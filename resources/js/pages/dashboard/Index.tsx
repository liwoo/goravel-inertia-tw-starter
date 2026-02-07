import React from "react";
import { Head, usePage } from "@inertiajs/react";
import Admin from "@/layouts/Admin";
import type { SharedData } from "@/types/app";
import { usePermissions } from "@/contexts/PermissionsContext";
import { useTranslation } from 'react-i18next';
import {
  RecentActivitiesWidget,
  RecentActivity,
} from "@/components/widgets";

// Dashboard page props
interface DashboardPageProps extends SharedData {
  recentActivities?: RecentActivity[];
  [key: string]: unknown; // Index signature for PageProps compatibility
}

const DashboardPage: React.FC = () => {
  const { props } = usePage<DashboardPageProps>();
  const { t } = useTranslation('dashboard');
  const user = props.auth?.user;
  const { isSuperAdmin } = usePermissions();

  const recentActivities = props.recentActivities || [];

  return (
    <Admin>
      <Head>
        <title>{props.pageTitle || t('title')}</title>
      </Head>

      <div className="p-4 md:p-6 flex flex-col gap-6">
        {/* Welcome Header */}
        <div className="flex flex-col gap-1">
          {user ? (
            <>
              <h1 className="text-2xl font-semibold tracking-tight">
                {t('welcomeBack', { name: user.name })}
              </h1>
              <p className="text-muted-foreground">
                {t('systemOverview')}
              </p>
            </>
          ) : (
            <>
              <h1 className="text-2xl font-semibold tracking-tight">
                {t('title')}
              </h1>
              <p className="text-muted-foreground">
                {t('simpleOverview')}
              </p>
            </>
          )}
        </div>

        {/* Main Content */}
        <div className="flex flex-col xl:flex-row gap-6 items-start">
          <div className="flex-1 flex flex-col gap-6 min-w-0">
            {/* Placeholder for future dashboard content */}
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <p className="text-muted-foreground text-sm">
                {t('placeholder')}
              </p>
            </div>
          </div>

          {/* Right Sidebar - Widgets */}
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
