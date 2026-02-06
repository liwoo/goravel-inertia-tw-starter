"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, Cell, XAxis, YAxis } from "recharts";
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
} from "@/components/ui/chart";
import { ChartActions } from "@/components/ui/chart-actions";
import { DistributionPoint } from "@/types/sme";
import { SmeClassificationChart } from "./SmeClassificationChart";

// Props for the carousel component
export interface SmeBusinessChartsCarouselProps {
  sectorData: DistributionPoint[];
  categoryData: DistributionPoint[];
  classificationData: DistributionPoint[];
  isLoading?: boolean;
}

// ============================================================================
// Sector Chart Content (Bar Chart)
// ============================================================================

// Use purple shade variations for bar charts
const BAR_COLORS = [
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
];

const generateBarChartConfig = (data: DistributionPoint[], labelKey: string): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "MSMEs",
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

interface BarChartContentProps {
  data: DistributionPoint[];
  dataKey: string;
  labelWidth?: number;
  height?: number;
}

function BarChartContent({ data, dataKey, labelWidth = 90, height }: BarChartContentProps) {
  const chartData = React.useMemo(() => {
    return [...data]
      .sort((a, b) => b.value - a.value)
      .map((item, index) => ({
        [dataKey]: item.label,
        value: item.value,
        percentage: item.percentage,
        fill: BAR_COLORS[index % BAR_COLORS.length],
        shortLabel: item.label.length > 12
          ? item.label.substring(0, 9) + "..."
          : item.label,
      }));
  }, [data, dataKey]);

  const chartConfig = React.useMemo(() => generateBarChartConfig(data, dataKey), [data, dataKey]);

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center h-[300px]">
        <p className="text-muted-foreground">No data available</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden w-full">
      <ChartContainer
        config={chartConfig}
        className={height ? "w-full max-w-full" : "h-[300px] w-full max-w-full"}
        style={height ? { height: `${height}px` } : undefined}
      >
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
            width={labelWidth}
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
                        {data[dataKey]}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-muted-foreground gap-2">
                      <span>Count:</span>
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
          maxBarSize={32}
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

// ============================================================================
// Main Carousel Component
// ============================================================================

export function SmeBusinessChartsCarousel({
  sectorData,
  categoryData,
  classificationData,
  isLoading = false,
}: SmeBusinessChartsCarouselProps) {
  // Chart titles and descriptions for each slide
  const chartInfo = [
    {
      title: "Economic Sector",
      description: "MSMEs by economic sector",
    },
    {
      title: "Category Distribution",
      description: "MSMEs by business category",
    },
    {
      title: "Classification",
      description: "MSMEs by MSME classification",
    },
  ];

  const [currentIndex, setCurrentIndex] = React.useState(0);

  // Create refs for each chart slide
  const chartRef1 = React.useRef<HTMLDivElement>(null);
  const chartRef2 = React.useRef<HTMLDivElement>(null);
  const chartRef3 = React.useRef<HTMLDivElement>(null);

  // Prepare CSV-friendly data for each chart
  const csvData1 = React.useMemo(() => {
    return sectorData.map(item => ({
      "Sector": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [sectorData]);

  const csvData2 = React.useMemo(() => {
    return categoryData.map(item => ({
      "Category": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [categoryData]);

  const csvData3 = React.useMemo(() => {
    return classificationData.map(item => ({
      "Classification": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [classificationData]);

  // Create arrays to switch between charts based on currentIndex
  const chartRefs = [chartRef1, chartRef2, chartRef3];
  const csvDataSets = [csvData1, csvData2, csvData3];
  const filenames = ["sme-sector-distribution", "sme-category-distribution", "sme-classification-distribution"];

  // Render function for fullscreen charts
  const renderFullscreen = React.useCallback((width: number, height: number) => {
    const chartHeight = height - 20; // Leave some padding
    if (currentIndex === 0) {
      return (
        <div style={{ width, height: chartHeight }}>
          <BarChartContent data={sectorData} dataKey="sector" labelWidth={90} height={chartHeight} />
        </div>
      );
    } else if (currentIndex === 1) {
      return (
        <div style={{ width, height: chartHeight }}>
          <BarChartContent data={categoryData} dataKey="category" labelWidth={90} height={chartHeight} />
        </div>
      );
    } else {
      return (
        <div style={{ width, height: chartHeight }}>
          <SmeClassificationChart data={classificationData} height={chartHeight} />
        </div>
      );
    }
  }, [currentIndex, sectorData, categoryData, classificationData]);

  // Track current index for card header update
  const handleCarouselChange = (index: number) => {
    setCurrentIndex(index);
  };

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>Business Distribution</CardTitle>
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
          {/* Chart 1: Sector Distribution (Bar) */}
          <div ref={chartRef1} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <BarChartContent data={sectorData} dataKey="sector" labelWidth={90} />
          </div>

          {/* Chart 2: Category Distribution (Bar) */}
          <div ref={chartRef2} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <BarChartContent data={categoryData} dataKey="category" labelWidth={90} />
          </div>

          {/* Chart 3: Classification Distribution (Donut) */}
          <div ref={chartRef3} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <SmeClassificationChart data={classificationData} />
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

export default SmeBusinessChartsCarousel;
