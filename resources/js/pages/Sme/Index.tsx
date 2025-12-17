import React, { useState } from 'react';
import { Head, router, usePage } from '@inertiajs/react';
import { SharedData } from '@/types/app.d';
import { Download, Power, PowerOff } from 'lucide-react';
import { toast } from 'sonner';
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
import { ExportDialog } from '@/components/ExportDialog';
import { ExportField, ExportOptions, ExportColumn } from '@/types/export';
import { exportData, formatDateForExport } from '@/utils/exportUtils';
import { PageAction, BulkAction } from '@/types/crud';

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

// SME export field definitions
const smeExportFields: ExportField[] = [
  { id: 'usmeNumber', label: 'USME Number' },
  { id: 'name', label: 'Business Name' },
  { id: 'registrationNumber', label: 'Registration Number' },
  { id: 'taxIdentificationNumber', label: 'TIN' },
  { id: 'businessCategory', label: 'Category' },
  { id: 'sector', label: 'Sector' },
  { id: 'subSector', label: 'Sub-Sector' },
  { id: 'businessDescription', label: 'Description' },
  { id: 'contactPhone', label: 'Phone' },
  { id: 'contactEmail', label: 'Email' },
  { id: 'physicalAddress', label: 'Physical Address' },
  { id: 'postalAddress', label: 'Postal Address' },
  { id: 'website', label: 'Website' },
  { id: 'district', label: 'District' },
  { id: 'traditionalAuthority', label: 'Traditional Authority' },
  { id: 'operationalStartDate', label: 'Operational Start Date' },
  { id: 'formalisationScore', label: 'Formalisation Score' },
  { id: 'createdAt', label: 'Date Added' },
];

// Default fields to export
const defaultSmeExportFields = [
  'usmeNumber', 'name', 'businessCategory', 'sector',
  'district', 'contactPhone', 'contactEmail', 'registrationNumber'
];

export default function SmeIndex({
  data,
  filters,
  permissions,
  meta,
  stats
}: SmeIndexProps) {
  const isMobile = useIsMobile();
  const [showExportDialog, setShowExportDialog] = useState(false);
  const { props } = usePage<SharedData>();
  const currentUser = props.auth?.user;

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Handle export
  const handleExport = async (options: ExportOptions) => {
    // Convert ExportFields to ExportColumns
    const columns: ExportColumn[] = smeExportFields
      .filter(f => options.fields.includes(f.id))
      .map(f => ({
        id: f.id,
        label: f.label,
        formatter: f.id.includes('Date') || f.id === 'createdAt'
          ? (value: any) => formatDateForExport(value)
          : undefined,
        width: f.id === 'businessDescription' ? 40 : 20,
      }));

    // Prepare statistics if requested
    const statistics = options.includeStats && stats ? {
      'Total SMEs': stats.totalSmes,
      'New This Month': stats.newThisMonth,
      'New Last Month': stats.newLastMonth,
      'Has Registration': stats.hasRegistration,
      'Registration Rate': `${stats.hasRegistrationPercentage.toFixed(1)}%`,
    } : undefined;

    // Use current page data - for full export, would need API call
    const exportRows = data.data.map(sme => ({
      ...sme,
      usmeNumber: sme.usmeNumber || (sme as any).usme_number,
      businessCategory: sme.businessCategory || (sme as any).business_category,
      createdAt: sme.createdAt || (sme as any).created_at,
      formalisationScore:
        sme.businessFormalisation?.formalisationScore ??
        sme.businessFormalisation?.formalisation_score ??
        (sme as any).business_formalisation?.formalisationScore ??
        (sme as any).business_formalisation?.formalisation_score ??
        '',
    }));

    await exportData(
      {
        rows: exportRows,
        columns,
        statistics,
        title: 'SME Export',
      },
      {
        ...options,
        filename: `smes-export-${new Date().toISOString().split('T')[0]}`,
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

  // Bulk actions for SME management
  const bulkActions: BulkAction[] = [
    {
      key: 'deactivate',
      label: 'Deactivate',
      icon: <PowerOff className="h-4 w-4" />,
      variant: 'warning',
      confirm: true,
      confirmMessage: 'Are you sure you want to deactivate the selected SMEs?',
    },
    {
      key: 'activate',
      label: 'Activate',
      icon: <Power className="h-4 w-4" />,
      variant: 'default',
      confirm: true,
      confirmMessage: 'Are you sure you want to activate the selected SMEs?',
    },
  ];

  // Handle bulk actions
  const handleBulkAction = async (action: string, selectedIds: number[]) => {
    if (selectedIds.length === 0) return;

    try {
      const endpoint = action === 'deactivate'
        ? '/api/smes/bulk-deactivate'
        : '/api/smes/bulk-activate';

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ ids: selectedIds }),
      });

      const result = await response.json();

      if (response.ok) {
        toast.success(result.message || `${selectedIds.length} SME(s) ${action}d successfully`);
        router.reload({ only: ['data'] });
      } else {
        toast.error(result.message || `Failed to ${action} SMEs`);
      }
    } catch (error) {
      console.error(`Bulk ${action} error:`, error);
      toast.error(`Failed to ${action} SMEs`);
    }
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
            pageActions={pageActions}
            bulkActions={permissions.canEdit ? bulkActions : []}
            onBulkAction={handleBulkAction}
          />
        </div>
      </div>

      {/* Export Dialog */}
      <ExportDialog
        open={showExportDialog}
        onClose={() => setShowExportDialog(false)}
        onExport={handleExport}
        totalItems={data.data.length}
        availableFields={smeExportFields}
        defaultFields={defaultSmeExportFields}
        title="Export SMEs"
        description={`Export ${data.data.length.toLocaleString()} SME records from the current page.`}
        showStatsOption={!!stats}
        defaultFilename="smes-export"
        preparedBy={currentUser?.name}
      />
    </Admin>
  );
}
