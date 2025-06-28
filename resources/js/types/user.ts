// TypeScript type definitions for User
export interface User {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  is_super_admin: boolean;
  created_at: string;
  updated_at: string;
  roles?: Role[];
}

export interface Role {
  id: number;
  name: string;
  slug: string;
  description: string;
}

export interface UserListResponse {
  data: User[];
  total: number;
  currentPage: number;
  lastPage: number;
  perPage: number;
  from: number;
  to: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface UserListRequest {
  page?: number;
  pageSize?: number;
  search?: string;
  sort?: string;
  direction?: 'ASC' | 'DESC';
  filters?: Record<string, any>;
}

export interface UserStats {
  totalUsers: number;
  activeUsers: number;
  inactiveUsers: number;
  superAdmins: number;
}

export interface UserFormData {
  name: string;
  email: string;
  password: string;
  is_active: boolean;
  is_super_admin: boolean;
  role_id?: number;
}

export interface UserBulkOperation {
  action: 'delete' | 'activate' | 'deactivate';
  ids: number[];
}

export interface UserExportOptions {
  format: 'csv' | 'json' | 'excel';
  fields?: string[];
  includeStats?: boolean;
  filters?: UserListRequest;
}

// Props interfaces for React components
export interface UserIndexProps {
  data: UserListResponse;
  filters: UserListRequest;
  stats?: UserStats;
  roles?: Role[];
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canManage: boolean;
    canExport: boolean;
    canViewReports: boolean;
  };
}

export interface UserFormProps {
  user?: User;
  roles?: Role[];
  onSubmit: (data: UserFormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export interface UserDetailProps {
  user: User;
  onEdit: () => void;
  onDelete: () => void;
  canEdit: boolean;
  canDelete: boolean;
}