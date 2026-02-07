# Chart Components - Reference Examples

Concrete code from the Books dashboard chart implementations.

## Donut Chart: Gender Distribution

Source: `resources/js/pages/dashboard/sections/SmeGenderChart.tsx`

### Props & Color Setup

```typescript
interface SmeGenderChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Semantic colors with case-insensitive lookup
const GENDER_COLORS: Record<string, string> = {
  male: "hsl(var(--primary))",
  female: "var(--chart-contrast)",
  other: "hsl(var(--chart-5))",
  unknown: "hsl(var(--muted-foreground))",
};

const getGenderColor = (label: string): string | undefined =>
  GENDER_COLORS[label.toLowerCase()];

const FALLBACK_COLORS = [
  "hsl(var(--chart-1))", "hsl(var(--chart-2))",
  "hsl(var(--chart-3))", "hsl(var(--chart-4))", "hsl(var(--chart-5))",
];
```

### Dynamic Config Generation

```typescript
const generateChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = { value: { label: "MSME Owners" } };
  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_");
    config[key] = {
      label: item.label,
      color: getGenderColor(item.label) || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    };
  });
  return config;
};
```

### Data Transform with useMemo

```typescript
const chartData = React.useMemo(() =>
  data.map((item, index) => ({
    gender: item.label.toLowerCase().replace(/\s+/g, "_"),
    label: item.label,
    value: item.value,
    percentage: item.percentage,
    fill: getGenderColor(item.label) || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
  })),
[data]);

const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);
const totalOwners = React.useMemo(() =>
  data.reduce((acc, curr) => acc + curr.value, 0),
[data]);
```

### Full Donut JSX

```tsx
<Card className="flex flex-col">
  <CardHeader className="items-center pb-0">
    <CardTitle>Owner Gender Distribution</CardTitle>
    <CardDescription>Primary SME owners by gender</CardDescription>
  </CardHeader>
  <CardContent className="flex-1 pb-0">
    <ChartContainer config={chartConfig} className="mx-auto aspect-square max-h-[300px]">
      <PieChart>
        <ChartTooltip cursor={false} content={({ active, payload }) => {
          if (!active || !payload?.length) return null;
          const d = payload[0]?.payload;
          if (!d) return null;
          const pct = typeof d.percentage === 'number' ? d.percentage.toFixed(1) : '0';
          return (
            <div className="rounded-lg border bg-background p-2 shadow-sm">
              <div className="flex min-w-[130px] items-center text-xs text-muted-foreground">
                <div className="h-2.5 w-2.5 shrink-0 rounded-[2px] mr-2"
                  style={{ backgroundColor: d.fill }} />
                <span className="flex-1">{d.label}</span>
                <div className="ml-auto flex items-baseline gap-1 font-mono font-medium tabular-nums text-foreground">
                  {d.value}
                  <span className="font-normal text-muted-foreground">({pct}%)</span>
                </div>
              </div>
            </div>
          );
        }} />
        <Pie data={chartData} dataKey="value" nameKey="gender"
          innerRadius={60} outerRadius={100} strokeWidth={2}
          stroke="hsl(var(--background))">
          {chartData.map((entry, i) => <Cell key={`cell-${i}`} fill={entry.fill} />)}
          <Label content={({ viewBox }) => {
            if (viewBox && "cx" in viewBox && "cy" in viewBox) {
              return (
                <text x={viewBox.cx} y={viewBox.cy} textAnchor="middle" dominantBaseline="middle">
                  <tspan x={viewBox.cx} y={viewBox.cy}
                    className="fill-foreground text-3xl font-bold">
                    {totalOwners.toLocaleString()}
                  </tspan>
                  <tspan x={viewBox.cx} y={(viewBox.cy || 0) + 24}
                    className="fill-muted-foreground text-sm">
                    Total Owners
                  </tspan>
                </text>
              );
            }
          }} />
        </Pie>
        <ChartLegend content={<ChartLegendContent nameKey="gender" />}
          className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center" />
      </PieChart>
    </ChartContainer>
  </CardContent>
</Card>
```

---

## Horizontal Bar Chart: Sector Distribution

Source: `resources/js/pages/dashboard/sections/SmeSectorChart.tsx`

### Data Transform (sorted descending, label truncation)

```typescript
const BAR_COLORS = [
  "hsl(var(--chart-2))", "hsl(var(--chart-3))", "hsl(var(--chart-4))",
];

const chartData = React.useMemo(() =>
  [...data]
    .sort((a, b) => b.value - a.value)
    .map((item, index) => ({
      sector: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: BAR_COLORS[index % BAR_COLORS.length],
      shortLabel: item.label.length > 10
        ? item.label.substring(0, 8) + "..." : item.label,
    })),
[data]);
```

### Full Bar Chart JSX

```tsx
<Card className="flex flex-col overflow-hidden">
  <CardHeader>
    <CardTitle>Economic Sector Distribution</CardTitle>
    <CardDescription>MSMEs by economic sector (sorted by count)</CardDescription>
  </CardHeader>
  <CardContent className="overflow-hidden w-full">
    <ChartContainer config={chartConfig} className="h-[300px] w-full max-w-full">
      <BarChart data={chartData} layout="vertical"
        margin={{ left: 0, right: 4, top: 0, bottom: 0 }}>
        <CartesianGrid horizontal={false} strokeDasharray="3 3" />
        <YAxis dataKey="shortLabel" type="category" tickLine={false}
          tickMargin={2} axisLine={false} width={70} tick={{ fontSize: 9 }} />
        <XAxis dataKey="value" type="number" tickLine={false}
          axisLine={false} tickMargin={8} tick={{ fontSize: 10 }} />
        <ChartTooltip
          cursor={{ fill: "hsl(var(--muted))", opacity: 0.3 }}
          wrapperStyle={{ zIndex: 1000 }}
          allowEscapeViewBox={{ x: false, y: true }}
          content={({ active, payload }) => {
            if (!active || !payload?.length) return null;
            const d = payload[0]?.payload;
            if (!d) return null;
            const pct = typeof d.percentage === 'number' ? d.percentage.toFixed(1) : '0';
            return (
              <div className="rounded-lg border bg-background p-2 shadow-sm max-w-[200px]">
                <div className="flex flex-col gap-1 text-xs">
                  <div className="flex items-center gap-2">
                    <div className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                      style={{ backgroundColor: d.fill }} />
                    <span className="font-medium text-foreground truncate">{d.sector}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground gap-2">
                    <span>Count:</span>
                    <span className="font-mono font-medium tabular-nums text-foreground">{d.value}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground gap-2">
                    <span>Share:</span>
                    <span className="font-mono font-medium tabular-nums text-foreground">{pct}%</span>
                  </div>
                </div>
              </div>
            );
          }}
        />
        <Bar dataKey="value" radius={[0, 4, 4, 0]} maxBarSize={32}>
          {chartData.map((entry, i) => <Cell key={`cell-${i}`} fill={entry.fill} />)}
        </Bar>
      </BarChart>
    </ChartContainer>
  </CardContent>
</Card>
```

---

## Population Pyramid: Age & Gender

Source: `resources/js/pages/dashboard/sections/SmeAgePyramidChart.tsx`

### Data Transform (male values negated)

```typescript
const MALE_COLOR = "hsl(var(--primary))";
const FEMALE_COLOR = "var(--chart-contrast)";

const chartData = React.useMemo(() =>
  data.map(item => ({
    ageGroup: item.ageGroup,
    male: -item.male,            // Negative = left side
    female: item.female,         // Positive = right side
    maleActual: item.male,       // Keep original for tooltip
    femaleActual: item.female,
    malePercentage: item.malePercentage,
    femalePercentage: item.femalePercentage,
  })),
[data]);

// Symmetric axis
const maxValue = React.useMemo(() => {
  let max = 0;
  data.forEach(item => { max = Math.max(max, item.male, item.female); });
  return Math.ceil(max * 1.1);
}, [data]);
```

### Custom Legend (not using ChartLegend)

```tsx
<div className="flex items-center justify-center gap-6 mb-2">
  <div className="flex items-center gap-1.5">
    <div className="h-3 w-3 rounded-sm" style={{ backgroundColor: MALE_COLOR }} />
    <span className="text-xs text-muted-foreground">Male ({totals.male})</span>
  </div>
  <div className="flex items-center gap-1.5">
    <div className="h-3 w-3 rounded-sm" style={{ backgroundColor: FEMALE_COLOR }} />
    <span className="text-xs text-muted-foreground">Female ({totals.female})</span>
  </div>
</div>
```

### Multi-series Tooltip

```tsx
content={({ active, payload }) => {
  if (!active || !payload?.length) return null;
  const d = payload[0]?.payload;
  return (
    <div className="rounded-lg border bg-background p-2 shadow-sm min-w-[160px]">
      <div className="font-medium text-sm mb-1.5">{d.ageGroup}</div>
      <div className="flex items-center text-xs text-muted-foreground mb-1">
        <div className="h-2.5 w-2.5 rounded-[2px] mr-2" style={{ backgroundColor: MALE_COLOR }} />
        <span className="flex-1">Male</span>
        <span className="ml-auto font-mono font-medium tabular-nums text-foreground">
          {d.maleActual} <span className="font-normal text-muted-foreground">({d.malePercentage?.toFixed(1)}%)</span>
        </span>
      </div>
      <div className="flex items-center text-xs text-muted-foreground">
        <div className="h-2.5 w-2.5 rounded-[2px] mr-2" style={{ backgroundColor: FEMALE_COLOR }} />
        <span className="flex-1">Female</span>
        <span className="ml-auto font-mono font-medium tabular-nums text-foreground">
          {d.femaleActual} <span className="font-normal text-muted-foreground">({d.femalePercentage?.toFixed(1)}%)</span>
        </span>
      </div>
    </div>
  );
}}
```

---

## KPI Cards with Trend Indicator

Source: `resources/js/pages/dashboard/sections/DashboardKpiCards.tsx`

### Trend Calculation

```typescript
function calculatePercentageChange(current: number, previous: number): number {
  if (previous === 0) return current > 0 ? 100 : 0;
  return Math.round(((current - previous) / previous) * 100);
}
```

### Card with Trend

```tsx
<Card>
  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
    <CardTitle className="text-sm font-medium text-muted-foreground">New This Month</CardTitle>
    <Users className="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div className="flex items-center justify-between">
      <div className="text-2xl font-bold">{stats.newThisMonth.toLocaleString()}</div>
      <TrendIndicator value={monthOverMonthChange} />
    </div>
    <p className="text-xs text-muted-foreground mt-1">
      vs {stats.newLastMonth.toLocaleString()} last month
    </p>
  </CardContent>
</Card>
```

### Card with Badge Count

```tsx
<Card>
  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
    <CardTitle className="text-sm font-medium text-muted-foreground">Events</CardTitle>
    <Calendar className="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div className="flex items-center justify-between">
      <div className="text-2xl font-bold">{stats.totalEvents.toLocaleString()}</div>
      {stats.upcomingEvents > 0 && (
        <span className="text-xs font-medium text-emerald-600">
          {stats.upcomingEvents} upcoming
        </span>
      )}
    </div>
  </CardContent>
</Card>
```

### Loading Skeleton

```tsx
function KpiCardSkeleton() {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="h-4 w-24 animate-pulse rounded bg-muted" />
        <div className="h-4 w-4 animate-pulse rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="h-8 w-16 animate-pulse rounded bg-muted" />
        <div className="mt-2 h-3 w-32 animate-pulse rounded bg-muted" />
      </CardContent>
    </Card>
  );
}
```

---

## Carousel with ChartActions Integration

Source: `resources/js/pages/dashboard/sections/SmeOwnerChartsCarousel.tsx`

### CSV Data Preparation

```typescript
const csvData = React.useMemo(() =>
  genderData.map(item => ({
    "Gender": item.label,
    "Count": item.value,
    "Percentage (%)": item.percentage.toFixed(1),
  })),
[genderData]);
```

### Fullscreen Render Callback

```typescript
const renderFullscreen = React.useCallback((width: number, height: number) => {
  const chartHeight = height - 20;
  if (currentIndex === 0) {
    return <div style={{ width, height: chartHeight }}>
      <GenderChartContent data={genderData} height={chartHeight} />
    </div>;
  } else if (currentIndex === 1) {
    return <div style={{ width, height: chartHeight }}>
      <YouthChart data={youthData} height={chartHeight} />
    </div>;
  }
  // ...
}, [currentIndex, genderData, youthData]);
```

### Dynamic Radius for Fullscreen

Charts that accept a `height` prop scale their dimensions:

```typescript
const { innerRadius, outerRadius } = React.useMemo(() => {
  if (height && height > 400) {
    const baseSize = Math.min(height * 0.35, 250);
    return { innerRadius: Math.floor(baseSize * 0.6), outerRadius: Math.floor(baseSize) };
  }
  return { innerRadius: 60, outerRadius: 100 }; // Default
}, [height]);
```

## Source Files

- Dashboard page: `resources/js/pages/dashboard/Index.tsx`
- KPI cards: `resources/js/pages/dashboard/sections/DashboardKpiCards.tsx`
- Gender chart: `resources/js/pages/dashboard/sections/SmeGenderChart.tsx`
- Sector chart: `resources/js/pages/dashboard/sections/SmeSectorChart.tsx`
- Age pyramid: `resources/js/pages/dashboard/sections/SmeAgePyramidChart.tsx`
- Carousel: `resources/js/pages/dashboard/sections/SmeOwnerChartsCarousel.tsx`
- Chart UI: `resources/js/components/ui/chart.tsx`
- Chart actions: `resources/js/components/ui/chart-actions.tsx`
- Types: `resources/js/types/sme.ts`
