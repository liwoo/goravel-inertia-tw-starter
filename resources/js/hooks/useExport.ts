/**
 * useExport Hook
 *
 * Custom hook for easy integration of export functionality with any data table.
 */

import { useState, useCallback } from 'react';
import { ExportField, ExportOptions, ExportColumn } from '@/types/export';
import { exportData, fieldsToColumns } from '@/utils/exportUtils';

export interface UseExportOptions {
    /** Function to get the data to export */
    getData: () => Promise<Record<string, any>[]> | Record<string, any>[];
    /** Column/field definitions */
    fields: ExportField[];
    /** Default filename (without extension) */
    filename?: string;
    /** Title for the export */
    title?: string;
    /** Function to get statistics */
    getStatistics?: () => Record<string, any>;
    /** Callback after successful export */
    onSuccess?: () => void;
    /** Callback on export error */
    onError?: (error: Error) => void;
}

export interface UseExportReturn {
    /** Whether the export dialog is open */
    isOpen: boolean;
    /** Open the export dialog */
    openDialog: () => void;
    /** Close the export dialog */
    closeDialog: () => void;
    /** Handle export - pass this to ExportDialog's onExport prop */
    handleExport: (options: ExportOptions) => Promise<void>;
    /** Whether an export is currently in progress */
    isExporting: boolean;
    /** Any export error that occurred */
    error: Error | null;
}

export function useExport({
    getData,
    fields,
    filename = 'export',
    title,
    getStatistics,
    onSuccess,
    onError,
}: UseExportOptions): UseExportReturn {
    const [isOpen, setIsOpen] = useState(false);
    const [isExporting, setIsExporting] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const openDialog = useCallback(() => {
        setIsOpen(true);
        setError(null);
    }, []);

    const closeDialog = useCallback(() => {
        setIsOpen(false);
    }, []);

    const handleExport = useCallback(async (options: ExportOptions) => {
        setIsExporting(true);
        setError(null);

        try {
            // Get the data
            const rows = await getData();

            // Convert fields to columns
            const columns: ExportColumn[] = fieldsToColumns(fields);

            // Get statistics if available
            const statistics = getStatistics?.();

            // Perform export
            await exportData(
                {
                    rows,
                    columns,
                    statistics,
                    title,
                },
                {
                    ...options,
                    filename: options.filename || filename,
                }
            );

            onSuccess?.();
        } catch (err) {
            const error = err instanceof Error ? err : new Error('Export failed');
            setError(error);
            onError?.(error);
            throw error;
        } finally {
            setIsExporting(false);
        }
    }, [getData, fields, filename, title, getStatistics, onSuccess, onError]);

    return {
        isOpen,
        openDialog,
        closeDialog,
        handleExport,
        isExporting,
        error,
    };
}

export default useExport;
