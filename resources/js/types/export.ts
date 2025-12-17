/**
 * Export Types
 *
 * Shared type definitions for export functionality.
 */

export type ExportFormat = 'csv' | 'json' | 'excel' | 'pdf';

export interface ExportField {
    id: string;
    label: string;
    /** Optional: specify how to format the value */
    formatter?: (value: any) => string;
}

export interface ExportOptions {
    format: ExportFormat;
    fields: string[];
    includeStats: boolean;
    filename?: string;
    /** Maximum rows to export (default: 200) */
    maxRows?: number;
    /** Name of user who prepared/triggered the export */
    preparedBy?: string;
}

export interface ExportColumn {
    id: string;
    label: string;
    /** Optional formatter for the column value */
    formatter?: (value: any, row: Record<string, any>) => string;
    /** Width for PDF/Excel columns (in characters) */
    width?: number;
}

export interface ExportData {
    /** Array of data rows */
    rows: Record<string, any>[];
    /** Column definitions */
    columns: ExportColumn[];
    /** Optional statistics to include */
    statistics?: Record<string, any>;
    /** Title for the export (used in PDF header) */
    title?: string;
}
