package seeders

import (
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/models"
)

type ConfigSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *ConfigSeeder) Signature() string {
	return "ConfigSeeder"
}

// Run executes the seeder logic.
func (s *ConfigSeeder) Run() error {
	// Financing configurations
	financings := []models.Config{
		{Name: "None", ConfigType: "Financing"},
		{Name: "Grants", ConfigType: "Financing"},
		{Name: "Loan", ConfigType: "Financing"},
	}

	// Improvement Aspects configurations
	improvementAspects := []models.Config{
		{Name: "None", ConfigType: "Improvement Aspects"},
		{Name: "Financing", ConfigType: "Improvement Aspects"},
		{Name: "Human Resource Management", ConfigType: "Improvement Aspects"},
		{Name: "Networking", ConfigType: "Improvement Aspects"},
		{Name: "Customer Care", ConfigType: "Improvement Aspects"},
		{Name: "Exploring new markets", ConfigType: "Improvement Aspects"},
		{Name: "Employing skilled workforce", ConfigType: "Improvement Aspects"},
		{Name: "Embarking on value addition", ConfigType: "Improvement Aspects"},
		{Name: "Embracing IT in business", ConfigType: "Improvement Aspects"},
		{Name: "Operations and production management", ConfigType: "Improvement Aspects"},
		{Name: "Quality management", ConfigType: "Improvement Aspects"},
		{Name: "Record keeping", ConfigType: "Improvement Aspects"},
		{Name: "Business management system", ConfigType: "Improvement Aspects"},
	}

	// Business Categories configurations
	businessCategories := []models.Config{
		{Name: "Partnership (Joint Venture)", ConfigType: "Business Categories"},
		{Name: "Primary Cooperative", ConfigType: "Business Categories"},
		{Name: "Secondary Cooperative", ConfigType: "Business Categories"},
		{Name: "Sole Proprietorship", ConfigType: "Business Categories"},
		{Name: "Group SME", ConfigType: "Business Categories"},
		{Name: "Private Company", ConfigType: "Business Categories"},
		{Name: "Association", ConfigType: "Business Categories"},
	}

	// Industries configurations
	industries := []models.Config{
		{Name: "Microfinance", ConfigType: "Industries"},
		{Name: "Business Counselling and Training", ConfigType: "Industries"},
		{Name: "MSMEs Consultants", ConfigType: "Industries"},
		{Name: "Banks with SME Products", ConfigType: "Industries"},
		{Name: "SME Insurance Facilities", ConfigType: "Industries"},
		{Name: "Market Linkages", ConfigType: "Industries"},
		{Name: "Institutions that deal with MSMEs", ConfigType: "Industries"},
	}

	// Sectors configurations (Sensible sectors in Malawi)
	sectors := []models.Config{
		{Name: "Agriculture", ConfigType: "Sectors"},
		{Name: "Manufacturing", ConfigType: "Sectors"},
		{Name: "Services", ConfigType: "Sectors"},
		{Name: "Trade", ConfigType: "Sectors"},
		{Name: "Tourism and Hospitality", ConfigType: "Sectors"},
		{Name: "Construction", ConfigType: "Sectors"},
		{Name: "Mining", ConfigType: "Sectors"},
		{Name: "Energy", ConfigType: "Sectors"},
		{Name: "Transport", ConfigType: "Sectors"},
		{Name: "ICT and Technology", ConfigType: "Sectors"},
	}

	// Registration Status configurations
	registrationStatuses := []models.Config{
		{Name: "Registered", ConfigType: "Registration Status"},
		{Name: "Unregistered", ConfigType: "Registration Status"},
		{Name: "In Progress", ConfigType: "Registration Status"},
		{Name: "Pending", ConfigType: "Registration Status"},
	}

	// Development Partners configurations
	developmentPartners := []models.Config{
		{Name: "UNDP", ConfigType: "Development Partners"},
		{Name: "GIZ", ConfigType: "Development Partners"},
		{Name: "World Bank", ConfigType: "Development Partners"},
		{Name: "African Development Bank", ConfigType: "Development Partners"},
		{Name: "USAID", ConfigType: "Development Partners"},
		{Name: "EU", ConfigType: "Development Partners"},
		{Name: "JICA", ConfigType: "Development Partners"},
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

	// Insert all configs into database
	for _, config := range allConfigs {
		if err := facades.Orm().Query().Create(&config); err != nil {
			return err
		}
	}

	return nil
}
