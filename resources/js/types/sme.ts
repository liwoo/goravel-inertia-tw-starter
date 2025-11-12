// TypeScript interfaces for Sme entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Sme interface matching the backend model
export interface Sme extends BaseModel {
  usmeNumber: string;
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
  createdBy?: number;
  updatedBy?: number;
  deletedBy?: number;
  ipAddress?: string;
  userAgent?: string;
  businessImprovementAspects: string[]; // Array, not string
  businessAccessedFinancing: string[]; // Array, not string
  // Relationships (optional, loaded via eager loading)
  primaryBusinessOwner?: any; // PrimaryBusinessOwner type
  additionalBusinessMembers?: any[]; // AdditionalBusinessMember[] type
  businessFormalisation?: any; // BusinessFormalisation type
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

// Sme statistics (if provided by backend)
export interface SmeStats {
  totalsmes: number;
  // Add your custom stats here
}
