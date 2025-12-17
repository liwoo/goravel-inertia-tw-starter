import * as React from 'react';
import { SmeStats } from '@/types/sme';
import { SmeKpiCards } from './SmeKpiCards';
import { SmeRegistrationTrend } from './SmeRegistrationTrend';

interface SmeChartsContainerProps {
  stats: SmeStats;
}

export function SmeChartsContainer({ stats }: SmeChartsContainerProps) {
  return (
    <div className="flex flex-col gap-6 px-4 lg:px-6">
      {/* KPI Cards Row */}
      <SmeKpiCards stats={stats} />

      {/* Registration Trend Chart */}
      <SmeRegistrationTrend data={stats.registrationTrend} />
    </div>
  );
}
