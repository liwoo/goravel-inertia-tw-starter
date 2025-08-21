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
import { CrudPageProps, CrudAction, PageAction, SimpleFilter } from '@/types/crud';
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
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CrudDataTable } from './CrudDataTable';
import { SearchBar } from './SearchBar';
import { FilterPanel } from './FilterPanel';
import { DynamicFilterBuilder } from '@/components/Filters/DynamicFilterBuilder';
import { ActiveFilterBadges } from './ActiveFilterBadges';
import { CrudPagination } from './CrudPagination';
import { CrudDrawer } from './CrudDrawer';
import { useDebounce } from '@/hooks/useDebounce';
import { useCrudSelection } from '@/hooks/useCrudSelection';
import { usePageSize } from '@/hooks/usePageSize';
import { useFilterMetadata } from '@/hooks/useFilterMetadata';
import { usePermissions } from '@/contexts/PermissionsContext';
import { PermissionGate } from '@/components/Permissions/PermissionGate';
import { FilterCondition, CompoundFilter } from '@/types/filters';
import { toast } from 'sonner';

export function CrudPage<T extends { id: number }>({
  data,
  filters,
  title,
  resourceName,
  route,
  columns,
  actions = [],
  customFilters = [],
  pageActions = [],
  simpleFilters = [],
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
  // Guard against undefined data
  if (!data || !data.data) {
    console.error('CrudPage: data prop is undefined or missing data.data', { data, resourceName });
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">No data available</p>
      </div>
    );
  }
  
  // Use provided permissions or auto-detect based on resourceName
  const { canPerformAction } = usePermissions();
  
  // Use provided route or default to /admin/{resourceName}
  const baseRoute = route || `/admin/${resourceName}`;
  
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
  const searchInputRef = React.useRef<HTMLInputElement>(null);

  // Parse dynamic filters from URL and backend response
  const parsedDynamicFilter = React.useMemo(() => {
    console.log('🔍 Parsing dynamic filters, filters prop:', filters);
    console.log('🌐 Current URL:', window.location.href);
    
    // First, always check URL parameters for direct filter application
    const urlParams = new URLSearchParams(window.location.search);
    const filtersParam = urlParams.get('filters');
    console.log('📄 URL filters param (raw):', filtersParam);
    if (filtersParam) {
      try {
        const parsed = JSON.parse(filtersParam);
        console.log('✅ Successfully parsed URL filters param:', parsed);
        console.log('🎯 Returning parsed filter from URL');
        return parsed;
      } catch (e) {
        console.error('❌ Failed to parse filters from URL:', e);
      }
    } else {
      console.log('📭 No filters parameter found in URL');
    }
    
    // Fallback: Check if filters contains a 'filters' key with custom filters from backend response
    if (filters?.filters) {
      console.log('Found filters.filters:', filters.filters, 'type:', typeof filters.filters);
      
      // Check for __custom_filters key (backend stores parsed JSON here)
      if (filters.filters.__custom_filters) {
        console.log('Found __custom_filters:', filters.filters.__custom_filters);
        return filters.filters.__custom_filters;
      }
      
      if (typeof filters.filters === 'string') {
        try {
          const parsed = JSON.parse(filters.filters);
          console.log('Parsed filters.filters string:', parsed);
          return parsed;
        } catch (e) {
          console.error('Failed to parse filters:', e);
          return null;
        }
      } else if (typeof filters.filters === 'object') {
        // Check if this is a dynamic filter object (has field or logic properties)
        // vs a custom filters object (arbitrary key-value pairs)
        if (filters.filters.field || filters.filters.logic) {
          console.log('Found dynamic filter object:', filters.filters);
          return filters.filters;
        }
        // Otherwise it's custom filters, not dynamic filters
        console.log('Found custom filters object, not dynamic filter:', filters.filters);
        return null;
      }
    }
    
    console.log('No dynamic filters found');
    return null;
  }, [filters]);
  
  // Search and filters state
  const [searchTerm, setSearchTerm] = React.useState(filters?.search || '');
  const [activeFilters, setActiveFilters] = React.useState(() => {
    if (!filters?.filters) return {};
    // Exclude the __custom_filters key from activeFilters (that's for dynamic filters)
    const { __custom_filters, ...customFilters } = filters.filters;
    return customFilters;
  });
  const [showFilters, setShowFilters] = React.useState(false);
  const [isSearching, setIsSearching] = React.useState(false);
  const [isRefreshing, setIsRefreshing] = React.useState(false);
  const [appliedDynamicFilter, setAppliedDynamicFilter] = React.useState<any>(parsedDynamicFilter);
  
  // Update applied dynamic filter when URL changes
  React.useEffect(() => {
    setAppliedDynamicFilter(parsedDynamicFilter);
  }, [parsedDynamicFilter]);
  
  // Update active filters when filters prop changes (after navigation)
  React.useEffect(() => {
    if (!filters?.filters) {
      setActiveFilters({});
      return;
    }
    // Exclude the __custom_filters key from activeFilters (that's for dynamic filters)
    const { __custom_filters, ...customFilters } = filters.filters;
    setActiveFilters(customFilters);
  }, [filters]);
  
  // Determine active simple filter based on current filters
  const getActiveSimpleFilter = React.useCallback(() => {
    if (!filters || !simpleFilters) return undefined;
    
    // Check both direct filters and nested filters.filters
    const actualFilters = filters.filters || filters;
    
    console.log('getActiveSimpleFilter debug:', {
      filters,
      actualFilters,
      simpleFilters: simpleFilters.map(f => ({ 
        key: f.key, 
        value: f.value, 
        filterParams: f.filterParams 
      }))
    });
    
    // Check each simple filter to see if its conditions match the current URL state
    for (const simpleFilter of simpleFilters) {
      if (simpleFilter.filterParams) {
        // Check if all filterParams match current filters
        const allParamsMatch = Object.entries(simpleFilter.filterParams).every(([key, value]) => {
          const currentValue = actualFilters[key] !== undefined ? actualFilters[key] : filters[key];
          
          // Handle boolean comparisons (string "true"/"false" vs boolean true/false)
          if (typeof value === 'boolean' && typeof currentValue === 'string') {
            return currentValue === value.toString();
          }
          if (typeof value === 'string' && typeof currentValue === 'boolean') {
            return value === currentValue.toString();
          }
          
          // Handle numeric comparisons
          if (typeof value === 'number' && typeof currentValue === 'string') {
            return currentValue === value.toString();
          }
          if (typeof value === 'string' && typeof currentValue === 'number') {
            return value === currentValue.toString();
          }
          
          // Default strict comparison
          return currentValue === value;
        });
        
        if (allParamsMatch) {
          return simpleFilter.value.toString();
        }
      } else {
        // Default check: does the filter's key match its value in current filters?
        const filterValue = actualFilters[simpleFilter.key] || filters[simpleFilter.key];
        if (filterValue === simpleFilter.value) {
          return simpleFilter.value.toString();
        }
      }
    }
    
    return undefined;
  }, [filters, simpleFilters]);
  
  const [activeSimpleFilter, setActiveSimpleFilter] = React.useState<string | undefined>(
    getActiveSimpleFilter()
  );
  
  // Update active filter when URL filters change
  React.useEffect(() => {
    const activeFilter = getActiveSimpleFilter();
    console.log('Filter state update:', { 
      filters, 
      activeFilter,
      status: filters?.status 
    });
    setActiveSimpleFilter(activeFilter);
  }, [filters, getActiveSimpleFilter]);

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
  
  // Fetch filter metadata for dynamic filters
  const { metadata: filterMetadata, loading: filterMetadataLoading } = useFilterMetadata(resourceName);

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

  // Helper function to build navigation parameters that preserve both custom and dynamic filters
  const buildNavigationParams = React.useCallback((overrides: Record<string, any> = {}) => {
    const params: Record<string, any> = {
      page: 1,
      pageSize: pageSize,
      ...overrides, // Apply overrides first so they can be overridden below if needed
    };
    
    // Preserve core navigation parameters from current filters
    if (filters?.sort && !overrides.hasOwnProperty('sort')) params.sort = filters.sort;
    if (filters?.direction && !overrides.hasOwnProperty('direction')) params.direction = filters.direction;
    if (filters?.search && !overrides.hasOwnProperty('search')) params.search = filters.search;
    
    // Preserve dynamic filters if they exist
    if (appliedDynamicFilter && !overrides.hasOwnProperty('filters')) {
      params.filters = JSON.stringify(appliedDynamicFilter);
    }
    
    // Preserve custom filters (but not __custom_filters which is for dynamic filters)
    Object.keys(activeFilters).forEach(key => {
      if (key !== '__custom_filters' && activeFilters[key] !== undefined && activeFilters[key] !== '' && !overrides.hasOwnProperty(key)) {
        params[key] = activeFilters[key];
      }
    });
    
    console.log('🔧 Built navigation params:', params);
    return params;
  }, [filters, pageSize, appliedDynamicFilter, activeFilters]);

  // Re-enable search functionality
  React.useEffect(() => {
    if (debouncedSearchTerm !== (filters?.search || '')) {
      setIsSearching(true);
      
      const params = buildNavigationParams({
        search: debouncedSearchTerm || undefined,
        page: 1,
      });
      
      console.log('🔍 Search params:', params);
      
      router.get(baseRoute, params, {
        preserveState: true,
        preserveScroll: true,
        only: ['data', 'filters'],
        onFinish: () => {
          setIsSearching(false);
          // Refocus the search input after update with a small delay
          setTimeout(() => {
            if (searchInputRef.current) {
              searchInputRef.current.focus();
              // Restore cursor position to end
              const length = searchInputRef.current.value.length;
              searchInputRef.current.setSelectionRange(length, length);
            }
          }, 50);
        },
      });
    }
  }, [debouncedSearchTerm, buildNavigationParams, baseRoute, filters?.search]);

  // Handlers
  const handleRefresh = React.useCallback(() => {
    setIsRefreshing(true);
    onRefresh?.();
    setTimeout(() => setIsRefreshing(false), 1000);
  }, [onRefresh]);

  const handleSort = React.useCallback((field: string, direction?: 'asc' | 'desc') => {
    // If direction is explicitly provided, use it; otherwise toggle
    const newDirection = direction || 
      (filters?.sort === field && filters?.direction === 'asc' ? 'desc' : 'asc');
    
    const params = buildNavigationParams({
      page: 1,
      sort: field,
      direction: newDirection,
    });
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams, filters?.sort, filters?.direction]);

  const handlePageChange = React.useCallback((page: number) => {
    const params = buildNavigationParams({ page });
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams]);

  const handlePageSizeChange = React.useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    
    const params = buildNavigationParams({ 
      page: 1, // Reset to first page when changing page size
      pageSize: newPageSize 
    });
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams, setPageSize]);

  const handleFilterChange = React.useCallback((filterKey: string, value: any) => {
    const newFilters = { ...activeFilters };
    if (value === '' || value === null || value === undefined) {
      delete newFilters[filterKey];
    } else {
      newFilters[filterKey] = value;
    }
    
    setActiveFilters(newFilters);
    
    // Use helper but override custom filters with the new filters
    const params = buildNavigationParams({
      page: 1,
      ...newFilters, // Override any existing custom filters with the new ones
    });
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams, activeFilters]);

  const handleDynamicFilterApply = React.useCallback((filter: FilterCondition | CompoundFilter | null) => {
    // Store the applied filter to maintain state
    setAppliedDynamicFilter(filter);
    
    // Build parameters with the new filter
    const params = buildNavigationParams({
      page: 1,
      filters: filter ? JSON.stringify(filter) : undefined,
    });
    
    // Remove filters key if no filter (to clear dynamic filters)
    if (!filter && params.filters) {
      delete params.filters;
    }
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams]);

  const handleRemoveDynamicFilter = React.useCallback((index?: number) => {
    if (!appliedDynamicFilter) return;
    
    let newFilter: FilterCondition | CompoundFilter | null = null;
    
    if ('logic' in appliedDynamicFilter && index !== undefined) {
      // Remove specific condition from compound filter
      const newConditions = appliedDynamicFilter.conditions.filter((_, i) => i !== index);
      if (newConditions.length > 1) {
        newFilter = { ...appliedDynamicFilter, conditions: newConditions };
      } else if (newConditions.length === 1) {
        newFilter = newConditions[0] as FilterCondition;
      }
    }
    
    handleDynamicFilterApply(newFilter);
  }, [appliedDynamicFilter, handleDynamicFilterApply]);

  const handleClearAllFilters = React.useCallback(() => {
    setAppliedDynamicFilter(null);
    setActiveFilters({});
    
    // Build clean parameters with only core navigation (no filters)
    const params: Record<string, any> = {
      page: 1,
      pageSize: pageSize,
    };
    
    // Preserve only sort, direction, and search
    if (filters?.sort) params.sort = filters.sort;
    if (filters?.direction) params.direction = filters.direction;
    if (filters?.search) params.search = filters.search;
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, filters, pageSize]);

  const handleSimpleFilterChange = React.useCallback((filterValue: string | undefined) => {
    setActiveSimpleFilter(filterValue);
    
    // Start with base navigation params including dynamic filters
    const overrides: Record<string, any> = { page: 1 };
    
    if (filterValue === undefined) {
      // "All" was selected - clear simple filter-specific params but preserve everything else
      // Don't add any specific filter parameters
    } else {
      // Find the selected filter
      const selectedFilter = simpleFilters.find(f => f.value.toString() === filterValue);
      if (selectedFilter) {
        // Apply the filter's parameters
        if (selectedFilter.filterParams) {
          // Use the explicitly defined filter parameters
          Object.assign(overrides, selectedFilter.filterParams);
        } else {
          // Default behavior: use the filter's key and value
          overrides[selectedFilter.key] = selectedFilter.value;
        }
      }
    }
    
    const params = buildNavigationParams(overrides);
    
    router.get(baseRoute, params, {
      preserveState: true,
      preserveScroll: true,
      only: ['data', 'filters'],
    });
  }, [baseRoute, buildNavigationParams, simpleFilters]);

  const handleCreate = React.useCallback(() => {
    // Reset form refs before opening create drawer
    createFormRef.current = null;
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
          // Close drawer first if it's open
          if (drawerState.isOpen) {
            closeDrawer();
          }
          // Refresh the page data after a small delay to ensure drawer is closed
          setTimeout(() => {
            router.reload({ 
              only: ['data', 'filters', 'stats'],
              preserveState: false,
              preserveScroll: true 
            });
          }, 100);
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
    // Reset form refs to prevent stale form data
    createFormRef.current = null;
    editFormRef.current = null;
  }, []);

  const handleDrawerSuccess = React.useCallback((message?: string) => {
    closeDrawer();
    if (message) {
      toast.success(message);
    }
    // Refresh the page data
    router.reload({ 
      only: ['data', 'filters', 'stats'],
      preserveState: false,
      preserveScroll: true 
    });
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

  // Count active filters including both custom and dynamic filters
  const customFilterCount = Object.keys(activeFilters).filter(key => 
    activeFilters[key] !== undefined && activeFilters[key] !== '' && activeFilters[key] !== null && activeFilters[key] !== '__all__'
  ).length;
  
  // Count dynamic filter conditions
  const dynamicFilterCount = (() => {
    if (!appliedDynamicFilter) return 0;
    if ('logic' in appliedDynamicFilter) {
      // Compound filter - count all conditions
      return appliedDynamicFilter.conditions.filter(c => !('logic' in c)).length;
    }
    // Single condition
    return 1;
  })();
  
  const activeFilterCount = customFilterCount + dynamicFilterCount;
  
  // Debug logging for filter state
  React.useEffect(() => {
    console.log('🔍 Filter state debug:', {
      currentURL: window.location.href,
      parsedDynamicFilter,
      appliedDynamicFilter,
      activeFilters,
      customFilterCount,
      dynamicFilterCount,
      activeFilterCount,
      filterMetadata,
      filterMetadataLoading,
      shouldShowBadges: !!(appliedDynamicFilter || Object.keys(activeFilters).some(key => activeFilters[key] !== undefined && activeFilters[key] !== '' && activeFilters[key] !== '__all__'))
    });
  }, [parsedDynamicFilter, appliedDynamicFilter, activeFilters, customFilterCount, dynamicFilterCount, activeFilterCount, filterMetadata, filterMetadataLoading]);

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

            {/* Filters - show if we have custom filters OR filter metadata for dynamic filters */}
            {(customFilters.length > 0 || (filterMetadata && filterMetadata.filters.length > 0)) && (
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

            {/* Page Actions */}
            {pageActions.length > 0 && (
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
                  {pageActions.map((action) => (
                    <DropdownMenuItem
                      key={action.key}
                      onClick={action.handler}
                      className={action.className}
                    >
                      {action.icon && <span className="mr-2">{action.icon}</span>}
                      {action.label}
                    </DropdownMenuItem>
                  ))}
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
                <kbd className="ml-2 pointer-events-none hidden lg:inline-flex h-5 select-none items-center gap-1 rounded border px-1.5 font-mono text-[10px] font-medium opacity-100">
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
                ref={searchInputRef}
                placeholder={`Search ${resourceName}...`}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-9"
                autoFocus
              />
            </div>
          </div>

          {/* Active Filters Display */}
          {(appliedDynamicFilter || Object.keys(activeFilters).some(key => activeFilters[key] !== undefined && activeFilters[key] !== '' && activeFilters[key] !== '__all__')) && (
            <ActiveFilterBadges
              dynamicFilter={appliedDynamicFilter}
              customFilters={activeFilters}
              filterMetadata={filterMetadata}
              onRemoveDynamicFilter={handleRemoveDynamicFilter}
              onRemoveCustomFilter={(key) => handleFilterChange(key, '')}
              onClearAll={handleClearAllFilters}
            />
          )}

          {/* Simple Filters */}
          {simpleFilters.length > 0 && (
            <Tabs 
              value={activeSimpleFilter || "all"} 
              onValueChange={(value) => handleSimpleFilterChange(value === "all" ? undefined : value)}
              className="w-fit"
            >
              <TabsList className="**:data-[slot=badge]:bg-muted-foreground/30 **:data-[slot=badge]:size-5 **:data-[slot=badge]:rounded-full **:data-[slot=badge]:px-1 flex">
                <TabsTrigger value="all">All</TabsTrigger>
                {simpleFilters.slice(0, 5).map((filter) => (
                  <TabsTrigger key={filter.key} value={filter.value.toString()}>
                    {filter.icon && <span className="mr-2">{filter.icon}</span>}
                    {filter.label}
                    {filter.badge !== undefined && (
                      <Badge variant="secondary" className="ml-2">
                        {filter.badge}
                      </Badge>
                    )}
                  </TabsTrigger>
                ))}
              </TabsList>
            </Tabs>
          )}

          {/* Filter Panel - Use Dynamic Filter Builder if metadata available, otherwise fall back to custom filters */}
          {showFilters && (
            <>
              {filterMetadata && filterMetadata.filters.length > 0 ? (
                <DynamicFilterBuilder
                  metadata={filterMetadata}
                  loading={filterMetadataLoading}
                  onApply={handleDynamicFilterApply}
                  initialFilter={appliedDynamicFilter}
                />
              ) : customFilters.length > 0 ? (
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
                      
                      // Build a clean URL without any filter parameters
                      const cleanParams: Record<string, any> = {
                        page: 1,
                        pageSize: pageSize,
                      };
                      
                      // Only include non-filter parameters from current filters
                      if (filters?.sort) cleanParams.sort = filters.sort;
                      if (filters?.direction) cleanParams.direction = filters.direction;
                      if (filters?.search) cleanParams.search = filters.search;
                      
                      router.get(baseRoute, cleanParams, {
                        preserveState: true,
                        preserveScroll: true,
                        only: ['data', 'filters'],
                      });
                    }}
                  />
                </div>
              ) : null}
            </>
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
        key={`${drawerState.type}-${selectedItem?.id || 'new'}`}
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
        {drawerState.isOpen && drawerState.type === 'create' && CreateForm && (
          <CreateForm
            ref={createFormRef}
            onSuccess={handleDrawerSuccess}
            onError={handleDrawerError}
            onCancel={closeDrawer}
            setIsSaving={setIsSaving}
          />
        )}
        
        {drawerState.isOpen && drawerState.type === 'edit' && EditForm && selectedItem && (
          <EditForm
            ref={editFormRef}
            item={selectedItem}
            onSuccess={handleDrawerSuccess}
            onError={handleDrawerError}
            onCancel={closeDrawer}
            setIsSaving={setIsSaving}
          />
        )}
        
        {drawerState.isOpen && drawerState.type === 'view' && DetailView && selectedItem && (
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