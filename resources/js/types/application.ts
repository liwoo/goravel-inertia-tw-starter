// TypeScript interfaces for Application entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Application type enum
export type ApplicationType = 'signup' | 'amend_formalisation';

// Application status enum
export type ApplicationStatus = 'Pending' | 'Approved' | 'Rejected';

// Formalisation amendment data structure (matches backend)
export interface FormalisationAmendmentData {
  registration_number?: string;
  tax_identification_number?: string;
  has_bank_account?: boolean;
  has_tax_clarification?: boolean;
  is_registered_for_vat?: boolean;
  is_member_of_association?: boolean;
  is_affiliated?: boolean;
  has_export_license?: boolean;
  has_accessed_bds?: boolean;
  annual_turnover?: number;
  estimated_value_of_assets?: number;
  change_reason?: string;
}

// Amendment application data with current vs proposed
export interface AmendmentApplicationData {
  current: FormalisationAmendmentData;
  proposed: FormalisationAmendmentData;
  sme_name: string;
}

// Core Application interface matching the backend model
export interface Application extends BaseModel {
  // Application type and metadata
  type: ApplicationType;
  data?: string; // JSON string for amendment applications
  sme_id?: number;
  rejection_reason?: string;

  // Common fields
  sme?: string;
  registrant_name?: string;
  email?: string;
  phone?: string;
  sme_registration_number?: string;
  sme_tax_identification_number?: string;
  status: ApplicationStatus;

  // Primary Business Owner Details (for signup applications)
  // These fields are nullable to support amendment applications
  first_name?: string | null;
  last_name?: string | null;
  other_names?: string | null;
  nationality?: string | null;
  national_id_number?: string | null;
  date_of_birth?: string | null;
  gender?: 'MALE' | 'FEMALE' | null;
  education_level?: string | null;
  malawian_status?: string | null;
  has_special_needs?: boolean;
  landline_number?: string | null;
  physical_address?: string | null;
  postal_address?: string | null;
  region?: string | null;
  district?: string | null;
  traditional_authority?: string | null;
  alt_contact_name?: string | null;
  alt_contact_relationship?: string | null;
  alt_contact_phone?: string | null;

  // Nested SME record (when loaded)
  sme_record?: {
    id: number;
    name: string;
  };
}

// Helper to parse amendment data from application
export function parseAmendmentData(dataStr?: string): AmendmentApplicationData | null {
  if (!dataStr) return null;
  try {
    return JSON.parse(dataStr) as AmendmentApplicationData;
  } catch {
    return null;
  }
}

// Application creation data (matches ApplicationCreateRequest)
export interface ApplicationCreateData {
  sme: string;
  registrant_name: string;
  email: string;
  phone: string;
  sme_registration_number: string;
  sme_tax_identification_number: string;
  status: string;

  // Primary Business Owner Details
  first_name: string;
  last_name: string;
  other_names?: string;
  nationality: string;
  national_id_number: string;
  date_of_birth: string;
  gender: 'MALE' | 'FEMALE';
  education_level: string;
  malawian_status: string;
  has_special_needs: boolean;
  landline_number?: string;
  physical_address?: string;
  postal_address?: string;
  region?: string;
  district?: string;
  traditional_authority?: string;
  alt_contact_name?: string;
  alt_contact_relationship?: string;
  alt_contact_phone?: string;
}

// Application update data (matches ApplicationUpdateRequest - all optional)
export interface ApplicationUpdateData extends Partial<ApplicationCreateData> { }

// Application list response (matches service GetList response)
export interface ApplicationListResponse extends PaginatedResult<Application> { }

// Application list request (extends base ListRequest with application-specific filters)
export interface ApplicationListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface ApplicationFormErrors {
  [key: string]: string | undefined;
}

// Application statistics (if provided by backend)
export interface ApplicationStats {
  totalapplications: number;
  // Add your custom stats here
}
