import React from 'react';
import { ProcurementNotice } from '@/types/procurementnotice';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { FileText, Calendar, Hash, Building2, CheckCircle, XCircle } from 'lucide-react';

/**
 * Procurement Notice table columns configuration
 */
export const procurementNoticeColumns: CrudColumn<ProcurementNotice>[] = [
    {
        key: 'procured_by',
        label: 'Procured By',
        sortable: true,
        className: 'min-w-[150px]',
        render: (item) => (
            <div className="font-medium text-gray-900 flex items-center">
                <p className="font-medium text-foreground">{item.procured_by}</p>
            </div>
        ),
    },
    {
        key: 'organization',
        label: 'Organization',
        sortable: true,
        className: 'min-w-[150px]',
        render: (item) => (
            <div className="font-medium text-gray-900 flex items-center">
                <p className="font-medium text-foreground">{item.organization}</p>
            </div>
        ),
    },
    {
        key: 'ref_no',
        label: 'Ref No',
        sortable: true,
        className: 'min-w-[150px]',
        render: (item) => (
            <div className="font-medium text-gray-900 flex items-center">
                <p className="font-medium text-foreground">{item.ref_no}</p>
            </div>
        ),
    },
    {
        key: 'procurement_type',
        label: 'Type',
        sortable: true,
        className: 'min-w-[150px]',
        render: (item) => (
            <div className="font-medium text-gray-900 flex items-center">
                <p className="font-medium text-foreground">{item.procurement_type}</p>
            </div>
        ),
    },
    {
        key: 'open_date',
        label: 'Period',
        sortable: true,
        className: 'min-w-[150px]',
        render: (item) => (
            <div className="space-y-1">
                <div className="font-medium text-gray-900 flex items-center">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <p className="font-medium text-foreground">
                        {item.open_date ? new Date(item.open_date).toLocaleDateString() : '-'} -
                        {item.close_date ? new Date(item.close_date).toLocaleDateString() : '-'}
                    </p>
                </div>
            </div>
        ),
    },
    {
        key: 'is_published',
        label: 'Status',
        sortable: true,
        className: 'min-w-[100px]',
        render: (item) => (
            <Badge variant={item.is_published ? 'default' : 'secondary'}>
                {item.is_published ? (
                    <div className="flex items-center gap-1">
                        <CheckCircle className="h-3 w-3" /> Published
                    </div>
                ) : (
                    <div className="flex items-center gap-1">
                        <XCircle className="h-3 w-3" /> Draft
                    </div>
                )}
            </Badge>
        ),
    },
    {
        key: 'createdAt',
        label: 'Created',
        sortable: true,
        className: 'w-32',
        render: (item) => {
            const dateValue = item.createdAt || item.created_at;
            if (!dateValue) return <span className="text-sm text-muted-foreground">-</span>;

            const date = new Date(dateValue);
            return (
                <div className="text-sm text-muted-foreground">
                    {date.toLocaleDateString('en-US', {
                        month: 'short',
                        day: 'numeric',
                        year: 'numeric'
                    })}
                </div>
            );
        },
    },
];

/**
 * Compact columns for mobile/smaller screens
 */
export const procurementNoticeColumnsMobile: CrudColumn<ProcurementNotice>[] = [
    {
        key: 'combined',
        label: 'Notice',
        sortable: false,
        render: (item) => (
            <div className="space-y-2">
                <div className="font-medium text-foreground">{item.procured_by}</div>
                <div className="text-sm text-muted-foreground">
                    {item.ref_no} • {item.procurement_type}
                </div>
            </div>
        ),
    },
];

/**
 * Procurement Notice filters configuration
 */
export const procurementNoticeFilters: CrudFilter[] = [
    {
        key: 'procured_by',
        label: 'Procured By',
        type: 'text',
        placeholder: 'Enter procured by',
    },
    {
        key: 'procurement_type',
        label: 'Type',
        type: 'text',
        placeholder: 'Enter type',
    },
    {
        key: 'ref_no',
        label: 'Ref No',
        type: 'text',
        placeholder: 'Enter ref no',
    },
    {
        key: 'organization',
        label: 'Organization',
        type: 'text',
        placeholder: 'Enter organization',
    },
    {
        key: 'is_published',
        label: 'Published',
        type: 'select',
        options: [
            { label: 'Yes', value: 'true' },
            { label: 'No', value: 'false' },
        ],
    },
];
