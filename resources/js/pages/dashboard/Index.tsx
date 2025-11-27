import React from "react";
import { Head, usePage } from "@inertiajs/react";
import Admin from "@/layouts/Admin";
import type { SharedData } from "@/types/app";
import { DistributionPoint } from "@/types/sme";
import { usePermissions } from "@/contexts/PermissionsContext";
import { PermissionGate } from "@/components/Permissions/PermissionGate";
import {
  DashboardKpiCards,
  SmeSectorChart,
  SmeGenderChart,
} from "./sections";
import {
  SmeRegionChart,
  SmeCategoryChart,
} from "@/pages/Sme/sections";
import {
  UpcomingEventsWidget,
  UpcomingProcurementsWidget,
  RecentActivitiesWidget,
  UpcomingEvent,
  UpcomingProcurement,
  RecentActivity,
} from "@/components/widgets";

// Dashboard statistics interface
interface DashboardStats {
  totalSmes: number;
  newThisMonth: number;
  newLastMonth: number;
  totalEvents: number;
  upcomingEvents: number;
  totalProcurements: number;
  activeProcurements: number;
  // Distribution data
  byRegion: DistributionPoint[];
  byCategory: DistributionPoint[];
  bySector: DistributionPoint[];
  byGender: DistributionPoint[];
}

// Dashboard page props
interface DashboardPageProps extends SharedData {
  stats?: DashboardStats;
  upcomingEvents?: UpcomingEvent[];
  upcomingProcurements?: UpcomingProcurement[];
  recentActivities?: RecentActivity[];
  [key: string]: unknown; // Index signature for PageProps compatibility
}

// Default/mock data for when backend hasn't implemented stats yet
const defaultStats: DashboardStats = {
  totalSmes: 0,
  newThisMonth: 0,
  newLastMonth: 0,
  totalEvents: 0,
  upcomingEvents: 0,
  totalProcurements: 0,
  activeProcurements: 0,
  byRegion: [],
  byCategory: [],
  bySector: [],
  byGender: [],
};

const DashboardPage: React.FC = () => {
  const { props } = usePage<DashboardPageProps>();
  const user = props.auth?.user;
  const { canPerformAction } = usePermissions();

  // Use provided stats or default empty stats
  const stats = props.stats || defaultStats;
  const upcomingEvents = props.upcomingEvents || [];
  const upcomingProcurements = props.upcomingProcurements || [];
  const recentActivities = props.recentActivities || [];

  // Permission checks for different data sections
  const canViewSmes = canPerformAction("smes", "read");
  const canViewEvents = canPerformAction("events", "read");
  const canViewProcurements = canPerformAction("procurement_notices", "read");

  // Check if user has any SME-related permissions to show charts
  const canViewSmeCharts = canViewSmes;

  return (
    <Admin>
      <Head>
        <title>{props.pageTitle || "Dashboard"}</title>
      </Head>

      <div className="p-4 md:p-6 flex flex-col gap-6">
        {/* Welcome Header - Full width above the main content */}
        <div className="flex flex-col gap-1">
          {user ? (
            <>
              <h1 className="text-2xl font-semibold tracking-tight">
                Welcome back, {user.name}!
              </h1>
              <p className="text-muted-foreground">
                Here is an overview of your SME database and upcoming activities.
              </p>
            </>
          ) : (
            <>
              <h1 className="text-2xl font-semibold tracking-tight">
                Dashboard
              </h1>
              <p className="text-muted-foreground">
                Overview of your SME database and upcoming activities.
              </p>
            </>
          )}
        </div>

        {/* Main Content + Sidebar Row */}
        <div className="flex flex-col xl:flex-row gap-6 items-start">
          {/* Main Content Area - Left Side */}
          <div className="flex-1 flex flex-col gap-6 min-w-0">
            {/* KPI Cards Row */}
            <DashboardKpiCards
              stats={stats}
              canViewSmes={canViewSmes}
              canViewEvents={canViewEvents}
              canViewProcurements={canViewProcurements}
            />

            {/* Charts Section - Only show if user can view SMEs */}
            {canViewSmeCharts && (
              <>
                {/* First Row of Charts: Region (Donut) + Sector (Bar) */}
                <div className="grid gap-4 lg:grid-cols-2">
                  <SmeRegionChart data={stats.byRegion} />
                  <SmeSectorChart data={stats.bySector} />
                </div>

                {/* Second Row of Charts: Category (Bar) + Gender (Donut) */}
                <div className="grid gap-4 lg:grid-cols-2">
                  <SmeCategoryChart data={stats.byCategory} />
                  <SmeGenderChart data={stats.byGender} />
                </div>
              </>
            )}
          </div>

          {/* Right Sidebar - Widgets (top-aligned with KPI cards) */}
          <aside className="w-full xl:w-[380px] 2xl:w-[420px] flex flex-col gap-4 shrink-0">
            <RecentActivitiesWidget activities={recentActivities} />
            <PermissionGate service="events" action="read">
              <UpcomingEventsWidget events={upcomingEvents} />
            </PermissionGate>
            <PermissionGate service="procurement_notices" action="read">
              <UpcomingProcurementsWidget procurements={upcomingProcurements} />
            </PermissionGate>
          </aside>
        </div>
      </div>
    </Admin>
  );
};

export default DashboardPage;
