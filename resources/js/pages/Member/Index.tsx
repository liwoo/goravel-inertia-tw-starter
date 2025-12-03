import React, { forwardRef } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  Member,
  MemberListResponse,
  MemberListRequest
} from '@/types/member';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  MemberCreateForm,
  MemberEditForm,
  MemberDetailView,
  memberColumns,
  memberColumnsMobile,
  memberFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { PortalCards } from './sections/PortalCards';

// Props interface for the Member Index page
interface MemberIndexProps {
  data: MemberListResponse;
  filters: MemberListRequest;
  permissions: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
  };
  meta?: {
    pagination: {
      defaultPageSize: number;
      maxPageSize: number;
      allowedSizes: number[];
    };
  };
  stats?: {
    [key: string]: any;
  };
  events?: any[];
  procurements?: any[];
  formalisation?: any;
  smeId?: number;
}

export default function MemberIndex({
  data,
  filters,
  permissions,
  meta,
  events,
  procurements,
  formalisation,
  smeId
}: MemberIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Wrap MemberCreateForm with forwardRef to pass smeId while maintaining ref forwarding
  const CreateFormWithSmeId = forwardRef((props: any, ref: any) => (
    <MemberCreateForm {...props} smeId={smeId} ref={ref} />
  ));

  return (
    <Admin title={"Member"}>
      <Head title="Member - Management" />

      <div className="flex flex-col gap-4 p-4 md:p-6">
        <PortalCards
          events={events}
          procurements={procurements}
          formalisation={formalisation}
        />

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Member>
            data={data}
            filters={filters}
            title="Members"
            resourceName="additional_business_members"
            columns={isMobile ? memberColumnsMobile : memberColumns}
            customFilters={memberFilters}
            paginationConfig={meta?.pagination}
            createForm={CreateFormWithSmeId}
            editForm={MemberEditForm}
            detailView={MemberDetailView}
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
