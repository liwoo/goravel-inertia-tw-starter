import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
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
const simpleFilters = createSimpleFilters(procurementNoticeSimpleFilters([]));

export default function ProcurementNoticeIndex({
    data,
    filters,
    permissions,
    meta
}: ProcurementNoticeIndexProps) {
    const isMobile = useIsMobile();

    const handleRefresh = () => {
        router.reload({ only: ['data'] });
    };

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
                    />
                </div>
            </div>
        </Admin>
    );
}
