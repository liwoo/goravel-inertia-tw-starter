// TypeScript types for the generic CRUD system
import React from 'react';

export interface PaginatedResult<T> {
  data: T[];
  total: number;
  perPage: number;
  currentPage: number;
  lastPage: number;
  from?: number;
  to?: number;
  hasNext?: boolean;
  hasPrev?: boolean;
}

export interface ListRequest {
  page?: number;
  perPage?: number;
  search?: string;
  sort?: string;
  direction?: 'asc' | 'desc';
  filters?: Record<string, any>;
}

export interface BaseModel {
  id: number;
  createdAt: string;
  updatedAt: string;
}

export interface CrudColumn<T = any> {
  key: string;
  label: string;
  sortable?: boolean;
  render?: (item: T) => React.ReactNode;
  className?: string;
  width?: string;
}

export interface CrudAction<T = any> {
  key: string;
  label: string;
  icon?: React.ReactNode;
  onClick: (item: T) => void;
  className?: string;
  confirm?: boolean;
  confirmMessage?: string;
  disabled?: (item: T) => boolean;
}

export interface CrudFilter {
  key: string;
  label: string;
  type: 'text' | 'select' | 'date' | 'number' | 'boolean';
  options?: { value: string | number; label: string }[];
  placeholder?: string;
  multiple?: boolean;
}

export interface CrudPageProps<T = any> {
  // Data
  data: PaginatedResult<T>;
  filters: ListRequest;
  
  // Configuration
  title: string;
  resourceName: string; // e.g., 'teams', 'players'
  columns: CrudColumn<T>[];
  actions?: CrudAction<T>[];
  customFilters?: CrudFilter[];
  
  // Pagination metadata (optional - will fallback to defaults if not provided)
  paginationConfig?: {
    defaultPageSize: number;
    maxPageSize: number;
    allowedSizes: number[];
  };
  
  // Permissions
  canCreate?: boolean;
  canEdit?: boolean;
  canDelete?: boolean;
  canView?: boolean;
  
  // Custom Components
  createForm?: React.ComponentType<CrudFormProps>;
  editForm?: React.ComponentType<CrudEditFormProps<T>>;
  detailView?: React.ComponentType<CrudDetailViewProps<T>>;
  
  // Callbacks
  onRefresh?: () => void;
  onBulkAction?: (action: string, selectedIds: number[]) => void;
  
  // Styling
  className?: string;
  tableClassName?: string;
}

export interface CrudFormProps {
  onSuccess: (message?: string) => void;
  onError?: (errors: any) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export interface CrudEditFormProps<T> {
  item: T;
  onSuccess: (message?: string) => void;
  onError?: (errors: any) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export interface CrudDetailViewProps<T> {
  item: T;
  onEdit?: () => void;
  onClose: () => void;
  canEdit?: boolean;
}

// Data table specific types
export interface DataTableProps<T> {
  data: T[];
  columns: CrudColumn<T>[];
  actions: CrudAction<T>[];
  sortField?: string;
  sortDirection?: 'asc' | 'desc';
  onSort: (field: string) => void;
  selectedIds: number[];
  onSelectionChange: (ids: number[]) => void;
  enableSelection?: boolean;
  loading?: boolean;
  emptyMessage?: string;
  className?: string;
}

// Search and filter types
export interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  onSearch: (value: string) => void;
  placeholder?: string;
  debounceMs?: number;
  loading?: boolean;
  className?: string;
}

export interface FilterPanelProps {
  filters: CrudFilter[];
  values: Record<string, any>;
  onChange: (key: string, value: any) => void;
  onClear: () => void;
  className?: string;
}

// Action dropdown types
export interface ActionDropdownProps<T> {
  actions: CrudAction<T>[];
  item: T;
  className?: string;
  buttonClassName?: string;
}

// Pagination types
export interface PaginationProps {
  currentPage: number;
  lastPage: number;
  total: number;
  perPage: number;
  onPageChange: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
  allowedPageSizes?: number[];
  showInfo?: boolean;
  className?: string;
}

// Drawer types
export interface DrawerProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
  children: React.ReactNode;
  className?: string;
  overlayClassName?: string;
}

// Form field types
export interface FormFieldProps {
  children: React.ReactNode;
  label?: string;
  error?: string;
  required?: boolean;
  className?: string;
  description?: string;
  hint?: string;
}

// Status types for common use cases
export type StatusType = 'active' | 'inactive' | 'pending' | 'draft' | 'published' | 'archived';

export interface StatusBadgeProps {
  status: StatusType | string;
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info';
  className?: string;
}

// Bulk action types
export interface BulkAction {
  key: string;
  label: string;
  icon?: React.ReactNode;
  variant?: 'default' | 'danger' | 'warning';
  confirm?: boolean;
  confirmMessage?: string;
}

// Extended CRUD page with bulk actions
export interface ExtendedCrudPageProps<T = any> extends CrudPageProps<T> {
  bulkActions?: BulkAction[];
  enableBulkSelect?: boolean;
  maxBulkSelect?: number;
}

// Hook return types
export interface UseCrudFiltersReturn {
  filters: Record<string, any>;
  updateFilter: (key: string, value: any) => void;
  clearFilters: () => void;
  hasActiveFilters: boolean;
  applyFilters: () => void;
  hasPendingChanges?: boolean;
  appliedFilters?: Record<string, any>;
}

export interface UseCrudSelectionReturn<T> {
  selectedIds: number[];
  selectedItems: T[];
  isAllSelected: boolean;
  isSomeSelected: boolean;
  toggleSelection: (id: number) => void;
  toggleAllSelection: () => void;
  clearSelection: () => void;
  setSelection: (ids: number[]) => void;
  selectRange?: (startIndex: number, endIndex: number) => void;
  canSelectMore?: boolean;
  remainingSelections?: number | null;
  selectionCount?: number;
}

// Loading and error states
export interface CrudState {
  loading: boolean;
  error: string | null;
  data: any[] | null;
  total: number;
  currentPage: number;
}

// Form validation types
export interface ValidationError {
  field: string;
  message: string;
}

export interface FormState {
  data: Record<string, any>;
  errors: ValidationError[];
  touched: Record<string, boolean>;
  isSubmitting: boolean;
  isValid: boolean;
}