import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { Application, ApplicationUpdateData } from '@/types/application';
import { FileText, User, Mail, Phone, Tag, MapPin, Briefcase } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';

interface ApplicationEditFormProps extends CrudEditFormProps<Application> {
  setIsSaving?: (saving: boolean) => void;
}

export const ApplicationEditForm = forwardRef<any, ApplicationEditFormProps>(({
  item: application,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<ApplicationUpdateData>({
    sme: application.sme,
    registrant_name: application.registrant_name,
    email: application.email,
    phone: application.phone,
    sme_registration_number: application.sme_registration_number,
    sme_tax_identification_number: application.sme_tax_identification_number,
    status: application.status,

    // Owner Details - convert null to undefined for form compatibility
    first_name: application.first_name ?? undefined,
    last_name: application.last_name ?? undefined,
    other_names: application.other_names ?? undefined,
    nationality: application.nationality ?? undefined,
    national_id_number: application.national_id_number ?? undefined,
    date_of_birth: application.date_of_birth ? application.date_of_birth.slice(0, 10) : '',
    gender: application.gender ?? undefined,
    education_level: application.education_level ?? undefined,
    malawian_status: application.malawian_status ?? undefined,
    has_special_needs: application.has_special_needs,
    landline_number: application.landline_number ?? undefined,
    physical_address: application.physical_address ?? undefined,
    postal_address: application.postal_address ?? undefined,
    region: application.region ?? undefined,
    district: application.district ?? undefined,
    traditional_authority: application.traditional_authority ?? undefined,
    alt_contact_name: application.alt_contact_name ?? undefined,
    alt_contact_relationship: application.alt_contact_relationship ?? undefined,
    alt_contact_phone: application.alt_contact_phone ?? undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    const newErrors: Record<string, string> = {};

    // Basic validation
    if (!formData.sme?.trim()) newErrors.sme = 'SME is required';
    if (!formData.registrant_name?.trim()) newErrors.registrant_name = 'Registrant Name is required';
    if (!formData.email?.trim()) newErrors.email = 'Email is required';
    if (!formData.phone?.trim()) newErrors.phone = 'Phone is required';
    if (!formData.sme_registration_number?.trim()) newErrors.sme_registration_number = 'Registration Number is required';
    if (!formData.sme_tax_identification_number?.trim()) newErrors.sme_tax_identification_number = 'Tax ID is required';
    if (!formData.status?.trim()) newErrors.status = 'Status is required';

    // Owner validation
    if (!formData.first_name?.trim()) newErrors.first_name = 'First Name is required';
    if (!formData.last_name?.trim()) newErrors.last_name = 'Last Name is required';
    if (!formData.nationality?.trim()) newErrors.nationality = 'Nationality is required';
    if (!formData.national_id_number?.trim()) newErrors.national_id_number = 'National ID is required';
    if (!formData.date_of_birth) newErrors.date_of_birth = 'Date of Birth is required';
    if (!formData.education_level?.trim()) newErrors.education_level = 'Education Level is required';
    if (!formData.malawian_status?.trim()) newErrors.malawian_status = 'Malawian Status is required';

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`/api/applications/${application.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Application updated successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">

        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>Contact and registrant details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="sme">SME *</Label>
                <div className="relative">
                  <FileText className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="sme"
                    value={formData.sme}
                    onChange={(e) => setFormData({ ...formData, sme: e.target.value })}
                    placeholder="Enter SME name"
                    className={`pl-9 ${errors.sme ? 'border-destructive' : ''}`}
                  />
                </div>
                {errors.sme && <p className="text-sm text-destructive">{errors.sme}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="registrant_name">Registrant Name *</Label>
                <div className="relative">
                  <User className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="registrant_name"
                    value={formData.registrant_name}
                    onChange={(e) => setFormData({ ...formData, registrant_name: e.target.value })}
                    placeholder="Enter registrant name"
                    className={`pl-9 ${errors.registrant_name ? 'border-destructive' : ''}`}
                  />
                </div>
                {errors.registrant_name && <p className="text-sm text-destructive">{errors.registrant_name}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email *</Label>
                <div className="relative">
                  <Mail className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="email"
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    placeholder="Enter email address"
                    className={`pl-9 ${errors.email ? 'border-destructive' : ''}`}
                  />
                </div>
                {errors.email && <p className="text-sm text-destructive">{errors.email}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">Phone *</Label>
                <div className="relative">
                  <Phone className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="phone"
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    placeholder="Enter phone number"
                    className={`pl-9 ${errors.phone ? 'border-destructive' : ''}`}
                  />
                </div>
                {errors.phone && <p className="text-sm text-destructive">{errors.phone}</p>}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Owner Personal Information */}
        <Card>
          <CardHeader>
            <CardTitle>Owner Personal Information</CardTitle>
            <CardDescription>Primary business owner details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="first_name">First Name *</Label>
                <Input
                  id="first_name"
                  value={formData.first_name}
                  onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                  className={errors.first_name ? 'border-destructive' : ''}
                />
                {errors.first_name && <p className="text-sm text-destructive">{errors.first_name}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="last_name">Last Name *</Label>
                <Input
                  id="last_name"
                  value={formData.last_name}
                  onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                  className={errors.last_name ? 'border-destructive' : ''}
                />
                {errors.last_name && <p className="text-sm text-destructive">{errors.last_name}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="other_names">Other Names</Label>
                <Input
                  id="other_names"
                  value={formData.other_names || ''}
                  onChange={(e) => setFormData({ ...formData, other_names: e.target.value })}
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="nationality">Nationality *</Label>
                <Input
                  id="nationality"
                  value={formData.nationality}
                  onChange={(e) => setFormData({ ...formData, nationality: e.target.value })}
                  className={errors.nationality ? 'border-destructive' : ''}
                />
                {errors.nationality && <p className="text-sm text-destructive">{errors.nationality}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="national_id_number">National ID *</Label>
                <Input
                  id="national_id_number"
                  value={formData.national_id_number}
                  onChange={(e) => setFormData({ ...formData, national_id_number: e.target.value })}
                  className={errors.national_id_number ? 'border-destructive' : ''}
                />
                {errors.national_id_number && <p className="text-sm text-destructive">{errors.national_id_number}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="date_of_birth">Date of Birth *</Label>
                <Input
                  id="date_of_birth"
                  type="date"
                  value={formData.date_of_birth}
                  onChange={(e) => setFormData({ ...formData, date_of_birth: e.target.value })}
                  className={errors.date_of_birth ? 'border-destructive' : ''}
                />
                {errors.date_of_birth && <p className="text-sm text-destructive">{errors.date_of_birth}</p>}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="gender">Gender *</Label>
                <Select
                  value={formData.gender}
                  onValueChange={(value: 'MALE' | 'FEMALE') => setFormData({ ...formData, gender: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select gender" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="MALE">Male</SelectItem>
                    <SelectItem value="FEMALE">Female</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="education_level">Education Level *</Label>
                <Input
                  id="education_level"
                  value={formData.education_level}
                  onChange={(e) => setFormData({ ...formData, education_level: e.target.value })}
                  className={errors.education_level ? 'border-destructive' : ''}
                />
                {errors.education_level && <p className="text-sm text-destructive">{errors.education_level}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="malawian_status">Malawian Status *</Label>
                <Input
                  id="malawian_status"
                  value={formData.malawian_status}
                  onChange={(e) => setFormData({ ...formData, malawian_status: e.target.value })}
                  className={errors.malawian_status ? 'border-destructive' : ''}
                />
                {errors.malawian_status && <p className="text-sm text-destructive">{errors.malawian_status}</p>}
              </div>
            </div>

            <div className="flex items-center space-x-2 pt-2">
              <Switch
                id="has_special_needs"
                checked={formData.has_special_needs}
                onCheckedChange={(checked) => setFormData({ ...formData, has_special_needs: checked })}
              />
              <Label htmlFor="has_special_needs">Has Special Needs</Label>
            </div>
          </CardContent>
        </Card>

        {/* Owner Contact & Address */}
        <Card>
          <CardHeader>
            <CardTitle>Owner Contact & Address</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="landline_number">Landline Number</Label>
                <Input
                  id="landline_number"
                  value={formData.landline_number || ''}
                  onChange={(e) => setFormData({ ...formData, landline_number: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="region">Region</Label>
                <Input
                  id="region"
                  value={formData.region || ''}
                  onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="district">District</Label>
                <Input
                  id="district"
                  value={formData.district || ''}
                  onChange={(e) => setFormData({ ...formData, district: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="traditional_authority">Traditional Authority</Label>
                <Input
                  id="traditional_authority"
                  value={formData.traditional_authority || ''}
                  onChange={(e) => setFormData({ ...formData, traditional_authority: e.target.value })}
                />
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="physical_address">Physical Address</Label>
                <Textarea
                  id="physical_address"
                  value={formData.physical_address || ''}
                  onChange={(e) => setFormData({ ...formData, physical_address: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="postal_address">Postal Address</Label>
                <Textarea
                  id="postal_address"
                  value={formData.postal_address || ''}
                  onChange={(e) => setFormData({ ...formData, postal_address: e.target.value })}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Alternate Contact */}
        <Card>
          <CardHeader>
            <CardTitle>Alternate Contact</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="alt_contact_name">Name</Label>
                <Input
                  id="alt_contact_name"
                  value={formData.alt_contact_name || ''}
                  onChange={(e) => setFormData({ ...formData, alt_contact_name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="alt_contact_relationship">Relationship</Label>
                <Input
                  id="alt_contact_relationship"
                  value={formData.alt_contact_relationship || ''}
                  onChange={(e) => setFormData({ ...formData, alt_contact_relationship: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="alt_contact_phone">Phone</Label>
                <Input
                  id="alt_contact_phone"
                  value={formData.alt_contact_phone || ''}
                  onChange={(e) => setFormData({ ...formData, alt_contact_phone: e.target.value })}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Registration Details */}
        <Card>
          <CardHeader>
            <CardTitle>Registration Details</CardTitle>
            <CardDescription>Official registration and tax information</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="sme_registration_number">SME Registration Number *</Label>
                <Input
                  id="sme_registration_number"
                  value={formData.sme_registration_number}
                  onChange={(e) => setFormData({ ...formData, sme_registration_number: e.target.value })}
                  placeholder="Enter registration number"
                  className={errors.sme_registration_number ? 'border-destructive' : ''}
                />
                {errors.sme_registration_number && <p className="text-sm text-destructive">{errors.sme_registration_number}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="sme_tax_identification_number">Tax Identification Number *</Label>
                <Input
                  id="sme_tax_identification_number"
                  value={formData.sme_tax_identification_number}
                  onChange={(e) => setFormData({ ...formData, sme_tax_identification_number: e.target.value })}
                  placeholder="Enter tax ID"
                  className={errors.sme_tax_identification_number ? 'border-destructive' : ''}
                />
                {errors.sme_tax_identification_number && <p className="text-sm text-destructive">{errors.sme_tax_identification_number}</p>}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Status */}
        <Card>
          <CardHeader>
            <CardTitle>Application Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 max-w-md">
              <Label htmlFor="status">Status *</Label>
              <Select
                value={formData.status}
                onValueChange={(value) => setFormData({ ...formData, status: value })}
              >
                <SelectTrigger className={errors.status ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="Approved">Approved</SelectItem>
                  <SelectItem value="Rejected">Rejected</SelectItem>
                </SelectContent>
              </Select>
              {errors.status && <p className="text-sm text-destructive">{errors.status}</p>}
            </div>
          </CardContent>
        </Card>

        {/* Metadata Section */}
        <div className="px-1">
          <h3 className="text-sm font-medium mb-2 text-muted-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">ID</p>
              <p className="font-medium text-foreground">#{application.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">
                {(application.createdAt || application.created_at) ? new Date(application.createdAt || application.created_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">
                {(application.updatedAt || application.updated_at) ? new Date(application.updatedAt || application.updated_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
          </div>
        </div>

      </div>
    </form>
  );
});

ApplicationEditForm.displayName = 'ApplicationEditForm';
