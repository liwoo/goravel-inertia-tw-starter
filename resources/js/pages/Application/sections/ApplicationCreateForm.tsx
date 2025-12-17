import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { ApplicationCreateData } from '@/types/application';
import { Tag, FileText, User, Mail, Phone, MapPin, Briefcase } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { DISTRICT_NAMES } from '@/constants/districts';
import { NATIONALITY_OPTIONS } from '@/types/nationalities';
import { EDUCATION_OPTIONS } from '@/types/education';
import { MALAWIAN_STATUS_OPTIONS } from '@/types/malawian-status';
import {
  formatMalawiPhone,
  validateMalawiPhone,
  validateBusinessRegistration,
  validateTIN,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';

interface ApplicationCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
}

export const ApplicationCreateForm = forwardRef<any, ApplicationCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<ApplicationCreateData>({
    sme: '',
    registrant_name: '',
    email: '',
    phone: '',
    sme_registration_number: '',
    sme_tax_identification_number: '',
    status: 'Pending',

    // Owner Details
    first_name: '',
    last_name: '',
    other_names: '',
    nationality: '',
    national_id_number: '',
    date_of_birth: '',
    gender: 'MALE',
    education_level: '',
    malawian_status: '',
    has_special_needs: false,
    landline_number: '',
    physical_address: '',
    postal_address: '',
    region: '',
    district: '',
    traditional_authority: '',
    alt_contact_name: '',
    alt_contact_relationship: '',
    alt_contact_phone: '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    const newErrors: Record<string, string> = {};

    // Basic validation
    if (!formData.sme?.trim()) newErrors.sme = 'SME is required';
    if (!formData.registrant_name?.trim()) newErrors.registrant_name = 'Registrant Name is required';
    if (!formData.email?.trim()) newErrors.email = 'Email is required';

    // Phone validation with Malawi format
    if (!formData.phone?.trim()) {
      newErrors.phone = 'Phone is required';
    } else if (!validateMalawiPhone(formData.phone)) {
      newErrors.phone = VALIDATION_MESSAGES.PHONE;
    }

    // Business Registration Number is optional, but validate format if provided
    if (formData.sme_registration_number?.trim() && !validateBusinessRegistration(formData.sme_registration_number)) {
      newErrors.sme_registration_number = VALIDATION_MESSAGES.BUSINESS_REG;
    }

    // Tax ID is optional, but validate format if provided
    if (formData.sme_tax_identification_number?.trim() && !validateTIN(formData.sme_tax_identification_number)) {
      newErrors.sme_tax_identification_number = VALIDATION_MESSAGES.TIN;
    }

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
      const response = await fetch('/api/applications', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Application created successfully');
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
                    onChange={(e) => {
                      const formatted = formatMalawiPhone(e.target.value);
                      setFormData({ ...formData, phone: formatted });
                    }}
                    placeholder={EXAMPLE_FORMATS.PHONE}
                    className={`pl-9 ${errors.phone ? 'border-destructive' : ''}`}
                  />
                </div>
                {errors.phone && <p className="text-sm text-destructive">{errors.phone}</p>}
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
                <Label htmlFor="sme_registration_number">Business Registration Number</Label>
                <Input
                  id="sme_registration_number"
                  value={formData.sme_registration_number}
                  onChange={(e) => setFormData({ ...formData, sme_registration_number: e.target.value.toUpperCase() })}
                  placeholder={EXAMPLE_FORMATS.BUSINESS_REG}
                  className={errors.sme_registration_number ? 'border-destructive' : ''}
                />
                {errors.sme_registration_number && <p className="text-sm text-destructive">{errors.sme_registration_number}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="sme_tax_identification_number">Tax Identification Number</Label>
                <Input
                  id="sme_tax_identification_number"
                  value={formData.sme_tax_identification_number}
                  onChange={(e) => setFormData({ ...formData, sme_tax_identification_number: e.target.value })}
                  placeholder={EXAMPLE_FORMATS.TIN}
                  className={errors.sme_tax_identification_number ? 'border-destructive' : ''}
                />
                {errors.sme_tax_identification_number && <p className="text-sm text-destructive">{errors.sme_tax_identification_number}</p>}
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
                <Select
                  value={formData.nationality}
                  onValueChange={(value) => setFormData({ ...formData, nationality: value })}
                >
                  <SelectTrigger className={errors.nationality ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Select nationality" />
                  </SelectTrigger>
                  <SelectContent>
                    {NATIONALITY_OPTIONS.map((nationality) => (
                      <SelectItem key={nationality.value} value={nationality.value}>
                        {nationality.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.nationality && <p className="text-sm text-destructive">{errors.nationality}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="national_id_number">National ID *</Label>
                <Input
                  id="national_id_number"
                  value={formData.national_id_number}
                  onChange={(e) => setFormData({ ...formData, national_id_number: e.target.value.toUpperCase() })}
                  placeholder={EXAMPLE_FORMATS.NATIONAL_ID}
                  maxLength={8}
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
                <Select
                  value={formData.education_level}
                  onValueChange={(value) => setFormData({ ...formData, education_level: value })}
                >
                  <SelectTrigger className={errors.education_level ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Select education level" />
                  </SelectTrigger>
                  <SelectContent>
                    {EDUCATION_OPTIONS.map((education) => (
                      <SelectItem key={education.value} value={education.value}>
                        {education.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.education_level && <p className="text-sm text-destructive">{errors.education_level}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="malawian_status">Malawian Status *</Label>
                <Select
                  value={formData.malawian_status}
                  onValueChange={(value) => setFormData({ ...formData, malawian_status: value })}
                >
                  <SelectTrigger className={errors.malawian_status ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    {MALAWIAN_STATUS_OPTIONS.map((status) => (
                      <SelectItem key={status.value} value={status.value}>
                        {status.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
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
                  onChange={(e) => {
                    const formatted = formatMalawiPhone(e.target.value);
                    setFormData({ ...formData, landline_number: formatted });
                  }}
                  placeholder={EXAMPLE_FORMATS.PHONE}
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
                <Select
                  value={formData.district || ''}
                  onValueChange={(value) => setFormData({ ...formData, district: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select district" />
                  </SelectTrigger>
                  <SelectContent>
                    {DISTRICT_NAMES.map((district) => (
                      <SelectItem key={district} value={district}>
                        {district}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
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
                  onChange={(e) => {
                    const formatted = formatMalawiPhone(e.target.value);
                    setFormData({ ...formData, alt_contact_phone: formatted });
                  }}
                  placeholder={EXAMPLE_FORMATS.PHONE}
                />
              </div>
            </div>
          </CardContent>
        </Card>

      </div>
    </form>
  );
});

ApplicationCreateForm.displayName = 'ApplicationCreateForm';
