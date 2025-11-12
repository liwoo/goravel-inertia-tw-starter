import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { Sme, SmeUpdateData } from '@/types/sme';
import { FolderOpen, Phone, Mail, MapPin, FileText, User } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface SmeEditFormProps extends CrudEditFormProps<Sme> {
  setIsSaving?: (saving: boolean) => void;
}

export const SmeEditForm = forwardRef<any, SmeEditFormProps>(({
  item: sme,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<SmeUpdateData>({
    // usmeNumber is auto-generated and cannot be updated - displayed as read-only
    name: sme.name,
    registrationNumber: sme.registrationNumber || '',
    taxIdentificationNumber: sme.taxIdentificationNumber || '',
    operationalStartDate: sme.operationalStartDate || '',
    businessCategory: sme.businessCategory,
    sector: sme.sector,
    subSector: sme.subSector || '',
    businessDescription: sme.businessDescription || '',
    contactPhone: sme.contactPhone,
    contactEmail: sme.contactEmail,
    physicalAddress: sme.physicalAddress || '',
    postalAddress: sme.postalAddress || '',
    website: sme.website || '',
    // region is inferred from district
    district: sme.district || '',
    traditionalAuthority: sme.traditionalAuthority || '',
    ipAddress: sme.ipAddress || '',
    userAgent: sme.userAgent || '',
    businessImprovementAspects: sme.businessImprovementAspects,
    businessAccessedFinancing: sme.businessAccessedFinancing,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    // usmeNumber is auto-generated and cannot be updated - no validation needed
    if (!formData.name?.trim()) {
      newErrors.name = 'Name is required';
    }
    if (!formData.businessCategory?.trim()) {
      newErrors.businessCategory = 'Business Category is required';
    }
    if (!formData.sector?.trim()) {
      newErrors.sector = 'Sector is required';
    }
    if (!formData.contactPhone?.trim()) {
      newErrors.contactPhone = 'Contact Phone is required';
    }
    if (!formData.contactEmail?.trim()) {
      newErrors.contactEmail = 'Contact Email is required';
    }
    if (!formData.businessImprovementAspects?.trim()) {
      newErrors.businessImprovementAspects = 'Business Improvement Aspects is required';
    }
    if (!formData.businessAccessedFinancing?.trim()) {
      newErrors.businessAccessedFinancing = 'Business Accessed Financing is required';
    }
    

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`/api/smes/${sme.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Sme updated successfully');
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

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Sme Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="usmeNumber">USME Number (Auto-generated)</Label>
                <Input
                  id="usmeNumber"
                  value={sme.usmeNumber}
                  readOnly
                  disabled
                  className="bg-muted cursor-not-allowed"
                />
                <p className="text-xs text-muted-foreground">This field is auto-generated and cannot be edited</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <User className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter name"
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="registrationNumber">Registration Number</Label>
                <Input
                  id="registrationNumber"
                  value={formData.registrationNumber}
                  onChange={(e) => setFormData({ ...formData, registrationNumber: e.target.value })}
                  placeholder="Enter registration number"
                  className={errors.registrationNumber ? 'border-destructive' : ''}
                />
                {errors.registrationNumber && (
                  <p className="text-sm text-destructive">{errors.registrationNumber}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="taxIdentificationNumber">Tax Identification Number</Label>
                <Input
                  id="taxIdentificationNumber"
                  value={formData.taxIdentificationNumber}
                  onChange={(e) => setFormData({ ...formData, taxIdentificationNumber: e.target.value })}
                  placeholder="Enter tax identification number"
                  className={errors.taxIdentificationNumber ? 'border-destructive' : ''}
                />
                {errors.taxIdentificationNumber && (
                  <p className="text-sm text-destructive">{errors.taxIdentificationNumber}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="operationalStartDate">Operational Start Date</Label>
                <Input
                  id="operationalStartDate"
                  value={formData.operationalStartDate}
                  onChange={(e) => setFormData({ ...formData, operationalStartDate: e.target.value })}
                  placeholder="Enter operational start date"
                  className={errors.operationalStartDate ? 'border-destructive' : ''}
                />
                {errors.operationalStartDate && (
                  <p className="text-sm text-destructive">{errors.operationalStartDate}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FolderOpen className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="businessCategory">Business Category *</Label>
                <Input
                  id="businessCategory"
                  value={formData.businessCategory}
                  onChange={(e) => setFormData({ ...formData, businessCategory: e.target.value })}
                  placeholder="Enter business category"
                  className={errors.businessCategory ? 'border-destructive' : ''}
                />
                {errors.businessCategory && (
                  <p className="text-sm text-destructive">{errors.businessCategory}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="sector">Sector *</Label>
                <Input
                  id="sector"
                  value={formData.sector}
                  onChange={(e) => setFormData({ ...formData, sector: e.target.value })}
                  placeholder="Enter sector"
                  className={errors.sector ? 'border-destructive' : ''}
                />
                {errors.sector && (
                  <p className="text-sm text-destructive">{errors.sector}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="subSector">Sub Sector</Label>
                <Input
                  id="subSector"
                  value={formData.subSector}
                  onChange={(e) => setFormData({ ...formData, subSector: e.target.value })}
                  placeholder="Enter sub sector"
                  className={errors.subSector ? 'border-destructive' : ''}
                />
                {errors.subSector && (
                  <p className="text-sm text-destructive">{errors.subSector}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="businessDescription">Business Description</Label>
                <Textarea
                  id="businessDescription"
                  value={formData.businessDescription}
                  onChange={(e) => setFormData({ ...formData, businessDescription: e.target.value })}
                  placeholder="Enter business description"
                  rows={4}
                  className={errors.businessDescription ? 'border-destructive' : ''}
                />
                {errors.businessDescription && (
                  <p className="text-sm text-destructive">{errors.businessDescription}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Phone className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="contactPhone">Contact Phone *</Label>
                <Input
                  id="contactPhone"
                  value={formData.contactPhone}
                  onChange={(e) => setFormData({ ...formData, contactPhone: e.target.value })}
                  placeholder="Enter phone number"
                  className={errors.contactPhone ? 'border-destructive' : ''}
                />
                {errors.contactPhone && (
                  <p className="text-sm text-destructive">{errors.contactPhone}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Mail className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="contactEmail">Contact Email *</Label>
                <Input
                  id="contactEmail"
                  value={formData.contactEmail}
                  onChange={(e) => setFormData({ ...formData, contactEmail: e.target.value })}
                  placeholder="Enter email address"
                  className={errors.contactEmail ? 'border-destructive' : ''}
                />
                {errors.contactEmail && (
                  <p className="text-sm text-destructive">{errors.contactEmail}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="physicalAddress">Physical Address</Label>
                <Textarea
                  id="physicalAddress"
                  value={formData.physicalAddress}
                  onChange={(e) => setFormData({ ...formData, physicalAddress: e.target.value })}
                  placeholder="Enter address"
                  rows={4}
                  className={errors.physicalAddress ? 'border-destructive' : ''}
                />
                {errors.physicalAddress && (
                  <p className="text-sm text-destructive">{errors.physicalAddress}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="postalAddress">Postal Address</Label>
                <Textarea
                  id="postalAddress"
                  value={formData.postalAddress}
                  onChange={(e) => setFormData({ ...formData, postalAddress: e.target.value })}
                  placeholder="Enter address"
                  rows={4}
                  className={errors.postalAddress ? 'border-destructive' : ''}
                />
                {errors.postalAddress && (
                  <p className="text-sm text-destructive">{errors.postalAddress}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="website">Website</Label>
                <Input
                  id="website"
                  value={formData.website}
                  onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                  placeholder="Enter website"
                  className={errors.website ? 'border-destructive' : ''}
                />
                {errors.website && (
                  <p className="text-sm text-destructive">{errors.website}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="region">Region (Inferred from District)</Label>
                <Input
                  id="region"
                  value={sme.region || ''}
                  readOnly
                  disabled
                  className="bg-muted cursor-not-allowed"
                />
                <p className="text-xs text-muted-foreground">This field is inferred from the district and cannot be edited</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="district">District</Label>
                <Input
                  id="district"
                  value={formData.district}
                  onChange={(e) => setFormData({ ...formData, district: e.target.value })}
                  placeholder="Enter district"
                  className={errors.district ? 'border-destructive' : ''}
                />
                {errors.district && (
                  <p className="text-sm text-destructive">{errors.district}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="traditionalAuthority">Traditional Authority</Label>
                <Input
                  id="traditionalAuthority"
                  value={formData.traditionalAuthority}
                  onChange={(e) => setFormData({ ...formData, traditionalAuthority: e.target.value })}
                  placeholder="Enter traditional authority"
                  className={errors.traditionalAuthority ? 'border-destructive' : ''}
                />
                {errors.traditionalAuthority && (
                  <p className="text-sm text-destructive">{errors.traditionalAuthority}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="ipAddress">Ip Address</Label>
                <Textarea
                  id="ipAddress"
                  value={formData.ipAddress}
                  onChange={(e) => setFormData({ ...formData, ipAddress: e.target.value })}
                  placeholder="Enter address"
                  rows={4}
                  className={errors.ipAddress ? 'border-destructive' : ''}
                />
                {errors.ipAddress && (
                  <p className="text-sm text-destructive">{errors.ipAddress}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="userAgent">User Agent</Label>
                <Input
                  id="userAgent"
                  value={formData.userAgent}
                  onChange={(e) => setFormData({ ...formData, userAgent: e.target.value })}
                  placeholder="Enter user agent"
                  className={errors.userAgent ? 'border-destructive' : ''}
                />
                {errors.userAgent && (
                  <p className="text-sm text-destructive">{errors.userAgent}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="businessImprovementAspects">Business Improvement Aspects *</Label>
                <Input
                  id="businessImprovementAspects"
                  value={formData.businessImprovementAspects}
                  onChange={(e) => setFormData({ ...formData, businessImprovementAspects: e.target.value })}
                  placeholder="Enter business improvement aspects"
                  className={errors.businessImprovementAspects ? 'border-destructive' : ''}
                />
                {errors.businessImprovementAspects && (
                  <p className="text-sm text-destructive">{errors.businessImprovementAspects}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="businessAccessedFinancing">Business Accessed Financing *</Label>
                <Input
                  id="businessAccessedFinancing"
                  value={formData.businessAccessedFinancing}
                  onChange={(e) => setFormData({ ...formData, businessAccessedFinancing: e.target.value })}
                  placeholder="Enter business accessed financing"
                  className={errors.businessAccessedFinancing ? 'border-destructive' : ''}
                />
                {errors.businessAccessedFinancing && (
                  <p className="text-sm text-destructive">{errors.businessAccessedFinancing}</p>
                )}
              </div>
            </div>

          </div>
        </div>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">ID</p>
              <p className="font-medium text-foreground">#{sme.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">
                {(sme.createdAt || sme.created_at) ? new Date(sme.createdAt || sme.created_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">
                {(sme.updatedAt || sme.updated_at) ? new Date(sme.updatedAt || sme.updated_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
});

SmeEditForm.displayName = 'SmeEditForm';
