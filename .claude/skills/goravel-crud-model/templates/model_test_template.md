# GORM Model Interaction Test Template

Use this template when generating tests after model creation.

## Template

Replace `<ModelName>`, `<model_name>`, and field placeholders with actual values.

```go
package unit

import (
    "testing"

    "github.com/goravel/framework/facades"
    "github.com/stretchr/testify/suite"

    "smedi-sme-db/app/models"
    "smedi-sme-db/tests"
)

type <ModelName>ModelTestSuite struct {
    suite.Suite
    tests.TestCase
}

func Test<ModelName>ModelTestSuite(t *testing.T) {
    suite.Run(t, &<ModelName>ModelTestSuite{})
}

func (s *<ModelName>ModelTestSuite) SetupTest() {
    s.RefreshDatabase()
}

func (s *<ModelName>ModelTestSuite) TestCreate() {
    model := &models.<ModelName>{
        // Set all required non-null fields here
    }
    err := facades.Orm().Query().Create(model)
    s.Nil(err)
    s.NotZero(model.ID)
}

func (s *<ModelName>ModelTestSuite) TestFindByID() {
    model := &models.<ModelName>{
        // Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    var found models.<ModelName>
    err := facades.Orm().Query().Where("id", model.ID).First(&found)
    s.Nil(err)
    s.Equal(model.ID, found.ID)
}

func (s *<ModelName>ModelTestSuite) TestUpdate() {
    model := &models.<ModelName>{
        // Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    _, err := facades.Orm().Query().Model(&models.<ModelName>{}).
        Where("id", model.ID).
        Update(map[string]interface{}{
            // "field": "new_value",
        })
    s.Nil(err)

    var updated models.<ModelName>
    s.Nil(facades.Orm().Query().Where("id", model.ID).First(&updated))
    // Assert updated field value
}

func (s *<ModelName>ModelTestSuite) TestSoftDelete() {
    model := &models.<ModelName>{
        // Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    _, err := facades.Orm().Query().Delete(model)
    s.Nil(err)

    // Should not find with normal query
    var notFound models.<ModelName>
    err = facades.Orm().Query().Where("id", model.ID).First(&notFound)
    s.NotNil(err)

    // Should find with trashed query
    var found models.<ModelName>
    err = facades.Orm().Query().WithTrashed().Where("id", model.ID).First(&found)
    s.Nil(err)
}
```

## Running the Test

```bash
APP_ENV=testing go test -v ./tests/unit -run Test<ModelName>ModelTestSuite
```

## Key Points

- Initialize ALL array/JSON fields (even empty: `[]string{}`, `[]int{}`)
- Use pointer for optional fields (`*int`, `*string`)
- carbon.DateTime non-pointer: `*carbon.NewDateTime(carbon.Parse("2025-01-01"))`
- Set `CreatedBy` if audit fields have NOT NULL constraint: `&createdByInt`
