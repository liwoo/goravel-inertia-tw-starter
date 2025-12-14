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
import { Carousel } from "@/components/ui/carousel";
import { DistributionPoint } from "@/types/sme";
import { SmeYouthChart } from "./SmeYouthChart";
import { SmeAgePyramidChart, AgeGenderDistributionPoint } from "./SmeAgePyramidChart";

// Props for the carousel component
export interface SmeOwnerChartsCarouselProps {
  genderData: DistributionPoint[];
  youthData: DistributionPoint[];
  ageGenderData: AgeGenderDistributionPoint[];
  isLoading?: boolean;
}

// ============================================================================
// Gender Chart (Chart 1) - Reused logic from existing SmeGenderChart
// ============================================================================

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

const FALLBACK_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
  "hsl(var(--chart-5))",
];

const generateGenderChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "SME Owners",
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

interface GenderChartContentProps {
  data: DistributionPoint[];
}

function GenderChartContent({ data }: GenderChartContentProps) {
  const chartData = React.useMemo(() => {
    return data.map((item, index) => ({
      gender: item.label.toLowerCase().replace(/\s+/g, "_"),
      label: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: getGenderColor(item.label) || FALLBACK_COLORS[index % FALLBACK_COLORS.length],
    }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateGenderChartConfig(data), [data]);

  const totalOwners = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.value, 0);
  }, [data]);

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center pb-0">
        <p className="text-muted-foreground">No data available</p>
      </div>
    );
  }

  return (
    <ChartContainer
      config={chartConfig}
      className="mx-auto aspect-square h-full max-h-[300px] min-h-[250px]"
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
  );
}

// ============================================================================
// Main Carousel Component
// ============================================================================

export function SmeOwnerChartsCarousel({
  genderData,
  youthData,
  ageGenderData,
  isLoading = false,
}: SmeOwnerChartsCarouselProps) {
  // Chart titles and descriptions for each slide
  const chartInfo = [
    {
      title: "Owner Gender Distribution",
      description: "Primary SME owners by gender",
    },
    {
      title: "Owner Youth Distribution",
      description: "SME owners by age group (18-35 = Youth)",
    },
    {
      title: "Age & Gender Distribution",
      description: "Population pyramid by age and gender",
    },
  ];

  const [currentIndex, setCurrentIndex] = React.useState(0);

  // Track current index for card header update
  const handleCarouselChange = (index: number) => {
    setCurrentIndex(index);
  };

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Owner Demographics</CardTitle>
          <CardDescription>Loading...</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 items-center justify-center pb-0">
          <div className="h-[300px] w-full animate-pulse rounded-md bg-muted" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="flex flex-col overflow-hidden">
      <CardHeader className="items-center pb-0">
        <CardTitle>{chartInfo[currentIndex].title}</CardTitle>
        <CardDescription>{chartInfo[currentIndex].description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-4 overflow-hidden">
        <CarouselWithCallback onIndexChange={handleCarouselChange}>
          {/* Chart 1: Gender Distribution */}
          <div className="pt-2 pb-10 h-[340px] overflow-hidden">
            <GenderChartContent data={genderData} />
          </div>

          {/* Chart 2: Youth Distribution */}
          <div className="pt-2 pb-10 h-[340px] overflow-hidden">
            <SmeYouthChart data={youthData} />
          </div>

          {/* Chart 3: Age/Gender Pyramid */}
          <div className="pt-2 pb-10 h-[340px] overflow-hidden">
            <SmeAgePyramidChart data={ageGenderData} />
          </div>
        </CarouselWithCallback>
      </CardContent>
    </Card>
  );
}

// Helper component to track carousel index changes
interface CarouselWithCallbackProps {
  children: React.ReactNode[];
  onIndexChange: (index: number) => void;
}

function CarouselWithCallback({ children, onIndexChange }: CarouselWithCallbackProps) {
  const [internalIndex, setInternalIndex] = React.useState(0);

  React.useEffect(() => {
    onIndexChange(internalIndex);
  }, [internalIndex, onIndexChange]);

  return (
    <CarouselTracked
      autoRotate={true}
      autoRotateInterval={7000}
      showArrows={true}
      showDots={true}
      currentIndex={internalIndex}
      onIndexChange={setInternalIndex}
    >
      {children}
    </CarouselTracked>
  );
}

// Extended Carousel component with index tracking
interface CarouselTrackedProps {
  children: React.ReactNode[];
  autoRotate?: boolean;
  autoRotateInterval?: number;
  showArrows?: boolean;
  showDots?: boolean;
  currentIndex: number;
  onIndexChange: (index: number) => void;
}

function CarouselTracked({
  children,
  autoRotate = false,
  autoRotateInterval = 7000,
  showArrows = true,
  showDots = true,
  currentIndex,
  onIndexChange,
}: CarouselTrackedProps) {
  const [isPaused, setIsPaused] = React.useState(false);
  const itemCount = React.Children.count(children);

  const goToNext = React.useCallback(() => {
    onIndexChange((currentIndex + 1) % itemCount);
  }, [currentIndex, itemCount, onIndexChange]);

  const goToPrevious = React.useCallback(() => {
    onIndexChange((currentIndex - 1 + itemCount) % itemCount);
  }, [currentIndex, itemCount, onIndexChange]);

  // Auto-rotation effect
  React.useEffect(() => {
    if (!autoRotate || isPaused || itemCount <= 1) return;

    const interval = setInterval(goToNext, autoRotateInterval);
    return () => clearInterval(interval);
  }, [autoRotate, autoRotateInterval, isPaused, goToNext, itemCount]);

  return (
    <div
      className="relative"
      onMouseEnter={() => setIsPaused(true)}
      onMouseLeave={() => setIsPaused(false)}
    >
      {/* Content */}
      <div className="overflow-hidden">
        <div
          className="flex transition-transform duration-300 ease-in-out"
          style={{ transform: `translateX(-${currentIndex * 100}%)` }}
        >
          {React.Children.map(children, (child, index) => (
            <div key={index} className="w-full flex-shrink-0 min-h-[350px]">
              {child}
            </div>
          ))}
        </div>
      </div>

      {/* Navigation Arrows */}
      {showArrows && itemCount > 1 && (
        <>
          <button
            className="absolute left-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/80 shadow-sm hover:bg-background flex items-center justify-center border border-border/50"
            onClick={goToPrevious}
            aria-label="Previous chart"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="m15 18-6-6 6-6" />
            </svg>
          </button>
          <button
            className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/80 shadow-sm hover:bg-background flex items-center justify-center border border-border/50"
            onClick={goToNext}
            aria-label="Next chart"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="m9 18 6-6-6-6" />
            </svg>
          </button>
        </>
      )}

      {/* Dots Navigation */}
      {showDots && itemCount > 1 && (
        <div className="absolute bottom-0 left-1/2 -translate-x-1/2 flex gap-1.5 py-2">
          {Array.from({ length: itemCount }).map((_, index) => (
            <button
              key={index}
              className={`h-2 rounded-full transition-all duration-200 ${
                index === currentIndex
                  ? "bg-primary w-4"
                  : "bg-muted-foreground/30 hover:bg-muted-foreground/50 w-2"
              }`}
              onClick={() => onIndexChange(index)}
              aria-label={`Go to chart ${index + 1}`}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default SmeOwnerChartsCarousel;
