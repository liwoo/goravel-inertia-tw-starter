import React from 'react';
import { Head, Link } from '@inertiajs/react';
import { ArrowLeft, Edit, Users, Shield, Calendar, Badge as BadgeIcon } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import Admin from '@/layouts/Admin';

interface Permission {
  id: number;
  name: string;
  slug: string;
  description: string;
  category: string;
  action: string;
}

interface User {
  id: number;
  name: string;
  email: string;
  assigned_at: string;
  is_active: boolean;
}

interface Role {
  id: number;
  name: string;
  slug: string;
  description: string;
  level: number;
  is_active: boolean;
  users_count: number;
  created_at: string;
  updated_at: string;
}

interface Props {
  role: Role;
  users: User[];
  permissions: Permission[];
  auth: {
    user: any;
  };
}

export default function RoleShow({ role, users, permissions, auth }: Props) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const groupPermissionsByCategory = (permissions: Permission[]) => {
    const grouped: Record<string, Permission[]> = {};
    permissions.forEach((permission) => {
      if (!grouped[permission.category]) {
        grouped[permission.category] = [];
      }
      grouped[permission.category].push(permission);
    });
    return grouped;
  };

  const groupedPermissions = groupPermissionsByCategory(permissions);

  return (
    <Admin title={`Role: ${role.name}`}>
      <Head title={`Role: ${role.name} - Permissions`} />
      
      <div className="space-y-6 min-w-0 overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link
              href="/admin/permissions"
              className="flex items-center text-sm text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Roles
            </Link>
          </div>
          
          <div className="flex items-center space-x-2">
            <Link href={`/admin/permissions/roles/${role.id}/edit`}>
              <Button variant="outline" size="sm">
                <Edit className="mr-2 h-4 w-4" />
                Edit Role
              </Button>
            </Link>
          </div>
        </div>

        {/* Role Details Card */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <BadgeIcon className="h-8 w-8 text-primary" />
                <div>
                  <CardTitle className="text-2xl">{role.name}</CardTitle>
                  <p className="text-sm text-muted-foreground">
                    Level {role.level} • {role.slug}
                  </p>
                </div>
              </div>
              <Badge variant={role.is_active ? "default" : "secondary"}>
                {role.is_active ? "Active" : "Inactive"}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <h4 className="font-medium mb-2">Description</h4>
                <p className="text-sm text-muted-foreground">
                  {role.description || "No description provided"}
                </p>
              </div>
              
              <div>
                <h4 className="font-medium mb-2 flex items-center">
                  <Calendar className="mr-2 h-4 w-4" />
                  Created
                </h4>
                <p className="text-sm text-muted-foreground">
                  {formatDate(role.created_at)}
                </p>
              </div>
              
              <div>
                <h4 className="font-medium mb-2 flex items-center">
                  <Users className="mr-2 h-4 w-4" />
                  Users
                </h4>
                <p className="text-sm text-muted-foreground">
                  {role.users_count} user{role.users_count !== 1 ? 's' : ''} assigned
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Users with this Role */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Users className="mr-2 h-5 w-5" />
                Assigned Users ({users.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              {users.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-4">
                  No users assigned to this role
                </p>
              ) : (
                <div className="space-y-3">
                  {users.map((user) => (
                    <div
                      key={user.id}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div>
                        <p className="font-medium">{user.name}</p>
                        <p className="text-sm text-muted-foreground">{user.email}</p>
                        <p className="text-xs text-muted-foreground">
                          Assigned: {formatDate(user.assigned_at)}
                        </p>
                      </div>
                      <Badge variant={user.is_active ? "default" : "secondary"}>
                        {user.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Role Permissions */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Shield className="mr-2 h-5 w-5" />
                Permissions ({permissions.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              {permissions.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-4">
                  No permissions assigned to this role
                </p>
              ) : (
                <div className="space-y-4">
                  {Object.entries(groupedPermissions).map(([category, categoryPermissions]) => (
                    <div key={category}>
                      <h4 className="font-medium text-sm uppercase tracking-wide text-muted-foreground mb-2">
                        {category}
                      </h4>
                      <div className="space-y-2">
                        {categoryPermissions.map((permission) => (
                          <div
                            key={permission.id}
                            className="flex items-center justify-between p-2 bg-muted/30 rounded"
                          >
                            <div>
                              <p className="text-sm font-medium">{permission.name}</p>
                              <p className="text-xs text-muted-foreground">
                                {permission.description}
                              </p>
                            </div>
                            <Badge variant="outline" className="text-xs">
                              {permission.action}
                            </Badge>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Stats Summary */}
        <Card>
          <CardHeader>
            <CardTitle>Role Summary</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
              <div>
                <p className="text-2xl font-bold text-primary">{role.level}</p>
                <p className="text-sm text-muted-foreground">Authority Level</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-primary">{users.length}</p>
                <p className="text-sm text-muted-foreground">Active Users</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-primary">{permissions.length}</p>
                <p className="text-sm text-muted-foreground">Permissions</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-primary">
                  {Object.keys(groupedPermissions).length}
                </p>
                <p className="text-sm text-muted-foreground">Categories</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </Admin>
  );
}