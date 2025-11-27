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

interface SmeRegionChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Define chart colors - using CSS variables for theme support
const CHART_COLORS = [
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
      label: "SMEs",
    },
  };

  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_");
    config[key] = {
      label: item.label,
      color: CHART_COLORS[index % CHART_COLORS.length],
    };
  });

  return config;
};

export function SmeRegionChart({ data, isLoading = false }: SmeRegionChartProps) {
  // Transform data for Recharts - assign colors based on index
  const chartData = React.useMemo(() => {
    return data.map((item, index) => ({
      region: item.label.toLowerCase().replace(/\s+/g, "_"),
      label: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: CHART_COLORS[index % CHART_COLORS.length],
    }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  const totalSmes = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.value, 0);
  }, [data]);

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Region Distribution</CardTitle>
          <CardDescription>SMEs by region</CardDescription>
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
          <CardTitle>Region Distribution</CardTitle>
          <CardDescription>SMEs by region</CardDescription>
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
        <CardTitle>Region Distribution</CardTitle>
        <CardDescription>SMEs by geographic region</CardDescription>
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
              nameKey="region"
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
                          {totalSmes.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground text-sm"
                        >
                          Total SMEs
                        </tspan>
                      </text>
                    );
                  }
                }}
              />
            </Pie>
            <ChartLegend
              content={<ChartLegendContent nameKey="region" />}
              className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
            />
          </PieChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

export default SmeRegionChart;
