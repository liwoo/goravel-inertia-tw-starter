"use client";

import * as React from 'react';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { RegistrationTrendPoint } from '@/types/sme';

interface SmeRegistrationTrendProps {
  data: RegistrationTrendPoint[];
}

const chartConfig = {
  count: {
    label: 'New Registrations',
    color: 'hsl(var(--chart-1))',
  },
  cumulative: {
    label: 'Cumulative Total',
    color: 'hsl(var(--chart-2))',
  },
} satisfies ChartConfig;

type ViewMode = 'count' | 'cumulative';

export function SmeRegistrationTrend({ data }: SmeRegistrationTrendProps) {
  const [viewMode, setViewMode] = React.useState<ViewMode>('count');

  // Format period string for display - backend sends "Jan 2025" format
  const formatPeriod = (period: string): string => {
    // Period is already in "Jan 2025" format from backend
    return period || '';
  };

  // Format period for shorter axis labels (e.g., "Jan 2025" -> "Jan")
  const formatAxisLabel = (period: string): string => {
    // Extract just the month from "Jan 2025"
    const parts = period?.split(' ');
    return parts?.[0] || period || '';
  };

  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Registration Trend</CardTitle>
          <CardDescription>No trend data available</CardDescription>
        </CardHeader>
        <CardContent className="flex h-[250px] items-center justify-center">
          <p className="text-sm text-muted-foreground">
            Registration trend data will appear here once available.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="@container/card">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <div className="space-y-1">
          <CardTitle>Registration Trend</CardTitle>
          <CardDescription>
            {viewMode === 'count'
              ? 'Monthly new SME registrations'
              : 'Cumulative SME registrations over time'}
          </CardDescription>
        </div>
        <Select value={viewMode} onValueChange={(v) => setViewMode(v as ViewMode)}>
          <SelectTrigger className="w-[160px]" aria-label="Select view mode">
            <SelectValue />
          </SelectTrigger>
          <SelectContent className="rounded-xl">
            <SelectItem value="count" className="rounded-lg">
              New Registrations
            </SelectItem>
            <SelectItem value="cumulative" className="rounded-lg">
              Cumulative Total
            </SelectItem>
          </SelectContent>
        </Select>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <AreaChart
            data={data}
            margin={{ top: 10, right: 10, left: 0, bottom: 0 }}
          >
            <defs>
              <linearGradient id="fillCount" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-count)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-count)"
                  stopOpacity={0.1}
                />
              </linearGradient>
              <linearGradient id="fillCumulative" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-cumulative)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-cumulative)"
                  stopOpacity={0.1}
                />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} strokeDasharray="3 3" />
            <XAxis
              dataKey="period"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={formatAxisLabel}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => value.toLocaleString()}
            />
            <ChartTooltip
              cursor={false}
              content={({ active, payload, label }) => {
                if (!active || !payload || !payload.length) return null;
                const data = payload[0]?.payload;
                if (!data) return null;
                return (
                  <div className="rounded-lg border bg-background p-2 shadow-sm">
                    <div className="text-xs">
                      <div className="font-medium text-foreground mb-1">
                        {formatPeriod(data.period)}
                      </div>
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <div
                          className="h-2 w-2 rounded-full"
                          style={{ backgroundColor: viewMode === 'count' ? 'var(--color-count)' : 'var(--color-cumulative)' }}
                        />
                        <span>{viewMode === 'count' ? 'New' : 'Total'}:</span>
                        <span className="font-mono font-medium text-foreground">
                          {(viewMode === 'count' ? data.count : data.cumulative)?.toLocaleString() ?? 0}
                        </span>
                      </div>
                    </div>
                  </div>
                );
              }}
            />
            {viewMode === 'count' ? (
              <Area
                dataKey="count"
                type="monotone"
                fill="url(#fillCount)"
                stroke="var(--color-count)"
                strokeWidth={2}
              />
            ) : (
              <Area
                dataKey="cumulative"
                type="monotone"
                fill="url(#fillCumulative)"
                stroke="var(--color-cumulative)"
                strokeWidth={2}
              />
            )}
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
