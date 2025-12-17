package requests

type MalawianStatusType string

const (
	MalawianStatusCitizen           MalawianStatusType = "Citizen"
	MalawianStatusPermanentResident MalawianStatusType = "Permanent Resident"
	MalawianStatusWorkPermit        MalawianStatusType = "Work Permit"
	MalawianStatusStudentVisa       MalawianStatusType = "Student Visa"
	MalawianStatusVisitor           MalawianStatusType = "Visitor"
)
