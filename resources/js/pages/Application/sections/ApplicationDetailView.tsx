import React from 'react';
import { Building2, User } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Application } from '@/types/application';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { DetailRow } from '@/components/ui/details-row';
import { BooleanBadge } from '@/components/ui/boolean-badge';

export function ApplicationDetailView({
  item: application,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Application>) {
  const formatDate = (date: string | Date | null) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  return (
    <div className="space-y-6">
      <Tabs defaultValue="basic" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="basic" title="Basic Information">
            <Building2 className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Basic Info</span>
          </TabsTrigger>
          <TabsTrigger value="owner" title="Primary Owner">
            <User className="h-4 w-4 sm:mr-2" />
            <span className="hidden sm:inline">Primary Owner</span>
          </TabsTrigger>
        </TabsList>
        <TabsContent value="basic" className="space-y-6">
          {/* Basic Information */}
          <Card>
            <CardHeader>
              <CardTitle>Basic Information</CardTitle>
              <CardDescription>Contact and registrant details</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-start gap-3">
                  <div className="p-2 rounded-lg bg-muted">
                    <FileText className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">SME</p>
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
                  <p className="text-sm text-muted-foreground">SME Registration Number</p>
                  <p className="font-medium text-foreground">{application.sme_registration_number}</p>
                </div>

                <div className="flex items-start gap-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm text-muted-foreground">SME Tax Identification Number</p>
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
                Primary Business Owner
              </CardTitle>
              <CardDescription>Details of the primary business owner</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
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
                  <BooleanBadge value={application.has_special_needs} />
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
