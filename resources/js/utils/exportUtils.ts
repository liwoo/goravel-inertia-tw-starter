/**
 * Export Utilities
 *
 * Generic export functionality for CSV, JSON, Excel, and PDF formats.
 * Works with any dataset containing rows and column definitions.
 */

import {
    ExportFormat,
    ExportField,
    ExportOptions,
    ExportColumn,
    ExportData
} from '@/types/export';

// Re-export types for convenience
export type { ExportColumn, ExportData } from '@/types/export';

/** Default maximum rows for export */
export const DEFAULT_MAX_EXPORT_ROWS = 200;

/** Export payload passed to format handlers */
interface ExportPayload {
    rows: Record<string, any>[];
    columns: ExportColumn[];
    statistics?: Record<string, any>;
    title?: string;
    preparedBy?: string;
}

/**
 * Main export function - routes to appropriate format handler
 */
export async function exportData(
    data: ExportData,
    options: ExportOptions
): Promise<void> {
    const { format, fields, includeStats, filename = 'export', maxRows = DEFAULT_MAX_EXPORT_ROWS, preparedBy } = options;

    // Filter columns based on selected fields
    const selectedColumns = data.columns.filter(col => fields.includes(col.id));

    // Apply row limit
    const limitedRows = data.rows.slice(0, maxRows);

    // Filter row data to only include selected fields
    const filteredRows = limitedRows.map(row => {
        const filteredRow: Record<string, any> = {};
        selectedColumns.forEach(col => {
            const value = row[col.id];
            filteredRow[col.id] = col.formatter ? col.formatter(value, row) : value;
        });
        return filteredRow;
    });

    const exportPayload: ExportPayload = {
        rows: filteredRows,
        columns: selectedColumns,
        statistics: includeStats ? data.statistics : undefined,
        title: data.title,
        preparedBy,
    };

    switch (format) {
        case 'csv':
            await exportToCSV(exportPayload, filename);
            break;
        case 'json':
            await exportToJSON(exportPayload, filename);
            break;
        case 'excel':
            await exportToExcel(exportPayload, filename);
            break;
        case 'pdf':
            await exportToPDF(exportPayload, filename);
            break;
        default:
            throw new Error(`Unsupported export format: ${format}`);
    }
}

/**
 * Export to CSV format
 */
async function exportToCSV(
    data: ExportPayload,
    filename: string
): Promise<void> {
    const { rows, columns, statistics, title, preparedBy } = data;

    // Build CSV content
    const lines: string[] = [];

    // Add title if present
    if (title) {
        lines.push(escapeCSV(title));
        lines.push('');
    }

    // Add export metadata
    lines.push(`Exported on,${escapeCSV(new Date().toLocaleString())}`);
    if (preparedBy) {
        lines.push(`Prepared by,${escapeCSV(preparedBy)}`);
    }
    lines.push(`Total records,${rows.length}`);
    lines.push('');

    // Add headers
    const headers = columns.map(col => escapeCSV(col.label));
    lines.push(headers.join(','));

    // Add data rows
    rows.forEach(row => {
        const values = columns.map(col => escapeCSV(formatValue(row[col.id])));
        lines.push(values.join(','));
    });

    // Add statistics if present
    if (statistics) {
        lines.push('');
        lines.push('Statistics');
        Object.entries(statistics).forEach(([key, value]) => {
            lines.push(`${escapeCSV(key)},${escapeCSV(formatValue(value))}`);
        });
    }

    const csvContent = lines.join('\n');
    downloadFile(csvContent, `${filename}.csv`, 'text/csv;charset=utf-8;');
}

/**
 * Export to JSON format
 */
async function exportToJSON(
    data: ExportPayload,
    filename: string
): Promise<void> {
    const { rows, columns, statistics, title, preparedBy } = data;

    const exportObj: any = {
        exportedAt: new Date().toISOString(),
        totalRecords: rows.length,
    };

    if (title) {
        exportObj.title = title;
    }

    if (preparedBy) {
        exportObj.preparedBy = preparedBy;
    }

    // Include column metadata
    exportObj.columns = columns.map(col => ({
        id: col.id,
        label: col.label,
    }));

    // Include data
    exportObj.data = rows;

    // Include statistics if present
    if (statistics) {
        exportObj.statistics = statistics;
    }

    const jsonContent = JSON.stringify(exportObj, null, 2);
    downloadFile(jsonContent, `${filename}.json`, 'application/json');
}

/**
 * Export to Excel format using SheetJS (xlsx)
 */
async function exportToExcel(
    data: ExportPayload,
    filename: string
): Promise<void> {
    const { rows, columns, statistics, title, preparedBy } = data;

    // Dynamically import xlsx library
    const XLSX = await import('xlsx');

    // Create workbook
    const wb = XLSX.utils.book_new();

    // Prepare data for main sheet
    const wsData: any[][] = [];

    // Add title row if present
    if (title) {
        wsData.push([title]);
        wsData.push([]);
    }

    // Add export metadata
    wsData.push(['Exported on', new Date().toLocaleString()]);
    if (preparedBy) {
        wsData.push(['Prepared by', preparedBy]);
    }
    wsData.push(['Total records', rows.length]);
    wsData.push([]);

    // Add headers
    wsData.push(columns.map(col => col.label));

    // Add data rows
    rows.forEach(row => {
        wsData.push(columns.map(col => formatValue(row[col.id])));
    });

    // Create main worksheet
    const ws = XLSX.utils.aoa_to_sheet(wsData);

    // Set column widths
    ws['!cols'] = columns.map(col => ({
        wch: col.width || Math.max(col.label.length, 15),
    }));

    XLSX.utils.book_append_sheet(wb, ws, 'Data');

    // Add statistics sheet if present
    if (statistics) {
        const statsData: any[][] = [['Statistic', 'Value']];
        Object.entries(statistics).forEach(([key, value]) => {
            statsData.push([key, formatValue(value)]);
        });

        const statsWs = XLSX.utils.aoa_to_sheet(statsData);
        statsWs['!cols'] = [{ wch: 30 }, { wch: 20 }];
        XLSX.utils.book_append_sheet(wb, statsWs, 'Statistics');
    }

    // Generate and download
    const excelBuffer = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
    const blob = new Blob([excelBuffer], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    });
    downloadBlob(blob, `${filename}.xlsx`);
}

/**
 * Export to PDF format using jsPDF
 */
async function exportToPDF(
    data: ExportPayload,
    filename: string
): Promise<void> {
    const { rows, columns, statistics, title, preparedBy } = data;

    // Dynamically import jsPDF and autoTable
    const jsPDFModule = await import('jspdf');
    const jsPDF = jsPDFModule.default;
    const autoTableModule = await import('jspdf-autotable');
    const autoTable = autoTableModule.default;

    // Create PDF document
    const doc = new jsPDF({
        orientation: columns.length > 5 ? 'landscape' : 'portrait',
        unit: 'mm',
        format: 'a4',
    });

    let yPosition = 15;

    // Add title if present
    if (title) {
        doc.setFontSize(16);
        doc.setFont('helvetica', 'bold');
        doc.text(title, 14, yPosition);
        yPosition += 10;
    }

    // Add export info
    doc.setFontSize(10);
    doc.setFont('helvetica', 'normal');
    doc.setTextColor(100);
    const exportInfo = preparedBy
        ? `Exported on ${new Date().toLocaleString()} | Prepared by: ${preparedBy} | ${rows.length} records`
        : `Exported on ${new Date().toLocaleString()} | ${rows.length} records`;
    doc.text(exportInfo, 14, yPosition);
    yPosition += 8;

    // Reset text color
    doc.setTextColor(0);

    // Prepare table data
    const tableHeaders = columns.map(col => col.label);
    const tableData = rows.map(row =>
        columns.map(col => formatValue(row[col.id]))
    );

    // Add main data table
    autoTable(doc, {
        head: [tableHeaders],
        body: tableData,
        startY: yPosition,
        styles: {
            fontSize: 8,
            cellPadding: 2,
        },
        headStyles: {
            fillColor: [59, 130, 246], // Blue header
            textColor: 255,
            fontStyle: 'bold',
        },
        alternateRowStyles: {
            fillColor: [245, 247, 250],
        },
        margin: { left: 14, right: 14 },
        tableWidth: 'auto',
    });

    // Add statistics if present
    if (statistics) {
        const finalY = (doc as any).lastAutoTable.finalY || yPosition;

        // Check if we need a new page
        if (finalY > doc.internal.pageSize.height - 50) {
            doc.addPage();
            yPosition = 15;
        } else {
            yPosition = finalY + 10;
        }

        doc.setFontSize(12);
        doc.setFont('helvetica', 'bold');
        doc.text('Statistics Summary', 14, yPosition);
        yPosition += 6;

        const statsData = Object.entries(statistics).map(([key, value]) => [
            key,
            formatValue(value),
        ]);

        autoTable(doc, {
            head: [['Statistic', 'Value']],
            body: statsData,
            startY: yPosition,
            styles: {
                fontSize: 9,
                cellPadding: 3,
            },
            headStyles: {
                fillColor: [100, 116, 139], // Gray header
                textColor: 255,
            },
            margin: { left: 14, right: 14 },
            tableWidth: 'wrap',
        });
    }

    // Add page numbers
    const pageCount = doc.getNumberOfPages();
    for (let i = 1; i <= pageCount; i++) {
        doc.setPage(i);
        doc.setFontSize(8);
        doc.setTextColor(150);
        doc.text(
            `Page ${i} of ${pageCount}`,
            doc.internal.pageSize.width / 2,
            doc.internal.pageSize.height - 10,
            { align: 'center' }
        );
    }

    // Download PDF
    doc.save(`${filename}.pdf`);
}

/**
 * Helper: Escape CSV value
 */
function escapeCSV(value: string): string {
    if (value === null || value === undefined) {
        return '';
    }
    const str = String(value);
    // Escape quotes and wrap in quotes if contains comma, quote, or newline
    if (str.includes(',') || str.includes('"') || str.includes('\n')) {
        return `"${str.replace(/"/g, '""')}"`;
    }
    return str;
}

/**
 * Helper: Format value for display
 */
function formatValue(value: any): string {
    if (value === null || value === undefined) {
        return '';
    }
    if (value instanceof Date) {
        return value.toISOString();
    }
    if (typeof value === 'boolean') {
        return value ? 'Yes' : 'No';
    }
    if (typeof value === 'object') {
        return JSON.stringify(value);
    }
    return String(value);
}

/**
 * Helper: Download file from string content
 */
function downloadFile(content: string, filename: string, mimeType: string): void {
    const blob = new Blob([content], { type: mimeType });
    downloadBlob(blob, filename);
}

/**
 * Helper: Download blob as file
 */
function downloadBlob(blob: Blob, filename: string): void {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
}

/**
 * Utility: Convert ExportField array to ExportColumn array
 */
export function fieldsToColumns(fields: ExportField[]): ExportColumn[] {
    return fields.map(field => ({
        id: field.id,
        label: field.label,
        formatter: field.formatter,
    }));
}

/**
 * Utility: Create export handler for use with ExportDialog
 */
export function createExportHandler(
    getData: () => Promise<Record<string, any>[]> | Record<string, any>[],
    columns: ExportColumn[],
    options?: {
        title?: string;
        getStatistics?: () => Record<string, any>;
    }
) {
    return async (exportOptions: ExportOptions): Promise<void> => {
        const rows = await getData();
        const statistics = options?.getStatistics?.();

        await exportData(
            {
                rows,
                columns,
                statistics,
                title: options?.title,
            },
            exportOptions
        );
    };
}

/**
 * Utility: Format date for export
 */
export function formatDateForExport(date: string | Date | null | undefined): string {
    if (!date) return '';
    const d = typeof date === 'string' ? new Date(date) : date;
    if (isNaN(d.getTime())) return '';
    return d.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    });
}

/**
 * Utility: Format datetime for export
 */
export function formatDateTimeForExport(date: string | Date | null | undefined): string {
    if (!date) return '';
    const d = typeof date === 'string' ? new Date(date) : date;
    if (isNaN(d.getTime())) return '';
    return d.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
}

/**
 * Utility: Format currency for export
 */
export function formatCurrencyForExport(
    value: number | null | undefined,
    currency: string = 'MWK'
): string {
    if (value === null || value === undefined) return '';
    return new Intl.NumberFormat('en-MW', {
        style: 'currency',
        currency,
        minimumFractionDigits: 2,
    }).format(value);
}

/**
 * Utility: Format number for export
 */
export function formatNumberForExport(value: number | null | undefined): string {
    if (value === null || value === undefined) return '';
    return new Intl.NumberFormat('en-US').format(value);
}

/**
 * Utility: Format percentage for export
 */
export function formatPercentageForExport(value: number | null | undefined): string {
    if (value === null || value === undefined) return '';
    return `${(value * 100).toFixed(1)}%`;
}
