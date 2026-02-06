import React from 'react';
import { Building2, FileText, User, ArrowRight, AlertCircle, FileEdit, UserPlus, Info } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Application, ApplicationType, parseAmendmentData, FormalisationAmendmentData } from '@/types/application';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { DetailRow } from '@/components/ui/details-row';
import { BooleanBadge } from '@/components/ui/boolean-badge';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

// Helper to get type display label
const getTypeLabel = (type?: ApplicationType): string => {
  switch (type) {
    case 'amend_formalisation':
      return 'Formalisation Amendment';
    case 'signup':
    default:
      return 'Sign Up';
  }
};

// Helper to get type badge style
const getTypeBadgeStyle = (type?: ApplicationType): string => {
  switch (type) {
    case 'amend_formalisation':
      return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
    case 'signup':
    default:
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
  }
};

// Amendment comparison row component
interface AmendmentComparisonRowProps {
  label: string;
  current: string | boolean | number | undefined | null;
  proposed: string | boolean | number | undefined | null;
}

function AmendmentComparisonRow({ label, current, proposed }: AmendmentComparisonRowProps) {
  const formatValue = (value: string | boolean | number | undefined | null): string => {
    if (value === undefined || value === null) return 'Not set';
    if (typeof value === 'boolean') return value ? 'Yes' : 'No';
    if (typeof value === 'number') return value.toLocaleString();
    return value;
  };

  const currentStr = formatValue(current);
  const proposedStr = formatValue(proposed);
  const hasChanged = currentStr !== proposedStr;

  return (
    <div className={`grid grid-cols-3 gap-4 py-2 px-3 rounded-md ${hasChanged ? 'bg-amber-50 dark:bg-amber-950/30' : ''}`}>
      <div className="font-medium text-sm">{label}</div>
      <div className={`text-sm ${hasChanged ? 'text-muted-foreground line-through' : ''}`}>
        {currentStr}
      </div>
      <div className={`text-sm ${hasChanged ? 'text-amber-700 dark:text-amber-400 font-medium' : ''}`}>
        {hasChanged ? proposedStr : '—'}
      </div>
    </div>
  );
}

export function ApplicationDetailView({
  item: application,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Application>) {
  const formatDate = (date: string | Date | null | undefined) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const applicationType = (application.type || 'signup') as ApplicationType;
  const isAmendment = applicationType === 'amend_formalisation';
  const amendmentData = isAmendment ? parseAmendmentData(application.data) : null;
  const TypeIcon = isAmendment ? FileEdit : UserPlus;

  return (
    <div className="space-y-6">
      {/* Rejection reason alert */}
      {application.status === 'Rejected' && application.rejection_reason && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Application Rejected</AlertTitle>
          <AlertDescription>{application.rejection_reason}</AlertDescription>
        </Alert>
      )}

      <Tabs defaultValue={isAmendment ? "amendment" : "basic"} className="w-full">
        <TabsList className={`grid w-full ${isAmendment ? 'grid-cols-2' : 'grid-cols-2'}`}>
          {isAmendment ? (
            <TabsTrigger value="amendment" title="Amendment Details">
              <FileEdit className="h-4 w-4 sm:mr-2" />
              <span className="hidden sm:inline">Amendment</span>
            </TabsTrigger>
          ) : (
            <TabsTrigger value="basic" title="Basic Information">
              <Building2 className="h-4 w-4 sm:mr-2" />
              <span className="hidden sm:inline">Basic Info</span>
            </TabsTrigger>
          )}
          <TabsTrigger value="owner" title="Primary Owner">
            <User className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">{isAmendment ? 'MSME Info' : 'Primary Owner'}</span>
          </TabsTrigger>
        </TabsList>

        {/* Amendment Details Tab - for amendment applications */}
        {isAmendment && amendmentData && (
          <TabsContent value="amendment" className="space-y-6">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      <FileEdit className="h-5 w-5" />
                      Formalisation Amendment Request
                    </CardTitle>
                    <CardDescription>
                      Comparison of current values vs. proposed changes for {amendmentData.sme_name}
                    </CardDescription>
                  </div>
                  <Badge className={getTypeBadgeStyle(applicationType)}>
                    {getTypeLabel(applicationType)}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Header row */}
                <div className="grid grid-cols-3 gap-4 py-2 px-3 bg-muted/50 rounded-md font-medium text-sm">
                  <div>Field</div>
                  <div>Current Value</div>
                  <div>Proposed Value</div>
                </div>

                {/* Registration & Tax */}
                <div className="space-y-1">
                  <h4 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Registration Details</h4>
                  <AmendmentComparisonRow
                    label="Registration Number"
                    current={amendmentData.current?.registration_number}
                    proposed={amendmentData.proposed?.registration_number}
                  />
                  <AmendmentComparisonRow
                    label="Tax Identification Number"
                    current={amendmentData.current?.tax_identification_number}
                    proposed={amendmentData.proposed?.tax_identification_number}
                  />
                </div>

                <Separator />

                {/* Compliance Status */}
                <div className="space-y-1">
                  <h4 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Compliance Status</h4>
                  <AmendmentComparisonRow
                    label="Has Bank Account"
                    current={amendmentData.current?.has_bank_account}
                    proposed={amendmentData.proposed?.has_bank_account}
                  />
                  <AmendmentComparisonRow
                    label="Has Tax Clarification"
                    current={amendmentData.current?.has_tax_clarification}
                    proposed={amendmentData.proposed?.has_tax_clarification}
                  />
                  <AmendmentComparisonRow
                    label="Registered for VAT"
                    current={amendmentData.current?.is_registered_for_vat}
                    proposed={amendmentData.proposed?.is_registered_for_vat}
                  />
                  <AmendmentComparisonRow
                    label="Member of Association"
                    current={amendmentData.current?.is_member_of_association}
                    proposed={amendmentData.proposed?.is_member_of_association}
                  />
                  <AmendmentComparisonRow
                    label="Is Affiliated"
                    current={amendmentData.current?.is_affiliated}
                    proposed={amendmentData.proposed?.is_affiliated}
                  />
                  <AmendmentComparisonRow
                    label="Has Export License"
                    current={amendmentData.current?.has_export_license}
                    proposed={amendmentData.proposed?.has_export_license}
                  />
                  <AmendmentComparisonRow
                    label="Has Accessed BDS"
                    current={amendmentData.current?.has_accessed_bds}
                    proposed={amendmentData.proposed?.has_accessed_bds}
                  />
                </div>

                <Separator />

                {/* Financial Information */}
                <div className="space-y-1">
                  <h4 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Financial Information</h4>
                  <AmendmentComparisonRow
                    label="Annual Turnover"
                    current={amendmentData.current?.annual_turnover}
                    proposed={amendmentData.proposed?.annual_turnover}
                  />
                  <AmendmentComparisonRow
                    label="Estimated Value of Assets"
                    current={amendmentData.current?.estimated_value_of_assets}
                    proposed={amendmentData.proposed?.estimated_value_of_assets}
                  />
                </div>

                {/* Change Reason */}
                {amendmentData.proposed?.change_reason && (
                  <>
                    <Separator />
                    <div className="space-y-2">
                      <h4 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Reason for Changes</h4>
                      <p className="text-sm bg-muted/50 p-3 rounded-md">{amendmentData.proposed.change_reason}</p>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        )}

        <TabsContent value="basic" className="space-y-6">
          {/* Basic Information */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Basic Information</CardTitle>
                  <CardDescription>Contact and registrant details</CardDescription>
                </div>
                <Badge className={getTypeBadgeStyle(applicationType)}>
                  {getTypeLabel(applicationType)}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-muted">
                    <FileText className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">MSME</p>
                    <p className="font-medium text-foreground">{application.sme}</p>
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-muted">
                    <User className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">Registrant Name</p>
                    <p className="font-medium text-foreground">{application.registrant_name}</p>
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">Email</p>
                    <p className="font-medium text-foreground">{application.email}</p>
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">Phone</p>
                    <p className="font-medium text-foreground">{application.phone}</p>
                  </div>
                </div>
              </div>
              <Separator />

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">MSME Registration Number</p>
                  <p className="font-medium text-foreground">{application.sme_registration_number}</p>
                </div>

                <div className="flex items-start gap-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">MSME Tax Identification Number</p>
                    <p className="font-medium text-foreground">{application.sme_tax_identification_number}</p>
                  </div>
                </div>
              </div>
              <Separator />

              <div>
                <p className="text-sm text-muted-foreground mb-2">Status</p>
                <div className="flex items-start gap-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">Current Status</p>
                    <Badge variant={
                      application.status === 'Approved' ? 'default' :
                        application.status === 'Rejected' ? 'destructive' : 'secondary'
                    }>
                      {application.status}
                    </Badge>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="owner">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                {isAmendment ? 'MSME Contact Information' : 'Primary Business Owner'}
              </CardTitle>
              <CardDescription>
                {isAmendment
                  ? 'Contact details for the MSME submitting the amendment'
                  : 'Details of the primary business owner'}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {isAmendment ? (
                // For amendment applications, show MSME contact info
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="MSME Name" value={amendmentData?.sme_name || application.sme} />
                    <DetailRow label="Contact Email" value={application.email} />
                    <DetailRow label="Contact Phone" value={application.phone} />
                  </div>
                  <Alert>
                    <Info className="h-4 w-4" />
                    <AlertDescription>
                      Amendment applications do not include primary business owner details.
                      View the Amendment tab for the proposed changes.
                    </AlertDescription>
                  </Alert>
                </div>
              ) : (
                // For signup applications, show full owner details
                <>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <DetailRow label="First Name" value={application.first_name} />
                    <DetailRow label="Last Name" value={application.last_name} />
                    <DetailRow label="Other Names" value={application.other_names} />
                    <DetailRow label="National ID" value={application.national_id_number} />
                    <DetailRow label="Nationality" value={application.nationality} />
                    <DetailRow label="Date of Birth" value={formatDate(application.date_of_birth)} />
                    <DetailRow label="Gender" value={application.gender} />
                    <DetailRow label="Education Level" value={application.education_level} />
                    <DetailRow label="Malawian Status" value={application.malawian_status} />
                    <div className="space-y-1">
                      <p className="text-sm text-muted-foreground">Special Needs</p>
                      <BooleanBadge value={application.has_special_needs ?? false} />
                    </div>
                  </div>

                  <Separator />

                  <div className="space-y-4">
                    <h4 className="font-semibold">Contact Information</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <DetailRow label="Landline" value={application.landline_number} />
                      <DetailRow label="Email" value={application.email} />
                    </div>
                  </div>

                  <Separator />

                  <div className="space-y-4">
                    <h4 className="font-semibold">Location</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <DetailRow label="Region" value={application.region} />
                      <DetailRow label="District" value={application.district} />
                      <DetailRow label="Traditional Authority" value={application.traditional_authority} />
                      <DetailRow label="Physical Address" value={application.physical_address} />
                      <DetailRow label="Postal Address" value={application.postal_address} />
                    </div>
                  </div>

                  <Separator />

                  <div className="space-y-4">
                    <h4 className="font-semibold">Alternative Contact</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <DetailRow label="Contact Name" value={application.alt_contact_name} />
                      <DetailRow label="Relationship" value={application.alt_contact_relationship} />
                      <DetailRow label="Contact Phone" value={application.alt_contact_phone} />
                    </div>
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Metadata Section */}
      <div className="px-1">
        <h3 className="text-sm font-medium mb-2 text-muted-foreground">Metadata</h3>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-muted-foreground">Created</p>
            <p className="font-medium text-foreground">{formatDate(application.createdAt || application.created_at || null)}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Last Updated</p>
            <p className="font-medium text-foreground">{formatDate(application.updatedAt || application.updated_at || null)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
