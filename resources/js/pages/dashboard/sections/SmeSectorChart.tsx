"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, Cell, XAxis, YAxis } from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
} from "@/components/ui/chart";
import { DistributionPoint } from "@/types/sme";

interface SmeSectorChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Use subtle shade variations for bar charts (2-3 shades)
const BAR_COLORS = [
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
];

// Generate chart config dynamically from data
const generateChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "SMEs",
    },
  };

  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_");
    config[key] = {
      label: item.label,
      color: BAR_COLORS[index % BAR_COLORS.length],
    };
  });

  return config;
};

export function SmeSectorChart({ data, isLoading = false }: SmeSectorChartProps) {
  // Transform and sort data by value descending
  const chartData = React.useMemo(() => {
    return [...data]
      .sort((a, b) => b.value - a.value)
      .map((item, index) => ({
        sector: item.label,
        value: item.value,
        percentage: item.percentage,
        fill: BAR_COLORS[index % BAR_COLORS.length],
        // Truncate long labels for Y-axis display
        shortLabel: item.label.length > 20
          ? item.label.substring(0, 17) + "..."
          : item.label,
      }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle>Economic Sector Distribution</CardTitle>
          <CardDescription>SMEs by economic sector</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] w-full animate-pulse rounded-md bg-muted" />
        </CardContent>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle>Economic Sector Distribution</CardTitle>
          <CardDescription>SMEs by economic sector</CardDescription>
        </CardHeader>
        <CardContent className="flex h-[300px] items-center justify-center">
          <p className="text-muted-foreground">No data available</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="flex flex-col">
      <CardHeader>
        <CardTitle>Economic Sector Distribution</CardTitle>
        <CardDescription>SMEs by economic sector (sorted by count)</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[300px] w-full">
          <BarChart
            data={chartData}
            layout="vertical"
            margin={{
              left: 0,
              right: 16,
              top: 0,
              bottom: 0,
            }}
          >
            <CartesianGrid horizontal={false} strokeDasharray="3 3" />
            <YAxis
              dataKey="shortLabel"
              type="category"
              tickLine={false}
              tickMargin={8}
              axisLine={false}
              width={140}
              tick={{ fontSize: 11 }}
            />
            <XAxis
              dataKey="value"
              type="number"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
            />
            <ChartTooltip
              cursor={{ fill: "hsl(var(--muted))", opacity: 0.3 }}
              content={({ active, payload }) => {
                if (!active || !payload || !payload.length) return null;
                const data = payload[0]?.payload;
                if (!data) return null;
                const pct = typeof data.percentage === 'number'
                  ? data.percentage.toFixed(1)
                  : String(data.percentage ?? 0);
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <div className="flex min-w-[150px] flex-col gap-1 text-xs">
                      <div className="flex items-center gap-2">
                        <div
                          className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                          style={{ backgroundColor: data.fill }}
                        />
                        <span className="font-medium text-foreground">
                          {data.sector}
                        </span>
                      </div>
                      <div className="flex items-center justify-between text-muted-foreground">
                        <span>Count:</span>
                        <span className="font-mono font-medium tabular-nums text-foreground">
                          {data.value}
                        </span>
                      </div>
                      <div className="flex items-center justify-between text-muted-foreground">
                        <span>Share:</span>
                        <span className="font-mono font-medium tabular-nums text-foreground">
                          {pct}%
                        </span>
                      </div>
                    </div>
                  </div>
                );
              }}
            />
            <Bar
              dataKey="value"
              radius={[0, 4, 4, 0]}
              maxBarSize={32}
            >
              {chartData.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={entry.fill}
                />
              ))}
            </Bar>
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

export default SmeSectorChart;
