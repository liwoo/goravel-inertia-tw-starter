package services

import (
	"time"

	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/models"
)

// DashboardService provides methods for dashboard statistics and widgets
type DashboardService struct{}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// UpcomingEventDTO represents the data structure for upcoming events in the dashboard widget
type UpcomingEventDTO struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Venue    string `json:"venue"`
	District string `json:"district,omitempty"`
}

// UpcomingProcurementDTO represents the data structure for upcoming procurements in the dashboard widget
type UpcomingProcurementDTO struct {
	ID              uint   `json:"id"`
	Organization    string `json:"organization"`
	RefNo           string `json:"refNo"`
	CloseDate       string `json:"closeDate"`
	ProcurementType string `json:"procurementType,omitempty"`
}

// RecentActivityDTO represents the data structure for recent user activities
type RecentActivityDTO struct {
	ID           uint   `json:"id"`
	EntityType   string `json:"entityType"`   // "sme", "event", "procurement"
	EntityID     uint   `json:"entityId"`
	EntityName   string `json:"entityName"`
	Action       string `json:"action"`       // "created", "updated"
	UserID       uint   `json:"userId"`
	UserName     string `json:"userName"`
	Timestamp    string `json:"timestamp"`
}

// EventStatistics contains statistics about events
type EventStatistics struct {
	TotalEvents    int64 `json:"totalEvents"`
	UpcomingEvents int64 `json:"upcomingEvents"`
}

// ProcurementStatistics contains statistics about procurement notices
type ProcurementStatistics struct {
	TotalProcurements  int64 `json:"totalProcurements"`
	ActiveProcurements int64 `json:"activeProcurements"`
}

// GetEventStatistics returns statistics about events
func (s *DashboardService) GetEventStatistics() EventStatistics {
	now := time.Now()

	var totalEvents int64
	var upcomingEvents int64

	// Get total events count (soft deletes handled by GORM)
	totalEvents, _ = facades.Orm().Query().Model(&models.Event{}).Count()

	// Get upcoming events count (date > now)
	upcomingEvents, _ = facades.Orm().Query().Model(&models.Event{}).
		Where("date > ?", now).
		Count()

	return EventStatistics{
		TotalEvents:    totalEvents,
		UpcomingEvents: upcomingEvents,
	}
}

// GetProcurementStatistics returns statistics about procurement notices
func (s *DashboardService) GetProcurementStatistics() ProcurementStatistics {
	now := time.Now()

	var totalProcurements int64
	var activeProcurements int64

	// Get total procurement notices count (soft deletes handled by GORM)
	totalProcurements, _ = facades.Orm().Query().Model(&models.ProcurementNotice{}).Count()

	// Get active procurement notices count (is_published=true AND close_date > now)
	activeProcurements, _ = facades.Orm().Query().Model(&models.ProcurementNotice{}).
		Where("is_published = ? AND close_date > ?", true, now).
		Count()

	return ProcurementStatistics{
		TotalProcurements:  totalProcurements,
		ActiveProcurements: activeProcurements,
	}
}

// GetUpcomingEvents returns the next N upcoming events ordered by date
func (s *DashboardService) GetUpcomingEvents(limit int) []UpcomingEventDTO {
	now := time.Now()
	var events []models.Event

	// Get upcoming events where date > now, ordered by date ascending
	err := facades.Orm().Query().Model(&models.Event{}).
		Where("date > ?", now).
		Order("date ASC").
		Limit(limit).
		Find(&events)

	if err != nil {
		return []UpcomingEventDTO{}
	}

	// Convert to DTOs
	result := make([]UpcomingEventDTO, len(events))
	for i, event := range events {
		result[i] = UpcomingEventDTO{
			ID:       event.ID,
			Title:    event.Title,
			Date:     event.Date.ToDateTimeString(),
			Venue:    event.Venue,
			District: event.District,
		}
	}

	return result
}

// GetUpcomingProcurements returns the next N upcoming procurement notices ordered by close date
func (s *DashboardService) GetUpcomingProcurements(limit int) []UpcomingProcurementDTO {
	now := time.Now()
	var procurements []models.ProcurementNotice

	// Get active procurement notices (is_published=true AND close_date > now), ordered by close_date ascending
	err := facades.Orm().Query().Model(&models.ProcurementNotice{}).
		Where("is_published = ? AND close_date > ?", true, now).
		Order("close_date ASC").
		Limit(limit).
		Find(&procurements)

	if err != nil {
		return []UpcomingProcurementDTO{}
	}

	// Convert to DTOs
	result := make([]UpcomingProcurementDTO, len(procurements))
	for i, proc := range procurements {
		result[i] = UpcomingProcurementDTO{
			ID:              proc.ID,
			Organization:    proc.Organization,
			RefNo:           proc.RefNo,
			CloseDate:       proc.CloseDate.ToDateTimeString(),
			ProcurementType: proc.ProcurementType,
		}
	}

	return result
}

// GetRecentActivities returns the most recent user activities across all entities
func (s *DashboardService) GetRecentActivities(limit int) []RecentActivityDTO {
	var activities []RecentActivityDTO

	// Helper to get user name by ID
	getUserName := func(userID *int) string {
		if userID == nil {
			return "System"
		}
		var user models.User
		err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", *userID).First(&user)
		if err != nil {
			return "Unknown"
		}
		return user.Name
	}

	// Get recent SME activities (created and updated)
	var smes []models.Sme
	facades.Orm().Query().Model(&models.Sme{}).
		Order("updated_at DESC").
		Limit(limit).
		Find(&smes)

	for _, sme := range smes {
		action := "updated"
		userID := sme.UpdatedBy
		timestamp := sme.UpdatedAt.ToDateTimeString()

		// Check if this is a newly created record (created_at == updated_at approximately)
		if sme.CreatedAt.ToDateTimeString() == sme.UpdatedAt.ToDateTimeString() {
			action = "created"
			userID = sme.CreatedBy
			timestamp = sme.CreatedAt.ToDateTimeString()
		}

		activities = append(activities, RecentActivityDTO{
			ID:         sme.ID,
			EntityType: "sme",
			EntityID:   sme.ID,
			EntityName: sme.Name,
			Action:     action,
			UserID:     uint(safeInt(userID)),
			UserName:   getUserName(userID),
			Timestamp:  timestamp,
		})
	}

	// Get recent Event activities
	var events []models.Event
	facades.Orm().Query().Model(&models.Event{}).
		Order("updated_at DESC").
		Limit(limit).
		Find(&events)

	for _, event := range events {
		action := "updated"
		userID := event.UpdatedBy
		timestamp := event.UpdatedAt.ToDateTimeString()

		if event.CreatedAt.ToDateTimeString() == event.UpdatedAt.ToDateTimeString() {
			action = "created"
			userID = event.CreatedBy
			timestamp = event.CreatedAt.ToDateTimeString()
		}

		activities = append(activities, RecentActivityDTO{
			ID:         event.ID,
			EntityType: "event",
			EntityID:   event.ID,
			EntityName: event.Title,
			Action:     action,
			UserID:     uint(safeInt(userID)),
			UserName:   getUserName(userID),
			Timestamp:  timestamp,
		})
	}

	// Get recent Procurement Notice activities
	var procurements []models.ProcurementNotice
	facades.Orm().Query().Model(&models.ProcurementNotice{}).
		Order("updated_at DESC").
		Limit(limit).
		Find(&procurements)

	for _, proc := range procurements {
		action := "updated"
		userID := proc.UpdatedBy
		timestamp := proc.UpdatedAt.ToDateTimeString()

		if proc.CreatedAt.ToDateTimeString() == proc.UpdatedAt.ToDateTimeString() {
			action = "created"
			userID = proc.CreatedBy
			timestamp = proc.CreatedAt.ToDateTimeString()
		}

		activities = append(activities, RecentActivityDTO{
			ID:         proc.ID,
			EntityType: "procurement",
			EntityID:   proc.ID,
			EntityName: proc.Organization + " - " + proc.RefNo,
			Action:     action,
			UserID:     uint(safeInt(userID)),
			UserName:   getUserName(userID),
			Timestamp:  timestamp,
		})
	}

	// Sort all activities by timestamp (newest first) and limit
	sortActivitiesByTimestamp(activities)
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities
}

// safeInt converts *int to int safely
func safeInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// sortActivitiesByTimestamp sorts activities by timestamp descending
func sortActivitiesByTimestamp(activities []RecentActivityDTO) {
	for i := 0; i < len(activities)-1; i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].Timestamp < activities[j].Timestamp {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}
}

// GetDashboardStats returns all statistics for the dashboard
func (s *DashboardService) GetDashboardStats(smeService *SmeService) map[string]interface{} {
	// Get SME statistics
	smeStats, _ := smeService.GetSmeStatistics()

	// Get event statistics
	eventStats := s.GetEventStatistics()

	// Get procurement statistics
	procurementStats := s.GetProcurementStatistics()

	// Get additional distributions
	bySector := smeService.GetDistributionBySector()
	byGender := smeService.GetDistributionByGender()

	// Combine all statistics
	return map[string]interface{}{
		"totalSmes":          smeStats["totalSmes"],
		"newThisMonth":       smeStats["newThisMonth"],
		"newLastMonth":       smeStats["newLastMonth"],
		"totalEvents":        eventStats.TotalEvents,
		"upcomingEvents":     eventStats.UpcomingEvents,
		"totalProcurements":  procurementStats.TotalProcurements,
		"activeProcurements": procurementStats.ActiveProcurements,
		"byRegion":           smeStats["byRegion"],
		"byCategory":         smeStats["byCategory"],
		"bySector":           bySector,
		"byGender":           byGender,
	}
}
