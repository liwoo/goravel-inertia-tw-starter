import * as React from 'react';
import { SmeStats } from '@/types/sme';
import { SmeKpiCards } from './SmeKpiCards';
import { SmeRegistrationTrend } from './SmeRegistrationTrend';
import { SmeDistributionCharts } from './SmeDistributionCharts';

interface SmeChartsContainerProps {
  stats: SmeStats;
  isLoading?: boolean;
}

export function SmeChartsContainer({ stats, isLoading = false }: SmeChartsContainerProps) {
  return (
    <div className="flex flex-col gap-6 px-4 lg:px-6">
      {/* KPI Cards Row */}
      <SmeKpiCards stats={stats} />

      {/* Registration Trend Chart */}
      <SmeRegistrationTrend data={stats.registrationTrend} />

      {/* Distribution Charts (Region Donut + Category Bar) */}
      <SmeDistributionCharts
        byRegion={stats.byRegion}
        byCategory={stats.byCategory}
        isLoading={isLoading}
      />
    </div>
  );
}
