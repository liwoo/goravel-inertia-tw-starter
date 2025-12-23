"use client";

import * as React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ChartActions } from "@/components/ui/chart-actions";
import { DistributionPoint } from "@/types/sme";
import { BdspStatusChart } from "./BdspStatusChart";
import { BdspServicesChart } from "./BdspServicesChart";

// Props for the carousel component
export interface BdspChartsCarouselProps {
  statusData: DistributionPoint[];
  servicesData: DistributionPoint[];
  isLoading?: boolean;
}

// ============================================================================
// Main Carousel Component
// ============================================================================

export function BdspChartsCarousel({
  statusData,
  servicesData,
  isLoading = false,
}: BdspChartsCarouselProps) {
  // Chart titles and descriptions for each slide
  const chartInfo = [
    {
      title: "BDSP Registration Status",
      description: "Business Development Service Providers by status",
    },
    {
      title: "Top Services Offered",
      description: "Most common services provided by BDSPs",
    },
  ];

  const [currentIndex, setCurrentIndex] = React.useState(0);

  // Create refs for each chart slide
  const chartRef1 = React.useRef<HTMLDivElement>(null);
  const chartRef2 = React.useRef<HTMLDivElement>(null);

  // Prepare CSV-friendly data for each chart
  const csvData1 = React.useMemo(() => {
    return statusData.map(item => ({
      "Status": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [statusData]);

  const csvData2 = React.useMemo(() => {
    return servicesData.map(item => ({
      "Service": item.label,
      "Count": item.value,
      "Percentage (%)": item.percentage.toFixed(1),
    }));
  }, [servicesData]);

  // Create arrays to switch between charts based on currentIndex
  const chartRefs = [chartRef1, chartRef2];
  const csvDataSets = [csvData1, csvData2];
  const filenames = ["bdsp-status-distribution", "bdsp-services-distribution"];

  // Render function for fullscreen charts
  const renderFullscreen = React.useCallback((width: number, height: number) => {
    const chartHeight = height - 20; // Leave some padding
    if (currentIndex === 0) {
      return (
        <div style={{ width, height: chartHeight }}>
          <BdspStatusChart data={statusData} height={chartHeight} />
        </div>
      );
    } else {
      return (
        <div style={{ width, height: chartHeight }}>
          <BdspServicesChart data={servicesData} height={chartHeight} />
        </div>
      );
    }
  }, [currentIndex, statusData, servicesData]);

  // Track current index for card header update
  const handleCarouselChange = (index: number) => {
    setCurrentIndex(index);
  };

  if (isLoading) {
    return (
      <Card className="flex flex-col">
        <CardHeader className="items-center pb-0">
          <CardTitle>BDSP Overview</CardTitle>
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
          {/* Chart 1: Registration Status (Donut) */}
          <div ref={chartRef1} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <BdspStatusChart data={statusData} />
          </div>

          {/* Chart 2: Top Services (Bar) */}
          <div ref={chartRef2} className="pt-2 pb-10 h-[340px] overflow-hidden">
            <BdspServicesChart data={servicesData} />
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

export default BdspChartsCarousel;
