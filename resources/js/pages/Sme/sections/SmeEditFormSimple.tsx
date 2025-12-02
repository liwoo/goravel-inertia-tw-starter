import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { CrudEditFormProps } from '@/types/crud';
import { Sme, SmeUpdateData } from '@/types/sme';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Building2, User, FileText, Users } from 'lucide-react';
import { snakefiyKeys } from '@/lib/utils';
import {
  validateMalawiPhone,
  validateBusinessRegistration,
  validateTIN,
  validateNationalID,
  VALIDATION_MESSAGES
} from '@/lib/malawi-validators';
import { BusinessFormalisation } from '@/types/business_formalisation';
import { SmeEditBusinessInfoTab } from './edit-tabs/SmeEditBusinessInfoTab';
import { SmeEditPrimaryOwnerTab } from './edit-tabs/SmeEditPrimaryOwnerTab';
import { SmeEditFormalizationTab, FormalizationFormState } from './edit-tabs/SmeEditFormalizationTab';
import { SmeEditAdditionalMembersTab } from './edit-tabs/SmeEditAdditionalMembersTab';

interface PrimaryOwnerData {
  id?: number;
  firstName: string;
  lastName: string;
  otherNames?: string;
  nationality: string;
  nationalIdNumber: string;
  dateOfBirth?: string;
  gender: string;
  educationLevel: string;
  malawianStatus: string;
  hasSpecialNeeds: boolean;
  phoneNumber: string;
  landlineNumber?: string;
  email?: string;
  physicalAddress?: string;
  postalAddress?: string;
  district?: string;
  traditionalAuthority?: string;
  altContactName?: string;
  altContactRelationship?: string;
  altContactPhone?: string;
}

interface AdditionalMemberData {
  id?: number;
  firstName: string;
  lastName: string;
  otherNames?: string;
  nationality: string;
  nationalIdNumber: string;
  dateOfBirth?: string;
  gender: string;
  email?: string;
  phoneNumber: string;
  isIntern: boolean;
  isPartTime: boolean;
}

interface SmeEditFormSimpleProps extends CrudEditFormProps<Sme> {
  setIsSaving?: (saving: boolean) => void;
  primaryOwner?: PrimaryOwnerData;
  formalization?: BusinessFormalisation;
  additionalMembers?: AdditionalMemberData[];
}

export const SmeEditFormSimple = forwardRef<any, SmeEditFormSimpleProps>(({
  item: sme,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving,
  primaryOwner: initialPrimaryOwner,
  formalization: initialFormalization,
  additionalMembers: initialAdditionalMembers = []
}, ref) => {
  const [activeTab, setActiveTab] = useState('business-info');
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Business data
  const [businessData, setBusinessData] = useState<SmeUpdateData>({
    name: sme.name,
    registrationNumber: sme.registration_number || sme.registrationNumber || undefined,
    taxIdentificationNumber: sme.tax_identification_number || sme.taxIdentificationNumber || undefined,
    operationalStartDate: sme.operational_start_date || sme.operationalStartDate || '',
    businessCategory: sme.business_category || sme.businessCategory || '',
    sector: sme.sector,
    subSector: sme.sub_sector || sme.subSector || undefined,
    businessDescription: sme.business_description || sme.businessDescription || undefined,
    contactPhone: sme.contact_phone || sme.contactPhone,
    contactEmail: sme.contact_email || sme.contactEmail,
    physicalAddress: sme.physical_address || sme.physicalAddress || undefined,
    postalAddress: sme.postal_address || sme.postalAddress || undefined,
    website: sme.website || undefined,
    district: sme.district || undefined,
    traditionalAuthority: sme.traditional_authority || sme.traditionalAuthority || undefined,
    businessImprovementAspects: Array.isArray(sme.business_improvement_aspects || sme.businessImprovementAspects)
      ? (sme.business_improvement_aspects || sme.businessImprovementAspects)
      : ((sme.business_improvement_aspects || sme.businessImprovementAspects) ? String(sme.business_improvement_aspects || sme.businessImprovementAspects).split(',').filter(Boolean) : []),
    businessAccessedFinancing: Array.isArray(sme.business_accessed_financing || sme.businessAccessedFinancing)
      ? (sme.business_accessed_financing || sme.businessAccessedFinancing)
      : ((sme.business_accessed_financing || sme.businessAccessedFinancing) ? String(sme.business_accessed_financing || sme.businessAccessedFinancing).split(',').filter(Boolean) : []),
  });

  // Primary owner data
  const [primaryOwner, setPrimaryOwner] = useState<PrimaryOwnerData | undefined>(initialPrimaryOwner);

  // Formalization data
  const [formalization, setFormalization] = useState<FormalizationFormState | undefined>(
    initialFormalization || undefined
  );

  // Additional members - track both current and original to detect deletions
  const [additionalMembers, setAdditionalMembers] = useState<AdditionalMemberData[]>(initialAdditionalMembers);
  const [originalMemberIds] = useState<number[]>(initialAdditionalMembers.filter(m => m.id).map(m => m.id!));

  // Save business info
  const handleSaveBusinessInfo = async () => {
    const newErrors: Record<string, string> = {};

    if (!businessData.name?.trim()) newErrors.name = 'Business Name is required';
    if (!businessData.businessCategory) newErrors.businessCategory = 'Business Category is required';
    if (!businessData.sector) newErrors.sector = 'Sector is required';
    if (!businessData.contactPhone) {
      newErrors.contactPhone = 'Contact Phone is required';
    } else if (!validateMalawiPhone(businessData.contactPhone)) {
      newErrors.contactPhone = VALIDATION_MESSAGES.PHONE;
    }
    if (!businessData.contactEmail) newErrors.contactEmail = 'Contact Email is required';

    if (businessData.registrationNumber && !validateBusinessRegistration(businessData.registrationNumber)) {
      newErrors.registrationNumber = VALIDATION_MESSAGES.BUSINESS_REG;
    }

    if (businessData.taxIdentificationNumber && !validateTIN(businessData.taxIdentificationNumber)) {
      newErrors.taxIdentificationNumber = VALIDATION_MESSAGES.TIN;
    }

    if (!businessData.businessImprovementAspects || businessData.businessImprovementAspects.length === 0) {
      newErrors.businessImprovementAspects = 'At least one improvement aspect is required';
    }
    if (!businessData.businessAccessedFinancing || businessData.businessAccessedFinancing.length === 0) {
      newErrors.businessAccessedFinancing = 'At least one financing source is required';
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
        body: JSON.stringify(snakefiyKeys(businessData)),
      });

      if (response.ok) {
        onSuccess('Business information updated successfully');
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

  // Save primary owner
  const handleSavePrimaryOwner = async () => {
    if (!primaryOwner || !primaryOwner.id) {
      onError?.('Primary owner data not available');
      return;
    }

    const newErrors: Record<string, string> = {};
    if (!primaryOwner.firstName) newErrors.firstName = 'First Name is required';
    if (!primaryOwner.lastName) newErrors.lastName = 'Last Name is required';
    if (!primaryOwner.nationality) newErrors.nationality = 'Nationality is required';
    if (!primaryOwner.nationalIdNumber) {
      newErrors.nationalIdNumber = 'National ID is required';
    } else if (!validateNationalID(primaryOwner.nationalIdNumber)) {
      newErrors.nationalIdNumber = VALIDATION_MESSAGES.NATIONAL_ID;
    }
    if (!primaryOwner.dateOfBirth) newErrors.dateOfBirth = 'Date of Birth is required';
    if (!primaryOwner.gender) newErrors.gender = 'Gender is required';
    if (!primaryOwner.educationLevel) newErrors.educationLevel = 'Education Level is required';
    if (!primaryOwner.malawianStatus) newErrors.malawianStatus = 'Malawian Status is required';
    if (!primaryOwner.phoneNumber) {
      newErrors.phoneNumber = 'Phone Number is required';
    } else if (!validateMalawiPhone(primaryOwner.phoneNumber)) {
      newErrors.phoneNumber = VALIDATION_MESSAGES.PHONE;
    }

    if (primaryOwner.landlineNumber && !validateMalawiPhone(primaryOwner.landlineNumber)) {
      newErrors.landlineNumber = VALIDATION_MESSAGES.PHONE;
    }

    if (primaryOwner.altContactPhone && !validateMalawiPhone(primaryOwner.altContactPhone)) {
      newErrors.altContactPhone = VALIDATION_MESSAGES.PHONE;
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`/api/primary-business-owners/${primaryOwner.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(snakefiyKeys(primaryOwner)),
      });

      if (response.ok) {
        onSuccess('Primary owner updated successfully');
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

  // Save formalization
  const handleSaveFormalization = async () => {
    if (!formalization) {
      onError?.('Formalization data not available');
      return;
    }

    setIsSaving?.(true);

    try {
      // If formalization has an id, update it; otherwise create new
      const isUpdate = formalization.id !== undefined;
      const url = isUpdate
        ? `/api/business-formalisations/${formalization.id}`
        : `/api/business-formalisations`;
      const method = isUpdate ? 'PUT' : 'POST';

      // For POST, add sme_id
      const payload = isUpdate
        ? snakefiyKeys(formalization)
        : snakefiyKeys({ ...formalization, sme_id: sme.id });

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        const data = await response.json();
        // Update formalization with the newly created id if it was a create
        if (!isUpdate && data.data?.id) {
          setFormalization({ ...formalization, id: data.data.id });
        }
        onSuccess(isUpdate ? 'Formalization updated successfully' : 'Formalization created successfully');
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

  // Save additional members
  const handleSaveAdditionalMembers = async () => {
    setIsSaving?.(true);

    try {
      const currentMemberIds = additionalMembers.filter(m => m.id).map(m => m.id!);
      const deletedMemberIds = originalMemberIds.filter(id => !currentMemberIds.includes(id));

      // Delete removed members
      for (const memberId of deletedMemberIds) {
        await fetch(`/api/additional_business_members/${memberId}`, {
          method: 'DELETE',
          headers: {
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
        });
      }

      // Create or update members
      for (const member of additionalMembers) {
        const payload = snakefiyKeys({
          ...member,
          smeId: sme.id,
        });

        if (member.id) {
          // Update existing member
          await fetch(`/api/additional_business_members/${member.id}`, {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'X-Requested-With': 'XMLHttpRequest',
            },
            body: JSON.stringify(payload),
          });
        } else {
          // Create new member
          const response = await fetch('/api/additional_business_members', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'X-Requested-With': 'XMLHttpRequest',
            },
            body: JSON.stringify(payload),
          });

          // Update local state with the new id
          if (response.ok) {
            const data = await response.json();
            if (data.data?.id) {
              member.id = data.data.id;
            }
          }
        }
      }

      // Update state with new ids
      setAdditionalMembers([...additionalMembers]);
      onSuccess('Team members saved successfully');
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit - saves the active tab
  useImperativeHandle(ref, () => ({
    handleSubmit: () => {
      switch (activeTab) {
        case 'business-info':
          return handleSaveBusinessInfo();
        case 'primary-owner':
          return handleSavePrimaryOwner();
        case 'formalization':
          return handleSaveFormalization();
        case 'additional-members':
          return handleSaveAdditionalMembers();
        default:
          return handleSaveBusinessInfo();
      }
    }
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="w-full">
          <TabsTrigger value="business-info" className="flex-1">
            <Building2 className="h-4 w-4 mr-2" />
            Business Info
          </TabsTrigger>
          <TabsTrigger value="primary-owner" className="flex-1">
            <User className="h-4 w-4 mr-2" />
            Primary Owner
          </TabsTrigger>
          <TabsTrigger value="formalization" className="flex-1">
            <FileText className="h-4 w-4 mr-2" />
            Formalization
          </TabsTrigger>
          <TabsTrigger value="additional-members" className="flex-1">
            <Users className="h-4 w-4 mr-2" />
            Team Members
          </TabsTrigger>
        </TabsList>

        <TabsContent value="business-info">
          <SmeEditBusinessInfoTab
            businessData={businessData}
            onChange={setBusinessData}
            errors={errors}
          />
        </TabsContent>

        <TabsContent value="primary-owner">
          <SmeEditPrimaryOwnerTab
            primaryOwner={primaryOwner}
            onChange={setPrimaryOwner}
            errors={errors}
          />
        </TabsContent>

        <TabsContent value="formalization">
          <SmeEditFormalizationTab
            formalization={formalization}
            onChange={setFormalization}
            errors={errors}
          />
        </TabsContent>

        <TabsContent value="additional-members">
          <SmeEditAdditionalMembersTab
            members={additionalMembers}
            onChange={setAdditionalMembers}
          />
        </TabsContent>
      </Tabs>

      {/* Metadata */}
      <Card>
        <CardHeader>
          <CardTitle>Metadata</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">SME ID</p>
              <p className="font-medium">#{sme.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">USME Number</p>
              <p className="font-medium">{sme.usme_number || sme.usmeNumber || '-'}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </form>
  );
});

SmeEditFormSimple.displayName = 'SmeEditFormSimple';
