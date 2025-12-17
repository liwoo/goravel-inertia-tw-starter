package models

type RegistrationStatus string

const (
	RegistrationPending   = "Pending"
	RegistrationConfirmed = "Active"
	RegistrationRejected  = "Rejected"
	RegistrationSuspended = "Inactive"
)
