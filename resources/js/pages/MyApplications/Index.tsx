import React from 'react';
import { Head, router } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Badge } from '@/components/ui/badge';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  Application,
  ApplicationListResponse,
  ApplicationListRequest,
} from '@/types/application';
import {
  applicationColumns,
  applicationColumnsMobile,
  ApplicationDetailView,
} from '@/pages/Application/sections';
import { useIsMobile } from '@/hooks/use-mobile';

interface MyApplicationsIndexProps {
  data: ApplicationListResponse;
  filters: ApplicationListRequest;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
  };
  meta?: {
    pagination: {
      defaultPageSize: number;
      maxPageSize: number;
      allowedSizes: number[];
    };
  };
  userName?: string;
  smeName?: string;
  usmeNumber?: string;
  classification?: string;
  error?: string;
}

export default function MyApplicationsIndex({
  data,
  filters,
  permissions,
  meta,
  userName,
  smeName,
  usmeNumber,
  classification,
  error,
}: MyApplicationsIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Filter columns for my applications - remove some admin-specific columns
  const myApplicationColumns = applicationColumns.filter(col =>
    !['registrant_name', 'email', 'phone'].includes(col.key as string)
  );

  const myApplicationColumnsMobile = applicationColumnsMobile;

  return (
    <Admin title="My Applications">
      <Head title="My Applications" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Header with SME info */}
        <div className="px-4 lg:px-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">
              My Applications
            </h1>
            <p className="text-muted-foreground">
              <span className="hidden sm:inline">Track the status of your applications and formalisation requests</span>
              <span className="sm:hidden">Track your applications</span>
            </p>
          </div>
          {smeName && (
            <div className="flex flex-col items-end gap-2 shrink-0">
              <div className="text-lg font-semibold text-right">{smeName}</div>
              <div className="flex items-center gap-2">
                {usmeNumber && (
                  <Badge variant="outline" className="text-xs font-mono">
                    {usmeNumber}
                  </Badge>
                )}
                {classification && (
                  <Badge
                    variant="default"
                    className={`text-sm px-3 py-1 font-semibold ${
                      classification === 'Micro'
                        ? 'bg-blue-600 hover:bg-blue-700'
                        : classification === 'Small'
                        ? 'bg-emerald-600 hover:bg-emerald-700'
                        : classification === 'Medium'
                        ? 'bg-amber-600 hover:bg-amber-700'
                        : ''
                    }`}
                  >
                    {classification} Enterprise
                  </Badge>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Error message if any */}
        {error && (
          <div className="px-4 lg:px-6">
            <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-md">
              {error}
            </div>
          </div>
        )}

        {/* Main CRUD Component - Read Only */}
        <div className="px-0">
          <CrudPage<Application>
            data={data}
            filters={filters}
            title="Applications"
            resourceName="my-applications"
            columns={isMobile ? myApplicationColumnsMobile : myApplicationColumns}
            paginationConfig={meta?.pagination}
            detailView={ApplicationDetailView}
            onRefresh={handleRefresh}
            canView={true}
            readOnly={true}
          />
        </div>
      </div>
    </Admin>
  );
}
