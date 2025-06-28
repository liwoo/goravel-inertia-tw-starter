import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search, Plus, Eye, Edit, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

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

interface RolesIndexProps {
  auth: any;
  roles: Role[];
  title: string;
  subtitle: string;
}

export default function RolesIndex({ auth, roles, title, subtitle }: RolesIndexProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(false);

  const filteredRoles = roles.filter(role => 
    role.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    role.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleAddNewRole = () => {
    router.visit('/admin/permissions/roles/create');
  };

  const handleViewRole = (roleId: number) => {
    router.visit(`/admin/permissions/roles/${roleId}`);
  };

  const handleEditRole = (roleId: number) => {
    router.visit(`/admin/permissions/roles/${roleId}/edit`);
  };

  const handleDeleteRole = async (roleId: number) => {
    if (!confirm('Are you sure you want to delete this role?')) return;

    setLoading(true);
    try {
      const response = await fetch(`/api/roles/${roleId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include'
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to delete role');
      }

      toast.success('Role deleted successfully');
      router.reload({ only: ['roles'] });
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete role');
    } finally {
      setLoading(false);
    }
  };

  const formatUserCount = (count: number) => {
    if (count === 0) return '0 user';
    if (count === 1) return '1 user';
    return `${count} users`;
  };

  return (
    <Admin title={title}>
      <Head title={title} />
      
      <div className="space-y-6 min-w-0 overflow-hidden">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center justify-between mb-2">
              <div>
                <h1 className="text-3xl font-bold">{title}</h1>
                <p className="text-muted-foreground mt-1">{subtitle}</p>
              </div>
              <Button 
                onClick={handleAddNewRole}
                className="bg-primary hover:bg-primary/90 text-primary-foreground"
              >
                <Plus className="w-4 h-4 mr-2" />
                Add new role
              </Button>
            </div>
          </div>

          {/* Search Bar */}
          <div className="mb-6">
            <div className="relative max-w-md">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-5 h-5" />
              <Input
                type="search"
                placeholder="Search roles..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 bg-background border-border text-foreground placeholder-muted-foreground focus:border-ring"
              />
            </div>
          </div>

          {/* Roles Table */}
          <div className="bg-card rounded-lg border overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground uppercase tracking-wider">
                    NAME
                  </th>
                  <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground uppercase tracking-wider">
                    DESCRIPTION
                  </th>
                  <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground uppercase tracking-wider">
                    USERS
                  </th>
                  <th className="text-right px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {filteredRoles.map((role) => (
                  <tr key={role.id} className="hover:bg-muted/50 transition-colors">
                    <td className="px-6 py-4">
                      <div className="font-medium text-foreground">{role.name}</div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-muted-foreground">{role.description}</div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-muted-foreground">{formatUserCount(role.users_count)}</div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end space-x-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleViewRole(role.id)}
                          className="text-muted-foreground hover:text-foreground hover:bg-muted"
                          disabled={loading}
                        >
                          <Eye className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleEditRole(role.id)}
                          className="text-muted-foreground hover:text-foreground hover:bg-muted"
                          disabled={loading}
                        >
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDeleteRole(role.id)}
                          className="text-muted-foreground hover:text-foreground hover:bg-muted"
                          disabled={loading}
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {filteredRoles.length === 0 && (
              <div className="text-center py-12 text-muted-foreground">
                {searchQuery ? 'No roles found matching your search.' : 'No roles available.'}
              </div>
            )}
          </div>

          {/* Add new role row at bottom */}
          <div className="mt-4">
            <Button
              variant="ghost"
              onClick={handleAddNewRole}
              className="text-primary hover:text-primary/90 hover:bg-muted"
            >
              <Plus className="w-4 h-4 mr-2" />
              Add new role
            </Button>
          </div>
        </div>
      </div>
    </Admin>
  );
}