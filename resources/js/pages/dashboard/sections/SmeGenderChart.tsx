"use client";

import * as React from "react";
import { Cell, Label, Pie, PieChart } from "recharts";
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
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
} from "@/components/ui/chart";
import { DistributionPoint } from "@/types/sme";

interface SmeGenderChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Define specific colors for gender categories (lowercase keys for case-insensitive lookup)
const GENDER_COLORS: Record<string, string> = {
  male: "hsl(var(--primary))",
  female: "var(--chart-contrast)",
  other: "hsl(var(--chart-5))",
  unknown: "hsl(var(--muted-foreground))",
};

// Helper to get color with case-insensitive lookup
const getGenderColor = (label: string): string | undefined => {
  return GENDER_COLORS[label.toLowerCase()];
};

// Fallback colors for unknown labels
const FALLBACK_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
  "hsl(var(--chart-5))",
];

// Generate chart config dynamically from data
const generateChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "MSME Owners",
    },
  };

  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_");
    config[key] = {
      label: item.label,
      color: getGenderColor(item.label) || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    };
  });

  return config;
};

export function SmeGenderChart({ data, isLoading = false }: SmeGenderChartProps) {
  // Transform data for Recharts - assign colors based on gender
  const chartData = React.useMemo(() => {
    return data.map((item, index) => ({
      gender: item.label.toLowerCase().replace(/\s+/g, "_"),
      label: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: getGenderColor(item.label) || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  const totalOwners = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.value, 0);
  }, [data]);

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Owner Gender Distribution</CardTitle>
          <CardDescription>SME owners by gender</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center pb-0">
          <div className="h-[250px] w-full animate-pulse rounded-md bg-muted" />
        </CardContent>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Owner Gender Distribution</CardTitle>
          <CardDescription>SME owners by gender</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center pb-0">
          <p className="text-muted-foreground">No data available</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Owner Gender Distribution</CardTitle>
        <CardDescription>Primary SME owners by gender</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[300px]"
        >
          <PieChart>
            <ChartTooltip
              cursor={false}
              content={({ active, payload }) => {
                if (!active || !payload || !payload.length) return null;
                const data = payload[0]?.payload;
                if (!data) return null;
                const pct = typeof data.percentage === 'number'
                  ? data.percentage.toFixed(1)
                  : String(data.percentage ?? 0);
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <div className="flex min-w-[130px] items-center text-xs text-muted-foreground">
                      <div
                        className="h-2.5 w-2.5 shrink-0 rounded-[2px] mr-2"
                        style={{ backgroundColor: data.fill }}
                      />
                      <span className="flex-1">{data.label}</span>
                      <div className="ml-auto flex items-baseline gap-1 font-mono font-medium tabular-nums text-foreground">
                        {data.value}
                        <span className="font-normal text-muted-foreground">
                          ({pct}%)
                        </span>
                      </div>
                    </div>
                  </div>
                );
              }}
            />
            <Pie
              data={chartData}
              dataKey="value"
              nameKey="gender"
              innerRadius={60}
              outerRadius={100}
              strokeWidth={2}
              stroke="hsl(var(--background))"
            >
              {chartData.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={entry.fill}
                />
              ))}
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="fill-foreground text-3xl font-bold"
                        >
                          {totalOwners.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground text-sm"
                        >
                          Total Owners
                        </tspan>
                      </text>
                    );
                  }
                }}
              />
            </Pie>
            <ChartLegend
              content={<ChartLegendContent nameKey="gender" />}
              className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
            />
          </PieChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

export default SmeGenderChart;
