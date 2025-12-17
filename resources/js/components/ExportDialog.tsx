import { useState } from 'react';
import { Download, FileSpreadsheet, FileJson, FileText, File } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { ExportFormat, ExportField, ExportOptions } from '@/types/export';

// Re-export types for backward compatibility
export type { ExportFormat, ExportField, ExportOptions } from '@/types/export';

/** Default maximum rows for export */
export const DEFAULT_MAX_EXPORT_ROWS = 200;

export interface ExportDialogProps {
    /** Whether the dialog is open */
    open: boolean;
    /** Callback when dialog is closed */
    onClose: () => void;
    /** Callback when export is triggered - receives options and should return a promise */
    onExport: (options: ExportOptions) => Promise<void>;
    /** Total number of items to export */
    totalItems: number;
    /** Available fields that can be exported */
    availableFields: ExportField[];
    /** Default selected field IDs */
    defaultFields?: string[];
    /** Title for the dialog */
    title?: string;
    /** Description for the dialog */
    description?: string;
    /** Whether to show the "Include statistics" option */
    showStatsOption?: boolean;
    /** Default filename (without extension) */
    defaultFilename?: string;
    /** Maximum rows to export (default: 200) */
    maxRows?: number;
    /** Name of user preparing the export */
    preparedBy?: string;
}

const formatIcons: Record<ExportFormat, React.ReactNode> = {
    csv: <FileText className="w-4 h-4" />,
    json: <FileJson className="w-4 h-4" />,
    excel: <FileSpreadsheet className="w-4 h-4" />,
    pdf: <File className="w-4 h-4" />,
};

const formatLabels: Record<ExportFormat, string> = {
    csv: 'CSV (Comma Separated)',
    json: 'JSON',
    excel: 'Excel (.xlsx)',
    pdf: 'PDF Document',
};

export function ExportDialog({
    open,
    onClose,
    onExport,
    totalItems,
    availableFields,
    defaultFields,
    title = 'Export Data',
    description,
    showStatsOption = true,
    defaultFilename = 'export',
    maxRows = DEFAULT_MAX_EXPORT_ROWS,
    preparedBy,
}: ExportDialogProps) {
    const [format, setFormat] = useState<ExportFormat>('csv');
    const [fields, setFields] = useState<string[]>(
        defaultFields ?? availableFields.map(f => f.id)
    );
    const [includeStats, setIncludeStats] = useState(false);
    const [isExporting, setIsExporting] = useState(false);

    // Calculate actual rows to export
    const rowsToExport = Math.min(totalItems, maxRows);
    const isLimited = totalItems > maxRows;

    const toggleField = (fieldId: string) => {
        setFields(prev =>
            prev.includes(fieldId)
                ? prev.filter(f => f !== fieldId)
                : [...prev, fieldId]
        );
    };

    const selectAll = () => {
        setFields(availableFields.map(f => f.id));
    };

    const deselectAll = () => {
        setFields([]);
    };

    const handleExport = async () => {
        if (fields.length === 0) return;

        setIsExporting(true);
        try {
            await onExport({
                format,
                fields,
                includeStats,
                filename: defaultFilename,
                maxRows,
                preparedBy,
            });
            onClose();
        } catch (error) {
            console.error('Export failed:', error);
        } finally {
            setIsExporting(false);
        }
    };

    const dynamicDescription = description ?? (
        isLimited
            ? `Export up to ${maxRows.toLocaleString()} of ${totalItems.toLocaleString()} records to a file.`
            : `Export ${totalItems.toLocaleString()} record${totalItems !== 1 ? 's' : ''} to a file.`
    );

    return (
        <Dialog open={open} onOpenChange={(isOpen) => !isOpen && onClose()}>
            <DialogContent className="max-w-md">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Download className="w-5 h-5" />
                        {title}
                    </DialogTitle>
                    <DialogDescription>
                        {dynamicDescription}
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    {/* Format Selection */}
                    <div className="space-y-2">
                        <Label htmlFor="format">Export Format</Label>
                        <Select value={format} onValueChange={(value: ExportFormat) => setFormat(value)}>
                            <SelectTrigger id="format">
                                <SelectValue>
                                    <div className="flex items-center gap-2">
                                        {formatIcons[format]}
                                        <span>{formatLabels[format]}</span>
                                    </div>
                                </SelectValue>
                            </SelectTrigger>
                            <SelectContent>
                                {(Object.keys(formatLabels) as ExportFormat[]).map((fmt) => (
                                    <SelectItem key={fmt} value={fmt}>
                                        <div className="flex items-center gap-2">
                                            {formatIcons[fmt]}
                                            <span>{formatLabels[fmt]}</span>
                                        </div>
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Fields Selection */}
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label>Fields to Include</Label>
                            <div className="flex gap-2">
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    className="h-auto py-1 px-2 text-xs"
                                    onClick={selectAll}
                                >
                                    Select All
                                </Button>
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    className="h-auto py-1 px-2 text-xs"
                                    onClick={deselectAll}
                                >
                                    Deselect All
                                </Button>
                            </div>
                        </div>
                        <div className="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto border rounded-md p-3 bg-muted/30">
                            {availableFields.map((field) => (
                                <div key={field.id} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={`field-${field.id}`}
                                        checked={fields.includes(field.id)}
                                        onCheckedChange={() => toggleField(field.id)}
                                    />
                                    <Label
                                        htmlFor={`field-${field.id}`}
                                        className="text-sm font-normal cursor-pointer"
                                    >
                                        {field.label}
                                    </Label>
                                </div>
                            ))}
                        </div>
                        <p className="text-xs text-muted-foreground">
                            {fields.length} of {availableFields.length} fields selected
                        </p>
                    </div>

                    {/* Statistics Option */}
                    {showStatsOption && (
                        <div className="flex items-center space-x-2 pt-2 border-t">
                            <Checkbox
                                id="includeStats"
                                checked={includeStats}
                                onCheckedChange={(checked) => setIncludeStats(!!checked)}
                            />
                            <Label htmlFor="includeStats" className="text-sm font-normal cursor-pointer">
                                Include statistics summary
                            </Label>
                        </div>
                    )}

                    {/* Row limit warning */}
                    {isLimited && (
                        <div className="text-xs text-amber-600 bg-amber-50 dark:bg-amber-950/30 dark:text-amber-400 p-2 rounded-md border border-amber-200 dark:border-amber-800">
                            Limited to {maxRows.toLocaleString()} rows. Use filters to narrow down your data for a complete export.
                        </div>
                    )}
                </div>

                <DialogFooter className="gap-2 sm:gap-0">
                    <Button variant="outline" onClick={onClose} disabled={isExporting}>
                        Cancel
                    </Button>
                    <Button
                        onClick={handleExport}
                        disabled={fields.length === 0 || isExporting}
                    >
                        {isExporting ? (
                            <>
                                <span className="animate-spin mr-2">⏳</span>
                                Exporting...
                            </>
                        ) : (
                            <>
                                <Download className="w-4 h-4 mr-2" />
                                Export {rowsToExport.toLocaleString()} Record{rowsToExport !== 1 ? 's' : ''}
                            </>
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

export default ExportDialog;
