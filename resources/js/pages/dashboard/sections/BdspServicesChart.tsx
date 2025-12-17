"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, Cell, XAxis, YAxis } from "recharts";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
} from "@/components/ui/chart";
import { DistributionPoint } from "@/types/sme";

interface BdspServicesChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Use purple shade variations for bar charts
const BAR_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
];

// Generate chart config dynamically from data
const generateChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "BDSPs",
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

export function BdspServicesChart({ data, isLoading = false }: BdspServicesChartProps) {
  // Transform and sort data by value descending, limit to top 10
  const chartData = React.useMemo(() => {
    return [...data]
      .sort((a, b) => b.value - a.value)
      .slice(0, 10)
      .map((item, index) => ({
        service: item.label,
        value: item.value,
        percentage: item.percentage,
        fill: BAR_COLORS[index % BAR_COLORS.length],
        // Truncate long labels for Y-axis display
        shortLabel: item.label.length > 10
          ? item.label.substring(0, 8) + "..."
          : item.label,
      }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateChartConfig(data), [data]);

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center h-[300px]">
        <div className="h-[250px] w-full animate-pulse rounded-md bg-muted" />
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center h-[300px]">
        <p className="text-muted-foreground">No data available</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden w-full">
      <ChartContainer config={chartConfig} className="h-[300px] w-full max-w-full">
        <BarChart
          data={chartData}
          layout="vertical"
          margin={{
            left: 0,
            right: 4,
            top: 0,
            bottom: 0,
          }}
        >
          <CartesianGrid horizontal={false} strokeDasharray="3 3" />
          <YAxis
            dataKey="shortLabel"
            type="category"
            tickLine={false}
            tickMargin={2}
            axisLine={false}
            width={70}
            tick={{ fontSize: 9 }}
          />
          <XAxis
            dataKey="value"
            type="number"
            tickLine={false}
            axisLine={false}
            tickMargin={8}
            tick={{ fontSize: 10 }}
          />
          <ChartTooltip
            cursor={{ fill: "hsl(var(--muted))", opacity: 0.3 }}
            wrapperStyle={{ zIndex: 1000 }}
            allowEscapeViewBox={{ x: false, y: true }}
            content={({ active, payload }) => {
              if (!active || !payload || !payload.length) return null;
              const data = payload[0]?.payload;
              if (!data) return null;
              const pct = typeof data.percentage === 'number'
                ? data.percentage.toFixed(1)
                : String(data.percentage ?? 0);
              return (
                <div className="rounded-lg border bg-background p-2 shadow-sm max-w-[200px]">
                  <div className="flex flex-col gap-1 text-xs">
                    <div className="flex items-center gap-2">
                      <div
                        className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                        style={{ backgroundColor: data.fill }}
                      />
                      <span className="font-medium text-foreground truncate">
                        {data.service}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-muted-foreground gap-2">
                      <span>Offering:</span>
                      <span className="font-mono font-medium tabular-nums text-foreground">
                        {data.value}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-muted-foreground gap-2">
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
            maxBarSize={28}
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
    </div>
  );
}

export default BdspServicesChart;
