package models

type Config struct {
	BaseAuditableModel

	Name        string  `json:"name" db:"name"`
	Code        *string `json:"code" db:"code"`
	ConfigType  string  `json:"config_type" db:"config_type"`
	Description *string `json:"description" db:"description"`
}

func (r *Config) TableName() string {
	return "configs"
}
