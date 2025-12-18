import React from "react";
import { Head, usePage } from "@inertiajs/react";
import Admin from "@/layouts/Admin";
import type { SharedData } from "@/types/app";
import { usePermissions } from "@/contexts/PermissionsContext";
import {
  RecentActivitiesWidget,
  RecentActivity,
} from "@/components/widgets";

// Dashboard statistics interface
interface DashboardStats {
  total: number;
  newThisMonth: number;
  newLastMonth: number;
}

// Dashboard page props
interface DashboardPageProps extends SharedData {
  stats?: DashboardStats;
  recentActivities?: RecentActivity[];
  [key: string]: unknown; // Index signature for PageProps compatibility
}

// Default/mock data for when backend hasn't implemented stats yet
const defaultStats: DashboardStats = {
  total: 0,
  newThisMonth: 0,
  newLastMonth: 0,
};

const DashboardPage: React.FC = () => {
  const { props } = usePage<DashboardPageProps>();
  const user = props.auth?.user;
  const { canPerformAction, isSuperAdmin } = usePermissions();

  // Use provided stats or default empty stats
  const stats = props.stats || defaultStats;
  const recentActivities = props.recentActivities || [];

  // Permission checks for different data sections
  const canViewBooks = canPerformAction("books", "read");

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
                Here is an overview of your database and activities.
              </p>
            </>
          ) : (
            <>
              <h1 className="text-2xl font-semibold tracking-tight">
                Dashboard
              </h1>
              <p className="text-muted-foreground">
                Overview of your database and activities.
              </p>
            </>
          )}
        </div>

        {/* Main Content + Sidebar Row */}
        <div className="flex flex-col xl:flex-row gap-6 items-start">
          {/* Main Content Area - Left Side */}
          <div className="flex-1 flex flex-col gap-6 min-w-0">

          </div>

          {/* Right Sidebar - Widgets (top-aligned with KPI cards) */}
          <aside className="w-full xl:w-[380px] 2xl:w-[420px] flex flex-col gap-4 shrink-0">
            {isSuperAdmin() && (
              <RecentActivitiesWidget activities={recentActivities} />
            )}
          </aside>
        </div>
      </div>
    </Admin>
  );
};

export default DashboardPage;
