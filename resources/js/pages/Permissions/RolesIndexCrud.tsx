import React from 'react';
import { Head } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { CrudPage } from '@/components/Crud/CrudPage';
import { Role } from '@/types/permissions';
import { 
  roleColumns, 
  roleColumnsMobile, 
  roleFilters, 
  roleQuickFilters,
  RoleCreateForm,
  RoleEditForm
} from './sections';

interface RolesIndexCrudProps {
  data: {
    data: Role[];
    total: number;
    perPage: number;
    currentPage: number;
    lastPage: number;
  };
  filters: {
    search?: string;
    sort?: string;
    direction?: 'asc' | 'desc';
    filters?: Record<string, any>;
  };
  title: string;
  subtitle: string;
}

export default function RolesIndexCrud({ 
  data, 
  filters, 
  title, 
  subtitle 
}: RolesIndexCrudProps) {
  
  const isMobile = false; // Could use useIsMobile hook

  return (
    <Admin title={title}>
      <Head title={title} />
      
      <CrudPage<Role>
        data={data}
        filters={filters}
        title={title}
        resourceName="roles"
        columns={isMobile ? roleColumnsMobile : roleColumns}
        customFilters={roleFilters}
        createForm={RoleCreateForm}
        editForm={RoleEditForm}
        onRefresh={() => {
          // Refresh handled by CrudPage
        }}
      />
    </Admin>
  );
}