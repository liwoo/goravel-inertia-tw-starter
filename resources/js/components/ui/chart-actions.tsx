"use client";

import * as React from "react";
import html2canvas from "html2canvas";
import { Download, FileDown, Maximize2, MoreVertical } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

export interface ChartActionsProps {
  /**
   * Reference to the chart container element for PNG capture
   */
  chartRef: React.RefObject<HTMLElement | null>;

  /**
   * Data for CSV export - array of objects
   */
  data: Array<Record<string, any>>;

  /**
   * Filename prefix for downloads (e.g., "sme-district-chart")
   */
  filename: string;

  /**
   * Optional: chart title for display in fullscreen and exports
   */
  title?: string;

  /**
   * Optional: application name for export headers (defaults to "SMEDI Dashboard")
   */
  appName?: string;

  /**
   * Optional: custom class name for the trigger button
   */
  className?: string;

  /**
   * Render function for fullscreen content - receives container dimensions
   * This should render the chart component at the given size
   */
  renderFullscreen?: (width: number, height: number) => React.ReactNode;
}

/**
 * ChartActions - Reusable dropdown component for dashboard charts
 *
 * Provides three core actions:
 * 1. Download CSV - exports the chart's underlying data
 * 2. Download PNG - captures the chart as an image
 * 3. Fullscreen - expands the chart to fullscreen view
 */
export function ChartActions({
  chartRef,
  data,
  filename,
  title,
  appName = "SMEDI Dashboard",
  className,
  renderFullscreen,
}: ChartActionsProps) {
  const [isFullscreenOpen, setIsFullscreenOpen] = React.useState(false);
  const [isDownloading, setIsDownloading] = React.useState(false);
  // Initialize with viewport-based dimensions (95vw - padding, 90vh - header - padding)
  const [containerSize, setContainerSize] = React.useState(() => ({
    width: typeof window !== 'undefined' ? Math.floor(window.innerWidth * 0.95 - 48) : 800,
    height: typeof window !== 'undefined' ? Math.floor(window.innerHeight * 0.90 - 100) : 600,
  }));
  const fullscreenContainerRef = React.useRef<HTMLDivElement>(null);

  // Format current date-time for exports
  const getFormattedDateTime = React.useCallback(() => {
    const now = new Date();
    return now.toLocaleString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }, []);

  // Measure fullscreen container when dialog opens using ResizeObserver
  React.useEffect(() => {
    if (!isFullscreenOpen) return;

    // Immediately set viewport-based dimensions as a reliable fallback
    const viewportWidth = Math.floor(window.innerWidth * 0.95 - 48);
    const viewportHeight = Math.floor(window.innerHeight * 0.90 - 100);
    setContainerSize({ width: viewportWidth, height: viewportHeight });

    const container = fullscreenContainerRef.current;
    if (!container) return;

    const updateSize = () => {
      const rect = container.getBoundingClientRect();
      if (rect.width > 0 && rect.height > 0) {
        setContainerSize({ width: Math.floor(rect.width), height: Math.floor(rect.height) });
      }
    };

    // Use ResizeObserver for reliable measurement after initial render
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width, height } = entry.contentRect;
        if (width > 0 && height > 0) {
          setContainerSize({ width: Math.floor(width), height: Math.floor(height) });
        }
      }
    });
    resizeObserver.observe(container);

    // Also try to measure after dialog animation completes
    const timeoutId = setTimeout(updateSize, 100);

    return () => {
      resizeObserver.disconnect();
      clearTimeout(timeoutId);
    };
  }, [isFullscreenOpen]);

  /**
   * Converts an array of objects to CSV format
   */
  const convertToCSV = React.useCallback((data: Array<Record<string, any>>): string => {
    if (!data || data.length === 0) return "";

    // Extract headers from the first object
    const headers = Object.keys(data[0]);

    // Create CSV header row
    const csvHeader = headers.join(",");

    // Create CSV data rows
    const csvRows = data.map((row) => {
      return headers.map((header) => {
        const value = row[header];

        // Handle different value types
        if (value === null || value === undefined) {
          return "";
        }

        // Handle nested objects/arrays by stringifying
        if (typeof value === "object") {
          return `"${JSON.stringify(value).replace(/"/g, '""')}"`;
        }

        // Escape values containing commas or quotes
        const stringValue = String(value);
        if (stringValue.includes(",") || stringValue.includes('"') || stringValue.includes("\n")) {
          return `"${stringValue.replace(/"/g, '""')}"`;
        }

        return stringValue;
      }).join(",");
    });

    return [csvHeader, ...csvRows].join("\n");
  }, []);

  /**
   * Triggers a browser download with the given content
   */
  const triggerDownload = React.useCallback((content: Blob | string, filename: string) => {
    const blob = content instanceof Blob
      ? content
      : new Blob([content], { type: "text/csv;charset=utf-8;" });

    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }, []);

  /**
   * Handles CSV export with header context
   */
  const handleDownloadCSV = React.useCallback(() => {
    try {
      const csvData = convertToCSV(data);
      if (!csvData) {
        console.error("No data available to export");
        return;
      }

      // Create header with context
      const dateTime = getFormattedDateTime();
      const headerLines = [
        `"${appName}"`,
        `"${title || 'Chart Data'}"`,
        `"Generated: ${dateTime}"`,
        `""`, // Empty line separator
      ];

      const csvWithHeader = headerLines.join("\n") + "\n" + csvData;

      const timestamp = new Date().toISOString().split("T")[0];
      const csvFilename = `${filename}-${timestamp}.csv`;
      triggerDownload(csvWithHeader, csvFilename);
    } catch (error) {
      console.error("Error exporting CSV:", error);
    }
  }, [data, filename, title, appName, convertToCSV, triggerDownload, getFormattedDateTime]);

  /**
   * Handles PNG export using html2canvas
   * Creates a dark-themed image with app name, chart title, and date-time
   */
  const handleDownloadPNG = React.useCallback(async () => {
    if (!chartRef.current) {
      console.error("Chart reference is not available");
      return;
    }

    setIsDownloading(true);

    try {
      const dateTime = getFormattedDateTime();

      // Create a wrapper element with dark background
      const wrapper = document.createElement("div");
      wrapper.style.cssText = `
        background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
        padding: 32px;
        border-radius: 16px;
        display: inline-block;
        min-width: 600px;
      `;

      // Add header section with app name, chart title, and date
      const headerSection = document.createElement("div");
      headerSection.style.cssText = `
        margin-bottom: 24px;
        text-align: center;
      `;

      // App name (main title)
      const appNameElement = document.createElement("h1");
      appNameElement.textContent = appName;
      appNameElement.style.cssText = `
        color: #ffffff;
        font-size: 28px;
        font-weight: 700;
        margin: 0 0 8px 0;
        font-family: system-ui, -apple-system, sans-serif;
      `;
      headerSection.appendChild(appNameElement);

      // Chart title (subtitle)
      if (title) {
        const titleElement = document.createElement("h2");
        titleElement.textContent = title;
        titleElement.style.cssText = `
          color: rgba(255, 255, 255, 0.85);
          font-size: 18px;
          font-weight: 500;
          margin: 0 0 8px 0;
          font-family: system-ui, -apple-system, sans-serif;
        `;
        headerSection.appendChild(titleElement);
      }

      // Date-time
      const dateElement = document.createElement("p");
      dateElement.textContent = dateTime;
      dateElement.style.cssText = `
        color: rgba(255, 255, 255, 0.6);
        font-size: 14px;
        margin: 0;
        font-family: system-ui, -apple-system, sans-serif;
      `;
      headerSection.appendChild(dateElement);

      wrapper.appendChild(headerSection);

      // Clone the chart content
      const chartClone = chartRef.current.cloneNode(true) as HTMLElement;
      chartClone.style.cssText = `
        background: rgba(255, 255, 255, 0.95);
        border-radius: 12px;
        padding: 16px;
      `;
      wrapper.appendChild(chartClone);

      // Temporarily add to document for rendering
      wrapper.style.position = "absolute";
      wrapper.style.left = "-9999px";
      wrapper.style.top = "-9999px";
      document.body.appendChild(wrapper);

      // Capture the wrapper element as canvas
      const canvas = await html2canvas(wrapper, {
        backgroundColor: null,
        scale: 2, // Higher quality
        logging: false,
        useCORS: true,
      });

      // Clean up
      document.body.removeChild(wrapper);

      // Convert canvas to blob
      canvas.toBlob((blob) => {
        if (!blob) {
          console.error("Failed to create image blob");
          setIsDownloading(false);
          return;
        }

        const timestamp = new Date().toISOString().split("T")[0];
        const pngFilename = `${filename}-${timestamp}.png`;
        triggerDownload(blob, pngFilename);
        setIsDownloading(false);
      }, "image/png");
    } catch (error) {
      console.error("Error capturing chart as PNG:", error);
      setIsDownloading(false);
    }
  }, [chartRef, filename, title, appName, triggerDownload, getFormattedDateTime]);

  /**
   * Handles fullscreen toggle
   */
  const handleFullscreen = React.useCallback(() => {
    setIsFullscreenOpen(true);
  }, []);

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className={cn(
              "h-8 w-8 text-muted-foreground hover:text-foreground",
              className
            )}
            aria-label="Chart actions"
          >
            <MoreVertical className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-48">
          <DropdownMenuItem onClick={handleDownloadCSV}>
            <FileDown className="mr-2 h-4 w-4" />
            Download CSV
          </DropdownMenuItem>
          <DropdownMenuItem onClick={handleDownloadPNG} disabled={isDownloading}>
            <Download className="mr-2 h-4 w-4" />
            {isDownloading ? "Capturing..." : "Download PNG"}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={handleFullscreen}>
            <Maximize2 className="mr-2 h-4 w-4" />
            Fullscreen View
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Dialog open={isFullscreenOpen} onOpenChange={setIsFullscreenOpen}>
        <DialogContent className="!max-w-[95vw] w-[95vw] !max-h-[95vh] h-[90vh] p-6 flex flex-col sm:!max-w-[95vw]">
          <DialogHeader className="flex-shrink-0 pb-4">
            <DialogTitle className="text-2xl font-semibold">
              {title || "Chart View"}
            </DialogTitle>
            <DialogDescription className="sr-only">
              Fullscreen view of {title || "chart"}
            </DialogDescription>
          </DialogHeader>
          <div
            ref={fullscreenContainerRef}
            className="flex-1 min-h-0 w-full overflow-hidden"
            style={{ height: 'calc(90vh - 120px)' }}
          >
            {renderFullscreen ? (
              renderFullscreen(containerSize.width, containerSize.height)
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Fullscreen view not available for this chart
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default ChartActions;
