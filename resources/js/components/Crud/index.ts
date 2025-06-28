// Main CRUD components
export { CrudPage } from './CrudPage';
export { CrudDataTable } from './CrudDataTable';
export { CrudDrawer } from './CrudDrawer';
export { CrudPagination } from './CrudPagination';

// Search and Filter components
export { SearchBar } from './SearchBar';
export { FilterPanel } from './FilterPanel';

// Status and display components
export { StatusBadge } from './StatusBadge';

// Re-export types for convenience
export type {
  CrudPageProps,
  CrudColumn,
  CrudAction,
  CrudFilter,
  CrudFormProps,
  CrudEditFormProps,
  CrudDetailViewProps,
  DataTableProps,
  SearchBarProps,
  FilterPanelProps,
  PaginationProps,
  DrawerProps,
  StatusBadgeProps,
  PaginatedResult,
  ListRequest,
  BaseModel,
} from '@/types/crud';