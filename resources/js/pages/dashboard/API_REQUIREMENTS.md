# Dashboard API Requirements

This document outlines the API endpoints and data structures required for the Dashboard page.

## Overview

The Dashboard page (`/dashboard`) requires the following data from the backend:
1. Dashboard statistics (KPIs and chart data)
2. Upcoming events list
3. Upcoming procurement notices list

---

## 1. Dashboard Statistics

### Endpoint
The dashboard statistics should be passed as props to the Dashboard Inertia page.

### Data Structure

```typescript
interface DashboardStats {
  // KPI Card Data
  totalSmes: number;           // Total count of all SMEs
  newThisMonth: number;        // SMEs created in current month
  newLastMonth: number;        // SMEs created in previous month (for trend calculation)
  totalEvents: number;         // Total count of all events
  upcomingEvents: number;      // Count of events where date > now
  totalProcurements: number;   // Total count of all procurement notices
  activeProcurements: number;  // Count of published procurements where close_date > now

  // Distribution Data for Charts
  byRegion: DistributionPoint[];    // SME distribution by region
  byCategory: DistributionPoint[];  // SME distribution by business_category
  bySector: DistributionPoint[];    // SME distribution by sector field
  byGender: DistributionPoint[];    // SME owner gender distribution
}

interface DistributionPoint {
  label: string;      // Display label (e.g., "Northern", "Agriculture")
  value: number;      // Count of items
  percentage: number; // Percentage of total (0-100)
}
```

### Backend Implementation Notes

#### For `byRegion`:
- Group SMEs by their `region` field (inferred from district)
- Expected labels: "Northern", "Central", "Southern"

#### For `byCategory`:
- Group SMEs by their `businessCategory` field
- Return top categories sorted by count

#### For `bySector` (NEW - needs implementation):
- Group SMEs by their `sector` field
- Return all sectors with counts
- Example sectors from the SME model: Agriculture, Manufacturing, Services, etc.

#### For `byGender` (NEW - needs implementation):
- Query through the `PrimaryBusinessOwner` relationship
- Group by the gender field on the owner
- Expected labels: "Male", "Female", "Other" (if applicable)
- Example query pseudocode:
  ```
  SELECT gender, COUNT(*) as count
  FROM primary_business_owners
  JOIN smes ON smes.id = primary_business_owners.sme_id
  WHERE smes.deleted_at IS NULL
  GROUP BY gender
  ```

---

## 2. Upcoming Events

### Purpose
Display the next 5 upcoming events in the dashboard sidebar widget.

### Expected Props Data

```typescript
interface UpcomingEvent {
  id: number;
  title: string;
  date: string;        // ISO date string (YYYY-MM-DD or full ISO)
  venue: string;
  district?: string;   // Optional, for location display
}

// Pass as: upcomingEvents: UpcomingEvent[]
```

### Backend Query
```sql
SELECT id, title, date, venue, district
FROM events
WHERE date >= CURRENT_DATE
  AND deleted_at IS NULL
ORDER BY date ASC
LIMIT 5
```

---

## 3. Upcoming Procurement Notices

### Purpose
Display active procurement notices that are still open for applications.

### Expected Props Data

```typescript
interface UpcomingProcurement {
  id: number;
  organization: string;
  ref_no: string;
  close_date: string;       // ISO date string
  procurement_type?: string; // Optional, for display badge
  invitation?: string;       // Optional
}

// Pass as: upcomingProcurements: UpcomingProcurement[]
```

### Backend Query
```sql
SELECT id, organization, ref_no, close_date, procurement_type, invitation
FROM procurement_notices
WHERE is_published = true
  AND close_date >= CURRENT_DATE
  AND deleted_at IS NULL
ORDER BY close_date ASC
LIMIT 5
```

---

## 4. Complete Page Props Structure

The Dashboard controller should return these props to the Inertia page:

```php
// DashboardController.php

public function index()
{
    return Inertia::render('dashboard/Index', [
        'stats' => [
            'totalSmes' => Sme::count(),
            'newThisMonth' => Sme::whereMonth('created_at', now()->month)
                                ->whereYear('created_at', now()->year)
                                ->count(),
            'newLastMonth' => Sme::whereMonth('created_at', now()->subMonth()->month)
                                 ->whereYear('created_at', now()->subMonth()->year)
                                 ->count(),
            'totalEvents' => Event::count(),
            'upcomingEvents' => Event::where('date', '>=', now())->count(),
            'totalProcurements' => ProcurementNotice::count(),
            'activeProcurements' => ProcurementNotice::where('is_published', true)
                                                      ->where('close_date', '>=', now())
                                                      ->count(),
            'byRegion' => $this->getSmeDistributionByRegion(),
            'byCategory' => $this->getSmeDistributionByCategory(),
            'bySector' => $this->getSmeDistributionBySector(),
            'byGender' => $this->getSmeDistributionByGender(),
        ],
        'upcomingEvents' => Event::where('date', '>=', now())
                                  ->orderBy('date')
                                  ->take(5)
                                  ->get(['id', 'title', 'date', 'venue', 'district']),
        'upcomingProcurements' => ProcurementNotice::where('is_published', true)
                                                    ->where('close_date', '>=', now())
                                                    ->orderBy('close_date')
                                                    ->take(5)
                                                    ->get(['id', 'organization', 'ref_no', 'close_date', 'procurement_type']),
    ]);
}

private function getSmeDistributionByRegion(): array
{
    $total = Sme::count();
    if ($total === 0) return [];

    return Sme::select('region', DB::raw('count(*) as count'))
              ->groupBy('region')
              ->get()
              ->map(fn($item) => [
                  'label' => $item->region,
                  'value' => $item->count,
                  'percentage' => round(($item->count / $total) * 100, 1),
              ])
              ->toArray();
}

private function getSmeDistributionByCategory(): array
{
    $total = Sme::count();
    if ($total === 0) return [];

    return Sme::select('business_category', DB::raw('count(*) as count'))
              ->groupBy('business_category')
              ->orderByDesc('count')
              ->get()
              ->map(fn($item) => [
                  'label' => $item->business_category,
                  'value' => $item->count,
                  'percentage' => round(($item->count / $total) * 100, 1),
              ])
              ->toArray();
}

private function getSmeDistributionBySector(): array
{
    $total = Sme::count();
    if ($total === 0) return [];

    return Sme::select('sector', DB::raw('count(*) as count'))
              ->groupBy('sector')
              ->orderByDesc('count')
              ->get()
              ->map(fn($item) => [
                  'label' => $item->sector,
                  'value' => $item->count,
                  'percentage' => round(($item->count / $total) * 100, 1),
              ])
              ->toArray();
}

private function getSmeDistributionByGender(): array
{
    $total = PrimaryBusinessOwner::count();
    if ($total === 0) return [];

    return PrimaryBusinessOwner::select('gender', DB::raw('count(*) as count'))
                               ->groupBy('gender')
                               ->orderByDesc('count')
                               ->get()
                               ->map(fn($item) => [
                                   'label' => $item->gender ?? 'Unknown',
                                   'value' => $item->count,
                                   'percentage' => round(($item->count / $total) * 100, 1),
                               ])
                               ->toArray();
}
```

---

## Summary of New Data Requirements

| Data Point | Source | Status |
|------------|--------|--------|
| `stats.byRegion` | Existing SME stats | Already implemented |
| `stats.byCategory` | Existing SME stats | Already implemented |
| `stats.bySector` | SME sector field | **NEW - Needs implementation** |
| `stats.byGender` | PrimaryBusinessOwner.gender | **NEW - Needs implementation** |
| `upcomingEvents` | Events table | **NEW - Needs implementation** |
| `upcomingProcurements` | ProcurementNotices table | **NEW - Needs implementation** |
| KPI totals | Various tables | **NEW - Needs implementation** |

---

## Testing Data

For development/testing, the frontend will show "No data available" states when these props are not provided. Once the backend implements these endpoints, the dashboard will automatically populate with real data.
