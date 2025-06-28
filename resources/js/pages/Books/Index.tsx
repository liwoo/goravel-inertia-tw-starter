import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import { Download, Upload, FileText, BarChart3, BookOpen, Users } from 'lucide-react';
import { 
  Book, 
  BookListResponse, 
  BookListRequest, 
  BookStats,
  BookBulkOperation,
  BookExportOptions,
  BookImportData 
} from '@/types/book';
import { CrudPage } from '@/components/Crud/CrudPage';
import { bookColumns, bookColumnsMobile, bookFilters, bookQuickFilters } from '@/components/Books/BookColumns';
import { BookCreateForm, BookEditForm, BookDetailView } from '@/components/Books/BookForms';
import { 
  BulkStatusUpdateDialog, 
  BulkTagsDialog, 
  BookExportDialog, 
  BookImportDialog 
} from '@/components/Books/BookActions';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Books Index page
interface BooksIndexProps {
  data: BookListResponse;
  filters: BookListRequest;
  stats?: BookStats;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canBorrow: boolean;
    canManageLibrary: boolean;
    canViewReports: boolean;
  };
  meta?: {
    pagination: {
      defaultPageSize: number;
      maxPageSize: number;
      allowedSizes: number[];
    };
  };
}

export default function BooksIndex({ 
  data, 
  filters, 
  stats,
  permissions,
  meta
}: BooksIndexProps) {
  const isMobile = useIsMobile();
  
  // Debug logging
  console.log('BooksIndex - data:', data);
  console.log('BooksIndex - filters:', filters);
  console.log('BooksIndex - stats:', stats);
  console.log('BooksIndex - permissions:', permissions);
  
  // Dialog states
  const [showImportDialog, setShowImportDialog] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [showBulkStatusDialog, setShowBulkStatusDialog] = useState(false);
  const [showBulkTagsDialog, setShowBulkTagsDialog] = useState(false);
  const [selectedBooks, setSelectedBooks] = useState<Book[]>([]);

  // Handle bulk operations
  const handleBulkAction = async (action: string, selectedIds: number[]) => {
    if (selectedIds.length === 0) return;

    // Get selected book objects
    const selected = data.data.filter(book => selectedIds.includes(book.id));
    setSelectedBooks(selected);

    const operations: Record<string, () => void> = {
      delete: () => handleBulkDelete(selectedIds),
      updateStatus: () => setShowBulkStatusDialog(true),
      export: () => handleBulkExport(selectedIds),
      addTags: () => setShowBulkTagsDialog(true),
    };

    const operation = operations[action];
    if (operation) {
      operation();
    }
  };

  const handleBulkDelete = (bookIds: number[]) => {
    const confirmMessage = `Are you sure you want to delete ${bookIds.length} book(s)? This action cannot be undone.`;
    if (confirm(confirmMessage)) {
      router.delete('/api/books/bulk', {
        data: { bookIds },
        onSuccess: () => {
          // Refresh will be handled by the parent
        },
      });
    }
  };

  const handleBulkStatusUpdate = (bookIds: number[]) => {
    const status = prompt('Enter new status (AVAILABLE, BORROWED, MAINTENANCE):');
    if (status && ['AVAILABLE', 'BORROWED', 'MAINTENANCE'].includes(status)) {
      router.put('/api/books/bulk/status', {
        bookIds,
        status,
      });
    }
  };

  const handleBulkExport = (bookIds: number[]) => {
    const format = prompt('Export format (csv, json, excel):') || 'csv';
    const options: BookExportOptions = {
      format: format as any,
      filters: { ...filters, bookIds },
    };
    
    // Trigger download
    window.open(`/api/books/export?${new URLSearchParams(options as any).toString()}`);
  };

  const handleBulkAddTags = (bookIds: number[]) => {
    const tags = prompt('Enter tags to add (comma-separated):');
    if (tags) {
      const tagArray = tags.split(',').map(tag => tag.trim()).filter(Boolean);
      router.put('/api/books/bulk/tags', {
        bookIds,
        tags: tagArray,
        action: 'add',
      });
    }
  };

  // Handle import
  const handleImport = async (importData: BookImportData) => {
    const formData = new FormData();
    formData.append('file', importData.file);
    formData.append('format', importData.format);
    formData.append('skipErrors', importData.skipErrors ? 'true' : 'false');
    formData.append('updateExisting', importData.updateExisting ? 'true' : 'false');

    try {
      await router.post('/api/books/import', formData, {
        forceFormData: true,
        onSuccess: () => {
          // Handle success
        },
      });
    } catch (error) {
      console.error('Import failed:', error);
    }
  };

  // Handle export
  const handleExport = async (options: BookExportOptions) => {
    const params = new URLSearchParams({
      format: options.format,
      ...(options.fields && { fields: options.fields.join(',') }),
      ...(options.includeStats && { includeStats: 'true' }),
      ...Object.fromEntries(
        Object.entries(options.filters || {}).map(([key, value]) => [key, String(value)])
      ),
    });

    window.open(`/api/books/export?${params.toString()}`);
  };

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  return (
    <Admin title={"Books"}>
      <Head title="Books - Library Management" />
      
      <div className="space-y-6 min-w-0 overflow-hidden">
        {/* Statistics Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Books</CardTitle>
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalBooks}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.totalValue > 0 && `Worth ${new Intl.NumberFormat('en-US', {
                    style: 'currency',
                    currency: 'USD',
                  }).format(stats.totalValue)}`}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Available</CardTitle>
                <div className="h-4 w-4 bg-green-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.availableBooks}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.totalBooks > 0 && `${Math.round((stats.availableBooks / stats.totalBooks) * 100)}% available`}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Borrowed</CardTitle>
                <div className="h-4 w-4 bg-blue-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-blue-600">{stats.borrowedBooks}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.totalBooks > 0 && `${Math.round((stats.borrowedBooks / stats.totalBooks) * 100)}% borrowed`}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Maintenance</CardTitle>
                <div className="h-4 w-4 bg-orange-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-orange-600">{stats.maintenanceBooks}</div>
                <p className="text-xs text-muted-foreground">
                  Avg. price ${stats.averagePrice.toFixed(2)}
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Quick Filter Buttons */}
        <div className="flex flex-wrap gap-2">
          {bookQuickFilters.map((filter) => (
            <Button
              key={filter.key}
              variant={JSON.stringify(filters) === JSON.stringify(filter.filters) ? 'default' : 'outline'}
              size="sm"
              onClick={() => {
                router.get('/admin/books', filter.filters, {
                  preserveState: true,
                  preserveScroll: true,
                  only: ['data', 'filters', 'stats'],
                });
              }}
              className="flex items-center space-x-2"
            >
              {filter.icon}
              <span>{filter.label}</span>
            </Button>
          ))}
        </div>

        {/* Top Authors */}
        {stats?.topAuthors && stats.topAuthors.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Users className="h-5 w-5" />
                <span>Top Authors</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap gap-2">
                {stats.topAuthors.slice(0, 5).map((author, index) => (
                  <Badge
                    key={index}
                    variant="outline"
                    className="cursor-pointer hover:bg-blue-50"
                    onClick={() => {
                      router.get('/admin/books', { ...filters, author: author.name }, {
                        preserveState: true,
                        preserveScroll: true,
                        only: ['data', 'filters', 'stats'],
                      });
                    }}
                  >
                    {author.name} ({author.count})
                  </Badge>
                ))}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Management Actions */}
        {permissions.canManageLibrary && (
          <div className="flex flex-wrap gap-2">
            <Button 
              variant="outline" 
              size="sm"
              onClick={() => setShowImportDialog(true)}
            >
              <Upload className="h-4 w-4 mr-2" />
              Import Books
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={() => setShowExportDialog(true)}
            >
              <Download className="h-4 w-4 mr-2" />
              Export Books
            </Button>
            {permissions.canViewReports && (
              <Button variant="outline" size="sm">
                <BarChart3 className="h-4 w-4 mr-2" />
                View Reports
              </Button>
            )}
          </div>
        )}

        {/* Main CRUD Component */}
        <CrudPage<Book>
          data={data}
          filters={filters}
          title="My Books"
          resourceName="books"
          columns={isMobile ? bookColumnsMobile : bookColumns}
          customFilters={bookFilters}
          paginationConfig={meta?.pagination}
          createForm={BookCreateForm}
          editForm={BookEditForm}
          detailView={BookDetailView}
          onBulkAction={handleBulkAction}
          onRefresh={handleRefresh}
        />

        {/* Action Dialogs */}
        {showImportDialog && (
          <BookImportDialog
            onClose={() => setShowImportDialog(false)}
            onImport={handleImport}
          />
        )}

        {showExportDialog && (
          <BookExportDialog
            onClose={() => setShowExportDialog(false)}
            onExport={handleExport}
            totalBooks={data.total}
          />
        )}

        {showBulkStatusDialog && (
          <BulkStatusUpdateDialog
            selectedBooks={selectedBooks}
            onClose={() => setShowBulkStatusDialog(false)}
            onUpdate={async (status, reason) => {
              await router.put('/api/books/bulk/status', {
                bookIds: selectedBooks.map(b => b.id),
                status,
                reason,
              });
              handleRefresh();
            }}
          />
        )}

        {showBulkTagsDialog && (
          <BulkTagsDialog
            selectedBooks={selectedBooks}
            onClose={() => setShowBulkTagsDialog(false)}
            onUpdate={async (tags, action) => {
              await router.put('/api/books/bulk/tags', {
                bookIds: selectedBooks.map(b => b.id),
                tags,
                action,
              });
              handleRefresh();
            }}
          />
        )}
      </div>
    </Admin>
  );
}