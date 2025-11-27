"use client";

import * as React from "react";
import { SmeRegionChart } from "./SmeRegionChart";
import { SmeCategoryChart } from "./SmeCategoryChart";
import { DistributionPoint } from "@/types/sme";

interface SmeDistributionChartsProps {
  byRegion: DistributionPoint[];
  byCategory: DistributionPoint[];
  isLoading?: boolean;
}

export function SmeDistributionCharts({
  byRegion,
  byCategory,
  isLoading = false,
}: SmeDistributionChartsProps) {
  return (
    <div className="space-y-4">
      {/* Section Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold tracking-tight">SME Distribution</h2>
          <p className="text-sm text-muted-foreground">
            Geographic and category breakdown of registered SMEs
          </p>
        </div>
      </div>

      {/* Charts Grid - Side by side on desktop, stacked on mobile */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Region Distribution (Donut Chart) */}
        <SmeRegionChart data={byRegion} isLoading={isLoading} />

        {/* Category Distribution (Bar Chart) */}
        <SmeCategoryChart data={byCategory} isLoading={isLoading} />
      </div>
    </div>
  );
}

export default SmeDistributionCharts;
