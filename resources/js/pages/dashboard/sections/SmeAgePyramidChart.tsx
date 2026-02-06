"use client";

import * as React from "react";
import {
  Bar,
  BarChart,
  XAxis,
  YAxis,
  Cell,
  ReferenceLine,
  ResponsiveContainer,
} from "recharts";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
} from "@/components/ui/chart";

// Age group data point for pyramid chart
export interface AgeGenderDistributionPoint {
  ageGroup: string;
  male: number;
  female: number;
  malePercentage: number;
  femalePercentage: number;
}

interface SmeAgePyramidChartProps {
  data: AgeGenderDistributionPoint[];
  isLoading?: boolean;
  height?: number;
}

// Colors for male and female - using primary/contrast for binary distinction
const MALE_COLOR = "hsl(var(--primary))";
const FEMALE_COLOR = "var(--chart-contrast)";

const chartConfig: ChartConfig = {
  male: {
    label: "Male",
    color: MALE_COLOR,
  },
  female: {
    label: "Female",
    color: FEMALE_COLOR,
  },
};

export function SmeAgePyramidChart({ data, isLoading = false, height }: SmeAgePyramidChartProps) {
  // Transform data for pyramid chart
  // Males will be negative (extend left), females positive (extend right)
  const chartData = React.useMemo(() => {
    return data.map((item) => ({
      ageGroup: item.ageGroup,
      male: -item.male, // Negative for left side
      female: item.female, // Positive for right side
      maleActual: item.male,
      femaleActual: item.female,
      malePercentage: item.malePercentage,
      femalePercentage: item.femalePercentage,
    }));
  }, [data]);

  // Calculate max value for symmetric axis
  const maxValue = React.useMemo(() => {
    let max = 0;
    data.forEach((item) => {
      max = Math.max(max, item.male, item.female);
    });
    return Math.ceil(max * 1.1); // Add 10% padding
  }, [data]);

  // Calculate totals
  const totals = React.useMemo(() => {
    let totalMale = 0;
    let totalFemale = 0;
    data.forEach((item) => {
      totalMale += item.male;
      totalFemale += item.female;
    });
    return { male: totalMale, female: totalFemale, total: totalMale + totalFemale };
  }, [data]);

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
      {/* Legend */}
      <div className="flex items-center justify-center gap-6 mb-2">
        <div className="flex items-center gap-1.5">
          <div
            className="h-3 w-3 rounded-sm"
            style={{ backgroundColor: MALE_COLOR }}
          />
          <span className="text-xs text-muted-foreground">
            Male ({totals.male})
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <div
            className="h-3 w-3 rounded-sm"
            style={{ backgroundColor: FEMALE_COLOR }}
          />
          <span className="text-xs text-muted-foreground">
            Female ({totals.female})
          </span>
        </div>
      </div>

      <ChartContainer
        config={chartConfig}
        className={height ? "w-full" : "h-[260px] w-full min-h-[200px]"}
        style={height ? { height: `${height - 40}px` } : undefined}
      >
        <BarChart
          data={chartData}
          layout="vertical"
          margin={{ top: 10, right: 30, left: 30, bottom: 10 }}
          barCategoryGap="15%"
        >
          <XAxis
            type="number"
            domain={[-maxValue, maxValue]}
            tickFormatter={(value) => Math.abs(value).toString()}
            axisLine={false}
            tickLine={false}
            tick={{ fontSize: 10 }}
          />
          <YAxis
            type="category"
            dataKey="ageGroup"
            axisLine={false}
            tickLine={false}
            tick={{ fontSize: 11 }}
            width={50}
          />
          <ReferenceLine x={0} stroke="hsl(var(--border))" />
          <ChartTooltip
            cursor={{ fill: "hsl(var(--muted))", opacity: 0.3 }}
            wrapperStyle={{ zIndex: 1000 }}
            allowEscapeViewBox={{ x: false, y: true }}
            content={({ active, payload }) => {
              if (!active || !payload || !payload.length) return null;
              const data = payload[0]?.payload;
              if (!data) return null;
              return (
                <div className="rounded-lg border bg-background p-2 shadow-sm min-w-[160px] max-w-[200px]">
                  <div className="font-medium text-sm mb-1.5">{data.ageGroup}</div>
                  <div className="flex items-center text-xs text-muted-foreground mb-1">
                    <div
                      className="h-2.5 w-2.5 shrink-0 rounded-[2px] mr-2"
                      style={{ backgroundColor: MALE_COLOR }}
                    />
                    <span className="flex-1">Male</span>
                    <div className="ml-auto flex items-baseline gap-1 font-mono font-medium tabular-nums text-foreground">
                      {data.maleActual}
                      <span className="font-normal text-muted-foreground">
                        ({data.malePercentage?.toFixed(1) ?? 0}%)
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center text-xs text-muted-foreground">
                    <div
                      className="h-2.5 w-2.5 shrink-0 rounded-[2px] mr-2"
                      style={{ backgroundColor: FEMALE_COLOR }}
                    />
                    <span className="flex-1">Female</span>
                    <div className="ml-auto flex items-baseline gap-1 font-mono font-medium tabular-nums text-foreground">
                      {data.femaleActual}
                      <span className="font-normal text-muted-foreground">
                        ({data.femalePercentage?.toFixed(1) ?? 0}%)
                      </span>
                    </div>
                  </div>
                </div>
              );
            }}
          />
          <Bar dataKey="male" radius={[4, 0, 0, 4]}>
            {chartData.map((entry, index) => (
              <Cell key={`male-${index}`} fill={MALE_COLOR} />
            ))}
          </Bar>
          <Bar dataKey="female" radius={[0, 4, 4, 0]}>
            {chartData.map((entry, index) => (
              <Cell key={`female-${index}`} fill={FEMALE_COLOR} />
            ))}
          </Bar>
        </BarChart>
      </ChartContainer>
    </div>
  );
}

export default SmeAgePyramidChart;
