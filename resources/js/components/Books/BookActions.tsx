import React, { useState } from 'react';
import { router } from '@inertiajs/react';
import { 
  Download, 
  Upload, 
  FileText, 
  Trash2, 
  Tag, 
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Info
} from 'lucide-react';
import { Book, BookBulkOperation, BookExportOptions, BookImportData } from '@/types/book';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
 * Bulk Status Update Dialog
 */
interface BulkStatusUpdateProps {
  selectedBooks: Book[];
  onClose: () => void;
  onUpdate: (status: string, reason?: string) => void;
}

export function BulkStatusUpdateDialog({ 
  selectedBooks, 
  onClose, 
  onUpdate 
}: BulkStatusUpdateProps) {
  const [status, setStatus] = useState<string>('');
  const [reason, setReason] = useState('');
  const [isUpdating, setIsUpdating] = useState(false);

  const handleUpdate = async () => {
    if (!status) return;
    
    setIsUpdating(true);
    try {
      await onUpdate(status, reason);
      onClose();
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Update Book Status</DialogTitle>
          <DialogDescription>
            Update the status for {selectedBooks.length} selected book(s).
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div>
            <Label htmlFor="status">New Status</Label>
            <Select value={status} onValueChange={setStatus}>
              <SelectTrigger>
                <SelectValue placeholder="Select new status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="AVAILABLE">Available</SelectItem>
                <SelectItem value="BORROWED">Borrowed</SelectItem>
                <SelectItem value="MAINTENANCE">Maintenance</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Label htmlFor="reason">Reason (Optional)</Label>
            <Textarea
              id="reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Enter reason for status change"
              rows={3}
            />
          </div>

          <div className="bg-blue-50 p-3 rounded-lg">
            <h4 className="font-medium text-blue-900 mb-2">Selected Books:</h4>
            <div className="space-y-1 max-h-32 overflow-y-auto">
              {selectedBooks.map((book) => (
                <div key={book.id} className="text-sm text-blue-800">
                  {book.title} by {book.author}
                </div>
              ))}
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button 
            onClick={handleUpdate} 
            disabled={!status || isUpdating}
          >
            {isUpdating ? 'Updating...' : 'Update Status'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

/**
 * Bulk Tags Management Dialog
 */
interface BulkTagsProps {
  selectedBooks: Book[];
  onClose: () => void;
  onUpdate: (tags: string[], action: 'add' | 'remove' | 'replace') => void;
}

export function BulkTagsDialog({ 
  selectedBooks, 
  onClose, 
  onUpdate 
}: BulkTagsProps) {
  const [tagInput, setTagInput] = useState('');
  const [tags, setTags] = useState<string[]>([]);
  const [action, setAction] = useState<'add' | 'remove' | 'replace'>('add');
  const [isUpdating, setIsUpdating] = useState(false);

  const addTag = () => {
    if (tagInput.trim() && !tags.includes(tagInput.trim())) {
      setTags([...tags, tagInput.trim()]);
      setTagInput('');
    }
  };

  const removeTag = (index: number) => {
    setTags(tags.filter((_, i) => i !== index));
  };

  const handleUpdate = async () => {
    if (tags.length === 0) return;
    
    setIsUpdating(true);
    try {
      await onUpdate(tags, action);
      onClose();
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Manage Tags</DialogTitle>
          <DialogDescription>
            Manage tags for {selectedBooks.length} selected book(s).
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div>
            <Label htmlFor="action">Action</Label>
            <Select value={action} onValueChange={(value: any) => setAction(value)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="add">Add tags</SelectItem>
                <SelectItem value="remove">Remove tags</SelectItem>
                <SelectItem value="replace">Replace all tags</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Label htmlFor="tags">Tags</Label>
            <div className="space-y-2">
              <div className="flex space-x-2">
                <Input
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  placeholder="Add a tag"
                  onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
                />
                <Button type="button" onClick={addTag} variant="outline" size="sm">
                  Add
                </Button>
              </div>
              {tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {tags.map((tag, index) => (
                    <Badge key={index} variant="outline" className="flex items-center gap-1">
                      {tag}
                      <button
                        type="button"
                        onClick={() => removeTag(index)}
                        className="ml-1 text-gray-500 hover:text-red-500"
                      >
                        ×
                      </button>
                    </Badge>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button 
            onClick={handleUpdate} 
            disabled={tags.length === 0 || isUpdating}
          >
            {isUpdating ? 'Updating...' : `${action === 'add' ? 'Add' : action === 'remove' ? 'Remove' : 'Replace'} Tags`}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

/**
 * Book Export Dialog
 */
interface BookExportProps {
  onClose: () => void;
  onExport: (options: BookExportOptions) => void;
  totalBooks: number;
}

export function BookExportDialog({ 
  onClose, 
  onExport, 
  totalBooks 
}: BookExportProps) {
  const [format, setFormat] = useState<'csv' | 'json' | 'excel' | 'pdf'>('csv');
  const [fields, setFields] = useState<string[]>([
    'title', 'author', 'isbn', 'status', 'price'
  ]);
  const [includeStats, setIncludeStats] = useState(false);
  const [isExporting, setIsExporting] = useState(false);

  const availableFields = [
    { id: 'title', label: 'Title' },
    { id: 'author', label: 'Author' },
    { id: 'isbn', label: 'ISBN' },
    { id: 'description', label: 'Description' },
    { id: 'status', label: 'Status' },
    { id: 'price', label: 'Price' },
    { id: 'publishedAt', label: 'Published Date' },
    { id: 'tags', label: 'Tags' },
    { id: 'createdAt', label: 'Created Date' },
    { id: 'updatedAt', label: 'Updated Date' },
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
          <DialogTitle>Export Books</DialogTitle>
          <DialogDescription>
            Export {totalBooks} book(s) to file.
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

/**
 * Book Import Dialog
 */
interface BookImportProps {
  onClose: () => void;
  onImport: (data: BookImportData) => void;
}

export function BookImportDialog({ onClose, onImport }: BookImportProps) {
  const [file, setFile] = useState<File | null>(null);
  const [format, setFormat] = useState<'csv' | 'json' | 'excel'>('csv');
  const [skipErrors, setSkipErrors] = useState(true);
  const [updateExisting, setUpdateExisting] = useState(false);
  const [isImporting, setIsImporting] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      // Auto-detect format from file extension
      const extension = selectedFile.name.split('.').pop()?.toLowerCase();
      if (extension === 'xlsx' || extension === 'xls') {
        setFormat('excel');
      } else if (extension === 'json') {
        setFormat('json');
      } else {
        setFormat('csv');
      }
    }
  };

  const handleImport = async () => {
    if (!file) return;
    
    setIsImporting(true);
    try {
      await onImport({
        file,
        format,
        skipErrors,
        updateExisting,
      });
      onClose();
    } finally {
      setIsImporting(false);
    }
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Import Books</DialogTitle>
          <DialogDescription>
            Import books from a file. Supported formats: CSV, JSON, Excel.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div>
            <Label htmlFor="file">Select File</Label>
            <Input
              id="file"
              type="file"
              accept=".csv,.json,.xlsx,.xls"
              onChange={handleFileChange}
            />
            {file && (
              <div className="mt-2 text-sm text-gray-600">
                Selected: {file.name} ({(file.size / 1024).toFixed(1)} KB)
              </div>
            )}
          </div>

          <div>
            <Label htmlFor="format">File Format</Label>
            <Select value={format} onValueChange={(value: any) => setFormat(value)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="csv">CSV</SelectItem>
                <SelectItem value="json">JSON</SelectItem>
                <SelectItem value="excel">Excel</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-3">
            <div className="flex items-center space-x-2">
              <Checkbox
                id="skipErrors"
                checked={skipErrors}
                onCheckedChange={(checked) => setSkipErrors(!!checked)}
              />
              <Label htmlFor="skipErrors" className="text-sm">
                Skip rows with errors
              </Label>
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox
                id="updateExisting"
                checked={updateExisting}
                onCheckedChange={(checked) => setUpdateExisting(!!checked)}
              />
              <Label htmlFor="updateExisting" className="text-sm">
                Update existing books (match by ISBN)
              </Label>
            </div>
          </div>

          <div className="bg-blue-50 p-3 rounded-lg">
            <div className="flex items-start space-x-2">
              <Info className="h-4 w-4 text-blue-500 mt-0.5" />
              <div className="text-sm text-blue-700">
                <p className="font-medium">Expected format:</p>
                <p>CSV: title,author,isbn,description,price,status</p>
                <p>JSON: Array of book objects with the same fields</p>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button 
            onClick={handleImport} 
            disabled={!file || isImporting}
          >
            <Upload className="w-4 h-4 mr-2" />
            {isImporting ? 'Importing...' : 'Import'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}