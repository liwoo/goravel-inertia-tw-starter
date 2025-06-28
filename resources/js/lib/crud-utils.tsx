import { CrudColumn, CrudAction, StatusType } from '@/types/crud';

/**
 * Utility functions for CRUD operations
 */

/**
 * Creates a standard ID column for CRUD tables
 */
export function createIdColumn<T extends { id: number }>(): CrudColumn<T> {
  return {
    key: 'id',
    label: 'ID',
    sortable: true,
    className: 'w-20',
    render: (item) => `#${item.id}`,
  };
}

/**
 * Creates a standard name column with avatar support
 */
export function createNameColumn<T extends { name: string; email?: string; avatar?: string }>(
  options: {
    showEmail?: boolean;
    showAvatar?: boolean;
    avatarSize?: 'sm' | 'md' | 'lg';
  } = {}
): CrudColumn<T> {
  const { showEmail = true, showAvatar = false, avatarSize = 'sm' } = options;
  
  const sizeClasses = {
    sm: 'h-6 w-6',
    md: 'h-8 w-8',
    lg: 'h-10 w-10',
  };

  return {
    key: 'name',
    label: 'Name',
    sortable: true,
    render: (item) => (
      <div className="flex items-center space-x-3">
        {showAvatar && (
          <div className={`${sizeClasses[avatarSize]} rounded-full bg-gray-200 flex items-center justify-center overflow-hidden`}>
            {item.avatar ? (
              <img src={item.avatar} alt={item.name} className="w-full h-full object-cover" />
            ) : (
              <span className="text-xs font-medium text-gray-600">
                {item.name.charAt(0).toUpperCase()}
              </span>
            )}
          </div>
        )}
        <div>
          <div className="font-medium text-gray-900">{item.name}</div>
          {showEmail && item.email && (
            <div className="text-sm text-gray-500">{item.email}</div>
          )}
        </div>
      </div>
    ),
  };
}

/**
 * Creates a standard status column with badge styling
 */
export function createStatusColumn<T extends { status: string }>(
  statusMap?: Record<string, { label: string; variant: StatusType }>
): CrudColumn<T> {
  return {
    key: 'status',
    label: 'Status',
    sortable: true,
    render: (item) => {
      const statusInfo = statusMap?.[item.status];
      return (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusClasses(item.status)}`}>
          {statusInfo?.label || item.status}
        </span>
      );
    },
  };
}

/**
 * Creates a standard date column with formatting
 */
export function createDateColumn<T>(
  key: keyof T,
  label: string,
  options: {
    format?: 'date' | 'datetime' | 'relative';
    sortable?: boolean;
  } = {}
): CrudColumn<T> {
  const { format = 'date', sortable = true } = options;

  return {
    key: String(key),
    label,
    sortable,
    render: (item) => {
      const value = item[key];
      if (!value) return '-';
      
      const date = new Date(String(value));
      
      switch (format) {
        case 'datetime':
          return date.toLocaleString();
        case 'relative':
          return getRelativeTime(date);
        default:
          return date.toLocaleDateString();
      }
    },
  };
}

/**
 * Creates standard CRUD actions (view, edit, delete)
 */
export function createStandardActions<T extends { id: number }>(options: {
  onView?: (item: T) => void;
  onEdit?: (item: T) => void;
  onDelete?: (item: T) => void;
  canView?: boolean;
  canEdit?: boolean;
  canDelete?: boolean;
  resourceName?: string;
}): CrudAction<T>[] {
  const {
    onView,
    onEdit,
    onDelete,
    canView = true,
    canEdit = true,
    canDelete = true,
    resourceName = 'item',
  } = options;

  const actions: CrudAction<T>[] = [];

  if (canView && onView) {
    actions.push({
      key: 'view',
      label: 'View',
      onClick: onView,
    });
  }

  if (canEdit && onEdit) {
    actions.push({
      key: 'edit',
      label: 'Edit',
      onClick: onEdit,
    });
  }

  if (canDelete && onDelete) {
    actions.push({
      key: 'delete',
      label: 'Delete',
      onClick: onDelete,
      className: 'text-red-600 hover:text-red-800',
      confirm: true,
      confirmMessage: `Are you sure you want to delete this ${resourceName}?`,
    });
  }

  return actions;
}

/**
 * Gets CSS classes for status badges
 */
function getStatusClasses(status: string): string {
  const normalizedStatus = status.toLowerCase();
  
  const statusMap: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-red-100 text-red-800',
    pending: 'bg-yellow-100 text-yellow-800',
    draft: 'bg-gray-100 text-gray-800',
    published: 'bg-blue-100 text-blue-800',
    archived: 'bg-purple-100 text-purple-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    reviewing: 'bg-yellow-100 text-yellow-800',
  };

  return statusMap[normalizedStatus] || 'bg-gray-100 text-gray-800';
}

/**
 * Gets relative time string (e.g., "2 hours ago")
 */
function getRelativeTime(date: Date): string {
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return 'Just now';
  }

  const diffInMinutes = Math.floor(diffInSeconds / 60);
  if (diffInMinutes < 60) {
    return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`;
  }

  const diffInHours = Math.floor(diffInMinutes / 60);
  if (diffInHours < 24) {
    return `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`;
  }

  const diffInDays = Math.floor(diffInHours / 24);
  if (diffInDays < 30) {
    return `${diffInDays} day${diffInDays > 1 ? 's' : ''} ago`;
  }

  return date.toLocaleDateString();
}

/**
 * Formats file size in human readable format
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Truncates text to specified length with ellipsis
 */
export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength - 3) + '...';
}

/**
 * Creates a column for displaying truncated text with tooltip
 */
export function createTextColumn<T>(
  key: keyof T,
  label: string,
  options: {
    maxLength?: number;
    sortable?: boolean;
    showTooltip?: boolean;
  } = {}
): CrudColumn<T> {
  const { maxLength = 50, sortable = true, showTooltip = true } = options;

  return {
    key: String(key),
    label,
    sortable,
    render: (item) => {
      const value = item[key];
      if (!value) return '-';
      
      const text = String(value);
      const truncated = truncateText(text, maxLength);
      
      if (showTooltip && text.length > maxLength) {
        return (
          <span title={text} className="cursor-help">
            {truncated}
          </span>
        );
      }
      
      return truncated;
    },
  };
}

/**
 * Creates a numeric column with formatting
 */
export function createNumericColumn<T>(
  key: keyof T,
  label: string,
  options: {
    format?: 'integer' | 'decimal' | 'currency' | 'percentage';
    currency?: string;
    sortable?: boolean;
    className?: string;
  } = {}
): CrudColumn<T> {
  const { format = 'integer', currency = 'USD', sortable = true, className = 'text-right' } = options;

  return {
    key: String(key),
    label,
    sortable,
    className,
    render: (item) => {
      const value = item[key];
      if (value === null || value === undefined) return '-';
      
      const num = Number(value);
      if (isNaN(num)) return String(value);
      
      switch (format) {
        case 'currency':
          return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency,
          }).format(num);
        case 'percentage':
          return new Intl.NumberFormat('en-US', {
            style: 'percent',
            minimumFractionDigits: 0,
            maximumFractionDigits: 2,
          }).format(num);
        case 'decimal':
          return new Intl.NumberFormat('en-US', {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
          }).format(num);
        default:
          return new Intl.NumberFormat('en-US').format(num);
      }
    },
  };
}