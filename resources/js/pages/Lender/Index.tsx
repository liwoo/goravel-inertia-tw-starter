import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Lender,
  LenderListResponse,
  LenderListRequest
} from '@/types/lender';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  LenderCreateForm,
  LenderEditForm,
  LenderDetailView,
  lenderColumns,
  lenderColumnsMobile,
  lenderFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Lender Index page
interface LenderIndexProps {
  data: LenderListResponse;
  filters: LenderListRequest;
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

export default function LenderIndex({
  data,
  filters,
  permissions,
  meta
}: LenderIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  return (
    <Admin title={"Lender"}>
      <Head title="Lender - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Lender>
            data={data}
            filters={filters}
            title="lenders"
            resourceName="lenders"
            columns={isMobile ? lenderColumnsMobile : lenderColumns}
            paginationConfig={meta?.pagination}
            createForm={LenderCreateForm}
            editForm={LenderEditForm}
            detailView={LenderDetailView}
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
