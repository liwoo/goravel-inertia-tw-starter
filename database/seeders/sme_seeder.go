package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
)

type SmeSeeder struct{}

func (s *SmeSeeder) Signature() string {
	return "SmeSeeder"
}

func (s *SmeSeeder) Run() error {
	// Check if SMEs already exist - skip if we have data
	existingCount, err := facades.Orm().Query().Model(&models.Sme{}).Count()
	if err != nil {
		return fmt.Errorf("failed to count existing MSMEs: %v", err)
	}

	if existingCount > 0 {
		facades.Log().Infof("SmeSeeder: Skipping - MSMEs table already has %d records", existingCount)
		fmt.Printf("SmeSeeder: Skipping - MSMEs table already has %d records\n", existingCount)
		return nil
	}

	fmt.Println("Seeding MSMEs...")

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Get config options
	businessCategories := getConfigValues("Business Categories")
	sectors := getConfigValues("Sectors")
	improvementAspects := getConfigValues("Improvement Aspects")
	financingSources := getConfigValues("Financing")

	// Districts (Malawi districts)
	districts := []string{
		"Balaka", "Blantyre", "Chikwawa", "Chiradzulu", "Chitipa",
		"Dedza", "Dowa", "Karonga", "Kasungu", "Likoma",
		"Lilongwe", "Machinga", "Mangochi", "Mchinji", "Mulanje",
		"Mwanza", "Mzimba", "Neno", "Nkhata Bay", "Nkhotakota",
		"Nsanje", "Ntcheu", "Ntchisi", "Phalombe", "Rumphi",
		"Salima", "Thyolo", "Zomba",
	}

	// Education levels
	educationLevels := []string{
		"Primary", "Secondary", "Diploma", "Bachelor's Degree",
		"Master's Degree", "PhD", "Vocational Training", "No Formal Education",
	}

	// Malawian statuses
	malawianStatuses := []string{
		"Citizen", "Resident", "Work Permit Holder", "Non-Resident",
	}

	// Genders (must match database enum: MALE, FEMALE)
	genders := []string{"MALE", "FEMALE"}

	// Create 100 SMEs
	for i := 1; i <= 100; i++ {
		fmt.Printf("Creating MSME %d/100...\n", i)

		// Generate unique USME number
		usmeNumber := fmt.Sprintf("USME-%06d", i)

		// Generate SME data
		businessName := faker.Word() + " " + randomBusinessType()
		district := districts[rand.Intn(len(districts))]
		region := getRegionFromDistrict(district)

		// Select random improvement aspects (1-3)
		selectedImprovements := selectRandomStrings(improvementAspects, 1+rand.Intn(3))

		// Select random financing sources (1-3)
		selectedFinancing := selectRandomStrings(financingSources, 1+rand.Intn(3))

		// Create the SME first
		sme := models.Sme{
			UsmeNumber:                 usmeNumber,
			Name:                       businessName,
			RegistrationNumber:         stringPtr(faker.UUIDDigit()),
			TaxIdentificationNumber:    stringPtr(fmt.Sprintf("TIN%09d", rand.Intn(1000000000))),
			OperationalStartDate:       carbonPtr(randomDateBetween(2010, 2023)),
			BusinessCategory:           randomString(businessCategories),
			Sector:                     randomString(sectors),
			SubSector:                  stringPtr(faker.Word()),
			BusinessDescription:        stringPtr(faker.Sentence()),
			ContactPhone:               generateMalawiPhone(),
			ContactEmail:               faker.Email(),
			PhysicalAddress:            stringPtr(faker.GetRealAddress().Address),
			PostalAddress:              stringPtr(fmt.Sprintf("P.O. Box %d, %s", rand.Intn(10000), district)),
			Website:                    stringPtr(faker.URL()),
			Region:                     stringPtr(region),
			District:                   stringPtr(district),
			TraditionalAuthority:       stringPtr(faker.LastName() + " TA"),
			BusinessImprovementAspects: selectedImprovements,
			BusinessAccessedFinancing:  selectedFinancing,
		}

		if err := facades.Orm().Query().Create(&sme); err != nil {
			return fmt.Errorf("failed to create SME: %v", err)
		}

		// Create Primary Business Owner
		nationality := requests.AllNationalities()[rand.Intn(len(requests.AllNationalities()))]
		owner := models.PrimaryBusinessOwner{
			FirstName:              faker.FirstName(),
			LastName:               faker.LastName(),
			OtherNames:             stringPtr(faker.FirstName()),
			Nationality:            string(nationality),
			NationalIdNumber:       generateNationalID(),
			DateOfBirth:            randomDateBetween(1950, 2000),
			Gender:                 genders[rand.Intn(len(genders))],
			EducationLevel:         educationLevels[rand.Intn(len(educationLevels))],
			MalawianStatus:         malawianStatuses[rand.Intn(len(malawianStatuses))],
			HasSpecialNeeds:        rand.Float32() < 0.1, // 10% have special needs
			PhoneNumber:            generateMalawiPhone(),
			LandlineNumber:         randomStringPtr(faker.Phonenumber()),
			Email:                  stringPtr(faker.Email()),
			PhysicalAddress:        stringPtr(faker.GetRealAddress().Address),
			PostalAddress:          stringPtr(fmt.Sprintf("P.O. Box %d, %s", rand.Intn(10000), district)),
			Region:                 stringPtr(region),
			District:               stringPtr(district),
			TraditionalAuthority:   stringPtr(faker.LastName() + " TA"),
			AltContactName:         randomStringPtr(faker.Name()),
			AltContactRelationship: randomStringPtr(randomRelationship()),
			AltContactPhone:        randomStringPtr(generateMalawiPhone()),
			SmeID:                  int(sme.ID),
		}

		if err := facades.Orm().Query().Create(&owner); err != nil {
			return fmt.Errorf("failed to create Primary Business Owner: %v", err)
		}

		// Create 0-5 Additional Business Members
		numMembers := rand.Intn(6)
		for j := 0; j < numMembers; j++ {
			memberNationality := requests.AllNationalities()[rand.Intn(len(requests.AllNationalities()))]
			member := models.AdditionalBusinessMember{
				FirstName:        faker.FirstName(),
				LastName:         faker.LastName(),
				OtherNames:       randomStringPtr(faker.FirstName()),
				Nationality:      string(memberNationality),
				NationalIdNumber: generateNationalID(),
				DateOfBirth:      carbonPtrNullable(randomDateBetween(1960, 2005)),
				Email:            randomStringPtr(faker.Email()),
				PhoneNumber:      generateMalawiPhone(),
				IsIntern:         rand.Float32() < 0.2, // 20% are interns
				IsPartTime:       rand.Float32() < 0.3, // 30% are part-time
				SmeId:            int(sme.ID),
			}

			if err := facades.Orm().Query().Create(&member); err != nil {
				return fmt.Errorf("failed to create Additional Business Member: %v", err)
			}
		}
	}

	fmt.Println("MSME seeding completed successfully!")
	return nil
}

// Helper functions

func getConfigValues(configType string) []string {
	var configs []models.Config
	facades.Orm().Query().Where("config_type = ?", configType).Find(&configs)

	values := make([]string, 0)
	for _, config := range configs {
		values = append(values, config.Name)
	}

	// Fallback defaults if no configs exist
	if len(values) == 0 {
		switch configType {
		case "Business Categories":
			return []string{"Retail", "Manufacturing", "Services", "Agriculture", "Technology"}
		case "Sectors":
			return []string{"Private", "Public", "NGO", "Cooperative"}
		case "Improvement Aspects":
			return []string{"Marketing", "Finance", "Operations", "Technology", "Human Resources"}
		case "Financing":
			return []string{"Bank Loan", "Microfinance", "Personal Savings", "Government Grant", "Angel Investor"}
		}
	}

	return values
}

func randomBusinessType() string {
	types := []string{
		"Enterprises", "Trading", "Services", "Solutions", "Industries",
		"Group", "Company", "Corporation", "Ventures", "Holdings",
		"Shop", "Store", "Market", "Farm", "Works",
	}
	return types[rand.Intn(len(types))]
}

func getRegionFromDistrict(district string) string {
	// Malawi regions by district
	northernDistricts := map[string]bool{
		"Chitipa": true, "Karonga": true, "Likoma": true,
		"Mzimba": true, "Nkhata Bay": true, "Rumphi": true,
	}

	centralDistricts := map[string]bool{
		"Dedza": true, "Dowa": true, "Kasungu": true,
		"Lilongwe": true, "Mchinji": true, "Nkhotakota": true,
		"Ntcheu": true, "Ntchisi": true, "Salima": true,
	}

	if northernDistricts[district] {
		return "Northern Region"
	} else if centralDistricts[district] {
		return "Central Region"
	}
	return "Southern Region"
}

func generateMalawiPhone() string {
	// Malawi phone numbers: +265 followed by 9 digits
	prefixes := []string{"88", "99", "77", "84", "85"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	return fmt.Sprintf("+265%s%07d", prefix, rand.Intn(10000000))
}

func generateNationalID() string {
	// Format: 2 letters + 8 digits
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id := string(letters[rand.Intn(len(letters))]) + string(letters[rand.Intn(len(letters))])
	id += fmt.Sprintf("%08d", rand.Intn(100000000))
	return id
}

func randomRelationship() string {
	relationships := []string{
		"Spouse", "Sibling", "Parent", "Child",
		"Friend", "Colleague", "Business Partner",
	}
	return relationships[rand.Intn(len(relationships))]
}

func randomDateBetween(startYear, endYear int) carbon.DateTime {
	year := startYear + rand.Intn(endYear-startYear+1)
	month := 1 + rand.Intn(12)
	day := 1 + rand.Intn(28)
	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return *carbon.NewDateTime(carbon.Parse(dateStr))
}

func randomString(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	return slice[rand.Intn(len(slice))]
}

func selectRandomStrings(slice []string, count int) []string {
	if len(slice) == 0 {
		return []string{}
	}

	if count > len(slice) {
		count = len(slice)
	}

	// Create a copy and shuffle
	shuffled := make([]string, len(slice))
	copy(shuffled, slice)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:count]
}

func stringPtr(s string) *string {
	if rand.Float32() < 0.1 { // 10% chance of nil
		return nil
	}
	return &s
}

func randomStringPtr(s string) *string {
	if rand.Float32() < 0.3 { // 30% chance of nil
		return nil
	}
	return &s
}

func carbonPtr(dt carbon.DateTime) *carbon.DateTime {
	return &dt
}

func carbonPtrNullable(dt carbon.DateTime) *carbon.DateTime {
	if rand.Float32() < 0.2 { // 20% chance of nil
		return nil
	}
	return &dt
}
