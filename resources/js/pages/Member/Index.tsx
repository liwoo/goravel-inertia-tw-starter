import React, { forwardRef, useEffect, useState } from 'react';
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
import {
  UpcomingEventsWidget,
  UpcomingProcurementsWidget,
  FormalisationScoreWidget,
  UpcomingEvent,
  UpcomingProcurement,
  FormalisationData
} from '@/components/widgets';
import axios from 'axios';
import { Badge } from '@/components/ui/badge';

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
  events?: UpcomingEvent[];
  procurements?: UpcomingProcurement[];
  formalisation?: FormalisationData;
  smeId?: number;
  userName?: string;
  smeName?: string;
}

export default function MemberIndex({
  data,
  filters,
  permissions,
  meta,
  events = [],
  procurements = [],
  formalisation,
  smeId,
  userName,
  smeName
}: MemberIndexProps) {
  const isMobile = useIsMobile();
  const [primaryOwner, setPrimaryOwner] = useState<any>(null);
  const [additionalMembers, setAdditionalMembers] = useState<any[]>([]);
  const [employeeSummary, setEmployeeSummary] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (smeId) {
      const fetchRelatedData = async () => {
        try {
          const [ownerRes, membersRes, summaryRes] = await Promise.all([
            axios.get(`/api/smes/${smeId}/primary_business_owner`).catch(() => ({ data: { data: null } })),
            axios.get(`/api/smes/${smeId}/additional_business_members`).catch(() => ({ data: { data: [] } })),
            axios.get(`/api/smes/${smeId}/business_employee_summary`).catch(() => ({ data: { data: null } })),
          ]);

          setPrimaryOwner(ownerRes.data.data);
          setAdditionalMembers(membersRes.data.data || []);
          setEmployeeSummary(summaryRes.data.data);
        } catch (error) {
          console.error('Error fetching related data:', error);
        } finally {
          setLoading(false);
        }
      };

      fetchRelatedData();
    } else {
      setLoading(false);
    }
  }, [smeId]);

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  // Wrap MemberCreateForm with forwardRef to pass smeId while maintaining ref forwarding
  const CreateFormWithSmeId = forwardRef((props: any, ref: any) => (
    <MemberCreateForm {...props} smeId={smeId} ref={ref} />
  ));

  return (
    <Admin title={"SME Portal"}>
      <Head title="SME Portal" />

      <div className="p-4 md:p-6 flex flex-col gap-6">
        {/* Welcome Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">
              Welcome{userName ? `, ${userName.split(' ')[0]}` : ''}
            </h1>
            <p className="text-muted-foreground">
              Manage your team, track opportunities, and monitor your business progress.
            </p>
          </div>
          {smeName && (
            <Badge variant="secondary" className="text-lg px-4 py-2 font-medium shrink-0">
              {smeName}
            </Badge>
          )}
        </div>

        {/* Main Content + Sidebar Row */}
        <div className="flex flex-col xl:flex-row gap-6 items-start">
          {/* Main Content Area - Left Side */}
          <div className="flex-1 flex flex-col gap-6 min-w-0">
            {/* Formalisation Score Widget */}
            <FormalisationScoreWidget
              formalisation={formalisation}
              primaryOwner={primaryOwner}
              additionalMembers={additionalMembers}
              employeeSummary={employeeSummary}
              isLoading={loading}
            />

            {/* Main CRUD Component */}
            <div className="px-0">
              <CrudPage<Member>
                data={data}
                filters={filters}
                title="Team Members"
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

          {/* Right Sidebar - Widgets (top-aligned with content) */}
          <aside className="w-full xl:w-[380px] 2xl:w-[420px] flex flex-col gap-4 shrink-0">
            <UpcomingEventsWidget events={events} />
            <UpcomingProcurementsWidget procurements={procurements} />
          </aside>
        </div>
      </div>
    </Admin>
  );
}
