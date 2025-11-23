import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Bdsp,
  BdspListResponse,
  BdspListRequest,
  BdspStats
} from '@/types/bdsp';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  BdspCreateForm,
  BdspEditForm,
  BdspDetailView,
  bdspColumns,
  bdspColumnsMobile,
  bdspFilters,
  bdspStatsConfigs,
  getBdspPageActions
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { renderStatsCards, createPageActions } from '@/lib/crud-page-utils';
import { BdspExportDialog, BdspExportOptions } from '@/components/Bdsp/BdspActions';

// Props interface for the Bdsp Index page
interface BdspIndexProps {
  data: BdspListResponse;
  filters: BdspListRequest;
  stats?: BdspStats;
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

export default function BdspIndex({
  data,
  filters,
  stats,
  permissions,
  meta
}: BdspIndexProps) {
  const isMobile = useIsMobile();

  const [showExportDialog, setShowExportDialog] = useState(false);

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  // Handle export
  const handleExport = async (options: BdspExportOptions) => {
    const params = new URLSearchParams({
      format: options.format,
      ...(options.fields && { fields: options.fields.join(',') }),
      ...(options.includeStats && { includeStats: 'true' }),
      ...Object.fromEntries(
        Object.entries(filters.filters || {}).map(([key, value]) => [key, String(value)])
      ),
    });

    window.open(`/api/bdsps/export?${params.toString()}`);
  };

  const pageActions = createPageActions(
    getBdspPageActions(permissions, {
      onExport: () => setShowExportDialog(true),
    })
  );

  return (
    <Admin title={"Bdsp"}>
      <Head title="Bdsp - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, bdspStatsConfigs)}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Bdsp>
            data={data}
            filters={filters}
            title="BDSPs Management"
            resourceName="bdsps"
            columns={isMobile ? bdspColumnsMobile : bdspColumns}
            customFilters={bdspFilters}
            pageActions={pageActions}
            paginationConfig={meta?.pagination}
            createForm={BdspCreateForm}
            editForm={BdspEditForm}
            detailView={BdspDetailView}
            onRefresh={handleRefresh}
            canCreate={permissions.canCreate}
            canEdit={permissions.canEdit}
            canDelete={permissions.canDelete}
            canView={true}
          />
        </div>

        {/* Action Dialogs */}
        {showExportDialog && (
          <BdspExportDialog
            onClose={() => setShowExportDialog(false)}
            onExport={handleExport}
            totalItems={data.total}
          />
        )}
      </div>
    </Admin>
  );
}
