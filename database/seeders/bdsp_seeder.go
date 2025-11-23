package seeders

import (
	"fmt"
	"math/rand"
	"smedi-sme-db/app/models"

	"github.com/go-faker/faker/v4"
	"github.com/goravel/framework/facades"
)

type BdspSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *BdspSeeder) Signature() string {
	return "BdspSeeder"
}

// Run executes the seeder logic.
func (s *BdspSeeder) Run() error {
	fmt.Println("Seeding BDSPs...")

	// Create 100 BDSPs
	status := models.RegistrationConfirmed
	for i := 1; i <= 100; i++ {
		fmt.Printf("Creating BDSP %d/100...\n", i)
		ubdspNumber := fmt.Sprintf("UBDSP-%06d", i)

		bdsp := models.Bdsp{
			UbdspNumber:        ubdspNumber,
			Name:               faker.Word(),
			PhysicalAddress:    stringPtr(faker.GetRealAddress().Address),
			PostalAddress:      stringPtr(fmt.Sprintf("P.O. Box %d, %s", rand.Intn(10000), faker.GetRealAddress().City)),
			RegistrationStatus: &status,
			Partners:           generateRandomWords(rand.Intn(5) + 1),
			ProductTypes:       generateRandomWords(rand.Intn(5) + 1),
			ServiceList:        generateRandomServices(rand.Intn(5) + 1),
		}

		if err := facades.Orm().Query().Create(&bdsp); err != nil {
			return fmt.Errorf("failed to create BDSP: %v", err)
		}
	}

	fmt.Println("BDSPs seeded successfully!")

	return nil
}

func generateRandomWords(count int) []string {
	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = faker.Word()
	}
	return words
}

func generateRandomServices(count int) []models.BdspService {
	services := make([]models.BdspService, count)
	for i := 0; i < count; i++ {
		services[i] = models.BdspService{
			Name:     faker.Word(),
			Cost:     float64(rand.Intn(100000) + 1),
			Duration: faker.Word(),
		}
	}
	return services
}
