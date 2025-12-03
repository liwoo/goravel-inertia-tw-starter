import React from 'react';
import { Calendar, User, Mail, Phone, CreditCard, Globe, Briefcase } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { CrudDetailViewProps } from '@/types/crud';
import { Member } from '@/types/member';

export function MemberDetailView({
  item: member,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Member>) {
  const formatDate = (date: string | Date | null | undefined) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const DetailRow = ({ label, value }: { label: string; value: any }) => (
    <div className="space-y-1">
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className="font-medium text-foreground">{value || '-'}</p>
    </div>
  );

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="h-5 w-5" />
            Personal Information
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <DetailRow label="First Name" value={member.first_name} />
            <DetailRow label="Last Name" value={member.last_name} />
            <DetailRow label="Other Names" value={member.other_names} />
            <DetailRow label="Gender" value={member.gender} />
            <DetailRow label="Nationality" value={member.nationality} />
            <DetailRow label="National ID Number" value={member.national_id_number} />
            <DetailRow label="Date of Birth" value={formatDate(member.date_of_birth)} />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Mail className="h-5 w-5" />
            Contact Information
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <DetailRow label="Email" value={member.email} />
            <DetailRow label="Phone Number" value={member.phone_number} />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Briefcase className="h-5 w-5" />
            Employment Status
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Intern:</span>
              <Badge variant={member.is_intern ? 'default' : 'secondary'}>
                {member.is_intern ? 'Yes' : 'No'}
              </Badge>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Part Time:</span>
              <Badge variant={member.is_part_time ? 'default' : 'secondary'}>
                {member.is_part_time ? 'Yes' : 'No'}
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <Separator />

      {/* Metadata Section */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Created</p>
            <p className="font-medium text-sm text-foreground">{formatDate(member.createdAt || member.created_at || null)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Last Updated</p>
            <p className="font-medium text-sm text-foreground">{formatDate(member.updatedAt || member.updated_at || null)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
