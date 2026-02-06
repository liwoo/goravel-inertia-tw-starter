import React, { useState } from 'react';
import { Head, router, usePage } from '@inertiajs/react';
import { SharedData } from '@/types/app.d';
import { Download, Globe, GlobeLock } from 'lucide-react';
import axios from 'axios';
import {
    ProcurementNotice,
    ProcurementNoticeListResponse,
    ProcurementNoticeListRequest
} from '@/types/procurementnotice';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
    ProcurementNoticeCreateForm,
    ProcurementNoticeEditForm,
    ProcurementNoticeDetailView,
    procurementNoticeColumns,
    procurementNoticeColumnsMobile,
    procurementNoticeFilters,
    procurementNoticeSimpleFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { createSimpleFilters } from '@/lib/crud-page-utils';
import { ExportDialog } from '@/components/ExportDialog';
import { ExportField, ExportOptions, ExportColumn } from '@/types/export';
import { exportData, formatDateForExport } from '@/utils/exportUtils';
import { PageAction, CrudAction } from '@/types/crud';
import { toast } from 'sonner';

// Props interface for the ProcurementNotice Index page
interface ProcurementNoticeIndexProps {
    data: ProcurementNoticeListResponse;
    filters: ProcurementNoticeListRequest;
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

// Procurement notice export field definitions
const procurementExportFields: ExportField[] = [
    { id: 'id', label: 'ID' },
    { id: 'ref_no', label: 'Reference Number' },
    { id: 'invitation', label: 'Invitation' },
    { id: 'procured_by', label: 'Procured By' },
    { id: 'organization', label: 'Organization' },
    { id: 'procurement_type', label: 'Procurement Type' },
    { id: 'market_approach', label: 'Market Approach' },
    { id: 'open_date', label: 'Open Date' },
    { id: 'close_date', label: 'Close Date' },
    { id: 'qualifying_districts', label: 'Qualifying Districts' },
    { id: 'classification', label: 'Classification' },
    { id: 'partners', label: 'Partners' },
    { id: 'is_published', label: 'Published' },
    { id: 'minimum_qualifying_score', label: 'Min. Qualifying Score' },
    { id: 'details', label: 'Details' },
    { id: 'createdAt', label: 'Date Added' },
];

// Default fields to export
const defaultProcurementExportFields = [
    'ref_no', 'invitation', 'procured_by', 'organization',
    'procurement_type', 'open_date', 'close_date', 'is_published'
];

const simpleFilters = createSimpleFilters(procurementNoticeSimpleFilters([]));

export default function ProcurementNoticeIndex({
    data,
    filters,
    permissions,
    meta
}: ProcurementNoticeIndexProps) {
    const isMobile = useIsMobile();
    const [showExportDialog, setShowExportDialog] = useState(false);
    const { props } = usePage<SharedData>();
    const currentUser = props.auth?.user;

    const handleRefresh = () => {
        router.reload({ only: ['data'] });
    };

    // Handle toggle publish
    const handleTogglePublish = async (item: ProcurementNotice) => {
        try {
            const response = await axios.post(`/api/procurement-notices/${item.id}/toggle-publish`);
            const statusText = response.data?.data?.is_published ? 'published' : 'unpublished';
            toast.success(`Procurement notice ${statusText} successfully`);
            router.reload({ only: ['data'] });
        } catch (error: any) {
            const errorMessage = error.response?.data?.message || 'Failed to update publish status';
            toast.error(errorMessage);
        }
    };

    // Custom row actions
    const customActions: CrudAction<ProcurementNotice>[] = permissions.canEdit ? [
        {
            key: 'toggle-publish',
            label: 'Toggle Publish',
            icon: <Globe className="w-4 h-4" />,
            onClick: handleTogglePublish,
        },
    ] : [];

    // Handle export
    const handleExport = async (options: ExportOptions) => {
        // Convert ExportFields to ExportColumns
        const columns: ExportColumn[] = procurementExportFields
            .filter(f => options.fields.includes(f.id))
            .map(f => ({
                id: f.id,
                label: f.label,
                formatter: f.id === 'open_date' || f.id === 'close_date' || f.id === 'createdAt'
                    ? (value: any) => formatDateForExport(value)
                    : f.id === 'qualifying_districts' || f.id === 'classification' || f.id === 'partners'
                        ? (value: any) => Array.isArray(value) ? value.join(', ') : value
                        : f.id === 'is_published'
                            ? (value: any) => value ? 'Yes' : 'No'
                            : undefined,
                width: f.id === 'details' ? 40 : 20,
            }));

        // Prepare data for export
        const exportRows = data.data.map(notice => ({
            ...notice,
            createdAt: notice.createdAt || (notice as any).created_at,
        }));

        await exportData(
            {
                rows: exportRows,
                columns,
                title: 'Procurement Notices Export',
            },
            {
                ...options,
                filename: `procurement-notices-export-${new Date().toISOString().split('T')[0]}`,
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
        <Admin title={"Procurement Notice"}>
            <Head title="Procurement Notice - Management" />

            <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
                {/* Main CRUD Component */}
                <div className="px-0">
                    <CrudPage<ProcurementNotice>
                        data={data}
                        filters={filters}
                        title="Procurement Notices"
                        resourceName="procurement-notices"
                        columns={isMobile ? procurementNoticeColumnsMobile : procurementNoticeColumns}
                        customFilters={procurementNoticeFilters}
                        simpleFilters={simpleFilters}
                        paginationConfig={meta?.pagination}
                        createForm={ProcurementNoticeCreateForm}
                        editForm={ProcurementNoticeEditForm}
                        detailView={ProcurementNoticeDetailView}
                        onRefresh={handleRefresh}
                        canCreate={permissions.canCreate}
                        canEdit={permissions.canEdit}
                        canDelete={permissions.canDelete}
                        canView={true}
                        pageActions={pageActions}
                        actions={customActions}
                    />
                </div>
            </div>

            {/* Export Dialog */}
            <ExportDialog
                open={showExportDialog}
                onClose={() => setShowExportDialog(false)}
                onExport={handleExport}
                totalItems={data.data.length}
                availableFields={procurementExportFields}
                defaultFields={defaultProcurementExportFields}
                title="Export Procurement Notices"
                description={`Export ${data.data.length.toLocaleString()} procurement notice records from the current page.`}
                showStatsOption={false}
                defaultFilename="procurement-notices-export"
                preparedBy={currentUser?.name}
            />
        </Admin>
    );
}
