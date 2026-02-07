# Backend Aggregation - Reference Examples

How Go services compute dashboard statistics and distribution data for chart consumption.

## Service Architecture

```
DashboardController
  └─ DashboardService.GetDashboardStats(smeService)
       ├─ smeService.GetSmeStatistics()         → totalSmes, newThisMonth, newLastMonth, byRegion, byCategory
       ├─ smeService.GetDistributionBySector()   → [{label, value, percentage}]
       ├─ smeService.GetDistributionByDistrict() → [{label, value, percentage}]
       ├─ smeService.GetDistributionByGender()   → [{label, value, percentage}]
       ├─ smeService.GetDistributionByYouth()    → [{label, value, percentage}]
       ├─ smeService.GetAgeGenderDistribution()  → [{ageGroup, male, female, malePercentage, femalePercentage}]
       ├─ smeService.GetDistributionByClassification() → [{label, value, percentage}]
       ├─ bdspService.GetDistributionByStatus()  → [{label, value, percentage}]
       └─ bdspService.GetTopServicesDistribution() → [{label, value, percentage}]
```

## Standard Distribution Query Pattern

Every distribution method follows this template:

```go
func (s *SmeService) GetDistributionByGender() []map[string]interface{} {
    // 1. Define result struct
    var results []struct {
        Label string `gorm:"column:label"`
        Value int64  `gorm:"column:value"`
    }

    // 2. Execute raw SQL with GROUP BY
    facades.Orm().Query().Raw(`
        SELECT COALESCE(p.gender, 'Unknown') as label, COUNT(*) as value
        FROM smes s
        INNER JOIN primary_business_owner p ON p.sme_id = s.id
        WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL
        GROUP BY p.gender
        ORDER BY value DESC
    `).Scan(&results)

    // 3. Calculate total for percentages
    total := int64(0)
    for _, r := range results {
        total += r.Value
    }

    // 4. Build output with percentage
    output := make([]map[string]interface{}, len(results))
    for i, r := range results {
        percentage := float64(0)
        if total > 0 {
            percentage = float64(r.Value) / float64(total) * 100
        }
        output[i] = map[string]interface{}{
            "label":      r.Label,
            "value":      r.Value,
            "percentage": percentage,
        }
    }
    return output
}
```

Key patterns:
- **COALESCE** for null handling (`COALESCE(field, 'Unknown')`)
- **Soft delete filter**: Always include `WHERE deleted_at IS NULL`
- **Server-side percentage**: `(value / total) * 100` computed in Go
- **Return shape**: `[]map[string]interface{}` with `{label, value, percentage}`
- **ORDER BY value DESC**: Default sort by count

## Youth Distribution (Computed Age Group)

Uses SQL `EXTRACT(YEAR FROM AGE(...))` for age calculation:

```go
func (s *SmeService) GetDistributionByYouth() []map[string]interface{} {
    var results []struct {
        Label string `gorm:"column:label"`
        Value int64  `gorm:"column:value"`
    }

    facades.Orm().Query().Raw(`
        SELECT
            CASE
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 18 AND 35
                    THEN 'Youth'
                ELSE 'Non-Youth'
            END as label,
            COUNT(*) as value
        FROM smes s
        INNER JOIN primary_business_owner p ON p.sme_id = s.id
        WHERE s.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND p.date_of_birth IS NOT NULL
        GROUP BY label
        ORDER BY label
    `).Scan(&results)

    // ... calculate percentages (same pattern)
}
```

## Age-Gender Distribution (Crosstab)

Returns a different shape for the population pyramid chart:

```go
func (s *SmeService) GetAgeGenderDistribution() []map[string]interface{} {
    var results []struct {
        AgeGroup string `gorm:"column:age_group"`
        Male     int64  `gorm:"column:male"`
        Female   int64  `gorm:"column:female"`
    }

    facades.Orm().Query().Raw(`
        SELECT
            CASE
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 18 AND 25 THEN '18-25'
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 26 AND 35 THEN '26-35'
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 36 AND 45 THEN '36-45'
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 46 AND 55 THEN '46-55'
                WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 56 AND 65 THEN '56-65'
                ELSE '65+'
            END as age_group,
            SUM(CASE WHEN LOWER(p.gender) = 'male' THEN 1 ELSE 0 END) as male,
            SUM(CASE WHEN LOWER(p.gender) = 'female' THEN 1 ELSE 0 END) as female
        FROM smes s
        INNER JOIN primary_business_owner p ON p.sme_id = s.id
        WHERE s.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND p.date_of_birth IS NOT NULL
        GROUP BY age_group
        ORDER BY age_group
    `).Scan(&results)

    // Calculate percentages per row
    output := make([]map[string]interface{}, len(results))
    for i, r := range results {
        rowTotal := r.Male + r.Female
        malePct, femalePct := float64(0), float64(0)
        if rowTotal > 0 {
            malePct = float64(r.Male) / float64(rowTotal) * 100
            femalePct = float64(r.Female) / float64(rowTotal) * 100
        }
        output[i] = map[string]interface{}{
            "ageGroup":         r.AgeGroup,
            "male":             r.Male,
            "female":           r.Female,
            "malePercentage":   malePct,
            "femalePercentage": femalePct,
        }
    }
    return output
}
```

## Classification Distribution (Custom Sort Order)

Some distributions need a specific display order rather than by count:

```go
func (s *SmeService) GetDistributionByClassification() []map[string]interface{} {
    // ... query as normal

    // Define display order
    orderMap := map[string]int{
        "Micro":        1,
        "Small":        2,
        "Medium":       3,
        "Unclassified": 4,
    }

    // Sort by defined priority
    sort.Slice(output, func(i, j int) bool {
        oi := orderMap[output[i]["label"].(string)]
        oj := orderMap[output[j]["label"].(string)]
        return oi < oj
    })

    return output
}
```

## District Distribution (Top N)

Limit results to top N for readability:

```go
func (s *SmeService) GetDistributionByDistrict() []map[string]interface{} {
    // ... raw SQL with GROUP BY s.district ORDER BY value DESC LIMIT 15
}
```

## BDSP Services Distribution (JSON Array Unpacking)

When data is stored as a JSON array, unpack in SQL:

```go
func (s *BdspService) GetTopServicesDistribution() []map[string]interface{} {
    var results []struct {
        Label string `gorm:"column:label"`
        Value int64  `gorm:"column:value"`
    }

    facades.Orm().Query().Raw(`
        SELECT service as label, COUNT(*) as value
        FROM bdsps,
            jsonb_array_elements_text(service_list_json) as service
        WHERE deleted_at IS NULL
            AND service_list_json IS NOT NULL
        GROUP BY service
        ORDER BY value DESC
        LIMIT 10
    `).Scan(&results)

    // ... calculate percentages
}
```

## Event & Procurement Statistics

Scalar counts (not distributions):

```go
func (s *DashboardService) GetEventStatistics() EventStatistics {
    now := time.Now()
    var totalEvents, upcomingEvents int64

    totalEvents, _ = facades.Orm().Query().Model(&models.Event{}).Count()
    upcomingEvents, _ = facades.Orm().Query().Model(&models.Event{}).
        Where("date > ?", now).Count()

    return EventStatistics{
        TotalEvents:    totalEvents,
        UpcomingEvents: upcomingEvents,
    }
}
```

## Widget DTOs

Dashboard widgets use specific DTOs:

```go
type UpcomingEventDTO struct {
    ID       uint   `json:"id"`
    Title    string `json:"title"`
    Date     string `json:"date"`
    Venue    string `json:"venue"`
    District string `json:"district,omitempty"`
}

type UpcomingProcurementDTO struct {
    ID              uint   `json:"id"`
    Organization    string `json:"organization"`
    RefNo           string `json:"refNo"`
    CloseDate       string `json:"closeDate"`
    ProcurementType string `json:"procurementType,omitempty"`
}

type RecentActivityDTO struct {
    ID         uint   `json:"id"`
    EntityType string `json:"entityType"`   // "sme", "event", "procurement"
    EntityID   uint   `json:"entityId"`
    EntityName string `json:"entityName"`
    Action     string `json:"action"`       // "created", "updated"
    UserID     uint   `json:"userId"`
    UserName   string `json:"userName"`
    Timestamp  string `json:"timestamp"`
}
```

## Controller Pattern

Aggregate everything into a single Inertia render:

```go
func (r *DashboardController) Show(ctx http.Context) http.Response {
    stats := r.dashboardService.GetDashboardStats(r.smeService)
    upcomingEvents := r.dashboardService.GetUpcomingEvents(5)
    upcomingProcurements := r.dashboardService.GetUpcomingProcurements(5)
    recentActivities := r.dashboardService.GetRecentActivities(10)

    return inertia.Render(ctx, "dashboard/Index", map[string]interface{}{
        "pageTitle":            "Dashboard",
        "stats":                stats,
        "upcomingEvents":       upcomingEvents,
        "upcomingProcurements": upcomingProcurements,
        "recentActivities":     recentActivities,
    })
}
```

## Adding a New Distribution

To add a new chart to the dashboard:

1. **Add the service method** in the relevant service file:

```go
func (s *SmeService) GetDistributionByNewField() []map[string]interface{} {
    var results []struct {
        Label string `gorm:"column:label"`
        Value int64  `gorm:"column:value"`
    }
    facades.Orm().Query().Raw(`
        SELECT COALESCE(s.new_field, 'Unknown') as label, COUNT(*) as value
        FROM smes s
        WHERE s.deleted_at IS NULL
        GROUP BY s.new_field
        ORDER BY value DESC
    `).Scan(&results)
    // ... calculate percentages
    return output
}
```

2. **Wire into GetDashboardStats**:

```go
func (s *DashboardService) GetDashboardStats(smeService *SmeService) map[string]interface{} {
    // ... existing calls
    byNewField := smeService.GetDistributionByNewField()

    return map[string]interface{}{
        // ... existing fields
        "byNewField": byNewField,
    }
}
```

3. **Add TypeScript type**:

```typescript
interface DashboardStats {
    // ... existing fields
    byNewField: DistributionPoint[];
}
```

4. **Create the chart component** (see chart-components.md)

5. **Add to the dashboard page** (grid or carousel)

## Source Files

- Dashboard service: `app/services/dashboard_service.go`
- SME service (distributions): `app/services/sme_service.go`
- BDSP service (distributions): `app/services/bdsp_service.go`
- Dashboard controller: `app/http/controllers/dashboard_controller.go`
- Inertia renderer: `app/http/inertia/`
