package helpers

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
)

// CrudHelper provides common CRUD operations
type CrudHelper struct {
	tableName string
	db        orm.Query
}

// NewCrudHelper creates a new CRUD helper
func NewCrudHelper(tableName string) *CrudHelper {
	return &CrudHelper{
		tableName: tableName,
		db:        facades.Orm().Query().Table(tableName),
	}
}

// Query returns the base query for custom operations
func (c *CrudHelper) Query() orm.Query {
	return facades.Orm().Query().Table(c.tableName)
}

// TableName returns the table name
func (c *CrudHelper) TableName() string {
	return c.tableName
}