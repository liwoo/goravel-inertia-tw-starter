// TypeScript interfaces for Config entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Config type enum matching backend (app/http/requests/config_type.go)
export enum ConfigType {
  Example = 'Example',
}

// Helper array for dropdowns
export const CONFIG_TYPES = Object.values(ConfigType);

// Config type display configuration
export const CONFIG_TYPE_CONFIG = {
  [ConfigType.Example]: {
    label: 'Example',
    description: 'Example options and methods',
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
