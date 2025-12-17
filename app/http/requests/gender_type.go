package requests

type GenderType string

const (
	GenderMale   GenderType = "MALE"
	GenderFemale GenderType = "FEMALE"
)

// GetAllGenders returns all valid gender values
func GetAllGenders() []GenderType {
	return []GenderType{
		GenderMale,
		GenderFemale,
	}
}
