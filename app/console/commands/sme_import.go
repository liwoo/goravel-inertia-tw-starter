package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// SmeImport is the artisan command for importing SME data from Excel files
type SmeImport struct {
}

// Signature The name and signature of the console command.
func (receiver *SmeImport) Signature() string {
	return "sme:import"
}

// Description The console command description.
func (receiver *SmeImport) Description() string {
	return "Import SME data from an Excel file"
}

// Extend The console command extend.
func (receiver *SmeImport) Extend() command.Extend {
	return command.Extend{
		Category: "sme",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to the Excel file to import",
				Required: true,
			},
			&command.IntFlag{
				Name:    "batch-size",
				Aliases: []string{"b"},
				Usage:   "Number of records to process per batch (default: 10)",
				Value:   10,
			},
			&command.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "Run without making database changes (for testing)",
			},
			&command.BoolFlag{
				Name:    "truncate",
				Aliases: []string{"t"},
				Usage:   "Delete ALL existing SMEs before import (use with caution!)",
			},
			&command.IntFlag{
				Name:    "start-row",
				Aliases: []string{"r"},
				Usage:   "Row number to start importing from (1-based, skips header). Use to resume interrupted imports.",
				Value:   1,
			},
			// Database connection flags
			&command.StringFlag{
				Name:    "dsn",
				Usage:   "Full database connection string (e.g., postgres://user:pass@host:port/dbname?sslmode=disable)",
			},
			&command.StringFlag{
				Name:    "db-host",
				Usage:   "Database host (default: uses app config)",
			},
			&command.IntFlag{
				Name:    "db-port",
				Usage:   "Database port (default: 5432)",
				Value:   5432,
			},
			&command.StringFlag{
				Name:    "db-user",
				Usage:   "Database user",
			},
			&command.StringFlag{
				Name:    "db-password",
				Usage:   "Database password",
			},
			&command.StringFlag{
				Name:    "db-name",
				Usage:   "Database name",
			},
			&command.StringFlag{
				Name:    "db-sslmode",
				Usage:   "Database SSL mode (disable, require, verify-ca, verify-full)",
				Value:   "disable",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *SmeImport) Handle(ctx console.Context) error {
	filePath := ctx.Option("file")
	batchSize := ctx.OptionInt("batch-size")
	dryRun := ctx.OptionBool("dry-run")
	truncate := ctx.OptionBool("truncate")
	startRow := ctx.OptionInt("start-row")

	// Database connection flags
	dsn := ctx.Option("dsn")
	dbHost := ctx.Option("db-host")
	dbPort := ctx.OptionInt("db-port")
	dbUser := ctx.Option("db-user")
	dbPassword := ctx.Option("db-password")
	dbName := ctx.Option("db-name")
	dbSSLMode := ctx.Option("db-sslmode")

	if filePath == "" {
		ctx.Error("File path is required. Use --file=/path/to/file.xlsx")
		return fmt.Errorf("file path is required")
	}

	if batchSize <= 0 {
		batchSize = 10
	}

	if startRow < 1 {
		startRow = 1
	}

	// Set up database connection
	var customDB *gorm.DB
	var err error

	if dsn != "" {
		// Use full connection string
		ctx.Info(fmt.Sprintf("Using custom database connection (DSN provided)"))
		customDB, err = createCustomDBConnection(dsn)
		if err != nil {
			ctx.Error(fmt.Sprintf("Failed to connect to custom database: %v", err))
			return err
		}
		defer closeCustomDB(customDB)
	} else if dbHost != "" {
		// Build connection string from individual parameters
		if dbUser == "" || dbName == "" {
			ctx.Error("When using --db-host, --db-user and --db-name are required")
			return fmt.Errorf("missing required database parameters")
		}

		builtDSN := buildDSN(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
		ctx.Info(fmt.Sprintf("Using custom database connection: %s@%s:%d/%s", dbUser, dbHost, dbPort, dbName))
		customDB, err = createCustomDBConnection(builtDSN)
		if err != nil {
			ctx.Error(fmt.Sprintf("Failed to connect to custom database: %v", err))
			return err
		}
		defer closeCustomDB(customDB)
	}

	ctx.Info(fmt.Sprintf("Starting SME import from: %s", filePath))
	ctx.Info(fmt.Sprintf("Batch size: %d", batchSize))
	if dryRun {
		ctx.Warning("DRY RUN MODE - No database changes will be made")
	}

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to open Excel file: %v", err))
		return err
	}
	defer f.Close()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		ctx.Error("No sheets found in the Excel file")
		return fmt.Errorf("no sheets found in the Excel file")
	}

	ctx.Info(fmt.Sprintf("Reading from sheet: %s", sheetName))

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to read rows: %v", err))
		return err
	}

	if len(rows) < 2 {
		ctx.Error("Excel file must have at least a header row and one data row")
		return fmt.Errorf("insufficient data in Excel file")
	}

	// Parse header row to get column indices
	headerRow := rows[0]
	columnMap := buildColumnMap(headerRow)

	totalDataRows := len(rows) - 1
	ctx.Info(fmt.Sprintf("Found %d data rows in file", totalDataRows))

	// Validate start row
	if startRow > totalDataRows {
		ctx.Error(fmt.Sprintf("Start row %d exceeds total data rows %d", startRow, totalDataRows))
		return fmt.Errorf("start row exceeds total data rows")
	}

	if startRow > 1 {
		ctx.Info(fmt.Sprintf("Starting from row %d (skipping first %d rows)", startRow, startRow-1))
	}

	// Delete existing SMEs only if --truncate flag is explicitly set
	if truncate {
		if startRow > 1 {
			ctx.Warning("Ignoring --truncate because --start-row is set (cannot truncate when resuming)")
		} else if !dryRun {
			ctx.Warning("TRUNCATE MODE: Deleting ALL existing SMEs and related data...")
			if err := deleteExistingSmes(ctx, customDB); err != nil {
				ctx.Error(fmt.Sprintf("Failed to delete existing SMEs: %v", err))
				return err
			}
			ctx.Success("Existing SMEs deleted successfully")
		} else {
			ctx.Warning("DRY RUN: Would delete all existing SMEs (--truncate flag set)")
		}
	} else {
		ctx.Info("Running in incremental mode (duplicates will be skipped)")
	}

	// Process rows in batches, starting from the specified row
	dataRows := rows[1:]
	totalRows := len(dataRows)
	successCount := 0
	errorCount := 0
	duplicateCount := 0
	skippedCount := startRow - 1
	var errors []string

	smeService := services.NewSmeService()

	// Start from the correct index (startRow is 1-based, array is 0-based)
	startIndex := startRow - 1
	rowsToProcess := totalRows - startIndex
	ctx.Info(fmt.Sprintf("Will process %d rows (from row %d to %d)", rowsToProcess, startRow, totalRows))

	for i := startIndex; i < totalRows; i += batchSize {
		end := i + batchSize
		if end > totalRows {
			end = totalRows
		}

		batch := dataRows[i:end]
		batchNum := ((i - startIndex) / batchSize) + 1
		ctx.Info(fmt.Sprintf("Processing batch %d (rows %d-%d of %d)...", batchNum, i+1, end, totalRows))

		for j, row := range batch {
			rowNum := i + j + 2 // +2 because Excel rows are 1-indexed and we skip header

			if dryRun {
				// Just validate the data in dry run mode
				_, err := parseRowData(row, columnMap, rowNum)
				if err != nil {
					errorCount++
					errMsg := fmt.Sprintf("Row %d: %v", rowNum, err)
					errors = append(errors, errMsg)
					ctx.Error(errMsg)
				} else {
					successCount++
					ctx.Info(fmt.Sprintf("Row %d: Would be imported successfully", rowNum))
				}
				continue
			}

			// Parse the row data
			smeData, err := parseRowData(row, columnMap, rowNum)
			if err != nil {
				errorCount++
				errMsg := fmt.Sprintf("Row %d: Parse error - %v", rowNum, err)
				errors = append(errors, errMsg)
				ctx.Error(errMsg)
				continue
			}

			// Check for duplicates before creating
			isDuplicate, matchReason := checkForDuplicate(smeData, customDB)
			if isDuplicate {
				duplicateCount++
				ctx.Warning(fmt.Sprintf("Row %d: Skipped (duplicate) - %s '%s' already exists (%s)",
					rowNum, "SME", smeData.BusinessName, matchReason))
				continue
			}

			// Create the SME with relationships in a transaction
			err = createSmeWithRelationships(ctx, smeData, smeService, customDB)
			if err != nil {
				errorCount++
				errMsg := fmt.Sprintf("Row %d: Create error - %v", rowNum, err)
				errors = append(errors, errMsg)
				ctx.Error(errMsg)
				continue
			}

			successCount++
			ctx.Success(fmt.Sprintf("Row %d: SME '%s' imported successfully", rowNum, smeData.BusinessName))
		}
	}

	// Print summary
	ctx.NewLine()
	ctx.Info("========== IMPORT SUMMARY ==========")
	ctx.Info(fmt.Sprintf("Total rows in file: %d", totalRows))
	if skippedCount > 0 {
		ctx.Info(fmt.Sprintf("Rows skipped (before start-row): %d", skippedCount))
	}
	ctx.Info(fmt.Sprintf("Rows processed: %d", rowsToProcess))
	ctx.Success(fmt.Sprintf("Successful imports: %d", successCount))
	if duplicateCount > 0 {
		ctx.Warning(fmt.Sprintf("Duplicates skipped: %d", duplicateCount))
	}
	if errorCount > 0 {
		ctx.Error(fmt.Sprintf("Failed imports: %d", errorCount))
		ctx.NewLine()
		ctx.Error("Errors encountered:")
		for _, errMsg := range errors {
			ctx.Error(fmt.Sprintf("  - %s", errMsg))
		}
	}
	ctx.Info("====================================")

	if dryRun {
		ctx.Warning("DRY RUN COMPLETE - No actual changes were made")
	}

	return nil
}

// SmeImportData holds parsed data from a single Excel row
type SmeImportData struct {
	// Primary Business Owner fields
	OwnerFullName           string
	OwnerFirstName          string
	OwnerLastName           string
	OwnerTitle              string
	OwnerNationality        string
	OwnerNationalID         string
	OwnerDateOfBirth        *carbon.DateTime
	OwnerGender             string
	OwnerEducation          string
	OwnerMalawianStatus     string
	OwnerHasSpecialNeeds    bool
	OwnerSpecialNeedsDesc   string
	OwnerPhone              string
	OwnerEmail              string
	OwnerPhysicalAddress    string
	OwnerPostalAddress      string
	OwnerRegion             string
	OwnerDistrict           string
	OwnerTA                 string
	OwnerNextOfKinName      string
	OwnerNextOfKinContact   string

	// SME fields
	BusinessName        string
	BusinessType        string
	IsRegistered        bool
	RegistrationNumber  string
	TaxID               string
	OperationalStartDate *carbon.DateTime
	Classification      string
	Sector              string
	SubSector           string
	BusinessDescription string
	ContactPhone        string
	ContactEmail        string
	PhysicalAddress     string
	PostalAddress       string
	Region              string
	District            string
	TraditionalAuthority string

	// Formalisation fields
	HasBankAccount      bool
	HasTaxClearance     bool
	IsRegisteredForVAT  bool
	IsMemberOfAssoc     bool
	IsAffiliated        bool
	HasExportLicense    bool
	HasAccessedBDS      bool
	AnnualTurnover      float64
	EstimatedAssets     float64

	// Additional fields
	AccessedFinancing   []string
	ImprovementAspects  []string
	OffersInternship    bool

	// Employee Summary fields
	FullTimeMales           int
	FullTimeFemales         int
	FullTimeWithContractMales   int
	FullTimeWithContractFemales int
	TemporaryMales          int
	TemporaryFemales        int
}

// buildColumnMap creates a mapping from column name to column index
func buildColumnMap(headerRow []string) map[string]int {
	columnMap := make(map[string]int)
	for i, col := range headerRow {
		// Normalize column names: uppercase and trim spaces
		normalizedCol := strings.ToUpper(strings.TrimSpace(col))
		columnMap[normalizedCol] = i
	}
	return columnMap
}

// getCellValue safely gets a cell value from a row
func getCellValue(row []string, columnMap map[string]int, columnName string) string {
	normalizedName := strings.ToUpper(strings.TrimSpace(columnName))
	if idx, ok := columnMap[normalizedName]; ok && idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

// parseYesNo converts YES/NO strings to boolean
func parseYesNo(value string) bool {
	upperVal := strings.ToUpper(strings.TrimSpace(value))
	return upperVal == "YES" || upperVal == "Y" || upperVal == "TRUE" || upperVal == "1"
}

// parseFloat parses a string to float64
func parseFloat(value string) float64 {
	// Remove commas and currency symbols
	cleaned := strings.ReplaceAll(value, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "MWK", "")
	cleaned = strings.ReplaceAll(cleaned, "MK", "")
	cleaned = strings.ReplaceAll(cleaned, "K", "")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return 0
	}

	f, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return f
}

// parseInt parses a string to int
func parseInt(value string) int {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return 0
	}

	i, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0
	}
	return i
}

// parseDate parses various date formats
func parseDate(value string) *carbon.DateTime {
	if value == "" {
		return nil
	}

	// Try multiple date formats
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2 Jan 2006",
	}

	for _, format := range formats {
		t, err := time.Parse(format, value)
		if err == nil {
			dateStr := t.Format("2006-01-02")
			dt := carbon.NewDateTime(carbon.Parse(dateStr))
			return dt
		}
	}

	return nil
}

// parseOwnerName parses full name into title, first name, and last name
func parseOwnerName(fullName string) (title, firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return "", "", ""
	}

	// Common titles to extract - ordered by length (longest first) to avoid
	// matching "Mr" when "Mrs" is intended (since "Mrs" starts with "Mr")
	titles := []string{
		"Prof.", "Prof",
		"Mrs.", "Mrs",
		"Miss",
		"Ms.", "Ms",
		"Mr.", "Mr",
		"Dr.", "Dr",
	}

	lowerFullName := strings.ToLower(fullName)
	for _, t := range titles {
		lowerTitle := strings.ToLower(t)
		if strings.HasPrefix(lowerFullName, lowerTitle) {
			// Make sure the title is followed by a space or end of string
			// to avoid matching "Mr" in "Mrambo" for example
			afterTitle := len(t)
			if afterTitle >= len(fullName) || fullName[afterTitle] == ' ' || fullName[afterTitle] == '.' {
				title = t
				fullName = strings.TrimSpace(fullName[afterTitle:])
				// Remove any leading period or space that might remain
				fullName = strings.TrimPrefix(fullName, ".")
				fullName = strings.TrimSpace(fullName)
				break
			}
		}
	}

	// Split remaining name into parts
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return title, "", ""
	} else if len(parts) == 1 {
		return title, parts[0], ""
	} else {
		firstName = parts[0]
		lastName = strings.Join(parts[1:], " ")
		return title, firstName, lastName
	}
}

// sanitizeName cleans up a name field, handling special characters common in Malawian names
func sanitizeName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Normalize multiple spaces to single space
	spaceRe := regexp.MustCompile(`\s+`)
	name = spaceRe.ReplaceAllString(name, " ")

	// Apostrophes in names like M'bwana, N'goma should be preserved
	// GORM/database drivers handle escaping via parameterized queries

	return name
}

// validateParsedName checks if a parsed name looks suspicious and returns a warning message
// Returns empty string if name looks valid
func validateParsedName(firstName, lastName, originalFullName string) string {
	// Check for suspiciously short first names (likely parsing errors)
	if len(firstName) == 1 && len(lastName) > 2 {
		// Single character first name with longer last name is suspicious
		// Likely a title parsing issue
		return fmt.Sprintf("Suspicious name parse: firstName='%s', lastName='%s' from '%s'",
			firstName, lastName, originalFullName)
	}

	// Check if first name looks like a leftover from title
	suspiciousFirstNames := []string{"s", "r", "rs", "iss", "rof"}
	for _, sus := range suspiciousFirstNames {
		if strings.ToLower(firstName) == sus {
			return fmt.Sprintf("Suspicious name parse (possible title remnant): firstName='%s', lastName='%s' from '%s'",
				firstName, lastName, originalFullName)
		}
	}

	return ""
}

// normalizeGender normalizes gender values to MALE/FEMALE
func normalizeGender(gender string) string {
	gender = strings.ToUpper(strings.TrimSpace(gender))
	if gender == "M" || gender == "MALE" || gender == "MAN" {
		return "MALE"
	}
	if gender == "F" || gender == "FEMALE" || gender == "WOMAN" {
		return "FEMALE"
	}
	return "MALE" // Default to MALE if unknown
}

// normalizeRegion normalizes region names
func normalizeRegion(region string) string {
	region = strings.TrimSpace(region)
	if region == "" {
		return ""
	}

	// Capitalize first letter of each word
	region = strings.ToLower(region)
	return strings.Title(region)
}

// parseRowData parses a single Excel row into SmeImportData
func parseRowData(row []string, columnMap map[string]int, rowNum int) (*SmeImportData, error) {
	data := &SmeImportData{}

	// Parse Primary Business Owner fields
	ownerFullName := sanitizeName(getCellValue(row, columnMap, "BUSINESS OWNER NAME"))
	data.OwnerFullName = ownerFullName
	data.OwnerTitle, data.OwnerFirstName, data.OwnerLastName = parseOwnerName(ownerFullName)

	// Sanitize the parsed names
	data.OwnerFirstName = sanitizeName(data.OwnerFirstName)
	data.OwnerLastName = sanitizeName(data.OwnerLastName)

	// Validate the parsed name and log warning if suspicious
	if warning := validateParsedName(data.OwnerFirstName, data.OwnerLastName, ownerFullName); warning != "" {
		facades.Log().Warning(fmt.Sprintf("Row %d: %s", rowNum, warning))
		// Attempt to fix common parsing issues
		if len(data.OwnerFirstName) <= 2 && len(data.OwnerLastName) > 2 {
			// Try re-parsing by treating the suspicious firstName as part of lastName
			parts := strings.Fields(ownerFullName)
			if len(parts) >= 2 {
				// Skip any detected title and use first real word as first name
				startIdx := 0
				if data.OwnerTitle != "" {
					// Find where the title ends in the parts
					for i, p := range parts {
						if strings.Contains(strings.ToLower(data.OwnerTitle), strings.ToLower(p)) ||
							strings.Contains(strings.ToLower(p), strings.ToLower(strings.TrimSuffix(data.OwnerTitle, "."))) {
							startIdx = i + 1
							break
						}
					}
				}
				if startIdx < len(parts) {
					data.OwnerFirstName = sanitizeName(parts[startIdx])
					if startIdx+1 < len(parts) {
						data.OwnerLastName = sanitizeName(strings.Join(parts[startIdx+1:], " "))
					} else {
						data.OwnerLastName = ""
					}
					facades.Log().Info(fmt.Sprintf("Row %d: Re-parsed name to firstName='%s', lastName='%s'",
						rowNum, data.OwnerFirstName, data.OwnerLastName))
				}
			}
		}
	}

	data.OwnerNationality = sanitizeName(getCellValue(row, columnMap, "NATIONALITY"))
	data.OwnerNationalID = strings.TrimSpace(getCellValue(row, columnMap, "NATIONAL ID NO"))
	data.OwnerDateOfBirth = parseDate(getCellValue(row, columnMap, "DATE OF BIRTH"))
	data.OwnerGender = normalizeGender(getCellValue(row, columnMap, "GENDER"))
	data.OwnerEducation = sanitizeName(getCellValue(row, columnMap, "EDUCATION"))
	data.OwnerMalawianStatus = sanitizeName(getCellValue(row, columnMap, "MALAWI STATUS"))
	data.OwnerHasSpecialNeeds = parseYesNo(getCellValue(row, columnMap, "HAS SPECIAL NEED?"))
	data.OwnerSpecialNeedsDesc = sanitizeName(getCellValue(row, columnMap, "SPECIAL NEEDS DESCRIPTION"))
	data.OwnerPhone = strings.TrimSpace(getCellValue(row, columnMap, "CONTACT"))
	data.OwnerEmail = strings.TrimSpace(getCellValue(row, columnMap, "EMAIL"))
	data.OwnerPhysicalAddress = sanitizeName(getCellValue(row, columnMap, "PHYSICAL ADDRESS"))
	data.OwnerPostalAddress = sanitizeName(getCellValue(row, columnMap, "POSTAL ADDRESS"))
	data.OwnerRegion = normalizeRegion(getCellValue(row, columnMap, "REGION"))
	data.OwnerDistrict = sanitizeName(getCellValue(row, columnMap, "DISTRICT"))
	data.OwnerTA = sanitizeName(getCellValue(row, columnMap, "TRADITIONAL AUTHORITY"))
	data.OwnerNextOfKinName = sanitizeName(getCellValue(row, columnMap, "NEXT OF KIN NAME"))
	data.OwnerNextOfKinContact = strings.TrimSpace(getCellValue(row, columnMap, "NEXT OF KIN CONTACT"))

	// Parse SME fields
	data.BusinessName = sanitizeName(getCellValue(row, columnMap, "BUSINESS NAME"))
	data.BusinessType = getCellValue(row, columnMap, "BUSINESS TYPE")
	data.IsRegistered = parseYesNo(getCellValue(row, columnMap, "REGISTERED"))
	data.RegistrationNumber = getCellValue(row, columnMap, "BUSINESS REGISTRATION NO")
	data.TaxID = getCellValue(row, columnMap, "MRA TIN")
	data.OperationalStartDate = parseDate(getCellValue(row, columnMap, "OPERATION START DATE"))
	data.Classification = normalizeClassification(getCellValue(row, columnMap, "CATEGORY"))
	data.Sector = getCellValue(row, columnMap, "SECTOR")
	data.SubSector = getCellValue(row, columnMap, "SUB SECTOR")
	data.BusinessDescription = getCellValue(row, columnMap, "SERVICES DESCRIPTION")
	data.ContactPhone = getCellValue(row, columnMap, "BUSINESS CONTACT")
	data.ContactEmail = getCellValue(row, columnMap, "BUSINESS EMAIL")
	data.PhysicalAddress = getCellValue(row, columnMap, "BUSINESS PHYSICAL ADDRESS")
	data.PostalAddress = getCellValue(row, columnMap, "BUSINESS POSTAL ADDRESS")
	data.Region = normalizeRegion(getCellValue(row, columnMap, "BUSINESS REGION"))
	data.District = getCellValue(row, columnMap, "BUSINESS DISTRICT")
	data.TraditionalAuthority = getCellValue(row, columnMap, "BUSINESS TRADITIONAL AUTHORITY")

	// Parse Formalisation fields
	data.HasBankAccount = parseYesNo(getCellValue(row, columnMap, "HAS BANK ACCOUNT?"))
	data.HasTaxClearance = parseYesNo(getCellValue(row, columnMap, "HAS TAX CLEARANCE CERTIFICATE?"))
	data.IsRegisteredForVAT = parseYesNo(getCellValue(row, columnMap, "IS REGISTERED FOR VAT?"))
	data.IsMemberOfAssoc = parseYesNo(getCellValue(row, columnMap, "IS ASSOCIATION MEMBER?"))
	data.IsAffiliated = parseYesNo(getCellValue(row, columnMap, "IS AFFLIATED?"))
	data.HasExportLicense = parseYesNo(getCellValue(row, columnMap, "DOES EXPORTS?"))
	data.HasAccessedBDS = parseYesNo(getCellValue(row, columnMap, "HAS ACCESS BUSINESS DEVELOPMENT SERVICES?"))
	data.AnnualTurnover = parseFloat(getCellValue(row, columnMap, "APPROX ANNUAL TURNOVER"))
	data.EstimatedAssets = parseFloat(getCellValue(row, columnMap, "VALUE ASSETS ESTIMATE"))

	// Parse additional fields
	accessedFinancing := getCellValue(row, columnMap, "ACCESSED FINANCING")
	if accessedFinancing != "" {
		// Split by comma or semicolon
		data.AccessedFinancing = splitAndTrim(accessedFinancing)
	} else {
		data.AccessedFinancing = []string{}
	}

	improvementAspects := getCellValue(row, columnMap, "IMPROVEMENT ASPECTS")
	otherImprovements := getCellValue(row, columnMap, "OTHER IMPROVEMENT ASPECTS")
	if improvementAspects != "" {
		data.ImprovementAspects = splitAndTrim(improvementAspects)
	} else {
		data.ImprovementAspects = []string{}
	}
	if otherImprovements != "" {
		data.ImprovementAspects = append(data.ImprovementAspects, splitAndTrim(otherImprovements)...)
	}

	data.OffersInternship = parseYesNo(getCellValue(row, columnMap, "OFFERS INTERNSHIP"))

	// Parse Employee Summary fields
	data.FullTimeMales = parseInt(getCellValue(row, columnMap, "FULL TIME PAID MALES"))
	data.FullTimeFemales = parseInt(getCellValue(row, columnMap, "FULL TIME PAID FEMALES"))
	data.FullTimeWithContractMales = parseInt(getCellValue(row, columnMap, "FULL TIME PAID MALES WITH CONTRACTS"))
	data.FullTimeWithContractFemales = parseInt(getCellValue(row, columnMap, "FULL TIME PAID FEMALES WITH CONTRACTS"))
	data.TemporaryMales = parseInt(getCellValue(row, columnMap, "TEMPORARY MALES"))
	data.TemporaryFemales = parseInt(getCellValue(row, columnMap, "TEMPORARY FEMALES"))

	// Validate required fields
	if data.BusinessName == "" {
		return nil, fmt.Errorf("business name is required")
	}

	// Use fallback values for required fields
	if data.ContactPhone == "" {
		data.ContactPhone = data.OwnerPhone
	}
	if data.ContactEmail == "" && data.OwnerEmail != "" {
		data.ContactEmail = data.OwnerEmail
	}
	if data.ContactEmail == "" {
		// Generate a placeholder email
		data.ContactEmail = generatePlaceholderEmail(data.BusinessName)
	}
	if data.Region == "" {
		data.Region = data.OwnerRegion
	}
	if data.District == "" {
		data.District = data.OwnerDistrict
	}

	// Default business type if not provided
	if data.BusinessType == "" {
		data.BusinessType = "Sole Proprietorship"
	}

	// Default sector if not provided
	if data.Sector == "" {
		data.Sector = "Other"
	}

	return data, nil
}

// splitAndTrim splits a string by comma or semicolon and trims each part
func splitAndTrim(s string) []string {
	// Replace semicolons with commas
	s = strings.ReplaceAll(s, ";", ",")
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// normalizeClassification normalizes classification values
func normalizeClassification(classification string) string {
	classification = strings.TrimSpace(classification)
	if classification == "" {
		return models.ClassificationUnclassified
	}

	lower := strings.ToLower(classification)
	if strings.Contains(lower, "micro") {
		return models.ClassificationMicro
	}
	if strings.Contains(lower, "small") {
		return models.ClassificationSmall
	}
	if strings.Contains(lower, "medium") {
		return models.ClassificationMedium
	}
	return models.ClassificationUnclassified
}

// generatePlaceholderEmail generates a placeholder email from business name
func generatePlaceholderEmail(businessName string) string {
	// Remove special characters and spaces
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	cleaned := re.ReplaceAllString(businessName, "")
	cleaned = strings.ToLower(cleaned)
	if len(cleaned) > 20 {
		cleaned = cleaned[:20]
	}
	// Use the constant from services package for consistency
	return fmt.Sprintf("%s%s", cleaned, services.PlaceholderEmailDomain)
}

// checkForDuplicate checks if an SME with matching criteria already exists in the database
// Returns (isDuplicate bool, matchReason string)
// Matching criteria (in order of priority):
// 1. National ID number (most unique identifier)
// 2. Business registration number (if registered)
// 3. Business name + Owner name combination
func checkForDuplicate(data *SmeImportData, customDB *gorm.DB) (bool, string) {
	var count int64

	// 1. Check by National ID (most reliable unique identifier)
	if data.OwnerNationalID != "" {
		if customDB != nil {
			customDB.Model(&models.PrimaryBusinessOwner{}).
				Where("national_id_number = ?", data.OwnerNationalID).
				Where("deleted_at IS NULL").
				Count(&count)
		} else {
			count, _ = facades.Orm().Query().
				Model(&models.PrimaryBusinessOwner{}).
				Where("national_id_number = ?", data.OwnerNationalID).
				Where("deleted_at IS NULL").
				Count()
		}
		if count > 0 {
			return true, fmt.Sprintf("National ID '%s' already registered", data.OwnerNationalID)
		}
	}

	// 2. Check by Business Registration Number (if registered)
	if data.RegistrationNumber != "" {
		if customDB != nil {
			customDB.Model(&models.Sme{}).
				Where("registration_number = ?", data.RegistrationNumber).
				Where("deleted_at IS NULL").
				Count(&count)
		} else {
			count, _ = facades.Orm().Query().
				Model(&models.Sme{}).
				Where("registration_number = ?", data.RegistrationNumber).
				Where("deleted_at IS NULL").
				Count()
		}
		if count > 0 {
			return true, fmt.Sprintf("Registration number '%s' already registered", data.RegistrationNumber)
		}
	}

	// 3. Check by Business Name + Owner Name combination
	// This catches cases where national ID or registration number might be missing
	if data.BusinessName != "" && (data.OwnerFirstName != "" || data.OwnerLastName != "") {
		ownerFirstName := strings.TrimSpace(data.OwnerFirstName)
		ownerLastName := strings.TrimSpace(data.OwnerLastName)
		businessName := strings.TrimSpace(data.BusinessName)

		if customDB != nil {
			customDB.Model(&models.Sme{}).
				Joins("INNER JOIN primary_business_owner ON primary_business_owner.sme_id = smes.id").
				Where("LOWER(smes.name) = LOWER(?)", businessName).
				Where("LOWER(primary_business_owner.first_name) = LOWER(?)", ownerFirstName).
				Where("LOWER(primary_business_owner.last_name) = LOWER(?)", ownerLastName).
				Where("smes.deleted_at IS NULL").
				Where("primary_business_owner.deleted_at IS NULL").
				Count(&count)
		} else {
			count, _ = facades.Orm().Query().
				Model(&models.Sme{}).
				Join("INNER JOIN primary_business_owner ON primary_business_owner.sme_id = smes.id").
				Where("LOWER(smes.name) = LOWER(?)", businessName).
				Where("LOWER(primary_business_owner.first_name) = LOWER(?)", ownerFirstName).
				Where("LOWER(primary_business_owner.last_name) = LOWER(?)", ownerLastName).
				Where("smes.deleted_at IS NULL").
				Where("primary_business_owner.deleted_at IS NULL").
				Count()
		}
		if count > 0 {
			return true, fmt.Sprintf("Business '%s' with owner '%s %s' already exists",
				businessName, ownerFirstName, ownerLastName)
		}
	}

	return false, ""
}

// deleteExistingSmes deletes all existing SMEs and related data
func deleteExistingSmes(ctx console.Context, customDB *gorm.DB) error {
	// Helper function to execute delete queries
	execDelete := func(query string) error {
		if customDB != nil {
			return customDB.Exec(query).Error
		}
		_, err := facades.Orm().Query().Exec(query)
		return err
	}

	// Delete in order due to foreign key constraints
	// 1. Delete applications (references smes via sme_id)
	if err := execDelete("DELETE FROM applications"); err != nil {
		return fmt.Errorf("failed to delete applications: %w", err)
	}
	ctx.Info("  Deleted applications records")

	// 2. Delete business_employee_summary
	if err := execDelete("DELETE FROM business_employee_summary"); err != nil {
		return fmt.Errorf("failed to delete business_employee_summary: %w", err)
	}
	ctx.Info("  Deleted business_employee_summary records")

	// 3. Delete business_formalisation
	if err := execDelete("DELETE FROM business_formalisation"); err != nil {
		return fmt.Errorf("failed to delete business_formalisation: %w", err)
	}
	ctx.Info("  Deleted business_formalisation records")

	// 4. Delete additional_business_members
	if err := execDelete("DELETE FROM additional_business_members"); err != nil {
		return fmt.Errorf("failed to delete additional_business_members: %w", err)
	}
	ctx.Info("  Deleted additional_business_members records")

	// 5. Delete primary_business_owner
	if err := execDelete("DELETE FROM primary_business_owner"); err != nil {
		return fmt.Errorf("failed to delete primary_business_owner: %w", err)
	}
	ctx.Info("  Deleted primary_business_owner records")

	// 6. Delete users associated with SMEs (SME portal users)
	if err := execDelete("DELETE FROM users WHERE id IN (SELECT user_id FROM user_sme)"); err != nil {
		// This might fail if there are no user_sme records, which is fine
		ctx.Warning(fmt.Sprintf("  Warning deleting SME users: %v", err))
	} else {
		ctx.Info("  Deleted SME users")
	}

	// 7. Delete user_sme pivot table
	if err := execDelete("DELETE FROM user_sme"); err != nil {
		ctx.Warning(fmt.Sprintf("  Warning deleting user_sme: %v", err))
	} else {
		ctx.Info("  Deleted user_sme records")
	}

	// 8. Delete smes
	if err := execDelete("DELETE FROM smes"); err != nil {
		return fmt.Errorf("failed to delete smes: %w", err)
	}
	ctx.Info("  Deleted smes records")

	return nil
}

// createSmeWithRelationships creates an SME with all its related records
func createSmeWithRelationships(ctx console.Context, data *SmeImportData, smeService *services.SmeService, customDB *gorm.DB) error {
	// Use either custom GORM DB or default Goravel ORM
	if customDB != nil {
		return createSmeWithRelationshipsGorm(ctx, data, customDB)
	}
	return createSmeWithRelationshipsDefault(ctx, data, smeService)
}

// createSmeWithRelationshipsGorm creates SME using custom GORM connection
func createSmeWithRelationshipsGorm(ctx console.Context, data *SmeImportData, db *gorm.DB) error {
	var smeID uint
	var bfID uint

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Create the SME record
		sme := &models.Sme{
			Name:               data.BusinessName,
			BusinessCategory:   data.BusinessType,
			Sector:             data.Sector,
			ContactPhone:       data.ContactPhone,
			ContactEmail:       data.ContactEmail,
			Classification:     data.Classification,
			IsActive:           true,
		}

		// Set optional fields
		if data.RegistrationNumber != "" {
			sme.RegistrationNumber = &data.RegistrationNumber
		}
		if data.TaxID != "" {
			sme.TaxIdentificationNumber = &data.TaxID
		}
		if data.OperationalStartDate != nil {
			sme.OperationalStartDate = data.OperationalStartDate
		}
		if data.SubSector != "" {
			sme.SubSector = &data.SubSector
		}
		if data.BusinessDescription != "" {
			sme.BusinessDescription = &data.BusinessDescription
		}
		if data.PhysicalAddress != "" {
			sme.PhysicalAddress = &data.PhysicalAddress
		}
		if data.PostalAddress != "" {
			sme.PostalAddress = &data.PostalAddress
		}
		if data.Region != "" {
			sme.Region = &data.Region
		}
		if data.District != "" {
			sme.District = &data.District
		}
		if data.TraditionalAuthority != "" {
			sme.TraditionalAuthority = &data.TraditionalAuthority
		}

		// Convert arrays to JSON for direct database insert
		if len(data.ImprovementAspects) > 0 {
			aspectsJSON, _ := json.Marshal(data.ImprovementAspects)
			sme.BusinessImprovementAspectJSON = string(aspectsJSON)
		} else {
			sme.BusinessImprovementAspectJSON = "[]"
		}
		if len(data.AccessedFinancing) > 0 {
			financingJSON, _ := json.Marshal(data.AccessedFinancing)
			sme.BusinessAccessedFinancingJSON = string(financingJSON)
		} else {
			sme.BusinessAccessedFinancingJSON = "[]"
		}

		// Generate USME number
		usmeData := map[string]interface{}{
			"district":          data.District,
			"region":            data.Region,
			"business_category": data.BusinessType,
		}
		usmeNumber, err := generateUsmeNumberForImportWithDB(usmeData, tx)
		if err != nil {
			return fmt.Errorf("failed to generate USME number: %w", err)
		}
		sme.UsmeNumber = usmeNumber

		if err := tx.Create(sme).Error; err != nil {
			return fmt.Errorf("failed to create SME: %w", err)
		}
		smeID = sme.ID

		// 2. Create Primary Business Owner
		pbo := buildPrimaryBusinessOwner(data, sme.ID)
		if err := tx.Create(pbo).Error; err != nil {
			return fmt.Errorf("failed to create primary business owner: %w", err)
		}

		// 3. Create Business Formalisation
		bf := &models.BusinessFormalisation{
			SmeID:                  int(sme.ID),
			HasBankAccount:         data.HasBankAccount,
			HasTaxClarification:    data.HasTaxClearance,
			IsRegisteredForVat:     data.IsRegisteredForVAT,
			IsMemberOfAssociation:  data.IsMemberOfAssoc,
			IsAffiliated:           data.IsAffiliated,
			HasExportLicense:       data.HasExportLicense,
			HasAccessedBds:         data.HasAccessedBDS,
			AnnualTurnover:         data.AnnualTurnover,
			EstimatedValueOfAssets: data.EstimatedAssets,
		}
		if err := tx.Create(bf).Error; err != nil {
			return fmt.Errorf("failed to create business formalisation: %w", err)
		}
		bfID = bf.ID

		// 4. Create Business Employee Summary
		bes := &models.BusinessEmployeeSummary{
			SmeID:                       int(sme.ID),
			FullTimeMales:               data.FullTimeMales,
			FullTimeFemales:             data.FullTimeFemales,
			FullTimeWithContractMales:   data.FullTimeWithContractMales,
			FullTimeWithContractFemales: data.FullTimeWithContractFemales,
			TemporaryMales:              data.TemporaryMales,
			TemporaryFemales:            data.TemporaryFemales,
		}
		if err := tx.Create(bes).Error; err != nil {
			return fmt.Errorf("failed to create business employee summary: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// After successful commit, calculate formalisation score and classification using GORM
	if calcErr := calculateFormalisationScoreWithDB(smeID, bfID, data, db); calcErr != nil {
		ctx.Warning(fmt.Sprintf("Warning: Failed to calculate formalisation score for SME %d: %v", smeID, calcErr))
	}

	if calcErr := calculateClassificationWithDB(smeID, data, db); calcErr != nil {
		ctx.Warning(fmt.Sprintf("Warning: Failed to calculate classification for SME %d: %v", smeID, calcErr))
	}

	return nil
}

// createSmeWithRelationshipsDefault creates SME using default Goravel ORM
func createSmeWithRelationshipsDefault(ctx console.Context, data *SmeImportData, smeService *services.SmeService) error {
	// Start a transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 1. Create the SME record
	sme := &models.Sme{
		Name:               data.BusinessName,
		BusinessCategory:   data.BusinessType,
		Sector:             data.Sector,
		ContactPhone:       data.ContactPhone,
		ContactEmail:       data.ContactEmail,
		Classification:     data.Classification,
		IsActive:           true,
	}

	// Set optional fields
	if data.RegistrationNumber != "" {
		sme.RegistrationNumber = &data.RegistrationNumber
	}
	if data.TaxID != "" {
		sme.TaxIdentificationNumber = &data.TaxID
	}
	if data.OperationalStartDate != nil {
		sme.OperationalStartDate = data.OperationalStartDate
	}
	if data.SubSector != "" {
		sme.SubSector = &data.SubSector
	}
	if data.BusinessDescription != "" {
		sme.BusinessDescription = &data.BusinessDescription
	}
	if data.PhysicalAddress != "" {
		sme.PhysicalAddress = &data.PhysicalAddress
	}
	if data.PostalAddress != "" {
		sme.PostalAddress = &data.PostalAddress
	}
	if data.Region != "" {
		sme.Region = &data.Region
	}
	if data.District != "" {
		sme.District = &data.District
	}
	if data.TraditionalAuthority != "" {
		sme.TraditionalAuthority = &data.TraditionalAuthority
	}

	// Set improvement aspects and financing (these will be converted to JSON by hooks)
	sme.BusinessImprovementAspects = data.ImprovementAspects
	sme.BusinessAccessedFinancing = data.AccessedFinancing

	// Convert arrays to JSON for direct database insert
	if len(data.ImprovementAspects) > 0 {
		aspectsJSON, _ := json.Marshal(data.ImprovementAspects)
		sme.BusinessImprovementAspectJSON = string(aspectsJSON)
	} else {
		sme.BusinessImprovementAspectJSON = "[]"
	}
	if len(data.AccessedFinancing) > 0 {
		financingJSON, _ := json.Marshal(data.AccessedFinancing)
		sme.BusinessAccessedFinancingJSON = string(financingJSON)
	} else {
		sme.BusinessAccessedFinancingJSON = "[]"
	}

	// Generate USME number using the service's logic
	usmeData := map[string]interface{}{
		"district":          data.District,
		"region":            data.Region,
		"business_category": data.BusinessType,
	}

	// We need to generate USME number - use the same logic as the service
	usmeNumber, err := generateUsmeNumberForImport(usmeData)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to generate USME number: %w", err)
	}
	sme.UsmeNumber = usmeNumber

	// Create SME record
	if err := tx.Create(sme); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create SME: %w", err)
	}

	// 2. Create Primary Business Owner
	pbo := buildPrimaryBusinessOwner(data, sme.ID)

	if err := tx.Create(pbo); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create primary business owner: %w", err)
	}

	// 3. Create Business Formalisation
	bf := &models.BusinessFormalisation{
		SmeID:                  int(sme.ID),
		HasBankAccount:         data.HasBankAccount,
		HasTaxClarification:    data.HasTaxClearance,
		IsRegisteredForVat:     data.IsRegisteredForVAT,
		IsMemberOfAssociation:  data.IsMemberOfAssoc,
		IsAffiliated:           data.IsAffiliated,
		HasExportLicense:       data.HasExportLicense,
		HasAccessedBds:         data.HasAccessedBDS,
		AnnualTurnover:         data.AnnualTurnover,
		EstimatedValueOfAssets: data.EstimatedAssets,
	}

	if err := tx.Create(bf); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create business formalisation: %w", err)
	}

	// 4. Create Business Employee Summary
	bes := &models.BusinessEmployeeSummary{
		SmeID:                       int(sme.ID),
		FullTimeMales:               data.FullTimeMales,
		FullTimeFemales:             data.FullTimeFemales,
		FullTimeWithContractMales:   data.FullTimeWithContractMales,
		FullTimeWithContractFemales: data.FullTimeWithContractFemales,
		TemporaryMales:              data.TemporaryMales,
		TemporaryFemales:            data.TemporaryFemales,
	}

	if err := tx.Create(bes); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create business employee summary: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// After successful commit, calculate formalisation score and classification
	// These use the service methods which handle their own queries
	_, err = smeService.CalculateFormalisationScore(sme.ID)
	if err != nil {
		ctx.Warning(fmt.Sprintf("Warning: Failed to calculate formalisation score for SME %d: %v", sme.ID, err))
	}

	_, err = smeService.CalculateClassification(sme.ID)
	if err != nil {
		ctx.Warning(fmt.Sprintf("Warning: Failed to calculate classification for SME %d: %v", sme.ID, err))
	}

	return nil
}

// buildPrimaryBusinessOwner creates a PrimaryBusinessOwner struct from import data
func buildPrimaryBusinessOwner(data *SmeImportData, smeID uint) *models.PrimaryBusinessOwner {
	pbo := &models.PrimaryBusinessOwner{
		SmeID:            int(smeID),
		FirstName:        data.OwnerFirstName,
		LastName:         data.OwnerLastName,
		Nationality:      data.OwnerNationality,
		NationalIdNumber: data.OwnerNationalID,
		Gender:           data.OwnerGender,
		EducationLevel:   data.OwnerEducation,
		MalawianStatus:   data.OwnerMalawianStatus,
		HasSpecialNeeds:  data.OwnerHasSpecialNeeds,
		PhoneNumber:      data.OwnerPhone,
	}

	// Set date of birth if provided
	if data.OwnerDateOfBirth != nil {
		pbo.DateOfBirth = *data.OwnerDateOfBirth
	} else {
		// Default to a reasonable date if not provided
		pbo.DateOfBirth = *carbon.NewDateTime(carbon.Parse("1990-01-01"))
	}

	// Set optional fields
	if data.OwnerSpecialNeedsDesc != "" {
		pbo.SpecialNeedsDescription = &data.OwnerSpecialNeedsDesc
	}
	if data.OwnerEmail != "" {
		pbo.Email = &data.OwnerEmail
	}
	if data.OwnerPhysicalAddress != "" {
		pbo.PhysicalAddress = &data.OwnerPhysicalAddress
	}
	if data.OwnerPostalAddress != "" {
		pbo.PostalAddress = &data.OwnerPostalAddress
	}
	if data.OwnerRegion != "" {
		pbo.Region = &data.OwnerRegion
	}
	if data.OwnerDistrict != "" {
		pbo.District = &data.OwnerDistrict
	}
	if data.OwnerTA != "" {
		pbo.TraditionalAuthority = &data.OwnerTA
	}
	if data.OwnerNextOfKinName != "" {
		pbo.AltContactName = &data.OwnerNextOfKinName
	}
	if data.OwnerNextOfKinContact != "" {
		pbo.AltContactPhone = &data.OwnerNextOfKinContact
	}

	// Default required fields if empty
	if pbo.FirstName == "" {
		pbo.FirstName = "Unknown"
	}
	if pbo.LastName == "" {
		pbo.LastName = "Owner"
	}
	if pbo.Nationality == "" {
		pbo.Nationality = "Malawian"
	}
	if pbo.NationalIdNumber == "" {
		pbo.NationalIdNumber = fmt.Sprintf("IMPORT-%d", smeID)
	}
	if pbo.EducationLevel == "" {
		pbo.EducationLevel = "Not Specified"
	}
	if pbo.MalawianStatus == "" {
		pbo.MalawianStatus = "citizen"
	}
	if pbo.PhoneNumber == "" {
		pbo.PhoneNumber = data.ContactPhone
	}

	return pbo
}

// generateUsmeNumberForImport generates a USME number for import
// This replicates the logic from sme_service.go generateUBI function
func generateUsmeNumberForImport(data map[string]interface{}) (string, error) {
	// Get current year
	year := time.Now().Year()

	// Get district code
	districtCode, err := getDistrictCodeForImport(data)
	if err != nil {
		return "", err
	}

	// Get category code
	categoryCode, err := getCategoryCodeForImport(data)
	if err != nil {
		return "", err
	}

	// Get next sequential number
	sequentialNumber, err := getNextSequentialNumberForImport(year, districtCode, categoryCode)
	if err != nil {
		return "", err
	}

	// Build UBI without check digit
	ubiWithoutCheck := fmt.Sprintf("MW-%04d-%s-%s-%06d", year, districtCode, categoryCode, sequentialNumber)

	// Calculate check digit
	checkDigit := calculateLuhnCheckDigitForImport(ubiWithoutCheck)

	// Return complete UBI
	return fmt.Sprintf("%s-%d", ubiWithoutCheck, checkDigit), nil
}

// getDistrictCodeForImport returns the 2-letter district code for a given district name
func getDistrictCodeForImport(data map[string]interface{}) (string, error) {
	districtCodes := map[string]string{
		// Northern Region
		"Chitipa":    "CT",
		"Karonga":    "KR",
		"Mzuzu":      "MZ",
		"Nkhata Bay": "NB",
		"Rumphi":     "RU",
		"Likoma":     "LK",
		"Mzimba":     "MH",

		// Central Region
		"Dedza":      "DE",
		"Dowa":       "DO",
		"Kasungu":    "KS",
		"Lilongwe":   "LI",
		"Mchinji":    "MC",
		"Nkhotakota": "NK",
		"Ntcheu":     "NU",
		"Ntchisi":    "NI",
		"Salima":     "SA",

		// Southern Region
		"Balaka":     "BA",
		"Blantyre":   "BT",
		"Chikwawa":   "CK",
		"Chiradzulu": "CR",
		"Machinga":   "MG",
		"Mangochi":   "MN",
		"Mulanje":    "MJ",
		"Mwanza":     "MW",
		"Nsanje":     "NS",
		"Thyolo":     "TH",
		"Phalombe":   "PH",
		"Zomba":      "ZO",
		"Neno":       "NE",
	}

	district, ok := data["district"]
	if !ok || district == nil {
		return "00", nil // Default district code if not provided
	}

	var districtStr string
	switch v := district.(type) {
	case string:
		districtStr = v
	case *string:
		if v != nil {
			districtStr = *v
		} else {
			return "00", nil
		}
	default:
		return "00", nil
	}

	// Try to find the district code
	code, exists := districtCodes[districtStr]
	if !exists {
		// Try case-insensitive match
		for k, v := range districtCodes {
			if strings.EqualFold(k, districtStr) {
				return v, nil
			}
		}
		return "00", nil // Default if district not found
	}

	return code, nil
}

// getCategoryCodeForImport returns the category code based on business category
func getCategoryCodeForImport(data map[string]interface{}) (string, error) {
	category, ok := data["business_category"]
	if !ok || category == nil {
		return "XX", nil // Default code
	}

	var categoryStr string
	switch v := category.(type) {
	case string:
		categoryStr = v
	case *string:
		if v != nil {
			categoryStr = *v
		} else {
			return "XX", nil
		}
	default:
		return "XX", nil
	}

	// Remove parentheses and their content
	re := regexp.MustCompile(`\s*\([^)]*\)\s*`)
	categoryStr = re.ReplaceAllString(categoryStr, " ")
	categoryStr = strings.TrimSpace(categoryStr)

	// Get first letters of words
	words := strings.Fields(categoryStr)

	var code string
	if len(words) >= 2 {
		code = strings.ToUpper(string(words[0][0]) + string(words[1][0]))
	} else if len(words) == 1 {
		if len(words[0]) >= 2 {
			code = strings.ToUpper(words[0][:2])
		} else {
			code = strings.ToUpper(string(words[0][0]) + "X")
		}
	} else {
		return "XX", nil
	}

	return code, nil
}

// getNextSequentialNumberForImport gets the next sequential number
func getNextSequentialNumberForImport(year int, districtCode, categoryCode string) (int, error) {
	pattern := fmt.Sprintf("MW-%04d-%s-%s-%%", year, districtCode, categoryCode)

	var sme models.Sme
	err := facades.Orm().Query().
		Where("usme_number LIKE ?", pattern).
		Order("usme_number DESC").
		First(&sme)

	if err != nil || sme.ID == 0 {
		return 1, nil
	}

	// Extract the sequential number from the usme_number
	parts := strings.Split(sme.UsmeNumber, "-")
	if len(parts) >= 5 {
		maxSequence, parseErr := strconv.Atoi(parts[4])
		if parseErr != nil {
			return 1, nil
		}
		return maxSequence + 1, nil
	}

	return 1, nil
}

// calculateLuhnCheckDigitForImport calculates the Luhn check digit
func calculateLuhnCheckDigitForImport(ubi string) int {
	digits := ""
	for _, char := range ubi {
		if char >= '0' && char <= '9' {
			digits += string(char)
		}
	}

	sum := 0
	alternate := false

	for i := len(digits) - 1; i >= 0; i-- {
		digit := int(digits[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

// createCustomDBConnection creates a GORM database connection from a DSN
func createCustomDBConnection(dsn string) (*gorm.DB, error) {
	// Parse the DSN to convert postgres:// URL format to key=value format if needed
	parsedDSN, err := parseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	db, err := gorm.Open(postgres.Open(parsedDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// parseDSN converts a postgres:// URL to a key=value DSN string
func parseDSN(dsn string) (string, error) {
	// If it's already in key=value format, return as-is
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		return dsn, nil
	}

	// Parse the URL
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	// Extract components
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "5432"
	}

	user := ""
	password := ""
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}

	dbname := strings.TrimPrefix(u.Path, "/")

	// Build key=value DSN
	parts := []string{
		fmt.Sprintf("host=%s", host),
		fmt.Sprintf("port=%s", port),
		fmt.Sprintf("user=%s", user),
		fmt.Sprintf("dbname=%s", dbname),
	}

	if password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", password))
	}

	// Add query parameters (like sslmode)
	for key, values := range u.Query() {
		if len(values) > 0 {
			parts = append(parts, fmt.Sprintf("%s=%s", key, values[0]))
		}
	}

	return strings.Join(parts, " "), nil
}

// closeCustomDB closes the custom database connection
func closeCustomDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

// buildDSN builds a PostgreSQL DSN from individual parameters
func buildDSN(host string, port int, user, password, dbname, sslmode string) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s", host, port, user, dbname)
	if password != "" {
		dsn += fmt.Sprintf(" password=%s", password)
	}
	if sslmode != "" {
		dsn += fmt.Sprintf(" sslmode=%s", sslmode)
	}
	return dsn
}

// generateUsmeNumberForImportWithDB generates a USME number using a custom GORM DB connection
func generateUsmeNumberForImportWithDB(data map[string]interface{}, db *gorm.DB) (string, error) {
	// Get current year
	year := time.Now().Year()

	// Get district code
	districtCode, err := getDistrictCodeForImport(data)
	if err != nil {
		return "", err
	}

	// Get category code
	categoryCode, err := getCategoryCodeForImport(data)
	if err != nil {
		return "", err
	}

	// Get next sequential number using the custom DB
	sequentialNumber, err := getNextSequentialNumberForImportWithDB(year, districtCode, categoryCode, db)
	if err != nil {
		return "", err
	}

	// Build UBI without check digit
	ubiWithoutCheck := fmt.Sprintf("MW-%04d-%s-%s-%06d", year, districtCode, categoryCode, sequentialNumber)

	// Calculate check digit
	checkDigit := calculateLuhnCheckDigitForImport(ubiWithoutCheck)

	// Return complete UBI
	return fmt.Sprintf("%s-%d", ubiWithoutCheck, checkDigit), nil
}

// getNextSequentialNumberForImportWithDB gets the next sequential number using custom GORM DB
func getNextSequentialNumberForImportWithDB(year int, districtCode, categoryCode string, db *gorm.DB) (int, error) {
	pattern := fmt.Sprintf("MW-%04d-%s-%s-%%", year, districtCode, categoryCode)

	var sme models.Sme
	err := db.Where("usme_number LIKE ?", pattern).
		Order("usme_number DESC").
		First(&sme).Error

	if err != nil || sme.ID == 0 {
		return 1, nil
	}

	// Extract the sequential number from the usme_number
	parts := strings.Split(sme.UsmeNumber, "-")
	if len(parts) >= 5 {
		maxSequence, parseErr := strconv.Atoi(parts[4])
		if parseErr != nil {
			return 1, nil
		}
		return maxSequence + 1, nil
	}

	return 1, nil
}

// calculateFormalisationScoreWithDB calculates the formalisation score using a GORM DB connection
// This mirrors the logic in SmeService.CalculateFormalisationScore but works with custom DB
func calculateFormalisationScoreWithDB(smeID uint, bfID uint, data *SmeImportData, db *gorm.DB) error {
	// Calculate score breakdown based on the import data
	// We have all the data in SmeImportData, so we can calculate directly

	// Compliance Score (70 points max) - from 7 boolean checkboxes
	complianceScore := 0
	if data.HasBankAccount {
		complianceScore += 10
	}
	if data.HasTaxClearance {
		complianceScore += 10
	}
	if data.IsRegisteredForVAT {
		complianceScore += 10
	}
	if data.IsMemberOfAssoc {
		complianceScore += 10
	}
	if data.IsAffiliated {
		complianceScore += 10
	}
	if data.HasExportLicense {
		complianceScore += 10
	}
	if data.HasAccessedBDS {
		complianceScore += 10
	}

	// Team Structure Score (20 points max)
	teamStructureScore := 0

	// Has primary business owner (we always create one): 5 points
	teamStructureScore += 5

	// We don't have additional members in import data, but check if we would have had any
	// For now, no additional members from import: 0 points for this

	// Calculate total team size
	totalTeamSize := 1 // Primary owner
	fullTimeEmployees := data.FullTimeMales + data.FullTimeFemales

	// Has full-time employees (1+): 5 points
	if fullTimeEmployees > 0 {
		teamStructureScore += 5
	}

	// Add all employees to team size
	totalTeamSize += data.FullTimeMales + data.FullTimeFemales
	totalTeamSize += data.TemporaryMales + data.TemporaryFemales

	// Team size > 5 people: 5 points
	if totalTeamSize > 5 {
		teamStructureScore += 5
	}

	// Financial Score (10 points max)
	financialScore := 0
	if data.AnnualTurnover > 0 {
		financialScore += 5
	}
	if data.EstimatedAssets > 0 {
		financialScore += 5
	}

	// Calculate total score
	totalScore := complianceScore + teamStructureScore + financialScore
	if totalScore > 100 {
		totalScore = 100
	}

	// Update the business_formalisation record with all score components
	return db.Model(&models.BusinessFormalisation{}).
		Where("id = ?", bfID).
		Updates(map[string]interface{}{
			"formalisation_score":  totalScore,
			"compliance_score":     complianceScore,
			"team_structure_score": teamStructureScore,
			"financial_score":      financialScore,
		}).Error
}

// calculateClassificationWithDB calculates and updates SME classification using a GORM DB connection
// This mirrors the logic in SmeService.CalculateClassification but works with custom DB
func calculateClassificationWithDB(smeID uint, data *SmeImportData, db *gorm.DB) error {
	// Calculate total employees from import data
	totalEmployees := 1 // Primary business owner

	// Add employees from BusinessEmployeeSummary data
	totalEmployees += data.FullTimeMales + data.FullTimeFemales +
		data.FullTimeWithContractMales + data.FullTimeWithContractFemales +
		data.TemporaryMales + data.TemporaryFemales

	// Get turnover and assets
	turnover := data.AnnualTurnover
	assets := data.EstimatedAssets

	// Determine classification using the same logic as DetermineClassification
	classification := determineClassificationForImport(totalEmployees, turnover, assets)

	// Update the SME record with the classification
	return db.Model(&models.Sme{}).
		Where("id = ?", smeID).
		Update("classification", classification).Error
}

// determineClassificationForImport applies the classification rules based on Malawi MSME Policy
// Classification is based on EITHER employees OR turnover meeting the criteria.
// This mirrors SmeService.DetermineClassification
func determineClassificationForImport(employees int, turnover, assets float64) string {
	// Classification thresholds (matching models package constants)
	const (
		// Micro thresholds
		microEmployeeMin = 1
		microEmployeeMax = 4
		microTurnoverMax = 5000000.0

		// Small thresholds
		smallEmployeeMin = 5
		smallEmployeeMax = 20
		smallTurnoverMin = 5000000.0
		smallTurnoverMax = 50000000.0

		// Medium thresholds
		mediumEmployeeMin = 21
		mediumEmployeeMax = 99
		mediumTurnoverMin = 50000000.0
		mediumTurnoverMax = 500000000.0
	)

	// Check Medium classification first (highest)
	// Employees: 21-99 OR Turnover: Above 50,000,000 - 500,000,000
	employeesMeetsMedium := employees >= mediumEmployeeMin && employees <= mediumEmployeeMax
	turnoverMeetsMedium := turnover > mediumTurnoverMin && turnover <= mediumTurnoverMax
	if employeesMeetsMedium || turnoverMeetsMedium {
		return models.ClassificationMedium
	}

	// Check Small classification
	// Employees: 5-20 OR Turnover: Above 5,000,000 - 50,000,000
	employeesMeetsSmall := employees >= smallEmployeeMin && employees <= smallEmployeeMax
	turnoverMeetsSmall := turnover > smallTurnoverMin && turnover <= smallTurnoverMax
	if employeesMeetsSmall || turnoverMeetsSmall {
		return models.ClassificationSmall
	}

	// Check Micro classification
	// Employees: 1-4 OR Turnover: Up to 5,000,000
	employeesMeetsMicro := employees >= microEmployeeMin && employees <= microEmployeeMax
	turnoverMeetsMicro := turnover > 0 && turnover <= microTurnoverMax
	if employeesMeetsMicro || turnoverMeetsMicro {
		return models.ClassificationMicro
	}

	// Default to Unclassified
	return models.ClassificationUnclassified
}
