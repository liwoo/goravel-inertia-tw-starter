import * as React from 'react';
import { TrendingUp, TrendingDown, Building2, FileCheck, CalendarPlus } from 'lucide-react';
import { Area, AreaChart, ResponsiveContainer } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { SmeStats, RegistrationTrendPoint } from '@/types/sme';
import { cn } from '@/lib/utils';

interface SmeKpiCardsProps {
  stats: SmeStats;
}

// Mini sparkline component for the trend card
function TrendSparkline({ data }: { data: RegistrationTrendPoint[] }) {
  return (
    <div className="h-12 w-24">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 0, right: 0, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="sparklineGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="hsl(var(--chart-1))" stopOpacity={0.4} />
              <stop offset="95%" stopColor="hsl(var(--chart-1))" stopOpacity={0} />
            </linearGradient>
          </defs>
          <Area
            type="monotone"
            dataKey="count"
            stroke="hsl(var(--chart-1))"
            strokeWidth={2}
            fill="url(#sparklineGradient)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}

// Calculate percentage change between two values
function calculatePercentageChange(current: number, previous: number): number {
  if (previous === 0) return current > 0 ? 100 : 0;
  return Math.round(((current - previous) / previous) * 100);
}

// Trend indicator component
function TrendIndicator({ value, suffix = '%' }: { value: number; suffix?: string }) {
  const isPositive = value >= 0;
  const Icon = isPositive ? TrendingUp : TrendingDown;

  return (
    <div
      className={cn(
        'flex items-center gap-1 text-xs font-medium',
        isPositive ? 'text-emerald-600' : 'text-red-600'
      )}
    >
      <Icon className="h-3 w-3" />
      <span>
        {isPositive ? '+' : ''}
        {value}
        {suffix}
      </span>
    </div>
  );
}

export function SmeKpiCards({ stats }: SmeKpiCardsProps) {
  const monthOverMonthChange = calculatePercentageChange(
    stats.newThisMonth,
    stats.newLastMonth
  );

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {/* Total MSMEs Card */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Total MSMEs
          </CardTitle>
          <Building2 className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {stats.totalSmes.toLocaleString()}
          </div>
          <p className="text-xs text-muted-foreground mt-1 hidden lg:block">
            Registered in the database
          </p>
        </CardContent>
      </Card>

      {/* New This Month Card */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            New This Month
          </CardTitle>
          <CalendarPlus className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold">
              {stats.newThisMonth.toLocaleString()}
            </div>
            <TrendIndicator value={monthOverMonthChange} />
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            <span className="hidden lg:inline">vs {stats.newLastMonth.toLocaleString()} last month</span>
            <span className="lg:hidden">vs {stats.newLastMonth.toLocaleString()}</span>
          </p>
        </CardContent>
      </Card>

      {/* Registered Businesses Card */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Formally Registered
          </CardTitle>
          <FileCheck className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold">
              {stats.hasRegistration.toLocaleString()}
            </div>
            <div className="flex items-center gap-1 text-xs font-medium text-muted-foreground">
              <span>{typeof stats.hasRegistrationPercentage === 'number' ? stats.hasRegistrationPercentage.toFixed(1) : stats.hasRegistrationPercentage}%</span>
              <span>of total</span>
            </div>
          </div>
          <div className="mt-2 h-1.5 w-full rounded-full bg-muted">
            <div
              className="h-full rounded-full bg-chart-2 transition-all"
              style={{ width: `${stats.hasRegistrationPercentage}%` }}
            />
          </div>
        </CardContent>
      </Card>

      {/* Registration Trend Sparkline Card */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Registration Trend
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {stats.registrationTrend.length > 0
                  ? stats.registrationTrend[stats.registrationTrend.length - 1].cumulative.toLocaleString()
                  : 0}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                Cumulative total
              </p>
            </div>
            {stats.registrationTrend.length > 0 && (
              <TrendSparkline data={stats.registrationTrend} />
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
