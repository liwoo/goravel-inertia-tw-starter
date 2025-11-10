import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Config,
  ConfigListResponse,
  ConfigListRequest
} from '@/types/config';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  ConfigCreateForm,
  ConfigEditForm,
  ConfigDetailView,
  configColumns,
  configColumnsMobile,
  configFilters,
  configSimpleFilters
} from './sections';
import { createSimpleFilters } from '@/lib/crud-page-utils';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Config Index page
interface ConfigIndexProps {
  data: ConfigListResponse;
  filters: ConfigListRequest;
  stats?: any; // Stats for config type counts
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

export default function ConfigIndex({
  data,
  filters,
  stats,
  permissions,
  meta
}: ConfigIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  // Create simple filters for config types
  const simpleFilters = createSimpleFilters(configSimpleFilters(stats));

  return (
    <Admin title={"Config"}>
      <Head title="Config - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Config>
            data={data}
            filters={filters}
            title="Config Management"
            resourceName="configs"
            columns={isMobile ? configColumnsMobile : configColumns}
            simpleFilters={simpleFilters}
            paginationConfig={meta?.pagination}
            createForm={ConfigCreateForm}
            editForm={ConfigEditForm}
            detailView={ConfigDetailView}
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
