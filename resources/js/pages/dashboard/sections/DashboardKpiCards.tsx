"use client";

import * as React from "react";
import {
  TrendingUp,
  TrendingDown,
  Building2,
  Calendar,
  FileText,
  Users,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface DashboardStats {
  totalSmes: number;
  newThisMonth: number;
  newLastMonth: number;
  totalEvents: number;
  upcomingEvents: number;
  totalProcurements: number;
  activeProcurements: number;
}

interface DashboardKpiCardsProps {
  stats: DashboardStats;
  isLoading?: boolean;
  canViewSmes?: boolean;
  canViewEvents?: boolean;
  canViewProcurements?: boolean;
}

// Calculate percentage change between two values
function calculatePercentageChange(current: number, previous: number): number {
  if (previous === 0) return current > 0 ? 100 : 0;
  return Math.round(((current - previous) / previous) * 100);
}

// Trend indicator component
function TrendIndicator({
  value,
  suffix = "%",
}: {
  value: number;
  suffix?: string;
}) {
  const isPositive = value >= 0;
  const Icon = isPositive ? TrendingUp : TrendingDown;

  return (
    <div
      className={cn(
        "flex items-center gap-1 text-xs font-medium",
        isPositive ? "text-emerald-600" : "text-red-600"
      )}
    >
      <Icon className="h-3 w-3" />
      <span>
        {isPositive ? "+" : ""}
        {value}
        {suffix}
      </span>
    </div>
  );
}

// Loading skeleton for KPI card
function KpiCardSkeleton() {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="h-4 w-24 animate-pulse rounded bg-muted" />
        <div className="h-4 w-4 animate-pulse rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="h-8 w-16 animate-pulse rounded bg-muted" />
        <div className="mt-2 h-3 w-32 animate-pulse rounded bg-muted" />
      </CardContent>
    </Card>
  );
}

export function DashboardKpiCards({
  stats,
  isLoading = false,
  canViewSmes = true,
  canViewEvents = true,
  canViewProcurements = true,
}: DashboardKpiCardsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <KpiCardSkeleton />
        <KpiCardSkeleton />
        <KpiCardSkeleton />
        <KpiCardSkeleton />
      </div>
    );
  }

  const monthOverMonthChange = calculatePercentageChange(
    stats.newThisMonth,
    stats.newLastMonth
  );

  // Collect visible cards
  const cards: React.ReactNode[] = [];

  // Total MSMEs Card - only if user can view SMEs
  if (canViewSmes) {
    cards.push(
      <Card key="total-smes">
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
          <p className="text-xs text-muted-foreground mt-1">
            Registered in the database
          </p>
        </CardContent>
      </Card>
    );

    // New This Month Card - only if user can view SMEs
    cards.push(
      <Card key="new-this-month">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            New This Month
          </CardTitle>
          <Users className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold">
              {stats.newThisMonth.toLocaleString()}
            </div>
            <TrendIndicator value={monthOverMonthChange} />
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            vs {stats.newLastMonth.toLocaleString()} last month
          </p>
        </CardContent>
      </Card>
    );
  }

  // Events Card - only if user can view events
  if (canViewEvents) {
    cards.push(
      <Card key="events">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Events
          </CardTitle>
          <Calendar className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold">
              {stats.totalEvents.toLocaleString()}
            </div>
            {stats.upcomingEvents > 0 && (
              <span className="text-xs font-medium text-emerald-600">
                {stats.upcomingEvents} upcoming
              </span>
            )}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Total events organized
          </p>
        </CardContent>
      </Card>
    );
  }

  // Procurement Notices Card - only if user can view procurements
  if (canViewProcurements) {
    cards.push(
      <Card key="procurements">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Procurements
          </CardTitle>
          <FileText className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="text-2xl font-bold">
              {stats.totalProcurements.toLocaleString()}
            </div>
            {stats.activeProcurements > 0 && (
              <span className="text-xs font-medium text-emerald-600">
                {stats.activeProcurements} active
              </span>
            )}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Total procurement notices
          </p>
        </CardContent>
      </Card>
    );
  }

  // If no cards are visible, show a message
  if (cards.length === 0) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          No data available. Contact your administrator for access.
        </CardContent>
      </Card>
    );
  }

  // Determine grid columns based on number of visible cards
  const gridCols = cards.length === 1
    ? "grid-cols-1"
    : cards.length === 2
    ? "grid-cols-1 sm:grid-cols-2"
    : cards.length === 3
    ? "grid-cols-1 sm:grid-cols-2 xl:grid-cols-3"
    : "grid-cols-1 sm:grid-cols-2 xl:grid-cols-4";

  return (
    <div className={`grid gap-4 ${gridCols}`}>
      {cards}
    </div>
  );
}

export default DashboardKpiCards;
