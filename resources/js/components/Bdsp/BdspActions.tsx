import React, { useState } from 'react';
import { Download } from 'lucide-react';
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

/**
 * BDSP Export Dialog
 */
export interface BdspExportOptions {
    format: 'csv' | 'json' | 'excel' | 'pdf';
    fields: string[];
    includeStats: boolean;
}

interface BdspExportProps {
    onClose: () => void;
    onExport: (options: BdspExportOptions) => void;
    totalItems: number;
}

export function BdspExportDialog({
    onClose,
    onExport,
    totalItems
}: BdspExportProps) {
    const [format, setFormat] = useState<'csv' | 'json' | 'excel' | 'pdf'>('csv');
    const [fields, setFields] = useState<string[]>([
        'name', 'email', 'phone', 'registration_status', 'postal_address'
    ]);
    const [includeStats, setIncludeStats] = useState(false);
    const [isExporting, setIsExporting] = useState(false);

    const availableFields = [
        { id: 'name', label: 'Name' },
        { id: 'email', label: 'Email' },
        { id: 'phone', label: 'Phone' },
        { id: 'postal_address', label: 'Postal Address' },
        { id: 'physical_address', label: 'Physical Address' },
        { id: 'registration_status', label: 'Status' },
        { id: 'created_at', label: 'Created Date' },
        { id: 'updated_at', label: 'Updated Date' },
    ];

    const toggleField = (fieldId: string) => {
        setFields(prev =>
            prev.includes(fieldId)
                ? prev.filter(f => f !== fieldId)
                : [...prev, fieldId]
        );
    };

    const handleExport = async () => {
        setIsExporting(true);
        try {
            await onExport({
                format,
                fields,
                includeStats,
            });
            onClose();
        } finally {
            setIsExporting(false);
        }
    };

    return (
        <Dialog open={true} onOpenChange={onClose}>
            <DialogContent className="max-w-md">
                <DialogHeader>
                    <DialogTitle>Export BDSPs</DialogTitle>
                    <DialogDescription>
                        Export {totalItems} BDSP(s) to file.
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    <div>
                        <Label htmlFor="format">Export Format</Label>
                        <Select value={format} onValueChange={(value: any) => setFormat(value)}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="csv">CSV</SelectItem>
                                <SelectItem value="json">JSON</SelectItem>
                                <SelectItem value="excel">Excel</SelectItem>
                                <SelectItem value="pdf">PDF</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div>
                        <Label>Fields to Include</Label>
                        <div className="grid grid-cols-2 gap-2 mt-2">
                            {availableFields.map((field) => (
                                <div key={field.id} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={field.id}
                                        checked={fields.includes(field.id)}
                                        onCheckedChange={() => toggleField(field.id)}
                                    />
                                    <Label htmlFor={field.id} className="text-sm">
                                        {field.label}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="flex items-center space-x-2">
                        <Checkbox
                            id="includeStats"
                            checked={includeStats}
                            onCheckedChange={(checked) => setIncludeStats(!!checked)}
                        />
                        <Label htmlFor="includeStats" className="text-sm">
                            Include statistics summary
                        </Label>
                    </div>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={onClose}>
                        Cancel
                    </Button>
                    <Button
                        onClick={handleExport}
                        disabled={fields.length === 0 || isExporting}
                    >
                        <Download className="w-4 h-4 mr-2" />
                        {isExporting ? 'Exporting...' : 'Export'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
