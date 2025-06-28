import React from 'react';
import { cn } from '@/lib/utils';
import { StatusBadgeProps, StatusType } from '@/types/crud';
import { Badge } from '@/components/ui/badge';

const statusVariantMap: Record<StatusType, string> = {
  active: 'bg-green-100 text-green-800 border-green-200',
  inactive: 'bg-red-100 text-red-800 border-red-200',
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  draft: 'bg-gray-100 text-gray-800 border-gray-200',
  published: 'bg-blue-100 text-blue-800 border-blue-200',
  archived: 'bg-purple-100 text-purple-800 border-purple-200',
};

export function StatusBadge({ 
  status, 
  variant, 
  className 
}: StatusBadgeProps) {
  const getStatusClass = () => {
    if (variant) {
      switch (variant) {
        case 'success':
          return 'bg-green-100 text-green-800 border-green-200';
        case 'warning':
          return 'bg-yellow-100 text-yellow-800 border-yellow-200';
        case 'danger':
          return 'bg-red-100 text-red-800 border-red-200';
        case 'info':
          return 'bg-blue-100 text-blue-800 border-blue-200';
        default:
          return 'bg-gray-100 text-gray-800 border-gray-200';
      }
    }

    // Auto-detect based on status value
    const normalizedStatus = status.toLowerCase() as StatusType;
    return statusVariantMap[normalizedStatus] || 'bg-gray-100 text-gray-800 border-gray-200';
  };

  const formatStatus = (status: string) => {
    return status
      .split(/[-_\s]/)
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  };

  return (
    <Badge
      variant="outline"
      className={cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border',
        getStatusClass(),
        className
      )}
    >
      {formatStatus(status)}
    </Badge>
  );
}