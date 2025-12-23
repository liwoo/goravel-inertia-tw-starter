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
  FormalisationScoreWidget,
  EventsCalendarWidget,
  FormalisationData,
  CalendarEvent
} from '@/components/widgets';
import { FormalisationEditSheet } from '@/components/FormalisationEditSheet';
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
  calendarEvents?: CalendarEvent[];
  formalisation?: FormalisationData;
  district?: string;
  smeId?: number;
  userName?: string;
  smeName?: string;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  classification?: string;
  usmeNumber?: string;
}

export default function MemberIndex({
  data,
  filters,
  permissions,
  meta,
  calendarEvents = [],
  formalisation,
  district,
  smeId,
  userName,
  smeName,
  registrationNumber,
  taxIdentificationNumber,
  classification,
  usmeNumber
}: MemberIndexProps) {
  const isMobile = useIsMobile();
  const [primaryOwner, setPrimaryOwner] = useState<any>(null);
  const [additionalMembers, setAdditionalMembers] = useState<any[]>([]);
  const [employeeSummary, setEmployeeSummary] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [hasPendingAmendment, setHasPendingAmendment] = useState(false);

  useEffect(() => {
    if (smeId) {
      const fetchRelatedData = async () => {
        try {
          const [ownerRes, membersRes, summaryRes, pendingRes] = await Promise.all([
            axios.get(`/api/smes/${smeId}/primary_business_owner`).catch(() => ({ data: { data: null } })),
            axios.get(`/api/smes/${smeId}/additional_business_members`).catch(() => ({ data: { data: [] } })),
            axios.get(`/api/smes/${smeId}/business_employee_summary`).catch(() => ({ data: { data: null } })),
            axios.get(`/api/applications/amendment/pending/${smeId}`).catch(() => ({ data: { data: { has_pending_amendment: false } } })),
          ]);

          setPrimaryOwner(ownerRes.data.data);
          setAdditionalMembers(membersRes.data.data || []);
          setEmployeeSummary(summaryRes.data.data);
          setHasPendingAmendment(pendingRes.data.data?.has_pending_amendment ?? false);
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
            <div className="flex flex-col items-end gap-2 shrink-0">
              <div className="text-lg font-semibold text-right">{smeName}</div>
              <div className="flex items-center gap-2">
                {usmeNumber && (
                  <Badge variant="outline" className="text-xs font-mono">
                    {usmeNumber}
                  </Badge>
                )}
                {classification && (
                  <Badge
                    variant="default"
                    className={`text-sm px-3 py-1 font-semibold ${
                      classification === 'Micro'
                        ? 'bg-blue-600 hover:bg-blue-700'
                        : classification === 'Small'
                        ? 'bg-emerald-600 hover:bg-emerald-700'
                        : classification === 'Medium'
                        ? 'bg-amber-600 hover:bg-amber-700'
                        : ''
                    }`}
                  >
                    {classification} Enterprise
                  </Badge>
                )}
              </div>
            </div>
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
              canEdit={!!smeId}
              onEditClick={() => setIsEditModalOpen(true)}
            />

            {/* Main CRUD Component */}
            <div className="px-0">
              <CrudPage<Member>
                data={data}
                filters={filters}
                title="Team Members"
                resourceName="additional_business_members"
                displayName="Additional Member"
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

          {/* Right Sidebar - Calendar Widget */}
          <aside className="w-full xl:w-[380px] 2xl:w-[420px] shrink-0">
            <EventsCalendarWidget events={calendarEvents} district={district} smeId={smeId} />
          </aside>
        </div>
      </div>

      {/* Formalisation Edit Sheet */}
      {smeId && (
        <FormalisationEditSheet
          isOpen={isEditModalOpen}
          onOpenChange={setIsEditModalOpen}
          smeId={smeId}
          hasPendingAmendment={hasPendingAmendment}
          currentData={{
            registrationNumber: registrationNumber,
            taxIdentificationNumber: taxIdentificationNumber,
            hasBankAccount: formalisation?.has_bank_account ?? formalisation?.hasBankAccount,
            hasTaxClarification: formalisation?.has_tax_clarification ?? formalisation?.hasTaxClarification,
            isRegisteredForVat: formalisation?.is_registered_for_vat ?? formalisation?.isRegisteredForVat,
            isMemberOfAssociation: formalisation?.is_member_of_association ?? formalisation?.isMemberOfAssociation,
            isAffiliated: formalisation?.is_affiliated ?? formalisation?.isAffiliated,
            hasExportLicense: formalisation?.has_export_license ?? formalisation?.hasExportLicense,
            hasAccessedBds: formalisation?.has_accessed_bds ?? formalisation?.hasAccessedBds,
            annualTurnover: formalisation?.annual_turnover ?? formalisation?.annualTurnover,
            estimatedValueOfAssets: formalisation?.estimated_value_of_assets ?? formalisation?.estimatedValueOfAssets,
          }}
          onSuccess={() => {
            // Reload the page to fetch updated data
            router.reload();
          }}
        />
      )}
    </Admin>
  );
}
