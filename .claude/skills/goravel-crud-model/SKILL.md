---
name: goravel-crud-model
description: Generate a Goravel model from an existing database table with post-generation fixes and automatic GORM test creation. Use after running migration.
argument-hint: "[TableName] [ModelName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Model Generator

Generate model for `$ARGUMENTS`.

## Step 1: Generate Model from Table

```bash
go run . artisan make:model --table=<table_name> <ModelName>
```

This creates `app/models/<model_name>.go`.

## Step 2: Post-Generation Fixes (CRITICAL)

The generator has known issues. Apply these fixes:

### Fix 1: Array/JSON Fields

Generator creates them as `string`. Change to proper types:

```go
// WRONG (generator output):
Partners string `json:"partners"`

// CORRECT:
Partners []string `json:"partners" db:"partners" gorm:"type:json;serializer:json"`
AttendingSmes []int `json:"attending_smes" db:"attending_smes" gorm:"type:json;serializer:json"`
```

### Fix 2: Audit Field Types

Ensure audit fields use pointer types:

```go
CreatedBy *int  `json:"created_by" db:"created_by"`
UpdatedBy *int  `json:"updated_by" db:"updated_by"`
DeletedBy *int  `json:"deleted_by" db:"deleted_by"`
```

### Fix 3: carbon.DateTime Fields

Verify date fields are non-pointer (unless nullable):

```go
// Non-nullable date:
Date carbon.DateTime `json:"date" db:"date"`

// Nullable date:
DateOfBirth *carbon.DateTime `json:"date_of_birth" db:"date_of_birth"`
```

### Fix 4: Add SearchFields Method

```go
func (m ModelName) SearchFields() []string {
    return []string{"field1", "field2", "field3"}
}
```

### Fix 5: Verify TableName Method

```go
func (m ModelName) TableName() string {
    return "table_name"
}
```

## Reference Pattern

See `app/models/book.go` for the canonical model implementation with:
- `BaseAuditableModel` embedding
- JSON field serialization hooks (`BeforeSave`/`AfterFind`)
- `SearchFields()` and `TableName()` methods

## Step 3: Auto-Generate GORM Interaction Test

**IMPORTANT**: After fixing the model, automatically create a GORM test at `tests/unit/<model_name>_model_test.go`:

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

func (s *<ModelName>ModelTestSuite) TestCreate<ModelName>() {
    model := &models.<ModelName>{
        // TODO: Set required fields based on the model
    }
    err := facades.Orm().Query().Create(model)
    s.Nil(err)
    s.NotZero(model.ID)
}

func (s *<ModelName>ModelTestSuite) TestFind<ModelName>() {
    model := &models.<ModelName>{
        // TODO: Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    var found models.<ModelName>
    err := facades.Orm().Query().Where("id", model.ID).First(&found)
    s.Nil(err)
    s.Equal(model.ID, found.ID)
}

func (s *<ModelName>ModelTestSuite) TestUpdate<ModelName>() {
    model := &models.<ModelName>{
        // TODO: Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    // TODO: Update a field
    _, err := facades.Orm().Query().Model(&models.<ModelName>{}).Where("id", model.ID).Update(map[string]interface{}{
        // "field": "new_value",
    })
    s.Nil(err)
}

func (s *<ModelName>ModelTestSuite) TestSoftDelete<ModelName>() {
    model := &models.<ModelName>{
        // TODO: Set required fields
    }
    s.Nil(facades.Orm().Query().Create(model))

    _, err := facades.Orm().Query().Delete(model)
    s.Nil(err)

    // Verify soft deleted
    var deleted models.<ModelName>
    err = facades.Orm().Query().WithTrashed().Where("id", model.ID).First(&deleted)
    s.Nil(err)
    s.NotNil(deleted.DeletedAt)
}
```

Fill in the TODO fields based on the actual model fields, then run:

```bash
APP_ENV=testing go test -v ./tests/unit -run Test<ModelName>ModelTestSuite
```

## Next Step

Run `/goravel-crud-service` to generate the service layer.
