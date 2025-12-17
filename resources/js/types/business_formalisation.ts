// TypeScript interfaces for BusinessFormalisation entities and operations
import { BaseModel } from './crud';

// Core BusinessFormalisation interface matching the backend model
export interface BusinessFormalisation extends BaseModel {
  smeId: number;
  hasBankAccount: boolean;
  hasTaxClarification: boolean;
  isRegisteredForVat: boolean;
  isMemberOfAssociation: boolean;
  isAffiliated: boolean;
  hasExportLicense: boolean;
  hasAccessedBds: boolean;
  annualTurnover: number;
  estimatedValueOfAssets: number;
  formalisationScore: number;
  complianceScore: number;
  teamStructureScore: number;
  financialScore: number;
  createdBy?: number;
  updatedBy?: number;
  deletedBy?: number;
  deletedAt?: string;
  deleted_at?: string;
  ipAddress?: string;
  userAgent?: string;
}

// BusinessFormalisation creation data (matches BusinessFormalisationCreateRequest)
export interface BusinessFormalisationCreateData {
  smeId?: number; // Optional since it will be set from the SME creation context
  hasBankAccount: boolean;
  hasTaxClarification: boolean;
  isRegisteredForVat: boolean;
  isMemberOfAssociation: boolean;
  isAffiliated: boolean;
  hasExportLicense: boolean;
  hasAccessedBds: boolean;
  annualTurnover: number;
  estimatedValueOfAssets: number;
  // formalisationScore is excluded - it's calculated internally
}

// BusinessFormalisation update data (matches BusinessFormalisationUpdateRequest)
export interface BusinessFormalisationUpdateData {
  hasBankAccount?: boolean;
  hasTaxClarification?: boolean;
  isRegisteredForVat?: boolean;
  isMemberOfAssociation?: boolean;
  isAffiliated?: boolean;
  hasExportLicense?: boolean;
  hasAccessedBds?: boolean;
  annualTurnover?: number;
  estimatedValueOfAssets?: number;
  // formalisationScore is excluded - it's calculated internally
}

// Form errors interface for validation
export interface BusinessFormalisationFormErrors {
  hasBankAccount?: string;
  hasTaxClarification?: string;
  isRegisteredForVat?: string;
  isMemberOfAssociation?: string;
  isAffiliated?: string;
  hasExportLicense?: string;
  hasAccessedBds?: string;
  annualTurnover?: string;
  estimatedValueOfAssets?: string;
}

// Helper type for form state in the SME creation wizard
export interface FormalisationFormData {
  hasBankAccount: boolean;
  hasTaxClarification: boolean;
  isRegisteredForVat: boolean;
  isMemberOfAssociation: boolean;
  isAffiliated: boolean;
  hasExportLicense: boolean;
  hasAccessedBds: boolean;
  annualTurnover: string; // String for form input, will be parsed to number
  estimatedValueOfAssets: string; // String for form input, will be parsed to number
}

// Default values for form initialization
export const DEFAULT_FORMALISATION_DATA: FormalisationFormData = {
  hasBankAccount: false,
  hasTaxClarification: false,
  isRegisteredForVat: false,
  isMemberOfAssociation: false,
  isAffiliated: false,
  hasExportLicense: false,
  hasAccessedBds: false,
  annualTurnover: '',
  estimatedValueOfAssets: '',
};