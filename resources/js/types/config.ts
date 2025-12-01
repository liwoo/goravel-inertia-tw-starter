// TypeScript interfaces for Config entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Config type enum matching backend (app/http/requests/config_type.go)
export enum ConfigType {
  Financing = 'Financing',
  ImprovementAspects = 'Improvement Aspects',
  BusinessCategories = 'Business Categories',
  Industries = 'Industries',
  Sectors = 'Sectors',
  RegistrationStatus = 'Registration Status',
  DevelopmentPartners = 'Development Partners',
  ProcurementType = 'Procurement Type',
  ProcuredBy = 'Procured By',
  Organization = 'Organization',
}

// Helper array for dropdowns
export const CONFIG_TYPES = Object.values(ConfigType);

// Config type display configuration
export const CONFIG_TYPE_CONFIG = {
  [ConfigType.Financing]: {
    label: 'Financing',
    description: 'Financing options and methods',
  },
  [ConfigType.ImprovementAspects]: {
    label: 'Improvement Aspects',
    description: 'Business improvement aspects',
  },
  [ConfigType.BusinessCategories]: {
    label: 'Business Categories',
    description: 'Categories of business operations',
  },
  [ConfigType.Industries]: {
    label: 'Industries',
    description: 'Industry classifications',
  },
  [ConfigType.Sectors]: {
    label: 'Sectors',
    description: 'Economic sectors',
  },
  [ConfigType.RegistrationStatus]: {
    label: 'Registration Status',
    description: 'Business registration statuses',
  },
  [ConfigType.DevelopmentPartners]: {
    label: 'Development Partners',
    description: 'Development partner organizations',
  },
  [ConfigType.ProcurementType]: {
    label: 'Procurement Type',
    description: 'Types of procurement',
  },
  [ConfigType.ProcuredBy]: {
    label: 'Procured By',
    description: 'Entities procuring goods/services',
  },
  [ConfigType.Organization]: {
    label: 'Organization',
    description: 'Organizations involved',
  },
} as const;

// Core Config interface matching the backend model
export interface Config extends BaseModel {
  name: string;
  code?: string;
  configType: ConfigType;
  config_type?: ConfigType; // Backend sends snake_case
  description?: string;
}

// Config creation data (matches ConfigCreateRequest)
export interface ConfigCreateData {
  name: string;
  code?: string;
  configType: ConfigType;
  description?: string;
}

// Config update data (matches ConfigUpdateRequest - all optional)
export interface ConfigUpdateData {
  name?: string;
  code?: string;
  configType?: ConfigType;
  description?: string;
}

// Config list response (matches service GetList response)
export interface ConfigListResponse extends PaginatedResult<Config> { }

// Config list request (extends base ListRequest with config-specific filters)
export interface ConfigListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface ConfigFormErrors {
  name?: string;
  code?: string;
  configType?: string;
  description?: string;
  general?: string;
}

// Config statistics (if provided by backend)
export interface ConfigStats {
  totalconfigs: number;
  // Add your custom stats here
}
