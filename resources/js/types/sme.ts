// TypeScript interfaces for Sme entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// SME Classification type based on Malawi MSME Policy
export type SmeClassification = 'Micro' | 'Small' | 'Medium' | 'Unclassified';

// Core Sme interface matching the backend model
// Includes both camelCase (frontend) and snake_case (API response) variants
export interface Sme extends BaseModel {
  // CamelCase properties (canonical)
  usmeNumber: string;
  name: string;
  isActive: boolean;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  operationalStartDate?: string;
  businessCategory: string;
  sector: string;
  subSector?: string;
  businessDescription?: string;
  contactPhone: string;
  contactEmail: string;
  physicalAddress?: string;
  postalAddress?: string;
  website?: string;
  district?: string;
  traditionalAuthority?: string;
  createdBy?: number;
  updatedBy?: number;
  deletedBy?: number;
  ipAddress?: string;
  userAgent?: string;
  businessImprovementAspects: string[];
  businessAccessedFinancing: string[];
  classification: SmeClassification;

  // Snake_case variants (from API responses)
  usme_number?: string;
  is_active?: boolean;
  registration_number?: string;
  tax_identification_number?: string;
  operational_start_date?: string;
  business_category?: string;
  sub_sector?: string;
  business_description?: string;
  contact_phone?: string;
  contact_email?: string;
  physical_address?: string;
  postal_address?: string;
  traditional_authority?: string;
  created_by?: number;
  updated_by?: number;
  deleted_by?: number;
  ip_address?: string;
  user_agent?: string;
  business_improvement_aspects?: string[];
  business_accessed_financing?: string[];
  // classification is same in both cases (single word)

  // Relationships (optional, loaded via eager loading)
  primaryBusinessOwner?: any; // PrimaryBusinessOwner type
  additionalBusinessMembers?: any[]; // AdditionalBusinessMember[] type
  businessFormalisation?: {
    id?: number;
    formalisationScore?: number;
    formalisation_score?: number;
  }; // BusinessFormalisation type
  business_formalisation?: {
    id?: number;
    formalisationScore?: number;
    formalisation_score?: number;
  }; // BusinessFormalisation snake_case variant
  businessEmployeeSummary?: any; // BusinessEmployeeSummary type
}

// Sme creation data (matches SmeCreateRequest)
export interface SmeCreateData {
  // usmeNumber is auto-generated - not needed in create request
  name: string;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  operationalStartDate?: string;
  businessCategory: string;
  sector: string;
  subSector?: string;
  businessDescription?: string;
  contactPhone: string;
  contactEmail: string;
  physicalAddress?: string;
  postalAddress?: string;
  website?: string;
  // region is inferred from district
  district?: string;
  traditionalAuthority?: string;
  businessImprovementAspects: string[]; // Array, not string
  businessAccessedFinancing: string[]; // Array, not string
}

// Sme update data (matches SmeUpdateRequest - all optional)
export interface SmeUpdateData {
  // usmeNumber is auto-generated and cannot be updated
  name?: string;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  operationalStartDate?: string;
  businessCategory?: string;
  sector?: string;
  subSector?: string;
  businessDescription?: string;
  contactPhone?: string;
  contactEmail?: string;
  physicalAddress?: string;
  postalAddress?: string;
  website?: string;
  // region is inferred from district
  district?: string;
  traditionalAuthority?: string;
  businessImprovementAspects?: string[]; // Array, not string
  businessAccessedFinancing?: string[]; // Array, not string
}

// Sme list response (matches service GetList response)
export interface SmeListResponse extends PaginatedResult<Sme> {}

// Sme list request (extends base ListRequest with sme-specific filters)
export interface SmeListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface SmeFormErrors {
  // usmeNumber is auto-generated - no validation needed
  name?: string;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  operationalStartDate?: string;
  businessCategory?: string;
  sector?: string;
  subSector?: string;
  businessDescription?: string;
  contactPhone?: string;
  contactEmail?: string;
  physicalAddress?: string;
  postalAddress?: string;
  website?: string;
  // region is inferred from district
  district?: string;
  traditionalAuthority?: string;
  businessImprovementAspects?: string; // Error message is string
  businessAccessedFinancing?: string; // Error message is string
  general?: string;
}

// Registration trend data point
export interface RegistrationTrendPoint {
  period: string;
  count: number;
  cumulative: number;
}

// Distribution data point (for region/category breakdowns)
export interface DistributionPoint {
  label: string;
  value: number;
  percentage: number;
}

// Age and gender distribution point (for population pyramid chart)
export interface AgeGenderDistributionPoint {
  ageGroup: string;
  male: number;
  female: number;
  malePercentage: number;
  femalePercentage: number;
}

// Sme statistics (for dashboard charts and KPIs)
export interface SmeStats {
  totalSmes: number;
  newThisMonth: number;
  newLastMonth: number;
  hasRegistration: number;
  hasRegistrationPercentage: number;
  registrationTrend: RegistrationTrendPoint[];
  byRegion: DistributionPoint[];
  byCategory: DistributionPoint[];
  byGender?: DistributionPoint[];
  byYouth?: DistributionPoint[];
  byAgeGender?: AgeGenderDistributionPoint[];
  byClassification?: DistributionPoint[];
}
