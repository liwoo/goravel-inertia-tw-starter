"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
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
  ChartTooltipContent,
} from "@/components/ui/chart";
import { DistributionPoint } from "@/types/sme";

interface SmeCategoryChartProps {
  data: DistributionPoint[];
  isLoading?: boolean;
}

// Chart configuration with primary theme color
const chartConfig = {
  value: {
    label: "SMEs",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

export function SmeCategoryChart({ data, isLoading = false }: SmeCategoryChartProps) {
  // Transform and sort data by value descending
  const chartData = React.useMemo(() => {
    return [...data]
      .sort((a, b) => b.value - a.value)
      .map((item) => ({
        category: item.label,
        value: item.value,
        percentage: item.percentage,
        // Truncate long labels for X-axis display
        shortLabel: item.label.length > 15
          ? item.label.substring(0, 12) + "..."
          : item.label,
      }));
  }, [data]);

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader>
          <CardTitle>Category Distribution</CardTitle>
          <CardDescription>SMEs by business category</CardDescription>
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
          <CardTitle>Category Distribution</CardTitle>
          <CardDescription>SMEs by business category</CardDescription>
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
        <CardTitle>Category Distribution</CardTitle>
        <CardDescription>SMEs by business category (sorted by count)</CardDescription>
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
              width={120}
              tick={{ fontSize: 12 }}
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
                      <div className="font-medium text-foreground">
                        {data.category}
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
              fill="var(--color-value)"
              radius={[0, 4, 4, 0]}
              maxBarSize={40}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

export default SmeCategoryChart;
