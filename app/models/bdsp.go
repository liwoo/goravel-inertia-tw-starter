package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Bdsp struct {
	BaseAuditableModel

	UbdspNumber        string   `json:"ubdsp_number" db:"ubdsp_number"`
	Name               string   `json:"name" db:"name"`
	PostalAddress      *string  `json:"postal_address" db:"postal_address_1"`
	PhysicalAddress    *string  `json:"physical_address" db:"physical_address"`
	RegistrationStatus *string  `json:"registration_status" db:"registration_status"`
	PartnersJSON       string   `json:"-" db:"partners_json" gorm:"column:partners_json"`
	Partners           []string `json:"partners" gorm:"-"`

	ProductTypesJSON string   `json:"-" db:"product_types_json" gorm:"column:product_types_json"`
	ProductTypes     []string `json:"product_types" gorm:"-"`

	ServiceListJSON string        `json:"-" db:"service_list_json" gorm:"column:service_list_json"`
	ServiceList     []BdspService `json:"service_list" gorm:"-"`
}

type BdspService struct {
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Duration string  `json:"duration"`
}

func (r *Bdsp) TableName() string {
	return "bdsps"
}

// BeforeSave hook to convert arrays to JSON
func (b *Bdsp) BeforeSave(tx *gorm.DB) error {
	// Convert Partners to JSON
	if len(b.Partners) > 0 {
		partnersBytes, err := json.Marshal(b.Partners)
		if err != nil {
			return err
		}
		b.PartnersJSON = string(partnersBytes)
	} else {
		b.PartnersJSON = ""
	}

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

	return nil
}

// AfterFind hook to convert JSON to arrays
func (b *Bdsp) AfterFind(tx *gorm.DB) error {
	// Convert PartnersJSON to array
	if b.PartnersJSON != "" {
		err := json.Unmarshal([]byte(b.PartnersJSON), &b.Partners)
		if err != nil {
			b.Partners = []string{}
		}
	} else {
		b.Partners = []string{}
	}

	// Convert ProductTypesJSON to array
	if b.ProductTypesJSON != "" {
		err := json.Unmarshal([]byte(b.ProductTypesJSON), &b.ProductTypes)
		if err != nil {
			b.ProductTypes = []string{}
		}
	} else {
		b.ProductTypes = []string{}
	}

	// Convert ServiceListJSON to array
	if b.ServiceListJSON != "" {
		err := json.Unmarshal([]byte(b.ServiceListJSON), &b.ServiceList)
		if err != nil {
			b.ServiceList = []BdspService{}
		}
	} else {
		b.ServiceList = []BdspService{}
	}

	return nil
}
