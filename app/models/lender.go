package models

type Lender struct {
	BaseAuditableModel

	Name    string  `json:"name" db:"name"`
	Email   string  `json:"email" db:"email"`
	Phone   *string `json:"phone" db:"phone"`
	Address *string `json:"address" db:"address"`
	Gender  *string `json:"gender" db:"gender"`
}

func (r *Lender) TableName() string {
	return "lenders"
}
