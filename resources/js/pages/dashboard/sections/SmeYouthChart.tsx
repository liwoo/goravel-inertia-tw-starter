"use client";

import * as React from "react";
import { Cell, Label, Pie, PieChart } from "recharts";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
} from "@/components/ui/chart";
import { DistributionPoint } from "@/types/sme";

interface SmeYouthChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Define specific colors for youth categories
const YOUTH_COLORS: Record<string, string> = {
  Youth: "hsl(var(--chart-3))",
  "Non-Youth": "hsl(var(--chart-4))",
  Unknown: "hsl(var(--chart-5))",
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
      label: "SME Owners",
    },
  };

  data.forEach((item, index) => {
    const key = item.label.toLowerCase().replace(/\s+/g, "_").replace(/-/g, "_");
    config[key] = {
      label: item.label,
      color: YOUTH_COLORS[item.label] || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    };
  });

  return config;
};

export function SmeYouthChart({ data, isLoading = false }: SmeYouthChartProps) {
  // Transform data for Recharts - assign colors based on youth status
  const chartData = React.useMemo(() => {
    return data.map((item, index) => ({
      category: item.label.toLowerCase().replace(/\s+/g, "_").replace(/-/g, "_"),
      label: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: YOUTH_COLORS[item.label] || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  const totalOwners = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.value, 0);
  }, [data]);

  // Find youth percentage for center label
  const youthData = data.find((d) => d.label === "Youth");
  const youthPercentage = youthData?.percentage ?? 0;

  if (isLoading) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex flex-1 items-center justify-center pb-0">
          <div className="h-[250px] w-full animate-pulse rounded-md bg-muted" />
        </div>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex flex-1 items-center justify-center pb-0">
          <p className="text-muted-foreground">No data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-hidden">
      <ChartContainer
        config={chartConfig}
        className="mx-auto aspect-square h-full max-h-[300px] min-h-[250px]"
      >
        <PieChart>
          <ChartTooltip
            cursor={false}
            wrapperStyle={{ zIndex: 1000 }}
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
            nameKey="category"
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
                        {youthPercentage.toFixed(0)}%
                      </tspan>
                      <tspan
                        x={viewBox.cx}
                        y={(viewBox.cy || 0) + 24}
                        className="fill-muted-foreground text-sm"
                      >
                        Youth Owners
                      </tspan>
                    </text>
                  );
                }
              }}
            />
          </Pie>
          <ChartLegend
            content={<ChartLegendContent nameKey="category" />}
            className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
          />
        </PieChart>
      </ChartContainer>
    </div>
  );
}

export default SmeYouthChart;
