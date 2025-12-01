import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Application,
  ApplicationListResponse,
  ApplicationListRequest
} from '@/types/application';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  ApplicationDetailView,
  applicationColumns,
  applicationColumnsMobile,
  applicationFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Application Index page
interface ApplicationIndexProps {
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
}

export default function ApplicationIndex({
  data,
  filters,
  permissions,
  meta
}: ApplicationIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  return (
    <Admin title={"Application"}>
      <Head title="Application - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Application>
            data={data}
            filters={filters}
            title="Applications"
            resourceName="applications"
            columns={isMobile ? applicationColumnsMobile : applicationColumns}
            customFilters={applicationFilters}
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
