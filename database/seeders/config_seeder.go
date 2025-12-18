package seeders

import (
	"starter-project/app/models"

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
	// strPtr := func(s string) *string { return &s }

	// Financing configurations
	// example := []models.Config{
	// 	{Name: "None", Code: strPtr("FIN_NONE"), ConfigType: "Financing"},
	// 	{Name: "Grants", Code: strPtr("FIN_GRANT"), ConfigType: "Financing"},
	// 	{Name: "Loan", Code: strPtr("FIN_LOAN"), ConfigType: "Financing"},
	// }

	// Combine all configs
	allConfigs := []models.Config{}
	// allConfigs = append(allConfigs, example...)

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
