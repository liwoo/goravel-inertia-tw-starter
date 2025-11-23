package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Event struct {
	orm.Model
	orm.SoftDeletes
	Title         string          `json:"title" db:"title"`
	Description   string          `json:"description" db:"description"`
	Date          carbon.DateTime `json:"date" db:"date"`
	Venue         string          `json:"venue" db:"venue"`
	Partners      []string        `json:"partners" db:"partners" gorm:"type:json;serializer:json"`
	District      string          `json:"district" db:"district"`
	AttendingSmes []int           `json:"attending_smes" db:"attending_smes" gorm:"type:json;serializer:json"`
	Notes         *string         `json:"notes" db:"notes"`
	CreatedBy     *int            `json:"created_by" db:"created_by"`
	UpdatedBy     *int            `json:"updated_by" db:"updated_by"`
	DeletedBy     *int            `json:"deleted_by" db:"deleted_by"`
	IpAddress     *string         `json:"ip_address" db:"ip_address"`
	UserAgent     *string         `json:"user_agent" db:"user_agent"`
}

func (r *Event) TableName() string {
	return "events"
}
