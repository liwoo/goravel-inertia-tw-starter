package feature

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
)

type SmeSuite struct {
	suite.Suite
	tests.TestCase
}

func TestSmeSuite(t *testing.T) {
	suite.Run(t, new(SmeSuite))
}

// SetupSuite runs once before all tests in the suite
func (s *SmeSuite) SetupSuite() {
	// Force a complete database reset before the suite starts
	s.RefreshDatabase()
}

// SetupTest will run before each test in the suite.
func (s *SmeSuite) SetupTest() {
	// Clear all SME-related tables to ensure clean state
	// This is more reliable than relying on RefreshDatabase between individual tests
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM business_employee_summary")
		orm.Query().Exec("DELETE FROM business_formalisation")
		orm.Query().Exec("DELETE FROM additional_business_members")
		orm.Query().Exec("DELETE FROM primary_business_owner")
		orm.Query().Exec("DELETE FROM smes")
	}
}

// TearDownTest will run after each test in the suite.
func (s *SmeSuite) TearDownTest() {
}

// TestCreateSmeWithAllRelationships is a comprehensive sanity test that creates an SME
// with all its relationships and verifies the data integrity
func (s *SmeSuite) TestCreateSmeWithAllRelationships() {
	// Prepare test data using the new Region/District types
	operationalDate := carbon.NewDateTime(carbon.Parse("2020-01-15"))
	region := string(requests.RegionCentral)
	district := string(requests.DistrictLilongwe)
	website := "https://test-sme.com"
	physicalAddress := "123 Main Street, Area 3, Lilongwe"
	postalAddress := "P.O. Box 12345, Lilongwe"
	tradAuth := "TA Chadza"
	regNumber := "REG-2020-001"
	tinNumber := "TIN-987654321"
	subSector := "Retail"
	businessDesc := "A test business selling agricultural products"

	// Step 1: Create the main SME record
	// Note: The BeforeSave hook will automatically convert the arrays to JSON
	sme := &models.Sme{
		UsmeNumber:                 "USME-2024-00001",
		Name:                       "Test SME Enterprise",
		RegistrationNumber:         &regNumber,
		TaxIdentificationNumber:    &tinNumber,
		OperationalStartDate:       operationalDate,
		BusinessCategory:           "Agriculture",
		Sector:                     "Agribusiness",
		SubSector:                  &subSector,
		BusinessDescription:        &businessDesc,
		ContactPhone:               "+265999123456",
		ContactEmail:               "info@testsme.com",
		PhysicalAddress:            &physicalAddress,
		PostalAddress:              &postalAddress,
		Website:                    &website,
		Region:                     &region,
		District:                   &district,
		TraditionalAuthority:       &tradAuth,
		BusinessImprovementAspects: []string{"Marketing", "Financial Management", "Product Quality"},
		BusinessAccessedFinancing:  []string{"Bank Loan", "Microfinance"},
	}

	err := facades.Orm().Query().Create(&sme)
	s.NoError(err, "SME creation should succeed")
	s.NotZero(sme.ID, "SME ID should be set after creation")
	s.Equal("USME-2024-00001", sme.UsmeNumber)
	s.Equal("Test SME Enterprise", sme.Name)

	// Verify region/district using our new types
	mappedRegion, exists := requests.GetRegionByDistrict(requests.DistrictLilongwe)
	s.True(exists, "Lilongwe should be a valid district")
	s.Equal(requests.RegionCentral, mappedRegion, "Lilongwe should be in Central Region")

	// Step 2: Create Primary Business Owner
	ownerDob := carbon.NewDateTime(carbon.Parse("1985-03-20"))
	ownerOtherNames := "John"
	ownerEmail := "owner@testsme.com"
	ownerLandline := "+2651234567"
	ownerRegion := string(requests.RegionCentral)
	ownerDistrict := string(requests.DistrictLilongwe)
	ownerTA := "TA Chadza"
	altContactName := "Jane Doe"
	altContactRel := "Spouse"
	altContactPhone := "+265888777666"
	ownerPhysicalAddr := "456 Owner Street, Area 10, Lilongwe"
	ownerPostalAddr := "P.O. Box 54321, Lilongwe"

	primaryOwner := &models.PrimaryBusinessOwner{
		FirstName:              "Michael",
		LastName:               "Banda",
		OtherNames:             &ownerOtherNames,
		Nationality:            "Malawian",
		NationalIdNumber:       "MN123456789",
		DateOfBirth:            *ownerDob,
		Gender:                 "MALE",
		EducationLevel:         "Tertiary",
		MalawianStatus:         "Citizen",
		HasSpecialNeeds:        false,
		PhoneNumber:            "+265999111222",
		LandlineNumber:         &ownerLandline,
		Email:                  &ownerEmail,
		PhysicalAddress:        &ownerPhysicalAddr,
		PostalAddress:          &ownerPostalAddr,
		Region:                 &ownerRegion,
		District:               &ownerDistrict,
		TraditionalAuthority:   &ownerTA,
		AltContactName:         &altContactName,
		AltContactRelationship: &altContactRel,
		AltContactPhone:        &altContactPhone,
		SmeID:                  int(sme.ID),
	}

	err = facades.Orm().Query().Create(&primaryOwner)
	s.NoError(err, "Primary business owner creation should succeed")
	s.NotZero(primaryOwner.ID, "Primary owner ID should be set")
	s.Equal("Michael", primaryOwner.FirstName)
	s.Equal("Banda", primaryOwner.LastName)
	s.Equal(int(sme.ID), primaryOwner.SmeID)

	// Step 3: Create Additional Business Members
	member1Dob := carbon.NewDateTime(carbon.Parse("1990-05-10"))
	member1Email := "member1@testsme.com"
	member1OtherNames := "Grace"

	member1 := &models.AdditionalBusinessMember{
		FirstName:        "Sarah",
		LastName:         "Phiri",
		OtherNames:       &member1OtherNames,
		Nationality:      "Malawian",
		NationalIdNumber: "MN987654321",
		DateOfBirth:      member1Dob,
		Email:            &member1Email,
		PhoneNumber:      "+265888222333",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(sme.ID),
	}

	member2Dob := carbon.NewDateTime(carbon.Parse("1995-08-25"))
	member2Email := "member2@testsme.com"
	member2OtherNames := "Peter"

	member2 := &models.AdditionalBusinessMember{
		FirstName:        "James",
		LastName:         "Mwale",
		OtherNames:       &member2OtherNames,
		Nationality:      "Malawian",
		NationalIdNumber: "MN456789123",
		DateOfBirth:      member2Dob,
		Email:            &member2Email,
		PhoneNumber:      "+265999333444",
		IsIntern:         true,
		IsPartTime:       false,
		SmeId:            int(sme.ID),
	}

	err = facades.Orm().Query().Create(&member1)
	s.NoError(err, "Additional member 1 creation should succeed")
	s.NotZero(member1.ID)

	err = facades.Orm().Query().Create(&member2)
	s.NoError(err, "Additional member 2 creation should succeed")
	s.NotZero(member2.ID)
	s.True(member2.IsIntern, "Member 2 should be marked as intern")

	// Step 4: Create Business Formalisation
	formalisation := &models.BusinessFormalisation{
		SmeID:                  int(sme.ID),
		HasBankAccount:         true,
		HasTaxClarification:    true,
		IsRegisteredForVat:     false,
		IsMemberOfAssociation:  true,
		IsAffiliated:           false,
		HasExportLicense:       false,
		HasAccessedBds:         true,
		AnnualTurnover:         250000.00,
		EstimatedValueOfAssets: 500000.00,
		FormalisationScore:     75,
	}

	err = facades.Orm().Query().Create(&formalisation)
	s.NoError(err, "Business formalisation creation should succeed")
	s.NotZero(formalisation.ID)
	s.Equal(75, formalisation.FormalisationScore)
	s.True(formalisation.HasBankAccount)

	// Step 5: Create Business Employee Summary
	employeeSummary := &models.BusinessEmployeeSummary{
		SmeID:           int(sme.ID),
		FullTimeMales:   5,
		FullTimeFemales: 3,
		PartTimeMales:   2,
		PartTimeFemales: 4,
		InternMales:     1,
		InternFemales:   2,
	}

	err = facades.Orm().Query().Create(&employeeSummary)
	s.NoError(err, "Business employee summary creation should succeed")
	s.NotZero(employeeSummary.ID)
	s.Equal(5, employeeSummary.FullTimeMales)
	s.Equal(3, employeeSummary.FullTimeFemales)

	// Step 6: Verify all relationships by retrieving the SME with related data
	var retrievedSme models.Sme
	err = facades.Orm().Query().
		With("PrimaryBusinessOwner").
		With("AdditionalBusinessMembers").
		With("BusinessFormalisation").
		With("BusinessEmployeeSummary").
		Where("id = ?", sme.ID).
		First(&retrievedSme)

	s.NoError(err, "Should retrieve SME with all relationships")
	s.Equal(sme.ID, retrievedSme.ID)
	s.Equal("Test SME Enterprise", retrievedSme.Name)

	// Verify Primary Business Owner relationship
	s.NotNil(retrievedSme.PrimaryBusinessOwner, "Primary business owner should be loaded")
	if retrievedSme.PrimaryBusinessOwner != nil {
		s.Equal("Michael", retrievedSme.PrimaryBusinessOwner.FirstName)
		s.Equal("Banda", retrievedSme.PrimaryBusinessOwner.LastName)
	}

	// Verify Additional Business Members relationship
	s.NotNil(retrievedSme.AdditionalBusinessMembers, "Additional business members should be loaded")
	s.Len(retrievedSme.AdditionalBusinessMembers, 2, "Should have 2 additional members")

	// Verify Business Formalisation relationship
	s.NotNil(retrievedSme.BusinessFormalisation, "Business formalisation should be loaded")
	if retrievedSme.BusinessFormalisation != nil {
		s.Equal(75, retrievedSme.BusinessFormalisation.FormalisationScore)
		s.Equal(250000.00, retrievedSme.BusinessFormalisation.AnnualTurnover)
	}

	// Verify Business Employee Summary relationship
	s.NotNil(retrievedSme.BusinessEmployeeSummary, "Business employee summary should be loaded")
	if retrievedSme.BusinessEmployeeSummary != nil {
		s.Equal(5, retrievedSme.BusinessEmployeeSummary.FullTimeMales)
		s.Equal(3, retrievedSme.BusinessEmployeeSummary.FullTimeFemales)
	}

	// Step 7: Verify count queries
	var smeCount int64
	smeCount, err = facades.Orm().Query().Model(&models.Sme{}).Count()
	s.NoError(err)
	s.Equal(int64(1), smeCount, "Should have exactly 1 SME")

	var ownerCount int64
	ownerCount, err = facades.Orm().Query().Model(&models.PrimaryBusinessOwner{}).Where("sme_id = ?", sme.ID).Count()
	s.NoError(err)
	s.Equal(int64(1), ownerCount, "Should have exactly 1 primary owner")

	var memberCount int64
	memberCount, err = facades.Orm().Query().Model(&models.AdditionalBusinessMember{}).Where("sme_id = ?", sme.ID).Count()
	s.NoError(err)
	s.Equal(int64(2), memberCount, "Should have exactly 2 additional members")

	var formalisationCount int64
	formalisationCount, err = facades.Orm().Query().Model(&models.BusinessFormalisation{}).Where("sme_id = ?", sme.ID).Count()
	s.NoError(err)
	s.Equal(int64(1), formalisationCount, "Should have exactly 1 formalisation record")

	var employeeCount int64
	employeeCount, err = facades.Orm().Query().Model(&models.BusinessEmployeeSummary{}).Where("sme_id = ?", sme.ID).Count()
	s.NoError(err)
	s.Equal(int64(1), employeeCount, "Should have exactly 1 employee summary")
}

// TestRegionDistrictValidation tests the region and district validation helpers
func (s *SmeSuite) TestRegionDistrictValidation() {
	// Test valid districts
	s.True(requests.ValidateDistrict("Lilongwe"), "Lilongwe should be valid")
	s.True(requests.ValidateDistrict("Blantyre"), "Blantyre should be valid")
	s.True(requests.ValidateDistrict("Mzimba"), "Mzimba should be valid")

	// Test invalid districts
	s.False(requests.ValidateDistrict("Invalid District"), "Invalid district should return false")
	s.False(requests.ValidateDistrict(""), "Empty district should return false")

	// Test valid regions
	s.True(requests.ValidateRegion("Northern Region"), "Northern Region should be valid")
	s.True(requests.ValidateRegion("Central Region"), "Central Region should be valid")
	s.True(requests.ValidateRegion("Southern Region"), "Southern Region should be valid")

	// Test invalid regions
	s.False(requests.ValidateRegion("Invalid Region"), "Invalid region should return false")

	// Test district to region mapping
	region, exists := requests.GetRegionByDistrict(requests.DistrictLilongwe)
	s.True(exists, "Lilongwe should exist")
	s.Equal(requests.RegionCentral, region, "Lilongwe should be in Central Region")

	region, exists = requests.GetRegionByDistrict(requests.DistrictBlantyre)
	s.True(exists, "Blantyre should exist")
	s.Equal(requests.RegionSouthern, region, "Blantyre should be in Southern Region")

	region, exists = requests.GetRegionByDistrict(requests.DistrictMzimba)
	s.True(exists, "Mzimba should exist")
	s.Equal(requests.RegionNorthern, region, "Mzimba should be in Northern Region")

	// Test getting all districts in a region
	centralDistricts := requests.GetDistrictsByRegion(requests.RegionCentral)
	s.Equal(9, len(centralDistricts), "Central Region should have 9 districts")

	northernDistricts := requests.GetDistrictsByRegion(requests.RegionNorthern)
	s.Equal(6, len(northernDistricts), "Northern Region should have 6 districts")

	southernDistricts := requests.GetDistrictsByRegion(requests.RegionSouthern)
	s.Equal(13, len(southernDistricts), "Southern Region should have 13 districts")

	// Test getting all regions and districts
	allRegions := requests.GetAllRegions()
	s.Equal(3, len(allRegions), "Should have exactly 3 regions")

	allDistricts := requests.GetAllDistricts()
	s.Equal(28, len(allDistricts), "Should have exactly 28 districts")
}
