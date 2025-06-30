"use client"

import * as React from 'react';
// @ts-ignore
import { Head, router } from '@inertiajs/react';
import { 
  Plus, 
  Filter, 
  MoreVertical, 
  Eye, 
  Edit, 
  Trash2, 
  Search,
  Settings2,
  RefreshCw,
  Download,
  Upload,
  X,
  Command,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { CrudPageProps, CrudAction } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Separator } from '@/components/ui/separator';
import { CrudDataTable } from './CrudDataTable';
import { SearchBar } from './SearchBar';
import { FilterPanel } from './FilterPanel';
import { CrudPagination } from './CrudPagination';
import { CrudDrawer } from './CrudDrawer';
import { useDebounce } from '@/hooks/useDebounce';
import { useCrudSelection } from '@/hooks/useCrudSelection';
import { usePageSize } from '@/hooks/usePageSize';
import { usePermissions } from '@/contexts/PermissionsContext';
import { PermissionGate } from '@/components/Permissions/PermissionGate';
import { toast } from 'sonner';

export function CrudPage<T extends { id: number }>({
  data,
  filters,
  title,
  resourceName,
  columns,
  actions = [],
  customFilters = [],
  paginationConfig,
  createForm: CreateForm,
  editForm: EditForm,
  detailView: DetailView,
  onRefresh,
  onBulkAction,
  className,
  tableClassName,
  // Accept permission props
  canCreate: propCanCreate,
  canEdit: propCanEdit,
  canDelete: propCanDelete,
  canView: propCanView,
}: CrudPageProps<T>) {
  
  // Use provided permissions or auto-detect based on resourceName
  const { canPerformAction } = usePermissions();
  
  const canCreate = propCanCreate !== undefined ? propCanCreate : canPerformAction(resourceName, 'create');
  const canEdit = propCanEdit !== undefined ? propCanEdit : canPerformAction(resourceName, 'update');
  const canDelete = propCanDelete !== undefined ? propCanDelete : canPerformAction(resourceName, 'delete');
  const canView = propCanView !== undefined ? propCanView : canPerformAction(resourceName, 'read');
  
  const canExport = canPerformAction(resourceName, 'export');
  const canBulkUpdate = canPerformAction(resourceName, 'bulk_update');
  const canBulkDelete = canPerformAction(resourceName, 'bulk_delete');
  // State management
  const [selectedItem, setSelectedItem] = React.useState<T | null>(null);
  const [drawerState, setDrawerState] = React.useState<{
    isOpen: boolean;
    type: 'create' | 'edit' | 'view' | undefined;
  }>({ isOpen: false, type: undefined });
  const [isSaving, setIsSaving] = React.useState(false);
  
  // Refs to form components
  const createFormRef = React.useRef<any>(null);
  const editFormRef = React.useRef<any>(null);

  // Search and filters
  const [searchTerm, setSearchTerm] = React.useState(filters?.search || '');
  const [activeFilters, setActiveFilters] = React.useState(filters?.filters || {});
  const [showFilters, setShowFilters] = React.useState(false);
  const [isSearching, setIsSearching] = React.useState(false);
  const [isRefreshing, setIsRefreshing] = React.useState(false);

  // Debounce search to avoid excessive requests
  const debouncedSearchTerm = useDebounce(searchTerm, 300);

  // Selection management for bulk actions
  const {
    selectedIds,
    selectedItems,
    toggleSelection,
    toggleAllSelection,
    clearSelection,
    setSelection,
  } = useCrudSelection(data.data);

  // Page size management with localStorage persistence
  const { pageSize, setPageSize, allowedSizes } = usePageSize(paginationConfig);

  // Generic error handler
  const handleError = React.useCallback((error: any, operation: string) => {
    console.error(`${operation} error:`, error);
    
    let errorMessage = 'Unknown error occurred';
    if (typeof error === 'string') {
      errorMessage = error;
    } else if (error?.message) {
      errorMessage = error.message;
    } else if (error?.errors?.validation_error) {
      errorMessage = error.errors.validation_error;
    }
    
    toast.error(`${operation} failed: ${errorMessage}`);
  }, []);

  // Re-enable search functionality
  React.useEffect(() => {
    if (debouncedSearchTerm !== (filters?.search || '')) {
      setIsSearching(true);
      
      router.get(`/admin/${resourceName}`, {
        ...(filters || {}),
        search: debouncedSearchTerm || undefined,
        page: 1,
        pageSize: pageSize,
      }, {
        preserveState: true,
        preserveScroll: true,
        only: ['data', 'filters'],
        onFinish: () => {
          setIsSearching(false);
        },
      });
    }
  }, [debouncedSearchTerm, resourceName, filters]);

  // Handlers
  const handleRefresh = React.useCallback(() => {
    setIsRefreshing(true);
    onRefresh?.();
    setTimeout(() => setIsRefreshing(false), 1000);
  }, [onRefresh]);

  const handleSort = React.useCallback((field: string) => {
    const newDirection = 
      filters?.sort === field && filters?.direction === 'asc' ? 'desc' : 'asc';
    
    router.get(`/admin/${resourceName}`, {
      ...(filters || {}),
      sort: field,
      direction: newDirection,
      pageSize: pageSize,
    }, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [resourceName, filters]);

  const handlePageChange = React.useCallback((page: number) => {
    router.get(`/admin/${resourceName}`, {
      ...(filters || {}),
      page,
      pageSize: pageSize,
    }, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [resourceName, filters, pageSize]);

  const handlePageSizeChange = React.useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    router.get(`/admin/${resourceName}`, {
      ...(filters || {}),
      page: 1, // Reset to first page when changing page size
      pageSize: newPageSize,
    }, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [resourceName, filters, setPageSize]);

  const handleFilterChange = React.useCallback((filterKey: string, value: any) => {
    const newFilters = { ...activeFilters };
    if (value === '' || value === null || value === undefined) {
      delete newFilters[filterKey];
    } else {
      newFilters[filterKey] = value;
    }
    
    setActiveFilters(newFilters);
    
    router.get(`/admin/${resourceName}`, {
      ...(filters || {}),
      filters: Object.keys(newFilters).length > 0 ? newFilters : undefined,
      page: 1,
      pageSize: pageSize,
    }, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [resourceName, filters, activeFilters, pageSize]);

  const handleCreate = React.useCallback(() => {
    setSelectedItem(null);
    setDrawerState({ isOpen: true, type: 'create' });
  }, []);

  const handleEdit = React.useCallback((item: T) => {
    setSelectedItem(item);
    setDrawerState({ isOpen: true, type: 'edit' });
  }, []);

  const handleView = React.useCallback((item: T) => {
    setSelectedItem(item);
    setDrawerState({ isOpen: true, type: 'view' });
  }, []);

  const handleDelete = React.useCallback(async (item: T) => {
    const confirmMessage = `Are you sure you want to delete this ${resourceName.slice(0, -1)}?`;
    if (confirm(confirmMessage)) {
      try {
        // Get CSRF token from meta tag
        const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
        
        const response = await fetch(`/api/${resourceName}/${item.id}`, {
          method: 'DELETE',
          headers: {
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
            'X-Inertia': 'true',
            'X-Inertia-Version': '1.0.0',
            ...(csrfToken && { 'X-CSRF-TOKEN': csrfToken }),
          },
        });

        if (response.ok) {
          toast.success(`${resourceName.slice(0, -1)} deleted successfully`);
          // Refresh the page data
          router.reload({ only: ['data', 'filters', 'stats'] });
          if (selectedIds.includes(item.id)) {
            clearSelection();
          }
        } else {
          const errorData = await response.json().catch(() => ({}));
          console.error('Delete error:', errorData);
          toast.error(`Failed to delete ${resourceName.slice(0, -1)}: ${errorData.message || 'Unknown error'}`);
        }
      } catch (error) {
        console.error('Delete error:', error);
        toast.error(`Failed to delete ${resourceName.slice(0, -1)}: Network error`);
      }
    }
  }, [resourceName, selectedIds, clearSelection]);

  const handleBulkDelete = React.useCallback(() => {
    if (selectedIds.length === 0) return;
    
    const confirmMessage = `Are you sure you want to delete ${selectedIds.length} item(s)?`;
    if (confirm(confirmMessage)) {
      try {
        onBulkAction?.('delete', selectedIds);
        toast.success(`${selectedIds.length} item(s) deleted successfully`);
        clearSelection();
      } catch (error) {
        handleError(error, 'Bulk delete');
      }
    }
  }, [selectedIds, onBulkAction, clearSelection, handleError]);

  const closeDrawer = React.useCallback(() => {
    setDrawerState({ isOpen: false, type: undefined });
    setSelectedItem(null);
  }, []);

  const handleDrawerSuccess = React.useCallback((message?: string) => {
    closeDrawer();
    if (message) {
      toast.success(message);
    }
    // Refresh the page data
    router.reload({ only: ['data', 'filters', 'stats'] });
  }, [closeDrawer]);

  const handleDrawerError = React.useCallback((errors: any) => {
    console.error('Drawer operation error:', errors);
    
    // Extract and show specific error message
    let errorMessage = 'Operation failed';
    console.log('Error object structure:', JSON.stringify(errors, null, 2));
    
    if (typeof errors === 'object' && errors !== null) {
      // Check for validation error first (higher priority)
      if ('errors' in errors && errors.errors && errors.errors.validation_error) {
        console.log('Found errors.errors:', errors.errors);
        console.log('Found validation_error, processing...');
        let validationError = errors.errors.validation_error;
        console.log('Parsing validation error:', validationError);
        
        if (typeof validationError === 'string' && validationError.includes('map[')) {
          // Extract field names and their error messages from the Go map format
          const fieldErrors = [];
          
          // Updated regex to handle the actual format with spaces in messages
          const fieldMatches = validationError.match(/(\w+):map\[(\w+):(.*?)\]/g);
          console.log('Field matches:', fieldMatches);
          
          if (fieldMatches) {
            for (const match of fieldMatches) {
              console.log('Processing match:', match);
              const fieldMatch = match.match(/(\w+):map\[(\w+):(.*?)\]/);
              
              if (fieldMatch) {
                const [, fieldName, errorType, errorMsg] = fieldMatch;
                console.log('Extracted:', { fieldName, errorType, errorMsg });
                fieldErrors.push(`${fieldName}: ${errorMsg}`);
              }
            }
            
            console.log('Final field errors:', fieldErrors);
            if (fieldErrors.length > 0) {
              errorMessage = `Validation failed:\n• ${fieldErrors.join('\n• ')}`;
            } else {
              // Fallback: just clean up the raw message a bit
              errorMessage = validationError.replace(/^validation errors: /, '').replace(/map\[|\]/g, '');
            }
          } else {
            // Fallback: just clean up the raw message a bit
            errorMessage = validationError.replace(/^validation errors: /, '').replace(/map\[|\]/g, '');
          }
        } else {
          errorMessage = validationError;
        }
      } else if (typeof errors.errors === 'string') {
        errorMessage = errors.errors;
      } else if ('message' in errors && typeof errors.message === 'string') {
        errorMessage = errors.message;
      }
    }
    
    // Create toast with ID so we can dismiss it specifically
    const toastId = toast.error(errorMessage, {
      duration: Infinity, // Make it persistent
      style: {
        whiteSpace: 'pre-line', // Preserve line breaks
      },
      action: {
        label: 'Close',
        onClick: () => {
          console.log('Toast close button clicked');
          toast.dismiss(toastId);
        },
      },
    });
    
    console.log('Created toast with ID:', toastId);
  }, []);

  // Build final actions - always include defaults + any additional custom actions
  const finalActions: CrudAction<T>[] = React.useMemo(() => {
    const defaultActions: CrudAction<T>[] = [];
    
    // Always include default actions if permissions allow
    if (canView && DetailView) {
      defaultActions.push({
        key: 'view',
        label: 'View',
        icon: <Eye className="w-4 h-4" />,
        onClick: handleView,
      });
    }
    
    if (canEdit && EditForm) {
      defaultActions.push({
        key: 'edit',
        label: 'Edit',
        icon: <Edit className="w-4 h-4" />,
        onClick: handleEdit,
      });
    }
    
    if (canDelete) {
      defaultActions.push({
        key: 'delete',
        label: 'Delete',
        icon: <Trash2 className="w-4 h-4" />,
        onClick: handleDelete,
        className: 'text-destructive focus:text-destructive',
        confirm: true,
        confirmMessage: `Are you sure you want to delete this ${resourceName.slice(0, -1)}?`,
      });
    }

    // Add any additional custom actions
    return [...defaultActions, ...(actions || [])];
  }, [canView, canEdit, canDelete, DetailView, EditForm, handleView, handleEdit, handleDelete, resourceName, actions]);

  const activeFilterCount = Object.keys(activeFilters).filter(key => 
    activeFilters[key] !== undefined && activeFilters[key] !== '' && activeFilters[key] !== null
  ).length;

  // Keyboard shortcuts
  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Cmd/Ctrl + N to create new
      if ((e.metaKey || e.ctrlKey) && e.key === 'n' && canCreate && CreateForm) {
        e.preventDefault();
        handleCreate();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [canCreate, CreateForm, handleCreate]);

  return (
    <>
      <Head title={title} />
      
      <div className={cn('flex flex-col gap-6 p-4 lg:p-6 min-w-0 overflow-hidden w-full', className)}>
        {/* Header */}
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between min-w-0">
          <div className="space-y-1 min-w-0">
            <h1 className="text-2xl font-semibold tracking-tight">
              {title}
            </h1>
            {data.total > 0 && (
              <p className="text-sm text-muted-foreground">
                {data.total} total {data.total === 1 ? 'item' : 'items'}
              </p>
            )}
          </div>
          
          <div className="flex items-center gap-2 flex-shrink-0">
            {/* Refresh Button */}
            <Button
              variant="outline"
              size="sm"
              onClick={handleRefresh}
              disabled={isRefreshing}
            >
              <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
              <span className="hidden lg:inline ml-2 whitespace-nowrap">Refresh</span>
            </Button>

            {/* Filters */}
            {customFilters.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowFilters(!showFilters)}
                className="relative"
              >
                <Filter className="h-4 w-4" />
                <span className="hidden lg:inline ml-2 whitespace-nowrap">Filters</span>
                {activeFilterCount > 0 && (
                  <Badge 
                    variant="secondary" 
                    className="absolute -top-2 -right-2 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs"
                  >
                    {activeFilterCount}
                  </Badge>
                )}
              </Button>
            )}

            {/* Export/Import Actions */}
            {(canExport || canBulkUpdate) && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="sm">
                    <Settings2 className="h-4 w-4" />
                    <span className="hidden lg:inline ml-2 whitespace-nowrap">Actions</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuLabel>Actions</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {canExport && (
                    <DropdownMenuItem>
                      <Download className="mr-2 h-4 w-4" />
                      Export Data
                    </DropdownMenuItem>
                  )}
                  {canBulkUpdate && (
                    <DropdownMenuItem>
                      <Upload className="mr-2 h-4 w-4" />
                      Import Data
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}

            {/* Create Button */}
            {canCreate && CreateForm && (
              <Button onClick={handleCreate} className="group">
                <Plus className="h-4 w-4" />
                <span className="hidden lg:inline ml-2 whitespace-nowrap">
                  Add {resourceName.slice(0, -1)}
                </span>
                <kbd className="ml-2 pointer-events-none hidden lg:inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 group-hover:bg-background">
                  <Command className="h-3 w-3" />N
                </kbd>
              </Button>
            )}
          </div>
        </div>

        {/* Search and Filters */}
        <div className="flex flex-col gap-4 min-w-0">
          {/* Search */}
          <div className="flex items-center gap-2 min-w-0">
            <div className="relative flex-1 max-w-sm min-w-0">
              {isSearching ? (
                <div className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                </div>
              ) : (
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              )}
              <Input
                placeholder={`Search ${resourceName}...`}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          {/* Filter Panel */}
          {showFilters && customFilters.length > 0 && (
            <div className="rounded-lg border bg-card p-4 min-w-0 overflow-hidden">
              <div className="flex items-center justify-between mb-4 min-w-0">
                <h3 className="text-sm font-medium">Filters</h3>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setShowFilters(false)}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
              <FilterPanel
                filters={customFilters}
                values={activeFilters}
                onChange={handleFilterChange}
                onClear={() => {
                  setActiveFilters({});
                  router.get(`/admin/${resourceName}`, {
                    ...(filters || {}),
                    filters: undefined,
                    page: 1,
                    pageSize: pageSize,
                  }, {
                    preserveState: true,
                    preserveScroll: true,
                    only: ['data', 'filters'],
                  });
                }}
              />
            </div>
          )}
        </div>

        {/* Bulk Actions */}
        {selectedIds.length > 0 && onBulkAction && (
          <div className="flex items-center justify-between rounded-lg border bg-muted/50 px-4 py-3 min-w-0">
            <div className="flex items-center gap-2 min-w-0">
              <Badge variant="secondary">
                {selectedIds.length} selected
              </Badge>
              <span className="text-sm text-muted-foreground">
                {selectedIds.length} of {data.data.length} row(s) selected
              </span>
            </div>
            <div className="flex items-center gap-2 flex-shrink-0">
              {canBulkDelete && (
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={handleBulkDelete}
                >
                  <Trash2 className="h-4 w-4 mr-1" />
                  Delete
                </Button>
              )}
              <Button
                variant="outline"
                size="sm"
                onClick={clearSelection}
              >
                Clear
              </Button>
            </div>
          </div>
        )}

        {/* Data Table */}
        <div className="min-w-0 overflow-hidden">
          <CrudDataTable
            data={data.data}
            columns={columns}
            actions={finalActions}
            sortField={filters?.sort}
            sortDirection={filters?.direction}
            onSort={handleSort}
            selectedIds={selectedIds}
            onSelectionChange={(ids) => {
              // Handle clear all case
              if (ids.length === 0) {
                clearSelection();
              } else if (ids.length === data.data.length) {
                toggleAllSelection();
              } else {
                // Set the selection directly to match the table state
                setSelection(ids);
              }
            }}
            enableSelection={!!onBulkAction}
            enableSearch={false} // We handle search externally
            enableColumnToggle={true}
            enablePagination={false} // We handle pagination externally
            className={tableClassName}
          />
        </div>

        {/* Pagination */}
        {data.lastPage > 1 && (
          <CrudPagination
            currentPage={data.currentPage}
            lastPage={data.lastPage}
            total={data.total}
            perPage={data.perPage}
            onPageChange={handlePageChange}
            onPageSizeChange={handlePageSizeChange}
            allowedPageSizes={allowedSizes}
          />
        )}
      </div>

      {/* Drawers for Create/Edit/View */}
      <CrudDrawer
        isOpen={drawerState.isOpen}
        onClose={closeDrawer}
        title={title}
        type={drawerState.type}
        resourceName={resourceName}
        size={drawerState.type === 'view' ? 'lg' : 'lg'}
        canEdit={drawerState.type === 'view' ? canEdit : false}
        canSave={drawerState.type === 'create' ? canCreate : drawerState.type === 'edit' ? canEdit : false}
        isSaving={isSaving}
        onEdit={drawerState.type === 'view' && canEdit ? () => {
          setDrawerState({ isOpen: true, type: 'edit' });
        } : undefined}
        onSave={
          drawerState.type === 'create' ? 
            () => createFormRef.current?.handleSubmit?.() :
          drawerState.type === 'edit' ?
            () => editFormRef.current?.handleSubmit?.() :
          undefined
        }
      >
        {drawerState.type === 'create' && CreateForm && (
          <CreateForm
            ref={createFormRef}
            onSuccess={handleDrawerSuccess}
            onError={handleDrawerError}
            onCancel={closeDrawer}
            setIsSaving={setIsSaving}
          />
        )}
        
        {drawerState.type === 'edit' && EditForm && selectedItem && (
          <EditForm
            ref={editFormRef}
            item={selectedItem}
            onSuccess={handleDrawerSuccess}
            onError={handleDrawerError}
            onCancel={closeDrawer}
            setIsSaving={setIsSaving}
          />
        )}
        
        {drawerState.type === 'view' && DetailView && selectedItem && (
          <DetailView
            item={selectedItem}
            onEdit={canEdit ? () => {
              setDrawerState({ isOpen: true, type: 'edit' });
            } : undefined}
            onClose={closeDrawer}
            canEdit={canEdit}
          />
        )}
      </CrudDrawer>
    </>
  );
}