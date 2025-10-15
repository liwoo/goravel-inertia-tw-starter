package feature

import (
	"github.com/stretchr/testify/suite"
	"players/tests"
	"testing"
)

type LenderTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestLenderTestSuite(t *testing.T) {
	suite.Run(t, new(LenderTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *LenderTestSuite) SetupTest() {
	//migrate fresh
	s.RefreshDatabase()
}

// TearDownTest will run after each test in the suite.
func (s *LenderTestSuite) TearDownTest() {
}

func (s *LenderTestSuite) TestIndex() {

}
