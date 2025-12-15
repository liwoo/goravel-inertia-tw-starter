import React, { forwardRef, useState, useEffect } from 'react';
import { CrudEditFormProps } from '@/types/crud';
import { Sme } from '@/types/sme';
import { SmeEditFormSimple } from './SmeEditFormSimple';
import axios from 'axios';
import { BusinessFormalisation } from '@/types/business_formalisation';
import { Loader2 } from 'lucide-react';

interface SmeEditFormProps extends CrudEditFormProps<Sme> {
  setIsSaving?: (saving: boolean) => void;
}

// Type definitions for related data
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

export const SmeEditForm = forwardRef<any, SmeEditFormProps>(({
  item: sme,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [primaryOwner, setPrimaryOwner] = useState<PrimaryOwnerData | undefined>();
  const [formalization, setFormalization] = useState<BusinessFormalisation | undefined>();
  const [additionalMembers, setAdditionalMembers] = useState<AdditionalMemberData[]>([]);
  const [loadingRelatedData, setLoadingRelatedData] = useState(true);

  useEffect(() => {
    const fetchRelatedData = async () => {
      setLoadingRelatedData(true);
      try {
        // Fetch all related data in parallel using the SME-specific endpoints
        const [ownerRes, formalizationRes, membersRes] = await Promise.all([
          // Fetch primary owner
          axios.get(`/api/smes/${sme.id}/primary_business_owner`)
            .catch(() => ({ data: { data: null } })),

          // Fetch formalization
          axios.get(`/api/smes/${sme.id}/business_formalisation`)
            .catch(() => ({ data: { data: null } })),

          // Fetch additional members
          axios.get(`/api/smes/${sme.id}/additional_business_members`)
            .catch(() => ({ data: { data: [] } }))
        ]);

        // Extract the data from responses (these are direct endpoints, not paginated)
        const ownerData = ownerRes.data.data;
        const formalizationData = formalizationRes.data.data;
        const membersData = membersRes.data.data || [];

        // Convert from snake_case to camelCase for primary owner
        if (ownerData) {
          setPrimaryOwner({
            id: ownerData.id,
            firstName: ownerData.first_name || '',
            lastName: ownerData.last_name || '',
            otherNames: ownerData.other_names || undefined,
            nationality: ownerData.nationality || '',
            nationalIdNumber: ownerData.national_id_number || '',
            dateOfBirth: ownerData.date_of_birth ? ownerData.date_of_birth.slice(0, 10) : undefined,
            gender: ownerData.gender || '',
            educationLevel: ownerData.education_level || '',
            malawianStatus: ownerData.malawian_status || '',
            hasSpecialNeeds: ownerData.has_special_needs || false,
            phoneNumber: ownerData.phone_number || '',
            landlineNumber: ownerData.landline_number || undefined,
            email: ownerData.email || undefined,
            physicalAddress: ownerData.physical_address || undefined,
            postalAddress: ownerData.postal_address || undefined,
            district: ownerData.district || undefined,
            traditionalAuthority: ownerData.traditional_authority || undefined,
            altContactName: ownerData.alt_contact_name || undefined,
            altContactRelationship: ownerData.alt_contact_relationship || undefined,
            altContactPhone: ownerData.alt_contact_phone || undefined,
          });
        }

        // Convert from snake_case to camelCase for formalization
        if (formalizationData) {
          setFormalization({
            id: formalizationData.id,
            smeId: formalizationData.sme_id,
            hasBankAccount: formalizationData.has_bank_account || false,
            hasTaxClarification: formalizationData.has_tax_clarification || false,
            isRegisteredForVat: formalizationData.is_registered_for_vat || false,
            isMemberOfAssociation: formalizationData.is_member_of_association || false,
            isAffiliated: formalizationData.is_affiliated || false,
            hasExportLicense: formalizationData.has_export_license || false,
            hasAccessedBds: formalizationData.has_accessed_bds || false,
            annualTurnover: formalizationData.annual_turnover || 0,
            estimatedValueOfAssets: formalizationData.estimated_value_of_assets || 0,
            formalisationScore: formalizationData.formalisation_score || 0,
            createdAt: formalizationData.created_at,
            updatedAt: formalizationData.updated_at,
            deletedAt: formalizationData.deleted_at,
          });
        }

        // Convert from snake_case to camelCase for additional members
        const convertedMembers = membersData.map((member: any) => ({
          id: member.id,
          firstName: member.first_name || '',
          lastName: member.last_name || '',
          otherNames: member.other_names || undefined,
          nationality: member.nationality || '',
          nationalIdNumber: member.national_id_number || '',
          dateOfBirth: member.date_of_birth ? member.date_of_birth.slice(0, 10) : undefined,
          gender: member.gender || '',
          email: member.email || undefined,
          phoneNumber: member.phone_number || '',
          isIntern: member.is_intern || false,
          isPartTime: member.is_part_time || false,
        }));
        setAdditionalMembers(convertedMembers);

      } catch (error) {
        console.error('Error fetching related data:', error);
      } finally {
        setLoadingRelatedData(false);
      }
    };

    if (sme?.id) {
      fetchRelatedData();
    }
  }, [sme?.id]);

  if (loadingRelatedData) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        <span className="ml-2 text-muted-foreground">Loading SME data...</span>
      </div>
    );
  }

  return (
    <SmeEditFormSimple
      ref={ref}
      item={sme}
      primaryOwner={primaryOwner}
      formalization={formalization}
      additionalMembers={additionalMembers}
      onSuccess={onSuccess}
      onError={onError}
      onCancel={onCancel}
      isLoading={isLoading}
      setIsSaving={setIsSaving}
    />
  );
});

SmeEditForm.displayName = 'SmeEditForm';