import React from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Globe, Users, User, X } from 'lucide-react';

export type PermissionScope = 'by_all' | 'by_my_role' | 'by_me' | 'none';

interface ScopeConfig {
  label: string;
  icon: React.ReactNode;
  color: string;
  description: string;
  value: PermissionScope;
}

const scopeConfig: Record<PermissionScope, ScopeConfig> = {
  by_all: {
    label: 'By All',
    icon: <Globe className="h-4 w-4" />,
    color: 'text-green-600 dark:text-green-400',
    description: 'Can access all resources',
    value: 'by_all'
  },
  by_my_role: {
    label: 'By My Role',
    icon: <Users className="h-4 w-4" />,
    color: 'text-blue-600 dark:text-blue-400',
    description: 'Can access resources created by same or lower role levels',
    value: 'by_my_role'
  },
  by_me: {
    label: 'By Me',
    icon: <User className="h-4 w-4" />,
    color: 'text-orange-600 dark:text-orange-400',
    description: 'Can only access own resources',
    value: 'by_me'
  },
  none: {
    label: 'None',
    icon: <X className="h-4 w-4" />,
    color: 'text-gray-500 dark:text-gray-400',
    description: 'No access',
    value: 'none'
  }
};

interface PermissionScopeSelectorProps {
  service: string;
  action: string;
  currentScope: PermissionScope;
  onChange: (scope: PermissionScope) => void;
  disabled?: boolean;
}

export function PermissionScopeSelector({
  service,
  action,
  currentScope = 'none',
  onChange,
  disabled = false
}: PermissionScopeSelectorProps) {
  const current = scopeConfig[currentScope];

  return (
    <Select value={currentScope} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger className="w-[140px] h-9">
        <SelectValue>
          <div className="flex items-center gap-2">
            <span className={current.color}>{current.icon}</span>
            <span className="text-sm">{current.label}</span>
          </div>
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {Object.entries(scopeConfig).map(([scope, config]) => (
          <SelectItem key={scope} value={scope}>
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-2">
                <span className={config.color}>{config.icon}</span>
                <span className="font-medium">{config.label}</span>
              </div>
              <span className="text-xs text-muted-foreground ml-6">
                {config.description}
              </span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

export function PermissionScopeBadge({ scope }: { scope: PermissionScope }) {
  const config = scopeConfig[scope];
  
  return (
    <Badge variant="outline" className="gap-1">
      <span className={config.color}>{config.icon}</span>
      <span>{config.label}</span>
    </Badge>
  );
}

export { scopeConfig };