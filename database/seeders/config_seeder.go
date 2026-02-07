package seeders

import (
	"books-database/app/models"
	"github.com/goravel/framework/facades"
)

type ConfigSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *ConfigSeeder) Signature() string {
	return "ConfigSeeder"
}

// Run executes the seeder logic.
// This seeder is idempotent - it will only insert configs that don't already exist.
// Existing configs (matched by Code) will be skipped to preserve any user modifications.
func (s *ConfigSeeder) Run() error {
	// Check if configs table already has data
	var existingCount int64
	existingCount, err := facades.Orm().Query().Model(&models.Config{}).Count()
	if err != nil {
		return err
	}

	// If we already have a significant number of configs, skip seeding entirely
	// This prevents re-seeding on every deployment while allowing initial seeding
	if existingCount > 50 {
		facades.Log().Infof("ConfigSeeder: Skipping - configs table already has data (%d records)", existingCount)
		return nil
	}

	// Helper function to create string pointer
	strPtr := func(s string) *string { return &s }

	// Financing configurations
	financings := []models.Config{
		{Name: "None", Code: strPtr("FIN_NONE"), ConfigType: "Financing"},
		{Name: "Grants", Code: strPtr("FIN_GRANT"), ConfigType: "Financing"},
		{Name: "Loan", Code: strPtr("FIN_LOAN"), ConfigType: "Financing"},
	}

	// Improvement Aspects configurations
	improvementAspects := []models.Config{
		{Name: "None", Code: strPtr("IMP_NONE"), ConfigType: "Improvement Aspects"},
		{Name: "Financing", Code: strPtr("IMP_FIN"), ConfigType: "Improvement Aspects"},
		{Name: "Human Resource Management", Code: strPtr("IMP_HRM"), ConfigType: "Improvement Aspects"},
		{Name: "Networking", Code: strPtr("IMP_NET"), ConfigType: "Improvement Aspects"},
		{Name: "Customer Care", Code: strPtr("IMP_CUST"), ConfigType: "Improvement Aspects"},
		{Name: "Exploring new markets", Code: strPtr("IMP_MKT"), ConfigType: "Improvement Aspects"},
		{Name: "Employing skilled workforce", Code: strPtr("IMP_WRKF"), ConfigType: "Improvement Aspects"},
		{Name: "Embarking on value addition", Code: strPtr("IMP_VAL"), ConfigType: "Improvement Aspects"},
		{Name: "Embracing IT in business", Code: strPtr("IMP_IT"), ConfigType: "Improvement Aspects"},
		{Name: "Operations and production management", Code: strPtr("IMP_OPS"), ConfigType: "Improvement Aspects"},
		{Name: "Quality management", Code: strPtr("IMP_QUAL"), ConfigType: "Improvement Aspects"},
		{Name: "Record keeping", Code: strPtr("IMP_REC"), ConfigType: "Improvement Aspects"},
		{Name: "Business management system", Code: strPtr("IMP_BMS"), ConfigType: "Improvement Aspects"},
	}

	// Business Categories configurations
	businessCategories := []models.Config{
		{Name: "Partnership (Joint Venture)", Code: strPtr("BC_PART"), ConfigType: "Business Categories"},
		{Name: "Primary Cooperative", Code: strPtr("BC_COOP1"), ConfigType: "Business Categories"},
		{Name: "Secondary Cooperative", Code: strPtr("BC_COOP2"), ConfigType: "Business Categories"},
		{Name: "Sole Proprietorship", Code: strPtr("BC_SOLE"), ConfigType: "Business Categories"},
		{Name: "Group Enterprise", Code: strPtr("BC_GROUP"), ConfigType: "Business Categories"},
		{Name: "Private Company", Code: strPtr("BC_PRIV"), ConfigType: "Business Categories"},
		{Name: "Association", Code: strPtr("BC_ASSOC"), ConfigType: "Business Categories"},
	}

	// Industries configurations
	industries := []models.Config{
		{Name: "Microfinance", Code: strPtr("IND_MFI"), ConfigType: "Industries"},
		{Name: "Business Counselling and Training", Code: strPtr("IND_BCT"), ConfigType: "Industries"},
		{Name: "Business Consultants", Code: strPtr("IND_CONS"), ConfigType: "Industries"},
		{Name: "Banks with SME Products", Code: strPtr("IND_BANK"), ConfigType: "Industries"},
		{Name: "Insurance Facilities", Code: strPtr("IND_INS"), ConfigType: "Industries"},
		{Name: "Market Linkages", Code: strPtr("IND_MKT"), ConfigType: "Industries"},
		{Name: "Industry Institutions", Code: strPtr("IND_INST"), ConfigType: "Industries"},
	}

	// Sectors configurations
	sectors := []models.Config{
		{Name: "Agriculture", Code: strPtr("SEC_AGR"), ConfigType: "Sectors"},
		{Name: "Manufacturing", Code: strPtr("SEC_MFG"), ConfigType: "Sectors"},
		{Name: "Services", Code: strPtr("SEC_SVC"), ConfigType: "Sectors"},
		{Name: "Trade", Code: strPtr("SEC_TRD"), ConfigType: "Sectors"},
		{Name: "Tourism and Hospitality", Code: strPtr("SEC_TRS"), ConfigType: "Sectors"},
		{Name: "Construction", Code: strPtr("SEC_CON"), ConfigType: "Sectors"},
		{Name: "Mining", Code: strPtr("SEC_MIN"), ConfigType: "Sectors"},
		{Name: "Energy", Code: strPtr("SEC_ENR"), ConfigType: "Sectors"},
		{Name: "Transport", Code: strPtr("SEC_TRN"), ConfigType: "Sectors"},
		{Name: "ICT and Technology", Code: strPtr("SEC_ICT"), ConfigType: "Sectors"},
	}

	// Registration Status configurations
	registrationStatuses := []models.Config{
		{Name: "Registered", Code: strPtr("REG_DONE"), ConfigType: "Registration Status"},
		{Name: "Unregistered", Code: strPtr("REG_NO"), ConfigType: "Registration Status"},
		{Name: "In Progress", Code: strPtr("REG_PROG"), ConfigType: "Registration Status"},
		{Name: "Pending", Code: strPtr("REG_PEND"), ConfigType: "Registration Status"},
	}

	// Development Partners configurations
	developmentPartners := []models.Config{
		{Name: "UNDP", Code: strPtr("DP_UNDP"), ConfigType: "Development Partners"},
		{Name: "GIZ", Code: strPtr("DP_GIZ"), ConfigType: "Development Partners"},
		{Name: "World Bank", Code: strPtr("DP_WB"), ConfigType: "Development Partners"},
		{Name: "African Development Bank", Code: strPtr("DP_AFDB"), ConfigType: "Development Partners"},
		{Name: "USAID", Code: strPtr("DP_USAID"), ConfigType: "Development Partners"},
		{Name: "EU", Code: strPtr("DP_EU"), ConfigType: "Development Partners"},
		{Name: "JICA", Code: strPtr("DP_JICA"), ConfigType: "Development Partners"},
	}

	// Combine all configs
	allConfigs := []models.Config{}
	allConfigs = append(allConfigs, financings...)
	allConfigs = append(allConfigs, improvementAspects...)
	allConfigs = append(allConfigs, businessCategories...)
	allConfigs = append(allConfigs, industries...)
	allConfigs = append(allConfigs, sectors...)
	allConfigs = append(allConfigs, registrationStatuses...)
	allConfigs = append(allConfigs, developmentPartners...)

	// Insert configs using upsert logic (skip existing)
	inserted := 0
	skipped := 0
	for _, config := range allConfigs {
		var existing models.Config
		err := facades.Orm().Query().Where("code = ?", *config.Code).First(&existing)
		if err == nil && existing.ID > 0 {
			// Config already exists, skip
			skipped++
			continue
		}
		// Config doesn't exist, create it
		if err := facades.Orm().Query().Create(&config); err != nil {
			return err
		}
		inserted++
	}

	facades.Log().Infof("ConfigSeeder: Inserted %d new configs, skipped %d existing", inserted, skipped)

	return nil
}
