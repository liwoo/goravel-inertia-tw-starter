import React from 'react';
import { 
  Shield, 
  Users, 
  Calendar,
  Hash,
  CheckCircle,
  XCircle,
  AlertCircle
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Role } from '@/types/permissions';

export function RoleDetailView({ 
  item: role,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Role>) {
  const formatDate = (date: string | Date) => {
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="space-y-8">
      {/* Role Icon and Status */}
      <div className="flex justify-center">
        <div className={`p-6 rounded-full ${role.is_active ? 'bg-teal-100 dark:bg-teal-900/30' : 'bg-gray-100 dark:bg-gray-900/30'}`}>
          <Shield className={`h-16 w-16 ${role.is_active ? 'text-teal-600 dark:text-teal-400' : 'text-gray-600 dark:text-gray-400'}`} />
        </div>
      </div>

      {/* Role Information */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Role Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Hash className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Role ID</p>
                <p className="font-medium text-foreground">#{role.id.toString().padStart(4, '0')}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Shield className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Role Name</p>
                <p className="font-medium text-foreground">{role.name}</p>
                <p className="text-xs text-muted-foreground">Slug: {role.slug}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <AlertCircle className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Description</p>
                <p className="text-sm text-foreground">{role.description || 'No description provided'}</p>
              </div>
            </div>
          </div>
        </div>

        <Separator className="my-6" />

        {/* Role Settings */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Role Settings</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Users className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Hierarchy Level</p>
                <Badge variant="secondary" className="font-mono">
                  Level {role.level}
                </Badge>
                <p className="text-xs text-muted-foreground">
                  Higher levels have more authority
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <CheckCircle className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Status</p>
                <div className="mt-1">
                  {role.is_active ? (
                    <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1 w-fit">
                      <CheckCircle className="h-3 w-3" />
                      Active
                    </Badge>
                  ) : (
                    <Badge className="bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1 w-fit">
                      <XCircle className="h-3 w-3" />
                      Inactive
                    </Badge>
                  )}
                </div>
              </div>
            </div>

            {role.parent && (
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">Parent Role</p>
                  <p className="font-medium text-foreground">{role.parent.name}</p>
                </div>
              </div>
            )}
          </div>
        </div>

        <Separator className="my-6" />

        {/* Timestamps */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Timestamps</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Created At</p>
                <p className="font-medium text-foreground">{formatDate(role.created_at)}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Last Updated</p>
                <p className="font-medium text-foreground">{formatDate(role.updated_at)}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}