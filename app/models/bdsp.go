package models

import (
	"encoding/json"

	"github.com/goravel/framework/database/orm"
	"gorm.io/gorm"
)

// BdspService represents a service offered by a BDSP
type BdspService struct {
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Duration string  `json:"duration"` // e.g., "1 hour", "2 days", etc.
}

type Bdsp struct {
	orm.Model
	orm.SoftDeletes

	// Basic Information
	Name            string  `json:"name" db:"name"`
	PostalAddress   *string `json:"postal_address" db:"postal_address"`
	PhysicalAddress *string `json:"physical_address" db:"physical_address"`

	// Registration Status - references config with type "Registration Status"
	RegistrationStatus string `json:"registration_status" db:"registration_status"`

	// Product Types - stored as JSON array, fetched from config with type "Product Types"
	ProductTypesJSON string   `json:"-" db:"product_types_json" gorm:"column:product_types_json;type:text"`
	ProductTypes     []string `json:"product_types" gorm:"-:all"`

	// Service List - stored as JSON array, each service has name, cost, and duration
	ServiceListJSON string        `json:"-" db:"service_list_json" gorm:"column:service_list_json;type:text"`
	ServiceList     []BdspService `json:"service_list" gorm:"-:all"`

	// Associated Partners - stored as JSON array, fetched from config with type "Development Partners"
	AssociatedPartnersJSON string   `json:"-" db:"associated_partners_json" gorm:"column:associated_partners_json;type:text"`
	AssociatedPartners     []string `json:"associated_partners" gorm:"-:all"`

	// Audit fields
	CreatedBy *int    `json:"created_by" db:"created_by"`
	UpdatedBy *int    `json:"updated_by" db:"updated_by"`
	DeletedBy *int    `json:"deleted_by" db:"deleted_by"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
}

func (r *Bdsp) TableName() string {
	return "bdsps"
}

// BeforeSave hook to convert arrays to JSON
func (b *Bdsp) BeforeSave(tx *gorm.DB) error {
	// Convert ProductTypes to JSON
	if len(b.ProductTypes) > 0 {
		productTypesBytes, err := json.Marshal(b.ProductTypes)
		if err != nil {
			return err
		}
		b.ProductTypesJSON = string(productTypesBytes)
	} else {
		b.ProductTypesJSON = ""
	}

	// Convert ServiceList to JSON
	if len(b.ServiceList) > 0 {
		serviceListBytes, err := json.Marshal(b.ServiceList)
		if err != nil {
			return err
		}
		b.ServiceListJSON = string(serviceListBytes)
	} else {
		b.ServiceListJSON = ""
	}

	// Convert AssociatedPartners to JSON
	if len(b.AssociatedPartners) > 0 {
		partnersBytes, err := json.Marshal(b.AssociatedPartners)
		if err != nil {
			return err
		}
		b.AssociatedPartnersJSON = string(partnersBytes)
	} else {
		b.AssociatedPartnersJSON = ""
	}

	return nil
}

// AfterFind hook to convert JSON to arrays
func (b *Bdsp) AfterFind(tx *gorm.DB) error {
	// Convert ProductTypesJSON to array
	if b.ProductTypesJSON != "" {
		err := json.Unmarshal([]byte(b.ProductTypesJSON), &b.ProductTypes)
		if err != nil {
			// If unmarshal fails, treat as empty array
			b.ProductTypes = []string{}
		}
	} else {
		b.ProductTypes = []string{}
	}

	// Convert ServiceListJSON to array
	if b.ServiceListJSON != "" {
		err := json.Unmarshal([]byte(b.ServiceListJSON), &b.ServiceList)
		if err != nil {
			// If unmarshal fails, treat as empty array
			b.ServiceList = []BdspService{}
		}
	} else {
		b.ServiceList = []BdspService{}
	}

	// Convert AssociatedPartnersJSON to array
	if b.AssociatedPartnersJSON != "" {
		err := json.Unmarshal([]byte(b.AssociatedPartnersJSON), &b.AssociatedPartners)
		if err != nil {
			// If unmarshal fails, treat as empty array
			b.AssociatedPartners = []string{}
		}
	} else {
		b.AssociatedPartners = []string{}
	}

	return nil
}
