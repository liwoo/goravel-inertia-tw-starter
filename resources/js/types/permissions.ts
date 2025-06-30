// TypeScript definitions for Permission Matrix system

export interface Permission {
  id: number;
  name: string;
  slug: string;
  description: string;
  category: string;
  resource: string;
  action: string;
  is_active: boolean;
  requires_ownership: boolean;
  can_delegate: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface Role {
  id: number;
  name: string;
  slug: string;
  description: string;
  level: number;
  is_active: boolean;
  parent_id?: number;
  parent?: Role;
  children?: Role[];
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface RoleWithPermissions extends Role {
  permission_ids: number[];
  permission_count: number;
  permissions?: Permission[];
}

export interface PermissionGrouped {
  category: string;
  permissions: Permission[];
}

export interface MatrixStats {
  total_roles: number;
  total_permissions: number;
  total_assignments: number;
  active_roles: number;
  active_permissions: number;
}

export interface PermissionMatrixData {
  roles: RoleWithPermissions[];
  permissions: PermissionGrouped[];
  matrix: Record<number, number[]>; // RoleID -> PermissionID[]
  stats: MatrixStats;
}

export interface BulkAssignmentRequest {
  role_id: number;
  permission_ids: number[];
  action: 'assign' | 'revoke';
}

export interface PermissionAssignmentRequest {
  role_id: number;
  permission_id: number;
}

export interface SyncPermissionsRequest {
  permission_ids: number[];
}

// API Response types
export interface PermissionResponse {
  success: boolean;
  message: string;
  data?: any;
  errors?: Record<string, string[]>;
}

export interface MatrixResponse extends PermissionResponse {
  data: PermissionMatrixData;
}

export interface AssignmentResponse extends PermissionResponse {
  data: {
    role_id: number;
    permission_id?: number;
    permission_count?: number;
    action: string;
  };
}

// Component props types
export interface PermissionMatrixProps {
  title: string;
  subtitle: string;
  matrixData: PermissionMatrixData;
  permissions: {
    canCreate: boolean;
    canView: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canManage: boolean;
  };
  user: any;
  breadcrumbs: Array<{
    label: string;
    href: string;
    active?: boolean;
  }>;
  stats: MatrixStats;
}

export interface PermissionGridProps {
  roles: RoleWithPermissions[];
  permissionGroups: PermissionGrouped[];
  matrix: Record<number, number[]>;
  onPermissionToggle: (roleId: number, permissionId: number, isAssigned: boolean) => void;
  onBulkToggle: (roleId: number, permissions: number[], action: 'assign' | 'revoke') => void;
  loading?: boolean;
  disabled?: boolean;
}

export interface RoleRowProps {
  role: RoleWithPermissions;
  permissions: Permission[];
  assignedPermissionIds: number[];
  onPermissionToggle: (permissionId: number, isAssigned: boolean) => void;
  onSelectAll: (permissions: number[], action: 'assign' | 'revoke') => void;
  loading?: boolean;
  disabled?: boolean;
}

export interface PermissionHeaderProps {
  permissionGroups: PermissionGrouped[];
  onSelectCategory: (category: string, permissions: number[], action: 'assign' | 'revoke') => void;
}

export interface MatrixStatsProps {
  stats: MatrixStats;
  loading?: boolean;
}

// Filter and search types
export interface PermissionFilters {
  search: string;
  category: string;
  roleLevel: string;
  showInactive: boolean;
}

export interface PermissionSort {
  field: 'name' | 'category' | 'action' | 'level' | 'created_at';
  direction: 'asc' | 'desc';
}

// UI State types
export interface PermissionMatrixState {
  data: PermissionMatrixData | null;
  loading: boolean;
  error: string | null;
  filters: PermissionFilters;
  sort: PermissionSort;
  selectedRoles: number[];
  selectedPermissions: number[];
  bulkAction: 'assign' | 'revoke' | null;
}

// Form validation types
export interface PermissionFormData {
  name: string;
  slug: string;
  description: string;
  category: string;
  action: string;
  resource?: string;
  requires_ownership: boolean;
  can_delegate: boolean;
  is_active: boolean;
}

export interface RoleFormData {
  name: string;
  slug: string;
  description: string;
  level: number;
  parent_id?: number;
  is_active: boolean;
}

// Validation schema types (for Zod)
export interface PermissionValidationRules {
  name: string;
  slug: string;
  description?: string;
  category: string;
  action: string;
  resource?: string;
}

export interface RoleValidationRules {
  name: string;
  slug: string;
  description?: string;
  level: number;
  parent_id?: number;
}

// Context types for permission matrix management
export interface PermissionMatrixContext {
  state: PermissionMatrixState;
  actions: {
    loadMatrix: () => Promise<void>;
    togglePermission: (roleId: number, permissionId: number) => Promise<void>;
    bulkAssign: (request: BulkAssignmentRequest) => Promise<void>;
    syncRolePermissions: (roleId: number, permissionIds: number[]) => Promise<void>;
    updateFilters: (filters: Partial<PermissionFilters>) => void;
    updateSort: (sort: PermissionSort) => void;
    clearSelection: () => void;
  };
}

// Standard permission categories and actions
export const PERMISSION_CATEGORIES = [
  'books',
  'users', 
  'roles',
  'system',
  'reports'
] as const;

export const PERMISSION_ACTIONS = [
  'create',
  'read', 
  'update',
  'delete',
  'manage',
  'assign',
  'borrow',
  'return',
  'export',
  'view',
  'configure',
  'monitor'
] as const;

export type PermissionCategory = typeof PERMISSION_CATEGORIES[number];
export type PermissionAction = typeof PERMISSION_ACTIONS[number];

// Service and Action types for permission matrix
export interface ServiceAction {
  service: string;
  actions: string[];
}