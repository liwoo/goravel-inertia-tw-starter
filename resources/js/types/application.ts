// TypeScript interfaces for Application entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Application interface matching the backend model
export interface Application extends BaseModel {
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
