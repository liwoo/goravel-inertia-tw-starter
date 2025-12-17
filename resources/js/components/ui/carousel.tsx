"use client";

import * as React from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface CarouselProps {
  children: React.ReactNode[];
  autoRotate?: boolean;
  autoRotateInterval?: number;
  className?: string;
  showArrows?: boolean;
  showDots?: boolean;
}

interface CarouselContextType {
  currentIndex: number;
  setCurrentIndex: (index: number) => void;
  itemCount: number;
  goToNext: () => void;
  goToPrevious: () => void;
}

const CarouselContext = React.createContext<CarouselContextType | null>(null);

export function useCarousel() {
  const context = React.useContext(CarouselContext);
  if (!context) {
    throw new Error("useCarousel must be used within a Carousel");
  }
  return context;
}

export function Carousel({
  children,
  autoRotate = false,
  autoRotateInterval = 7000,
  className,
  showArrows = true,
  showDots = true,
}: CarouselProps) {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [isPaused, setIsPaused] = React.useState(false);
  const itemCount = React.Children.count(children);

  const goToNext = React.useCallback(() => {
    setCurrentIndex((prev) => (prev + 1) % itemCount);
  }, [itemCount]);

  const goToPrevious = React.useCallback(() => {
    setCurrentIndex((prev) => (prev - 1 + itemCount) % itemCount);
  }, [itemCount]);

  // Auto-rotation effect
  React.useEffect(() => {
    if (!autoRotate || isPaused || itemCount <= 1) return;

    const interval = setInterval(goToNext, autoRotateInterval);
    return () => clearInterval(interval);
  }, [autoRotate, autoRotateInterval, isPaused, goToNext, itemCount]);

  const contextValue = React.useMemo(
    () => ({
      currentIndex,
      setCurrentIndex,
      itemCount,
      goToNext,
      goToPrevious,
    }),
    [currentIndex, itemCount, goToNext, goToPrevious]
  );

  return (
    <CarouselContext.Provider value={contextValue}>
      <div
        className={cn("relative", className)}
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
              <div key={index} className="w-full flex-shrink-0">
                {child}
              </div>
            ))}
          </div>
        </div>

        {/* Navigation Arrows */}
        {showArrows && itemCount > 1 && (
          <>
            <Button
              variant="ghost"
              size="icon"
              className="absolute left-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/80 shadow-sm hover:bg-background"
              onClick={goToPrevious}
              aria-label="Previous chart"
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-background/80 shadow-sm hover:bg-background"
              onClick={goToNext}
              aria-label="Next chart"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </>
        )}

        {/* Dots Navigation */}
        {showDots && itemCount > 1 && (
          <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
            {Array.from({ length: itemCount }).map((_, index) => (
              <button
                key={index}
                className={cn(
                  "h-2 w-2 rounded-full transition-all duration-200",
                  index === currentIndex
                    ? "bg-primary w-4"
                    : "bg-muted-foreground/30 hover:bg-muted-foreground/50"
                )}
                onClick={() => setCurrentIndex(index)}
                aria-label={`Go to chart ${index + 1}`}
              />
            ))}
          </div>
        )}
      </div>
    </CarouselContext.Provider>
  );
}

export function CarouselItem({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return <div className={cn("w-full", className)}>{children}</div>;
}
