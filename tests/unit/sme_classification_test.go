package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// SmeClassificationTestSuite contains all tests for the SME classification logic
type SmeClassificationTestSuite struct {
	suite.Suite
	smeService *services.SmeService
}

// SetupSuite runs once before all tests
func (suite *SmeClassificationTestSuite) SetupSuite() {
	// Note: We're testing the classification logic in isolation,
	// so we don't need a full database setup
}

// SetupTest runs before each test
func (suite *SmeClassificationTestSuite) SetupTest() {
	// Create a fresh SmeService for each test
	// Note: For unit tests, we're testing the determineClassification method
	// which is a pure function that doesn't require database access
	suite.smeService = &services.SmeService{}
}

// TestSmeClassificationTestSuite runs the test suite
func TestSmeClassificationTestSuite(t *testing.T) {
	suite.Run(t, new(SmeClassificationTestSuite))
}

// =============================================================================
// MICRO CLASSIFICATION TESTS
// =============================================================================

// TestMicroClassification_LowTurnoverAndAssets tests micro classification with
// 1 employee, low turnover, and low assets
func (suite *SmeClassificationTestSuite) TestMicroClassification_LowTurnoverAndAssets() {
	// 1 employee with low turnover (100,000) and low assets (500,000)
	// Both financial criteria are met within micro thresholds
	employees := 1
	turnover := 100000.0
	assets := 500000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMicro, classification,
		"1 employee with turnover 100,000 and assets 500,000 should be classified as Micro")
}

// TestMicroClassification_MaxEmployeesAtTurnoverThreshold tests micro classification
// at the maximum employee count and turnover threshold
func (suite *SmeClassificationTestSuite) TestMicroClassification_MaxEmployeesAtTurnoverThreshold() {
	// 4 employees at max turnover threshold (5,000,000)
	employees := 4
	turnover := 5000000.0
	assets := 100000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMicro, classification,
		"4 employees with turnover at max threshold (5,000,000) should be classified as Micro")
}

// TestMicroClassification_AssetsOnlyMeetCriteria tests micro classification when
// assets meet criteria but turnover exceeds the threshold
func (suite *SmeClassificationTestSuite) TestMicroClassification_AssetsOnlyMeetCriteria() {
	// 2 employees with assets meeting criteria (800,000) but turnover exceeding (6,000,000)
	// Since turnover exceeds micro threshold, it doesn't satisfy turnover criteria
	// But assets are within micro threshold, so assets criteria is satisfied
	employees := 2
	turnover := 6000000.0    // Exceeds micro threshold of 5,000,000
	assets := 800000.0       // Within micro threshold of 1,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMicro, classification,
		"2 employees with assets 800,000 (within micro threshold) should be Micro even if turnover exceeds threshold")
}

// TestMicroClassification_TurnoverOnlyMeetCriteria tests micro classification when
// only turnover criteria is met
func (suite *SmeClassificationTestSuite) TestMicroClassification_TurnoverOnlyMeetCriteria() {
	// 3 employees with turnover meeting criteria (3,000,000) but assets exceeding threshold
	employees := 3
	turnover := 3000000.0   // Within micro threshold of 5,000,000
	assets := 2000000.0     // Exceeds micro threshold of 1,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMicro, classification,
		"3 employees with turnover 3,000,000 (within micro threshold) should be Micro even if assets exceed threshold")
}

// TestMicroClassification_ExactBoundaryThresholds tests micro classification at exact boundaries
func (suite *SmeClassificationTestSuite) TestMicroClassification_ExactBoundaryThresholds() {
	// Test at exact thresholds
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		{
			name:       "Min employees (1), max turnover",
			employees:  1,
			turnover:   models.MicroTurnoverMax,
			assets:     0,
			expected:   models.ClassificationMicro,
		},
		{
			name:       "Max employees (4), max assets",
			employees:  4,
			turnover:   0,
			assets:     models.MicroAssetsMax,
			expected:   models.ClassificationMicro,
		},
		{
			name:       "Max employees (4), max turnover, max assets",
			employees:  4,
			turnover:   models.MicroTurnoverMax,
			assets:     models.MicroAssetsMax,
			expected:   models.ClassificationMicro,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// =============================================================================
// SMALL CLASSIFICATION TESTS
// =============================================================================

// TestSmallClassification_MediumTurnover tests small classification with
// 5 employees and medium turnover
func (suite *SmeClassificationTestSuite) TestSmallClassification_MediumTurnover() {
	// 5 employees with medium turnover (20,000,000)
	employees := 5
	turnover := 20000000.0   // Within small range (5M - 50M)
	assets := 5000000.0      // Within small threshold (20M)

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationSmall, classification,
		"5 employees with turnover 20,000,000 should be classified as Small")
}

// TestSmallClassification_MaxThresholds tests small classification at maximum thresholds
func (suite *SmeClassificationTestSuite) TestSmallClassification_MaxThresholds() {
	// 20 employees at max thresholds
	employees := 20
	turnover := 50000000.0   // Max small turnover threshold
	assets := 20000000.0     // Max small assets threshold

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationSmall, classification,
		"20 employees at max small thresholds should be classified as Small")
}

// TestSmallClassification_AssetsOnlyMeetCriteria tests small classification when
// assets meet criteria but turnover exceeds
func (suite *SmeClassificationTestSuite) TestSmallClassification_AssetsOnlyMeetCriteria() {
	// 10 employees with assets meeting criteria but turnover exceeding
	employees := 10
	turnover := 60000000.0   // Exceeds small threshold of 50,000,000
	assets := 15000000.0     // Within small threshold of 20,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationSmall, classification,
		"10 employees with assets 15,000,000 (within small threshold) should be Small even if turnover exceeds threshold")
}

// TestSmallClassification_TurnoverOnlyMeetCriteria tests small classification when
// only turnover criteria is met
func (suite *SmeClassificationTestSuite) TestSmallClassification_TurnoverOnlyMeetCriteria() {
	// 15 employees with turnover meeting criteria but assets exceeding
	employees := 15
	turnover := 30000000.0   // Within small range (5M - 50M)
	assets := 25000000.0     // Exceeds small threshold of 20,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationSmall, classification,
		"15 employees with turnover 30,000,000 should be Small even if assets exceed threshold")
}

// TestSmallClassification_MinEmployees tests small classification at minimum employee count
func (suite *SmeClassificationTestSuite) TestSmallClassification_MinEmployees() {
	// Exactly 5 employees (minimum for small)
	employees := 5
	turnover := 10000000.0   // Within small range
	assets := 5000000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationSmall, classification,
		"Exactly 5 employees with valid financial criteria should be Small")
}

// TestSmallClassification_ExactBoundaryThresholds tests small classification at exact boundaries
func (suite *SmeClassificationTestSuite) TestSmallClassification_ExactBoundaryThresholds() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		{
			name:       "Min employees (5), turnover just above min",
			employees:  5,
			turnover:   models.SmallTurnoverMin + 1,
			assets:     0,
			expected:   models.ClassificationSmall,
		},
		{
			name:       "Max employees (20), max turnover",
			employees:  20,
			turnover:   models.SmallTurnoverMax,
			assets:     0,
			expected:   models.ClassificationSmall,
		},
		{
			name:       "Mid employees (12), max assets only",
			employees:  12,
			turnover:   0,
			assets:     models.SmallAssetsMax,
			expected:   models.ClassificationSmall,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// =============================================================================
// MEDIUM CLASSIFICATION TESTS
// =============================================================================

// TestMediumClassification_HighTurnover tests medium classification with
// 21 employees and high turnover
func (suite *SmeClassificationTestSuite) TestMediumClassification_HighTurnover() {
	// 21 employees with high turnover (100,000,000)
	employees := 21
	turnover := 100000000.0   // Within medium range (50M - 500M)
	assets := 50000000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMedium, classification,
		"21 employees with turnover 100,000,000 should be classified as Medium")
}

// TestMediumClassification_MaxThresholds tests medium classification at maximum thresholds
func (suite *SmeClassificationTestSuite) TestMediumClassification_MaxThresholds() {
	// 99 employees at max thresholds
	employees := 99
	turnover := 500000000.0   // Max medium turnover threshold
	assets := 250000000.0     // Max medium assets threshold

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMedium, classification,
		"99 employees at max medium thresholds should be classified as Medium")
}

// TestMediumClassification_AssetsOnlyMeetCriteria tests medium classification when
// only assets meet criteria
func (suite *SmeClassificationTestSuite) TestMediumClassification_AssetsOnlyMeetCriteria() {
	// 50 employees with assets meeting criteria
	employees := 50
	turnover := 600000000.0   // Exceeds medium threshold of 500,000,000
	assets := 200000000.0     // Within medium threshold of 250,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMedium, classification,
		"50 employees with assets 200,000,000 (within medium threshold) should be Medium even if turnover exceeds threshold")
}

// TestMediumClassification_TurnoverOnlyMeetCriteria tests medium classification when
// only turnover criteria is met
func (suite *SmeClassificationTestSuite) TestMediumClassification_TurnoverOnlyMeetCriteria() {
	// 75 employees with turnover meeting criteria but assets exceeding
	employees := 75
	turnover := 300000000.0   // Within medium range (50M - 500M)
	assets := 300000000.0     // Exceeds medium threshold of 250,000,000

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMedium, classification,
		"75 employees with turnover 300,000,000 should be Medium even if assets exceed threshold")
}

// TestMediumClassification_MinEmployees tests medium classification at minimum employee count
func (suite *SmeClassificationTestSuite) TestMediumClassification_MinEmployees() {
	// Exactly 21 employees (minimum for medium)
	employees := 21
	turnover := 75000000.0   // Within medium range
	assets := 100000000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationMedium, classification,
		"Exactly 21 employees with valid financial criteria should be Medium")
}

// TestMediumClassification_ExactBoundaryThresholds tests medium classification at exact boundaries
func (suite *SmeClassificationTestSuite) TestMediumClassification_ExactBoundaryThresholds() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		{
			name:       "Min employees (21), turnover just above min",
			employees:  21,
			turnover:   models.MediumTurnoverMin + 1,
			assets:     0,
			expected:   models.ClassificationMedium,
		},
		{
			name:       "Max employees (99), max turnover",
			employees:  99,
			turnover:   models.MediumTurnoverMax,
			assets:     0,
			expected:   models.ClassificationMedium,
		},
		{
			name:       "Mid employees (60), max assets only",
			employees:  60,
			turnover:   0,
			assets:     models.MediumAssetsMax,
			expected:   models.ClassificationMedium,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// =============================================================================
// UNCLASSIFIED TESTS
// =============================================================================

// TestUnclassified_ZeroEmployees tests that 0 employees results in Unclassified
func (suite *SmeClassificationTestSuite) TestUnclassified_ZeroEmployees() {
	// 0 employees should be unclassified regardless of financial metrics
	employees := 0
	turnover := 5000000.0
	assets := 1000000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationUnclassified, classification,
		"0 employees should be classified as Unclassified regardless of financial criteria")
}

// TestUnclassified_TooManyEmployees tests that 100+ employees results in Unclassified
func (suite *SmeClassificationTestSuite) TestUnclassified_TooManyEmployees() {
	// 100+ employees should be unclassified (not an SME)
	testCases := []struct {
		name      string
		employees int
	}{
		{"100 employees", 100},
		{"150 employees", 150},
		{"500 employees", 500},
		{"1000 employees", 1000},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, 200000000.0, 100000000.0)
			suite.Equal(models.ClassificationUnclassified, classification,
				"%d employees should be classified as Unclassified", tc.employees)
		})
	}
}

// TestUnclassified_EmployeesInRangeButFinancialNotMet tests that employees in range
// but financial criteria not met results in Unclassified
func (suite *SmeClassificationTestSuite) TestUnclassified_EmployeesInRangeButFinancialNotMet() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
	}{
		// Micro employee range but no financial criteria met
		{
			name:      "Micro employees (3) with zero turnover and assets",
			employees: 3,
			turnover:  0,
			assets:    0,
		},
		// Small employee range but no financial criteria met
		{
			name:      "Small employees (10) with zero turnover and assets",
			employees: 10,
			turnover:  0,
			assets:    0,
		},
		// Medium employee range but no financial criteria met
		{
			name:      "Medium employees (50) with zero turnover and assets",
			employees: 50,
			turnover:  0,
			assets:    0,
		},
		// Micro employees but turnover is above micro threshold and assets exceed
		{
			name:      "Micro employees (2) with turnover and assets both exceeding micro thresholds",
			employees: 2,
			turnover:  6000000.0,    // Above micro threshold
			assets:    2000000.0,    // Above micro threshold
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(models.ClassificationUnclassified, classification, tc.name)
		})
	}
}

// TestUnclassified_NegativeEmployees tests edge case of negative employees
func (suite *SmeClassificationTestSuite) TestUnclassified_NegativeEmployees() {
	// Negative employees should be unclassified
	employees := -1
	turnover := 5000000.0
	assets := 1000000.0

	classification := suite.smeService.DetermineClassification(employees, turnover, assets)

	suite.Equal(models.ClassificationUnclassified, classification,
		"Negative employees should be classified as Unclassified")
}

// =============================================================================
// EDGE CASES
// =============================================================================

// TestEdgeCase_ExactlyAtBoundaryThresholds tests classification at exact boundary values
func (suite *SmeClassificationTestSuite) TestEdgeCase_ExactlyAtBoundaryThresholds() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		// Micro/Small boundary - 4 vs 5 employees
		{
			name:       "4 employees at micro/small boundary",
			employees:  4,
			turnover:   5000000.0,
			assets:     1000000.0,
			expected:   models.ClassificationMicro,
		},
		{
			name:       "5 employees at micro/small boundary",
			employees:  5,
			turnover:   10000000.0,
			assets:     5000000.0,
			expected:   models.ClassificationSmall,
		},
		// Small/Medium boundary - 20 vs 21 employees
		{
			name:       "20 employees at small/medium boundary",
			employees:  20,
			turnover:   50000000.0,
			assets:     20000000.0,
			expected:   models.ClassificationSmall,
		},
		{
			name:       "21 employees at small/medium boundary",
			employees:  21,
			turnover:   100000000.0,
			assets:     50000000.0,
			expected:   models.ClassificationMedium,
		},
		// Medium/Unclassified boundary - 99 vs 100 employees
		{
			name:       "99 employees at medium/unclassified boundary",
			employees:  99,
			turnover:   500000000.0,
			assets:     250000000.0,
			expected:   models.ClassificationMedium,
		},
		{
			name:       "100 employees at medium/unclassified boundary",
			employees:  100,
			turnover:   500000000.0,
			assets:     250000000.0,
			expected:   models.ClassificationUnclassified,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// TestEdgeCase_OnlyAssetsNoTurnover tests classification when only assets are provided
func (suite *SmeClassificationTestSuite) TestEdgeCase_OnlyAssetsNoTurnover() {
	testCases := []struct {
		name       string
		employees  int
		assets     float64
		expected   string
	}{
		{
			name:       "Micro: 3 employees, assets within threshold, no turnover",
			employees:  3,
			assets:     800000.0,
			expected:   models.ClassificationMicro,
		},
		{
			name:       "Small: 15 employees, assets within threshold, no turnover",
			employees:  15,
			assets:     15000000.0,
			expected:   models.ClassificationSmall,
		},
		{
			name:       "Medium: 50 employees, assets within threshold, no turnover",
			employees:  50,
			assets:     200000000.0,
			expected:   models.ClassificationMedium,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, 0, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// TestEdgeCase_OnlyTurnoverNoAssets tests classification when only turnover is provided
func (suite *SmeClassificationTestSuite) TestEdgeCase_OnlyTurnoverNoAssets() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		expected   string
	}{
		{
			name:       "Micro: 2 employees, turnover within threshold, no assets",
			employees:  2,
			turnover:   3000000.0,
			expected:   models.ClassificationMicro,
		},
		{
			name:       "Small: 10 employees, turnover within threshold, no assets",
			employees:  10,
			turnover:   25000000.0,
			expected:   models.ClassificationSmall,
		},
		{
			name:       "Medium: 75 employees, turnover within threshold, no assets",
			employees:  75,
			turnover:   300000000.0,
			expected:   models.ClassificationMedium,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, 0)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// TestEdgeCase_TurnoverAtExactThresholdBoundaries tests turnover at exact threshold values
func (suite *SmeClassificationTestSuite) TestEdgeCase_TurnoverAtExactThresholdBoundaries() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		// Micro: turnover "up to" 5,000,000 means <= 5,000,000
		{
			name:       "Micro: turnover exactly at 5,000,000 (max threshold)",
			employees:  3,
			turnover:   models.MicroTurnoverMax,
			assets:     0,
			expected:   models.ClassificationMicro,
		},
		// Small: turnover "above 5,000,000" means > 5,000,000
		// When assets=0, only turnover is considered
		{
			name:       "Small: turnover exactly at 5,000,000 (min threshold boundary), no assets",
			employees:  10,
			turnover:   models.SmallTurnoverMin,
			assets:     0,
			expected:   models.ClassificationUnclassified, // turnover must be > 5M for small, assets=0 doesn't help
		},
		{
			name:       "Small: turnover exactly at 5,000,000 but assets satisfy criteria",
			employees:  10,
			turnover:   models.SmallTurnoverMin,
			assets:     5000000.0, // Assets satisfy Small criteria (> 0 and <= 20M)
			expected:   models.ClassificationSmall, // Classified as Small because assets satisfy criteria
		},
		{
			name:       "Small: turnover just above 5,000,000",
			employees:  10,
			turnover:   models.SmallTurnoverMin + 1,
			assets:     0,
			expected:   models.ClassificationSmall,
		},
		// Medium: turnover "above 50,000,000" means > 50,000,000
		// When assets=0, only turnover is considered
		{
			name:       "Medium: turnover exactly at 50,000,000 (min threshold boundary), no assets",
			employees:  50,
			turnover:   models.MediumTurnoverMin,
			assets:     0,
			expected:   models.ClassificationUnclassified, // turnover must be > 50M for medium, assets=0 doesn't help
		},
		{
			name:       "Medium: turnover exactly at 50,000,000 but assets satisfy criteria",
			employees:  50,
			turnover:   models.MediumTurnoverMin,
			assets:     100000000.0, // Assets satisfy Medium criteria (> 0 and <= 250M)
			expected:   models.ClassificationMedium, // Classified as Medium because assets satisfy criteria
		},
		{
			name:       "Medium: turnover just above 50,000,000",
			employees:  50,
			turnover:   models.MediumTurnoverMin + 1,
			assets:     0,
			expected:   models.ClassificationMedium,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// TestEdgeCase_LargeValues tests classification with very large financial values
func (suite *SmeClassificationTestSuite) TestEdgeCase_LargeValues() {
	testCases := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		{
			name:       "Medium with very large turnover",
			employees:  80,
			turnover:   499999999.0, // Just under max
			assets:     100000000.0,
			expected:   models.ClassificationMedium,
		},
		{
			name:       "Medium with max turnover",
			employees:  80,
			turnover:   500000000.0, // At max
			assets:     100000000.0,
			expected:   models.ClassificationMedium,
		},
		{
			name:       "Unclassified with turnover exceeding all thresholds",
			employees:  80,
			turnover:   600000000.0, // Above max
			assets:     300000000.0, // Above max
			expected:   models.ClassificationUnclassified,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			classification := suite.smeService.DetermineClassification(tc.employees, tc.turnover, tc.assets)
			suite.Equal(tc.expected, classification, tc.name)
		})
	}
}

// =============================================================================
// STANDALONE TESTS (for additional coverage with standard testing pattern)
// =============================================================================

// TestDetermineClassification_AllCategories provides comprehensive coverage
func TestDetermineClassification_AllCategories(t *testing.T) {
	smeService := &services.SmeService{}

	tests := []struct {
		name       string
		employees  int
		turnover   float64
		assets     float64
		expected   string
	}{
		// Micro
		{"Micro - standard case", 3, 3000000.0, 500000.0, models.ClassificationMicro},
		{"Micro - min employees", 1, 1000000.0, 100000.0, models.ClassificationMicro},
		{"Micro - max employees", 4, 5000000.0, 1000000.0, models.ClassificationMicro},
		{"Micro - turnover only", 2, 4000000.0, 0, models.ClassificationMicro},
		{"Micro - assets only", 2, 0, 800000.0, models.ClassificationMicro},

		// Small
		{"Small - standard case", 12, 30000000.0, 15000000.0, models.ClassificationSmall},
		{"Small - min employees", 5, 10000000.0, 5000000.0, models.ClassificationSmall},
		{"Small - max employees", 20, 50000000.0, 20000000.0, models.ClassificationSmall},
		{"Small - turnover only", 10, 25000000.0, 0, models.ClassificationSmall},
		{"Small - assets only", 10, 0, 15000000.0, models.ClassificationSmall},

		// Medium
		{"Medium - standard case", 50, 200000000.0, 150000000.0, models.ClassificationMedium},
		{"Medium - min employees", 21, 75000000.0, 50000000.0, models.ClassificationMedium},
		{"Medium - max employees", 99, 500000000.0, 250000000.0, models.ClassificationMedium},
		{"Medium - turnover only", 60, 300000000.0, 0, models.ClassificationMedium},
		{"Medium - assets only", 60, 0, 200000000.0, models.ClassificationMedium},

		// Unclassified
		{"Unclassified - zero employees", 0, 5000000.0, 1000000.0, models.ClassificationUnclassified},
		{"Unclassified - 100 employees", 100, 200000000.0, 100000000.0, models.ClassificationUnclassified},
		{"Unclassified - negative employees", -5, 5000000.0, 1000000.0, models.ClassificationUnclassified},
		{"Unclassified - no financial data with micro employees", 3, 0, 0, models.ClassificationUnclassified},
		{"Unclassified - no financial data with small employees", 10, 0, 0, models.ClassificationUnclassified},
		{"Unclassified - no financial data with medium employees", 50, 0, 0, models.ClassificationUnclassified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := smeService.DetermineClassification(tt.employees, tt.turnover, tt.assets)
			assert.Equal(t, tt.expected, result, "Classification for %s should be %s", tt.name, tt.expected)
		})
	}
}

// TestClassificationConstants verifies the classification constants are correctly defined
func TestClassificationConstants(t *testing.T) {
	// Verify classification string constants
	assert.Equal(t, "Micro", models.ClassificationMicro)
	assert.Equal(t, "Small", models.ClassificationSmall)
	assert.Equal(t, "Medium", models.ClassificationMedium)
	assert.Equal(t, "Unclassified", models.ClassificationUnclassified)

	// Verify employee thresholds
	assert.Equal(t, 1, models.MicroEmployeeMin)
	assert.Equal(t, 4, models.MicroEmployeeMax)
	assert.Equal(t, 5, models.SmallEmployeeMin)
	assert.Equal(t, 20, models.SmallEmployeeMax)
	assert.Equal(t, 21, models.MediumEmployeeMin)
	assert.Equal(t, 99, models.MediumEmployeeMax)

	// Verify turnover thresholds
	assert.Equal(t, 5000000.0, models.MicroTurnoverMax)
	assert.Equal(t, 5000000.0, models.SmallTurnoverMin)
	assert.Equal(t, 50000000.0, models.SmallTurnoverMax)
	assert.Equal(t, 50000000.0, models.MediumTurnoverMin)
	assert.Equal(t, 500000000.0, models.MediumTurnoverMax)

	// Verify assets thresholds
	assert.Equal(t, 1000000.0, models.MicroAssetsMax)
	assert.Equal(t, 20000000.0, models.SmallAssetsMax)
	assert.Equal(t, 250000000.0, models.MediumAssetsMax)
}

// TestClassificationLogicConsistency verifies the classification logic is consistent
// across boundary transitions
func TestClassificationLogicConsistency(t *testing.T) {
	smeService := &services.SmeService{}

	// Test employee boundary transitions
	t.Run("Employee boundaries", func(t *testing.T) {
		// 4 -> 5 employees should transition from Micro to Small (with valid financials)
		microResult := smeService.DetermineClassification(4, 5000000.0, 1000000.0)
		assert.Equal(t, models.ClassificationMicro, microResult)

		smallResult := smeService.DetermineClassification(5, 10000000.0, 5000000.0)
		assert.Equal(t, models.ClassificationSmall, smallResult)

		// 20 -> 21 employees should transition from Small to Medium (with valid financials)
		smallResult2 := smeService.DetermineClassification(20, 50000000.0, 20000000.0)
		assert.Equal(t, models.ClassificationSmall, smallResult2)

		mediumResult := smeService.DetermineClassification(21, 100000000.0, 50000000.0)
		assert.Equal(t, models.ClassificationMedium, mediumResult)

		// 99 -> 100 employees should transition from Medium to Unclassified
		mediumResult2 := smeService.DetermineClassification(99, 500000000.0, 250000000.0)
		assert.Equal(t, models.ClassificationMedium, mediumResult2)

		unclassifiedResult := smeService.DetermineClassification(100, 500000000.0, 250000000.0)
		assert.Equal(t, models.ClassificationUnclassified, unclassifiedResult)
	})
}
