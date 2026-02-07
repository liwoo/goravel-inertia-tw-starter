import React, { useState } from 'react';
// @ts-ignore
import { Head, router } from '@inertiajs/react';
import { Download, Upload, FileText, BarChart3, BookOpen, Users } from 'lucide-react';
import { useTranslation } from 'react-i18next';
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
import {
  BookCreateForm,
  BookEditForm,
  BookDetailView,
  getBookColumns,
  getBookColumnsMobile,
  getBookFilters,
  getBookStatsConfigs,
  getBookSimpleFilters,
  getBookPageActions,
  getBookBulkActions
} from './sections';
import { 
  BulkStatusUpdateDialog, 
  BulkTagsDialog, 
  BookExportDialog, 
  BookImportDialog 
} from '@/components/Books/BookActions';
import { 
  renderStatsCards, 
  createPageActions, 
  createSimpleFilters 
} from '@/lib/crud-page-utils';
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
  const { t } = useTranslation('books');

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

    const bookBulkActions = getBookBulkActions(t);

    const operations: Record<string, () => void> = {
      delete: () => bookBulkActions.handleBulkDelete(selectedIds),
      updateStatus: () => setShowBulkStatusDialog(true),
      export: () => bookBulkActions.handleBulkExport(selectedIds, filters),
      addTags: () => setShowBulkTagsDialog(true),
    };

    const operation = operations[action];
    if (operation) {
      operation();
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

  // Use extracted configurations
  const simpleFilters = createSimpleFilters(getBookSimpleFilters(t, stats));

  const pageActions = createPageActions(
    getBookPageActions(t, permissions, {
      onImport: () => setShowImportDialog(true),
      onExport: () => setShowExportDialog(true),
      onReports: () => router.visit('/admin/books/reports'),
    })
  );

  return (
    <Admin title={t('page.title')}>
      <Head title={t('page.headTitle')} />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, getBookStatsConfigs(t))}


        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Book>
          data={data}
          filters={filters}
          title={t('page.myBooks')}
          resourceName="books"
          columns={isMobile ? getBookColumnsMobile(t) : getBookColumns(t)}
          customFilters={getBookFilters(t)}
          simpleFilters={simpleFilters}
          pageActions={pageActions}
          paginationConfig={meta?.pagination}
          createForm={BookCreateForm}
          editForm={BookEditForm}
          detailView={BookDetailView}
          onBulkAction={handleBulkAction}
          onRefresh={handleRefresh}
          // Pass the permissions from backend instead of auto-detecting
          canCreate={permissions.canCreate}
          canEdit={permissions.canEdit}
          canDelete={permissions.canDelete}
          canView={true} // If they're on this page, they can view
          />
        </div>

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