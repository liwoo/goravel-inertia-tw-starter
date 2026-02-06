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
import { DistributionPoint } from "@/types/sme";
import { SmeDistrictChart } from "./SmeDistrictChart";

// Props for the carousel component
export interface SmeLocationChartsCarouselProps {
  regionData: DistributionPoint[];
  districtData: DistributionPoint[];
  isLoading?: boolean;
}

// ============================================================================
// Region Chart (Chart 1) - Donut chart for regions
// ============================================================================

// Use purple shade variations for charts
const CHART_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
  "hsl(var(--chart-5))",
];

const generateRegionChartConfig = (data: DistributionPoint[]): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "MSMEs",
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

interface RegionChartContentProps {
  data: DistributionPoint[];
  height?: number;
}

function RegionChartContent({ data, height }: RegionChartContentProps) {
  const chartData = React.useMemo(() => {
    return data.map((item, index) => ({
      region: item.label.toLowerCase().replace(/\s+/g, "_"),
      label: item.label,
      value: item.value,
      percentage: item.percentage,
      fill: CHART_COLORS[index % CHART_COLORS.length],
    }));
  }, [data]);

  const chartConfig = React.useMemo(() => generateRegionChartConfig(data), [data]);

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
          nameKey="region"
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
                      Total MSMEs
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
  );
}

// ============================================================================
// Main Carousel Component
// ============================================================================

export function SmeLocationChartsCarousel({
  regionData,
  districtData,
  isLoading = false,
}: SmeLocationChartsCarouselProps) {
  // Chart titles and descriptions for each slide
  const chartInfo = [
    {
      title: "Region Distribution",
      description: "MSMEs by geographic region",
    },
    {
      title: "District Distribution",
      description: "Top 15 districts by MSME count",
    },
  ];

  const [currentIndex, setCurrentIndex] = React.useState(0);

  // Create refs for each chart slide
  const chartRef1 = React.useRef<HTMLDivElement>(null);
  const chartRef2 = React.useRef<HTMLDivElement>(null);

  // Prepare CSV-friendly data for each chart
  const csvData1 = React.useMemo(() => {
    return regionData.map(item => ({
      "Region": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [regionData]);

  const csvData2 = React.useMemo(() => {
    return districtData.map(item => ({
      "District": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [districtData]);

  // Create arrays to switch between charts based on currentIndex
  const chartRefs = [chartRef1, chartRef2];
  const csvDataSets = [csvData1, csvData2];
  const filenames = ["sme-region-distribution", "sme-district-distribution"];

  // Render function for fullscreen charts
  const renderFullscreen = React.useCallback((width: number, height: number) => {
    const chartHeight = height - 20; // Leave some padding
    if (currentIndex === 0) {
      return (
        <div style={{ width, height: chartHeight }}>
          <RegionChartContent data={regionData} height={chartHeight} />
        </div>
      );
    } else {
      return (
        <div style={{ width, height: chartHeight }}>
          <SmeDistrictChart data={districtData} height={chartHeight} />
        </div>
      );
    }
  }, [currentIndex, regionData, districtData]);

  // Track current index for card header update
  const handleCarouselChange = (index: number) => {
    setCurrentIndex(index);
  };

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Location Distribution</CardTitle>
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
          {/* Chart 1: Region Distribution (Donut) */}
          <div ref={chartRef1} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <RegionChartContent data={regionData} />
          </div>

          {/* Chart 2: District Distribution (Bar) */}
          <div ref={chartRef2} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <SmeDistrictChart data={districtData} />
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

export default SmeLocationChartsCarousel;
