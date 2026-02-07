import React from 'react';
// @ts-ignore
import { Head, router } from '@inertiajs/react';
import { useTranslation } from 'react-i18next';
import {
  Author,
  AuthorListResponse,
  AuthorListRequest,
  AuthorStats,
} from '@/types/author';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  AuthorCreateForm,
  AuthorEditForm,
  AuthorDetailView,
  getAuthorColumns,
  getAuthorColumnsMobile,
  getAuthorFilters,
  getAuthorStatsConfigs,
  getAuthorSimpleFilters,
} from './sections';
import {
  renderStatsCards,
  createSimpleFilters
} from '@/lib/crud-page-utils';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the Authors Index page
interface AuthorsIndexProps {
  data: AuthorListResponse;
  filters: AuthorListRequest;
  stats?: AuthorStats;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canManage: boolean;
  };
  meta?: {
    pagination: {
      defaultPageSize: number;
      maxPageSize: number;
      allowedSizes: number[];
    };
  };
}

export default function AuthorsIndex({
  data,
  filters,
  stats,
  permissions,
  meta
}: AuthorsIndexProps) {
  const isMobile = useIsMobile();
  const { t } = useTranslation('authors');

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  // Use extracted configurations
  const simpleFilters = createSimpleFilters(getAuthorSimpleFilters(t, stats));

  return (
    <Admin title={t('page.title')}>
      <Head title={t('page.headTitle')} />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, getAuthorStatsConfigs(t))}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Author>
            data={data}
            filters={filters}
            title={t('page.myAuthors')}
            resourceName="authors"
            columns={isMobile ? getAuthorColumnsMobile(t) : getAuthorColumns(t)}
            customFilters={getAuthorFilters(t)}
            simpleFilters={simpleFilters}
            paginationConfig={meta?.pagination}
            createForm={AuthorCreateForm}
            editForm={AuthorEditForm}
            detailView={AuthorDetailView}
            onRefresh={handleRefresh}
            canCreate={permissions.canCreate}
            canEdit={permissions.canEdit}
            canDelete={permissions.canDelete}
            canView={true}
          />
        </div>
      </div>
    </Admin>
  );
}
