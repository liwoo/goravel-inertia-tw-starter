import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Sme,
  SmeListResponse,
  SmeListRequest,
  SmeStats
} from '@/types/sme';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  SmeCreateForm,
  SmeEditForm,
  SmeDetailView,
  smeColumns,
  smeColumnsMobile,
  smeFilters,
  SmeChartsContainer
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Sme Index page
interface SmeIndexProps {
  data: SmeListResponse;
  filters: SmeListRequest;
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
  stats?: SmeStats;
}

export default function SmeIndex({
  data,
  filters,
  permissions,
  meta,
  stats
}: SmeIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  return (
    <Admin title={"SME Management"}>
      <Head title="Sme - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Charts and KPIs Section */}
        {stats && <SmeChartsContainer stats={stats} />}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Sme>
            data={data}
            filters={filters}
            title="SMEs Management"
            resourceName="smes"
            columns={isMobile ? smeColumnsMobile : smeColumns}
            customFilters={smeFilters}
            paginationConfig={meta?.pagination}
            createForm={SmeCreateForm}
            editForm={SmeEditForm}
            detailView={SmeDetailView}
            onRefresh={handleRefresh}
            canCreate={permissions.canCreate}
            canEdit={permissions.canEdit}
            canDelete={permissions.canDelete}
            canView={true}
          />
        </div>
      </div>
    </Admin>
  );
}
