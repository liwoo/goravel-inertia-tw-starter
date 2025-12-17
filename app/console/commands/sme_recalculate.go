package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// SmeRecalculate is the artisan command for recalculating SME classifications and formalisation scores
type SmeRecalculate struct {
	db *gorm.DB // Custom database connection (nil = use default)
}

// Signature The name and signature of the console command.
func (receiver *SmeRecalculate) Signature() string {
	return "sme:recalculate"
}

// Description The console command description.
func (receiver *SmeRecalculate) Description() string {
	return "Recalculate formalisation scores and classifications for all SMEs"
}

// Extend The console command extend.
func (receiver *SmeRecalculate) Extend() command.Extend {
	return command.Extend{
		Category: "sme",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "classification-only",
				Aliases: []string{"c"},
				Usage:   "Only recalculate classifications (skip formalisation scores)",
			},
			&command.BoolFlag{
				Name:    "score-only",
				Aliases: []string{"s"},
				Usage:   "Only recalculate formalisation scores (skip classifications)",
			},
			&command.IntFlag{
				Name:    "sme-id",
				Aliases: []string{"i"},
				Usage:   "Recalculate for a specific SME ID only",
			},
			&command.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "Show what would be recalculated without making changes",
			},
			// Remote database connection flags
			&command.StringFlag{
				Name:  "db-host",
				Usage: "Remote database host (e.g., db.example.com)",
			},
			&command.IntFlag{
				Name:  "db-port",
				Usage: "Remote database port (default: 5432)",
				Value: 5432,
			},
			&command.StringFlag{
				Name:  "db-name",
				Usage: "Remote database name",
			},
			&command.StringFlag{
				Name:  "db-user",
				Usage: "Remote database username",
			},
			&command.StringFlag{
				Name:  "db-password",
				Usage: "Remote database password",
			},
			&command.StringFlag{
				Name:  "db-sslmode",
				Usage: "SSL mode (disable, require, verify-ca, verify-full)",
				Value: "require",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *SmeRecalculate) Handle(ctx console.Context) error {
	classificationOnly := ctx.OptionBool("classification-only")
	scoreOnly := ctx.OptionBool("score-only")
	smeID := ctx.OptionInt("sme-id")
	dryRun := ctx.OptionBool("dry-run")

	// Remote database connection options
	dbHost := ctx.Option("db-host")
	dbPort := ctx.OptionInt("db-port")
	dbName := ctx.Option("db-name")
	dbUser := ctx.Option("db-user")
	dbPassword := ctx.Option("db-password")
	dbSSLMode := ctx.Option("db-sslmode")

	if classificationOnly && scoreOnly {
		ctx.Error("Cannot use --classification-only and --score-only together")
		return fmt.Errorf("conflicting options")
	}

	if dryRun {
		ctx.Warning("DRY RUN MODE - No changes will be made")
	}

	// Check if remote database connection is requested
	if dbHost != "" {
		if dbName == "" || dbUser == "" {
			ctx.Error("When using --db-host, you must also provide --db-name and --db-user")
			return fmt.Errorf("missing required database options")
		}

		ctx.Info(fmt.Sprintf("Connecting to remote database: %s@%s:%d/%s (sslmode=%s)", dbUser, dbHost, dbPort, dbName, dbSSLMode))

		db, err := receiver.connectToRemoteDB(dbHost, dbPort, dbName, dbUser, dbPassword, dbSSLMode)
		if err != nil {
			ctx.Error(fmt.Sprintf("Failed to connect to remote database: %v", err))
			return err
		}
		defer func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}()

		receiver.db = db
		ctx.Success("Connected to remote database successfully")
	}

	// Single SME mode
	if smeID > 0 {
		return receiver.recalculateSingle(ctx, uint(smeID), classificationOnly, scoreOnly, dryRun)
	}

	// All SMEs mode
	return receiver.recalculateAll(ctx, classificationOnly, scoreOnly, dryRun)
}

// connectToRemoteDB creates a direct GORM connection to a remote PostgreSQL database
func (receiver *SmeRecalculate) connectToRemoteDB(host string, port int, dbname, user, password, sslmode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getDB returns the appropriate database connection (remote or default)
func (receiver *SmeRecalculate) getDB() *gorm.DB {
	if receiver.db != nil {
		return receiver.db
	}
	// Use default Goravel ORM - need to get underlying GORM instance
	return nil
}

func (receiver *SmeRecalculate) recalculateSingle(ctx console.Context, smeID uint, classificationOnly, scoreOnly, dryRun bool) error {
	var sme models.Sme
	var err error

	if receiver.db != nil {
		// Use remote database
		err = receiver.db.
			Preload("BusinessFormalisation").
			Preload("BusinessEmployeeSummary").
			Preload("PrimaryBusinessOwner").
			Preload("AdditionalBusinessMembers").
			Where("id = ? AND deleted_at IS NULL", smeID).
			First(&sme).Error
	} else {
		// Use default Goravel ORM
		err = facades.Orm().Query().
			With("BusinessFormalisation").
			With("BusinessEmployeeSummary").
			With("PrimaryBusinessOwner").
			With("AdditionalBusinessMembers").
			Where("id = ?", smeID).
			First(&sme)
	}

	if err != nil || sme.ID == 0 {
		ctx.Error(fmt.Sprintf("SME with ID %d not found", smeID))
		return fmt.Errorf("SME not found")
	}

	ctx.Info(fmt.Sprintf("Processing SME: %s (ID: %d)", sme.Name, sme.ID))
	ctx.Info(fmt.Sprintf("  Current classification: %s", sme.Classification))

	if sme.BusinessFormalisation != nil {
		ctx.Info(fmt.Sprintf("  Current formalisation score: %d", sme.BusinessFormalisation.FormalisationScore))
	}

	if dryRun {
		ctx.Info("  [DRY RUN] Would recalculate scores and classification")
		return nil
	}

	var newScore int
	var newClassification string

	// Recalculate formalisation score
	if !classificationOnly {
		if receiver.db != nil {
			newScore, err = receiver.calculateFormalisationScoreRemote(smeID)
		} else {
			smeService := services.NewSmeService()
			newScore, err = smeService.CalculateFormalisationScore(smeID)
		}
		if err != nil {
			ctx.Error(fmt.Sprintf("  Failed to calculate formalisation score: %v", err))
		} else {
			ctx.Success(fmt.Sprintf("  New formalisation score: %d", newScore))
		}
	}

	// Recalculate classification
	if !scoreOnly {
		if receiver.db != nil {
			newClassification, err = receiver.calculateClassificationRemote(smeID)
		} else {
			smeService := services.NewSmeService()
			newClassification, err = smeService.CalculateClassification(smeID)
		}
		if err != nil {
			ctx.Error(fmt.Sprintf("  Failed to calculate classification: %v", err))
		} else {
			ctx.Success(fmt.Sprintf("  New classification: %s", newClassification))
		}
	}

	return nil
}

func (receiver *SmeRecalculate) recalculateAll(ctx console.Context, classificationOnly, scoreOnly, dryRun bool) error {
	var smes []models.Sme
	var err error

	if receiver.db != nil {
		err = receiver.db.
			Select("id", "name", "classification").
			Where("deleted_at IS NULL").
			Find(&smes).Error
	} else {
		err = facades.Orm().Query().
			Model(&models.Sme{}).
			Select("id", "name", "classification").
			Find(&smes)
	}

	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to fetch SMEs: %v", err))
		return err
	}

	totalCount := len(smes)
	ctx.Info(fmt.Sprintf("Found %d SMEs to process", totalCount))

	if dryRun {
		ctx.Info("[DRY RUN] Would recalculate for all SMEs")
		return nil
	}

	var scoreSuccessCount, scoreErrorCount int
	var classificationSuccessCount, classificationErrorCount int
	var classificationChanges []string

	smeService := services.NewSmeService()

	for i, sme := range smes {
		// Progress indicator every 100 SMEs
		if (i+1)%100 == 0 || i == 0 {
			ctx.Info(fmt.Sprintf("Processing %d/%d...", i+1, totalCount))
		}

		oldClassification := sme.Classification

		// Recalculate formalisation score
		if !classificationOnly {
			if receiver.db != nil {
				_, err = receiver.calculateFormalisationScoreRemote(sme.ID)
			} else {
				_, err = smeService.CalculateFormalisationScore(sme.ID)
			}
			if err != nil {
				scoreErrorCount++
			} else {
				scoreSuccessCount++
			}
		}

		// Recalculate classification
		if !scoreOnly {
			var newClassification string
			if receiver.db != nil {
				newClassification, err = receiver.calculateClassificationRemote(sme.ID)
			} else {
				newClassification, err = smeService.CalculateClassification(sme.ID)
			}
			if err != nil {
				classificationErrorCount++
			} else {
				classificationSuccessCount++
				if oldClassification != newClassification {
					classificationChanges = append(classificationChanges,
						fmt.Sprintf("  %s (ID: %d): %s -> %s", sme.Name, sme.ID, oldClassification, newClassification))
				}
			}
		}
	}

	// Print summary
	ctx.NewLine()
	ctx.Info("========== RECALCULATION SUMMARY ==========")
	ctx.Info(fmt.Sprintf("Total SMEs: %d", totalCount))

	if !classificationOnly {
		ctx.Success(fmt.Sprintf("Formalisation scores updated: %d", scoreSuccessCount))
		if scoreErrorCount > 0 {
			ctx.Error(fmt.Sprintf("Formalisation score errors: %d", scoreErrorCount))
		}
	}

	if !scoreOnly {
		ctx.Success(fmt.Sprintf("Classifications updated: %d", classificationSuccessCount))
		if classificationErrorCount > 0 {
			ctx.Error(fmt.Sprintf("Classification errors: %d", classificationErrorCount))
		}

		if len(classificationChanges) > 0 {
			ctx.NewLine()
			ctx.Info(fmt.Sprintf("Classification changes (%d):", len(classificationChanges)))
			for _, change := range classificationChanges {
				ctx.Info(change)
			}
		} else {
			ctx.Info("No classification changes detected")
		}
	}

	ctx.Info("============================================")

	return nil
}

// calculateFormalisationScoreRemote calculates formalisation score using direct GORM connection
func (receiver *SmeRecalculate) calculateFormalisationScoreRemote(smeID uint) (int, error) {
	var bf models.BusinessFormalisation
	err := receiver.db.Where("sme_id = ? AND deleted_at IS NULL", smeID).First(&bf).Error
	if err != nil {
		return 0, fmt.Errorf("failed to load business formalisation: %w", err)
	}

	// Calculate score using same logic as SmeService
	score := 0

	if bf.HasBankAccount {
		score += 15
	}
	if bf.HasTaxClarification {
		score += 20
	}
	if bf.IsRegisteredForVat {
		score += 15
	}
	if bf.IsMemberOfAssociation {
		score += 10
	}
	if bf.IsAffiliated {
		score += 10
	}
	if bf.HasExportLicense {
		score += 15
	}
	if bf.HasAccessedBds {
		score += 15
	}

	// Update the score
	err = receiver.db.Model(&models.BusinessFormalisation{}).
		Where("sme_id = ?", smeID).
		Update("formalisation_score", score).Error
	if err != nil {
		return score, fmt.Errorf("failed to update formalisation score: %w", err)
	}

	return score, nil
}

// calculateClassificationRemote calculates classification using direct GORM connection
func (receiver *SmeRecalculate) calculateClassificationRemote(smeID uint) (string, error) {
	var sme models.Sme
	err := receiver.db.
		Preload("BusinessFormalisation").
		Preload("BusinessEmployeeSummary").
		Preload("PrimaryBusinessOwner").
		Preload("AdditionalBusinessMembers").
		Where("id = ? AND deleted_at IS NULL", smeID).
		First(&sme).Error

	if err != nil {
		return "", fmt.Errorf("failed to load SME: %w", err)
	}

	// Calculate total employees
	totalEmployees := 0

	if sme.PrimaryBusinessOwner != nil && sme.PrimaryBusinessOwner.ID != 0 {
		totalEmployees += 1
	}

	totalEmployees += len(sme.AdditionalBusinessMembers)

	if sme.BusinessEmployeeSummary != nil {
		bes := sme.BusinessEmployeeSummary
		totalEmployees += bes.FullTimeMales + bes.FullTimeFemales +
			bes.PartTimeMales + bes.PartTimeFemales +
			bes.InternMales + bes.InternFemales +
			bes.FullTimeWithContractMales + bes.FullTimeWithContractFemales +
			bes.TemporaryMales + bes.TemporaryFemales
	}

	// Get turnover and assets
	var turnover, assets float64
	if sme.BusinessFormalisation != nil {
		turnover = sme.BusinessFormalisation.AnnualTurnover
		assets = sme.BusinessFormalisation.EstimatedValueOfAssets
	}

	// Determine classification using same logic as SmeService
	classification := determineClassification(totalEmployees, turnover, assets)

	// Update the classification
	err = receiver.db.Model(&models.Sme{}).
		Where("id = ?", smeID).
		Update("classification", classification).Error
	if err != nil {
		return classification, fmt.Errorf("failed to update classification: %w", err)
	}

	return classification, nil
}

// determineClassification applies the classification rules based on Malawi MSME Policy
func determineClassification(employees int, turnover, assets float64) string {
	// Medium: 21-99 employees
	if employees >= models.MediumEmployeeMin && employees <= models.MediumEmployeeMax {
		turnoverMet := turnover > models.MediumTurnoverMin && turnover <= models.MediumTurnoverMax
		assetsMet := assets <= models.MediumAssetsMax && assets > 0
		if turnoverMet || assetsMet {
			return models.ClassificationMedium
		}
	}

	// Small: 5-20 employees
	if employees >= models.SmallEmployeeMin && employees <= models.SmallEmployeeMax {
		turnoverMet := turnover > models.SmallTurnoverMin && turnover <= models.SmallTurnoverMax
		assetsMet := assets <= models.SmallAssetsMax && assets > 0
		if turnoverMet || assetsMet {
			return models.ClassificationSmall
		}
	}

	// Micro: 1-4 employees
	if employees >= models.MicroEmployeeMin && employees <= models.MicroEmployeeMax {
		turnoverMet := turnover > 0 && turnover <= models.MicroTurnoverMax
		assetsMet := assets > 0 && assets <= models.MicroAssetsMax
		if turnoverMet || assetsMet {
			return models.ClassificationMicro
		}
	}

	return models.ClassificationUnclassified
}
