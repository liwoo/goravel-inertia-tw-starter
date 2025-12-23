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

interface SmeClassificationChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
  height?: number;
}

// Colors matching the classification badges in SmeColumns.tsx
const CLASSIFICATION_COLORS: Record<string, string> = {
  micro: "hsl(217, 91%, 60%)",      // Blue for Micro
  small: "hsl(160, 84%, 39%)",       // Emerald for Small
  medium: "hsl(271, 91%, 65%)",      // Purple for Medium
  unclassified: "hsl(var(--muted-foreground))", // Gray for Unclassified
};

// Fallback colors for any unexpected values
const FALLBACK_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
];

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
      color: CLASSIFICATION_COLORS[key] || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    };
  });

  return config;
};

export function SmeClassificationChart({ data, isLoading = false, height }: SmeClassificationChartProps) {
  const chartData = React.useMemo(() => {
    // Define preferred order for classifications
    const order = ['micro', 'small', 'medium', 'unclassified'];

    return [...data]
      .sort((a, b) => {
        const aIndex = order.indexOf(a.label.toLowerCase());
        const bIndex = order.indexOf(b.label.toLowerCase());
        return aIndex - bIndex;
      })
      .map((item, index) => {
        const key = item.label.toLowerCase().replace(/\s+/g, "_");
        return {
          classification: key,
          label: item.label,
          value: item.value,
          percentage: item.percentage,
          fill: CLASSIFICATION_COLORS[key] || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
        };
      });
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  const totalSmes = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.value, 0);
  }, [data]);

  // Calculate dynamic radius based on height (scale proportionally)
  const { innerRadius, outerRadius } = React.useMemo(() => {
    if (height && height > 400) {
      // Scale up for fullscreen - use ~30% of height for outer radius
      const baseSize = Math.min(height * 0.35, 250);
      return {
        innerRadius: Math.floor(baseSize * 0.6),
        outerRadius: Math.floor(baseSize),
      };
    }
    // Default sizes for normal view
    return { innerRadius: 60, outerRadius: 100 };
  }, [height]);

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center pb-0">
        <div className="h-[250px] w-[250px] animate-pulse rounded-full bg-muted" />
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center pb-0">
        <p className="text-muted-foreground">No data available</p>
      </div>
    );
  }

  // Use dynamic height if provided, otherwise use default constraints
  const containerClass = height
    ? "mx-auto aspect-square"
    : "mx-auto aspect-square h-full max-h-[300px] min-h-[250px]";

  return (
    <ChartContainer
      config={chartConfig}
      className={containerClass}
      style={height ? { height: `${height}px` } : undefined}
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
          nameKey="classification"
          innerRadius={innerRadius}
          outerRadius={outerRadius}
          strokeWidth={2}
          stroke="hsl(var(--background))"
        >
          {chartData.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.fill} />
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
          content={<ChartLegendContent nameKey="classification" />}
          className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
        />
      </PieChart>
    </ChartContainer>
  );
}

export default SmeClassificationChart;
