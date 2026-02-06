package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/models"
)

// UpdateBusinessEmployeeSummary listener updates the BusinessEmployeeSummary aggregates
// when AdditionalBusinessMembers are created, updated, or deleted
type UpdateBusinessEmployeeSummary struct {
}

// Signature returns the listener's signature
func (listener *UpdateBusinessEmployeeSummary) Signature() string {
	return "update_business_employee_summary"
}

// Queue returns the queue details for async processing
func (listener *UpdateBusinessEmployeeSummary) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

// Handle processes the event asynchronously
func (listener *UpdateBusinessEmployeeSummary) Handle(args ...any) error {
	// Extract the smeId from the event args
	// Note: After JSON serialization in queue, numbers come back as float64
	smeId, ok := toInt(args[0])
	if !ok || smeId == 0 {
		facades.Log().Warning("UpdateBusinessEmployeeSummary: Invalid smeId")
		return nil
	}

	// Recalculate the employee summary for this SME
	return recalculateEmployeeSummary(smeId)
}

// recalculateEmployeeSummary recalculates and updates the BusinessEmployeeSummary for a given SME
func recalculateEmployeeSummary(smeId int) error {
	// Get all additional business members for this SME
	var members []models.AdditionalBusinessMember
	err := facades.Orm().Query().
		Where("sme_id = ?", smeId).
		Find(&members)

	if err != nil {
		facades.Log().Errorf("Failed to fetch members for SME %d: %v", smeId, err)
		return err
	}

	// Calculate aggregates
	var fullTimeMales, fullTimeFemales int
	var partTimeMales, partTimeFemales int
	var internMales, internFemales int

	for _, member := range members {
		// Determine if the member is male or female based on gender field
		isMale := member.Gender != nil && *member.Gender == "MALE"
		isFemale := member.Gender != nil && *member.Gender == "FEMALE"

		if member.IsIntern {
			if isMale {
				internMales++
			} else if isFemale {
				internFemales++
			}
		} else if member.IsPartTime {
			if isMale {
				partTimeMales++
			} else if isFemale {
				partTimeFemales++
			}
		} else {
			// Full-time employee
			if isMale {
				fullTimeMales++
			} else if isFemale {
				fullTimeFemales++
			}
		}
	}

	// Check if a summary already exists
	var summary models.BusinessEmployeeSummary
	err = facades.Orm().Query().
		Where("sme_id = ?", smeId).
		First(&summary)

	if err != nil || summary.ID == 0 {
		// Create new summary
		summary = models.BusinessEmployeeSummary{
			SmeID:           smeId,
			FullTimeMales:   fullTimeMales,
			FullTimeFemales: fullTimeFemales,
			PartTimeMales:   partTimeMales,
			PartTimeFemales: partTimeFemales,
			InternMales:     internMales,
			InternFemales:   internFemales,
		}

		err = facades.Orm().Query().Create(&summary)
		if err != nil {
			facades.Log().Errorf("Failed to create employee summary for SME %d: %v", smeId, err)
			return err
		}

		facades.Log().Infof("Created employee summary for SME %d", smeId)
	} else {
		// Update existing summary
		summary.FullTimeMales = fullTimeMales
		summary.FullTimeFemales = fullTimeFemales
		summary.PartTimeMales = partTimeMales
		summary.PartTimeFemales = partTimeFemales
		summary.InternMales = internMales
		summary.InternFemales = internFemales

		err = facades.Orm().Query().Save(&summary)

		if err != nil {
			facades.Log().Errorf("Failed to update employee summary for SME %d: %v", smeId, err)
			return err
		}

		facades.Log().Infof("Updated employee summary for SME %d", smeId)
	}

	return nil
}
