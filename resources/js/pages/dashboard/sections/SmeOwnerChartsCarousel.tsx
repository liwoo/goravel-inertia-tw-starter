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
import { ChartActions } from "@/components/ui/chart-actions";
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
  height?: number;
}

function GenderChartContent({ data, height }: GenderChartContentProps) {
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
          nameKey="gender"
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

  // Create refs for each chart slide
  const chartRef1 = React.useRef<HTMLDivElement>(null);
  const chartRef2 = React.useRef<HTMLDivElement>(null);
  const chartRef3 = React.useRef<HTMLDivElement>(null);

  // Prepare CSV-friendly data for each chart
  const csvData1 = React.useMemo(() => {
    return genderData.map(item => ({
      "Gender": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [genderData]);

  const csvData2 = React.useMemo(() => {
    return youthData.map(item => ({
      "Category": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [youthData]);

  const csvData3 = React.useMemo(() => {
    return ageGenderData.map(item => ({
      "Age Group": item.ageGroup,
      "Male": item.male,
      "Female": item.female,
    }));
  }, [ageGenderData]);

  // Create arrays to switch between charts based on currentIndex
  const chartRefs = [chartRef1, chartRef2, chartRef3];
  const csvDataSets = [csvData1, csvData2, csvData3];
  const filenames = ["sme-owner-gender-distribution", "sme-owner-youth-distribution", "sme-owner-age-gender-pyramid"];

  // Render function for fullscreen charts
  const renderFullscreen = React.useCallback((width: number, height: number) => {
    const chartHeight = height - 20; // Leave some padding
    if (currentIndex === 0) {
      return (
        <div style={{ width, height: chartHeight }}>
          <GenderChartContent data={genderData} height={chartHeight} />
        </div>
      );
    } else if (currentIndex === 1) {
      return (
        <div style={{ width, height: chartHeight }}>
          <SmeYouthChart data={youthData} height={chartHeight} />
        </div>
      );
    } else {
      return (
        <div style={{ width, height: chartHeight }}>
          <SmeAgePyramidChart data={ageGenderData} height={chartHeight} />
        </div>
      );
    }
  }, [currentIndex, genderData, youthData, ageGenderData]);

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
      <CardHeader className="items-center pb-0 relative">
        <div className="absolute right-4 top-4">
          <ChartActions
            chartRef={chartRefs[currentIndex]}
            data={csvDataSets[currentIndex]}
            filename={filenames[currentIndex]}
            title={chartInfo[currentIndex].title}
            renderFullscreen={renderFullscreen}
          />
        </div>
        <CardTitle>{chartInfo[currentIndex].title}</CardTitle>
        <CardDescription>{chartInfo[currentIndex].description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-4 overflow-hidden">
        <CarouselWithCallback onIndexChange={handleCarouselChange}>
          {/* Chart 1: Gender Distribution */}
          <div ref={chartRef1} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <GenderChartContent data={genderData} />
          </div>

          {/* Chart 2: Youth Distribution */}
          <div ref={chartRef2} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <SmeYouthChart data={youthData} />
          </div>

          {/* Chart 3: Age/Gender Pyramid */}
          <div ref={chartRef3} className="pt-2 pb-10 h-[340px] overflow-hidden">
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
