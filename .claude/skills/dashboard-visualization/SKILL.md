---
name: dashboard-visualization
description: Guide for building dashboards, charts, KPI cards, and data visualizations. Use when creating dashboard pages, adding charts to features, building stat displays, aggregating data for visual presentation, or implementing export (CSV/PNG/fullscreen) for charts. Based on the Books dashboard implementation using Recharts and shadcn/ui chart components.
---

# Dashboard & Visualization Pattern

This skill defines how to build dashboards, charts, and data visualizations in this codebase. Follow these patterns when adding visual analytics to any feature or building new dashboard pages.

## Architecture Overview

```
Go Backend                          React Frontend
┌─────────────────────┐             ┌──────────────────────────────┐
│ DashboardController  │──Inertia──▶│ DashboardPage (Index.tsx)    │
│   └─ DashboardService│            │   ├─ KPI Cards Row           │
│       ├─ SmeService   │            │   ├─ Charts Grid (2-col)     │
│       ├─ BdspService  │            │   │   ├─ Carousel (charts)   │
│       └─ Raw SQL agg  │            │   │   └─ Carousel (charts)   │
└─────────────────────┘             │   └─ Sidebar (widgets)       │
                                    └──────────────────────────────┘
```

Data flows via **Inertia props** (not REST API calls). The Go controller aggregates all data in `GetDashboardStats()` and passes it as page props. The frontend receives it via `usePage<Props>()`.

## Key Principles

1. **Server-side aggregation** - All counts, percentages, and distributions are computed in Go services using raw SQL. The frontend only renders.
2. **Shadcn chart wrapper** - All charts use `ChartContainer` + `ChartConfig` from `@/components/ui/chart` which wraps Recharts with theme support.
3. **Standard data shape** - Distribution data uses `DistributionPoint { label, value, percentage }`. All charts consume this shape.
4. **Every chart gets ChartActions** - CSV export, PNG capture, and fullscreen view via the reusable `ChartActions` component.
5. **Permission-gated sections** - Dashboard sections are wrapped in `PermissionGate` or conditional renders based on `canPerformAction()`.
6. **Carousel grouping** - Related charts are grouped into carousels (auto-rotating, pauseable) to save space.
7. **Responsive layout** - `xl:flex-row` for main + sidebar, `lg:grid-cols-2` for chart pairs, stacked on mobile.

## File Structure

```
resources/js/pages/dashboard/
  Index.tsx                           # Main dashboard page
  SmeDashboard.tsx                    # Role-specific dashboard variant
  sections/
    index.ts                          # Barrel exports
    DashboardKpiCards.tsx             # KPI stat cards row
    Sme{Name}Chart.tsx               # Individual chart components
    Sme{Name}ChartsCarousel.tsx       # Carousel wrappers grouping related charts
    Bdsp{Name}Chart.tsx               # BDSP-specific charts

app/services/
  dashboard_service.go                # Widget data, event/procurement stats
  sme_service.go                      # SME distribution queries (GetDistributionBy*)
  bdsp_service.go                     # BDSP distribution queries

app/http/controllers/
  dashboard_controller.go             # Aggregates all stats, renders via Inertia
```

## Standard Data Types

### DistributionPoint (the universal chart data shape)

```typescript
// resources/js/types/sme.ts
export interface DistributionPoint {
  label: string;      // Category name (e.g., "Male", "Agriculture")
  value: number;      // Count
  percentage: number;  // Pre-calculated server-side
}

export interface AgeGenderDistributionPoint {
  ageGroup: string;       // e.g., "18-25", "26-35"
  male: number;
  female: number;
  malePercentage: number;
  femalePercentage: number;
}
```

### DashboardStats (Inertia page props)

```typescript
interface DashboardStats {
  // Scalar KPIs
  totalSmes: number;
  newThisMonth: number;
  newLastMonth: number;
  totalEvents: number;
  upcomingEvents: number;
  totalProcurements: number;
  activeProcurements: number;
  // Distribution arrays
  byRegion: DistributionPoint[];
  byDistrict: DistributionPoint[];
  byCategory: DistributionPoint[];
  bySector: DistributionPoint[];
  byGender: DistributionPoint[];
  byYouth: DistributionPoint[];
  byAgeGender: AgeGenderDistributionPoint[];
  byClassification: DistributionPoint[];
  bdspByStatus: DistributionPoint[];
  bdspByServices: DistributionPoint[];
}
```

## Backend: Data Aggregation Patterns

### Service Method Pattern

Each distribution query follows this pattern in Go:

```go
// In sme_service.go
func (s *SmeService) GetDistributionByGender() []map[string]interface{} {
    var results []struct {
        Label string `gorm:"column:label"`
        Value int64  `gorm:"column:value"`
    }

    facades.Orm().Query().Raw(`
        SELECT COALESCE(p.gender, 'Unknown') as label, COUNT(*) as value
        FROM smes s
        INNER JOIN primary_business_owner p ON p.sme_id = s.id
        WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL
        GROUP BY p.gender
        ORDER BY value DESC
    `).Scan(&results)

    // Calculate percentages
    total := int64(0)
    for _, r := range results { total += r.Value }

    output := make([]map[string]interface{}, len(results))
    for i, r := range results {
        output[i] = map[string]interface{}{
            "label":      r.Label,
            "value":      r.Value,
            "percentage": float64(r.Value) / float64(total) * 100,
        }
    }
    return output
}
```

### Dashboard Controller Pattern

The controller aggregates all service calls and passes to Inertia:

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

### GetDashboardStats aggregation

```go
func (s *DashboardService) GetDashboardStats(smeService *SmeService) map[string]interface{} {
    smeStats, _ := smeService.GetSmeStatistics()
    eventStats := s.GetEventStatistics()
    procurementStats := s.GetProcurementStatistics()
    bySector := smeService.GetDistributionBySector()
    byGender := smeService.GetDistributionByGender()
    // ... more distributions

    return map[string]interface{}{
        "totalSmes":     smeStats["totalSmes"],
        "newThisMonth":  smeStats["newThisMonth"],
        "bySector":      bySector,
        "byGender":      byGender,
        // ... combine all stats
    }
}
```

## Frontend: Chart Components

### Component Imports

```typescript
// Recharts primitives
import { Bar, BarChart, CartesianGrid, Cell, XAxis, YAxis } from "recharts";
import { Pie, PieChart, Label } from "recharts";
import { ReferenceLine, ResponsiveContainer } from "recharts";

// Shadcn chart wrappers (theme-aware)
import {
  ChartConfig, ChartContainer, ChartTooltip,
  ChartLegend, ChartLegendContent, ChartTooltipContent,
} from "@/components/ui/chart";

// Export actions
import { ChartActions } from "@/components/ui/chart-actions";

// Data type
import { DistributionPoint } from "@/types/sme";
```

### ChartConfig (Theme-Aware Colors)

Every chart needs a `ChartConfig` object that maps data keys to labels and colors:

```typescript
const chartConfig: ChartConfig = {
  value: { label: "MSMEs" },
  male: { label: "Male", color: "hsl(var(--primary))" },
  female: { label: "Female", color: "var(--chart-contrast)" },
};
```

For dynamic data, generate the config from the data array:

```typescript
const generateChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = { value: { label: "MSMEs" } };
  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_");
    config[key] = {
      label: item.label,
      color: COLORS[key] || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    };
  });
  return config;
};
```

### Color System

| Variable | Purpose | Usage |
|----------|---------|-------|
| `hsl(var(--primary))` | Main brand blue | Male, primary series |
| `var(--chart-contrast)` | Pink/accent | Female, secondary series |
| `hsl(var(--chart-1))` through `--chart-5` | Palette cycle | Multi-category charts |
| `hsl(var(--muted-foreground))` | Gray | Unknown/Other |

Semantic color maps for known categories:

```typescript
const GENDER_COLORS: Record<string, string> = {
  male: "hsl(var(--primary))",
  female: "var(--chart-contrast)",
  other: "hsl(var(--chart-5))",
  unknown: "hsl(var(--muted-foreground))",
};

const STATUS_COLORS: Record<string, string> = {
  active: "hsl(160, 84%, 39%)",    // Emerald
  inactive: "hsl(var(--muted-foreground))",
  pending: "hsl(45, 93%, 47%)",     // Amber
  rejected: "hsl(0, 84%, 60%)",     // Red
};
```

Always use **case-insensitive lookup**: `COLORS[label.toLowerCase()]`.

## Chart Type Recipes

### Donut Chart (Pie with center label)

```tsx
<ChartContainer config={chartConfig} className="mx-auto aspect-square max-h-[300px]">
  <PieChart>
    <ChartTooltip cursor={false} content={<CustomTooltip />} />
    <Pie data={chartData} dataKey="value" nameKey="category"
      innerRadius={60} outerRadius={100} strokeWidth={2}
      stroke="hsl(var(--background))">
      {chartData.map((entry, i) => <Cell key={i} fill={entry.fill} />)}
      <Label content={({ viewBox }) => {
        if (viewBox && "cx" in viewBox && "cy" in viewBox) {
          return (
            <text x={viewBox.cx} y={viewBox.cy} textAnchor="middle" dominantBaseline="middle">
              <tspan x={viewBox.cx} y={viewBox.cy} className="fill-foreground text-3xl font-bold">
                {total.toLocaleString()}
              </tspan>
              <tspan x={viewBox.cx} y={(viewBox.cy || 0) + 24} className="fill-muted-foreground text-sm">
                Total
              </tspan>
            </text>
          );
        }
      }} />
    </Pie>
    <ChartLegend content={<ChartLegendContent nameKey="category" />}
      className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center" />
  </PieChart>
</ChartContainer>
```

### Horizontal Bar Chart (sorted descending)

```tsx
<ChartContainer config={chartConfig} className="h-[300px] w-full">
  <BarChart data={chartData} layout="vertical" margin={{ left: 0, right: 4 }}>
    <CartesianGrid horizontal={false} strokeDasharray="3 3" />
    <YAxis dataKey="shortLabel" type="category" tickLine={false} axisLine={false}
      width={70} tick={{ fontSize: 9 }} />
    <XAxis dataKey="value" type="number" tickLine={false} axisLine={false}
      tick={{ fontSize: 10 }} />
    <ChartTooltip cursor={{ fill: "hsl(var(--muted))", opacity: 0.3 }}
      content={<CustomTooltip />} />
    <Bar dataKey="value" radius={[0, 4, 4, 0]} maxBarSize={32}>
      {chartData.map((entry, i) => <Cell key={i} fill={entry.fill} />)}
    </Bar>
  </BarChart>
</ChartContainer>
```

### Population Pyramid (mirrored bar chart)

```tsx
// Male values are negated for left-side display
const chartData = data.map(item => ({
  ageGroup: item.ageGroup,
  male: -item.male,     // Negative = extends left
  female: item.female,  // Positive = extends right
  maleActual: item.male,
  femaleActual: item.female,
}));

<BarChart data={chartData} layout="vertical" barCategoryGap="15%">
  <XAxis type="number" domain={[-maxValue, maxValue]}
    tickFormatter={(v) => Math.abs(v).toString()} />
  <YAxis type="category" dataKey="ageGroup" width={50} />
  <ReferenceLine x={0} stroke="hsl(var(--border))" />
  <Bar dataKey="male" radius={[4, 0, 0, 4]}>
    {chartData.map((_, i) => <Cell key={i} fill={MALE_COLOR} />)}
  </Bar>
  <Bar dataKey="female" radius={[0, 4, 4, 0]}>
    {chartData.map((_, i) => <Cell key={i} fill={FEMALE_COLOR} />)}
  </Bar>
</BarChart>
```

## KPI Cards

```tsx
function DashboardKpiCards({ stats, canViewSmes, canViewEvents, canViewProcurements }) {
  const cards: React.ReactNode[] = [];

  if (canViewSmes) {
    cards.push(
      <Card key="total">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">Total MSMEs</CardTitle>
          <Building2 className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.totalSmes.toLocaleString()}</div>
          <p className="text-xs text-muted-foreground mt-1">Registered in the database</p>
        </CardContent>
      </Card>
    );
  }
  // ... more permission-gated cards

  // Dynamic grid columns based on visible card count
  const gridCols = cards.length <= 2 ? "sm:grid-cols-2" : "sm:grid-cols-2 xl:grid-cols-4";
  return <div className={`grid grid-cols-1 gap-4 ${gridCols}`}>{cards}</div>;
}
```

### Trend Indicator

```tsx
function TrendIndicator({ value }: { value: number }) {
  const isPositive = value >= 0;
  const Icon = isPositive ? TrendingUp : TrendingDown;
  return (
    <div className={cn("flex items-center gap-1 text-xs font-medium",
      isPositive ? "text-emerald-600" : "text-red-600")}>
      <Icon className="h-3 w-3" />
      <span>{isPositive ? "+" : ""}{value}%</span>
    </div>
  );
}

// Calculate: ((current - previous) / previous) * 100
const change = calculatePercentageChange(stats.newThisMonth, stats.newLastMonth);
```

## Custom Tooltip Pattern

All chart tooltips follow this structure:

```tsx
content={({ active, payload }) => {
  if (!active || !payload || !payload.length) return null;
  const data = payload[0]?.payload;
  if (!data) return null;
  const pct = typeof data.percentage === 'number'
    ? data.percentage.toFixed(1) : String(data.percentage ?? 0);
  return (
    <div className="rounded-lg border bg-background p-2 shadow-sm">
      <div className="flex min-w-[130px] items-center text-xs text-muted-foreground">
        <div className="h-2.5 w-2.5 shrink-0 rounded-[2px] mr-2"
          style={{ backgroundColor: data.fill }} />
        <span className="flex-1">{data.label}</span>
        <div className="ml-auto flex items-baseline gap-1 font-mono font-medium tabular-nums text-foreground">
          {data.value}
          <span className="font-normal text-muted-foreground">({pct}%)</span>
        </div>
      </div>
    </div>
  );
}}
```

## ChartActions (Export & Fullscreen)

Every chart card should include `ChartActions` in the header:

```tsx
const chartRef = React.useRef<HTMLDivElement>(null);

const csvData = React.useMemo(() =>
  data.map(item => ({
    "Category": item.label,
    "Count": item.value,
    "Percentage (%)": item.percentage.toFixed(1),
  })),
[data]);

const renderFullscreen = React.useCallback((width: number, height: number) => (
  <div style={{ width, height: height - 20 }}>
    <MyChartContent data={data} height={height - 20} />
  </div>
), [data]);

<CardHeader className="items-center pb-0 relative">
  <div className="absolute right-4 top-4">
    <ChartActions
      chartRef={chartRef}
      data={csvData}
      filename="my-chart-export"
      title="My Chart Title"
      renderFullscreen={renderFullscreen}
    />
  </div>
  <CardTitle>My Chart</CardTitle>
  <CardDescription>Description here</CardDescription>
</CardHeader>

<CardContent ref={chartRef}>
  {/* chart content */}
</CardContent>
```

**ChartActions provides:**
- **CSV Download** - Converts data array to CSV with app name/title/date header
- **PNG Download** - Uses `html2canvas` with dark gradient background, 2x scale
- **Fullscreen** - Opens in `Dialog` at 95vw x 90vh with `ResizeObserver`

## Carousel Pattern (Grouping Related Charts)

```tsx
function MyChartsCarousel({ dataA, dataB, dataC }) {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const chartRefs = [useRef(null), useRef(null), useRef(null)];

  const chartInfo = [
    { title: "Chart A", description: "Description A" },
    { title: "Chart B", description: "Description B" },
    { title: "Chart C", description: "Description C" },
  ];

  return (
    <Card className="flex flex-col overflow-hidden">
      <CardHeader className="items-center pb-0 relative">
        <div className="absolute right-4 top-4">
          <ChartActions chartRef={chartRefs[currentIndex]} data={csvData[currentIndex]}
            filename={filenames[currentIndex]} title={chartInfo[currentIndex].title}
            renderFullscreen={renderFullscreen} />
        </div>
        <CardTitle>{chartInfo[currentIndex].title}</CardTitle>
        <CardDescription>{chartInfo[currentIndex].description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-4 overflow-hidden">
        <CarouselTracked autoRotate autoRotateInterval={7000}
          showArrows showDots currentIndex={currentIndex} onIndexChange={setCurrentIndex}>
          <div ref={chartRefs[0]} className="h-[340px]"><ChartA data={dataA} /></div>
          <div ref={chartRefs[1]} className="h-[340px]"><ChartB data={dataB} /></div>
          <div ref={chartRefs[2]} className="h-[340px]"><ChartC data={dataC} /></div>
        </CarouselTracked>
      </CardContent>
    </Card>
  );
}
```

Carousel features: auto-rotate (7s), pause on hover, arrow + dot navigation, CSS `translateX` transitions.

## Dashboard Page Layout

```tsx
<div className="p-4 md:p-6 flex flex-col gap-6">
  {/* Welcome header */}
  <div className="flex flex-col gap-1">
    <h1 className="text-2xl font-semibold tracking-tight">Welcome back, {name}!</h1>
    <p className="text-muted-foreground hidden md:inline">Full description</p>
    <p className="text-muted-foreground md:hidden">Short description</p>
  </div>

  {/* Main + Sidebar */}
  <div className="flex flex-col xl:flex-row gap-6 items-start">
    {/* Main content */}
    <div className="flex-1 flex flex-col gap-6 min-w-0">
      <KpiCards stats={stats} />
      <div className="grid gap-4 lg:grid-cols-2">
        <ChartCarouselA />
        <ChartCarouselB />
      </div>
      <div className="grid gap-4 lg:grid-cols-2">
        <ChartCarouselC />
        <ChartCarouselD />
      </div>
    </div>

    {/* Sidebar */}
    <aside className="w-full xl:w-[380px] 2xl:w-[420px] flex flex-col gap-4 shrink-0">
      <PermissionGate service="events" action="read">
        <UpcomingEventsWidget events={events} />
      </PermissionGate>
      <PermissionGate service="procurement_notices" action="read">
        <UpcomingProcurementsWidget procurements={procurements} />
      </PermissionGate>
    </aside>
  </div>
</div>
```

## Loading & Empty States

```tsx
// Loading skeleton
if (isLoading) {
  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Chart Title</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-[250px] w-full animate-pulse rounded-md bg-muted" />
      </CardContent>
    </Card>
  );
}

// Empty state
if (!data || data.length === 0) {
  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Chart Title</CardTitle>
      </CardHeader>
      <CardContent className="flex h-[300px] items-center justify-center">
        <p className="text-muted-foreground">No data available</p>
      </CardContent>
    </Card>
  );
}
```

## Performance Patterns

- `React.useMemo()` for chart data transforms, config generation, and total calculations
- `React.useCallback()` for fullscreen render functions and carousel event handlers
- `ResizeObserver` for dynamic fullscreen sizing (not polling)
- Cleanup: `clearInterval` in carousel `useEffect` returns

## Checklist

When building new visualizations:

- [ ] Backend: Add distribution query in the relevant service (raw SQL, returns `{label, value, percentage}[]`)
- [ ] Backend: Wire into `GetDashboardStats()` or create new controller method
- [ ] Frontend: Create chart component accepting `data: DistributionPoint[]` + `isLoading?: boolean`
- [ ] Frontend: Generate `ChartConfig` dynamically from data with semantic + fallback colors
- [ ] Frontend: Transform data with `React.useMemo()` including `fill` color per entry
- [ ] Frontend: Implement loading skeleton (`animate-pulse`) and empty state
- [ ] Frontend: Add custom tooltip showing label, count, and percentage
- [ ] Frontend: Wrap in `Card` with `CardHeader` (title + description) and `CardContent`
- [ ] Frontend: Add `ChartActions` with CSV data, chart ref, filename, and fullscreen renderer
- [ ] Frontend: Group related charts into a carousel if 2+ charts share a theme
- [ ] Frontend: Gate visibility behind `PermissionGate` or `canPerformAction()`
- [ ] Frontend: Ensure responsive layout (`lg:grid-cols-2`, stacked on mobile)

For detailed code examples, see [examples/chart-components.md](examples/chart-components.md) and [examples/backend-aggregation.md](examples/backend-aggregation.md).
