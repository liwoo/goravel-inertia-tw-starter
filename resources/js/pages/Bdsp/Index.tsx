import React, { useState } from 'react';
import { Head, router, usePage } from '@inertiajs/react';
import { SharedData } from '@/types/app.d';
import { Download } from 'lucide-react';
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
  bdspStatsConfigs
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { renderStatsCards } from '@/lib/crud-page-utils';
import { ExportDialog } from '@/components/ExportDialog';
import { ExportField, ExportOptions, ExportColumn } from '@/types/export';
import { exportData } from '@/utils/exportUtils';
import { PageAction } from '@/types/crud';

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

// BDSP export field definitions
const bdspExportFields: ExportField[] = [
  { id: 'id', label: 'ID' },
  { id: 'ubdsp_number', label: 'UBDSP Number' },
  { id: 'name', label: 'Name' },
  { id: 'postal_address', label: 'Postal Address' },
  { id: 'physical_address', label: 'Physical Address' },
  { id: 'registration_status', label: 'Registration Status' },
  { id: 'partners', label: 'Partners' },
  { id: 'product_types', label: 'Product Types' },
  { id: 'service_list', label: 'Services' },
  { id: 'createdAt', label: 'Date Added' },
];

// Default fields to export
const defaultBdspExportFields = [
  'ubdsp_number', 'name', 'postal_address', 'registration_status', 'product_types'
];

export default function BdspIndex({
  data,
  filters,
  stats,
  permissions,
  meta
}: BdspIndexProps) {
  const isMobile = useIsMobile();
  const { props } = usePage<SharedData>();
  const currentUser = props.auth?.user;

  const [showExportDialog, setShowExportDialog] = useState(false);

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  // Handle export
  const handleExport = async (options: ExportOptions) => {
    // Convert ExportFields to ExportColumns
    const columns: ExportColumn[] = bdspExportFields
      .filter(f => options.fields.includes(f.id))
      .map(f => ({
        id: f.id,
        label: f.label,
        formatter: f.id === 'partners' || f.id === 'product_types'
          ? (value: any) => Array.isArray(value) ? value.join(', ') : value
          : f.id === 'service_list'
            ? (value: any) => Array.isArray(value)
              ? value.map((s: any) => `${s.name} (${s.cost} - ${s.duration})`).join('; ')
              : value
            : undefined,
        width: f.id === 'service_list' ? 40 : 20,
      }));

    // Prepare statistics if requested
    const statistics = options.includeStats && stats ? {
      'Total BDSPs': stats.totalBdsps,
      'Active BDSPs': stats.activeBdsps,
      'Pending BDSPs': stats.pendingBdsps,
      'Rejected BDSPs': stats.rejectedBdsps,
      'Suspended BDSPs': stats.suspendedBdsps,
    } : undefined;

    // Prepare data for export
    const exportRows = data.data.map(bdsp => ({
      ...bdsp,
      createdAt: bdsp.createdAt || (bdsp as any).created_at,
    }));

    await exportData(
      {
        rows: exportRows,
        columns,
        statistics,
        title: 'BDSP Export',
      },
      {
        ...options,
        filename: `bdsps-export-${new Date().toISOString().split('T')[0]}`,
      }
    );
  };

  // Page actions including export
  const pageActions: PageAction[] = [
    {
      key: 'export',
      label: 'Export Data',
      icon: <Download className="h-4 w-4" />,
      handler: () => setShowExportDialog(true),
    },
  ];

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

        {/* Export Dialog */}
        <ExportDialog
          open={showExportDialog}
          onClose={() => setShowExportDialog(false)}
          onExport={handleExport}
          totalItems={data.data.length}
          availableFields={bdspExportFields}
          defaultFields={defaultBdspExportFields}
          title="Export BDSPs"
          description={`Export ${data.data.length.toLocaleString()} BDSP records from the current page.`}
          showStatsOption={!!stats}
          defaultFilename="bdsps-export"
          preparedBy={currentUser?.name}
        />
      </div>
    </Admin>
  );
}
