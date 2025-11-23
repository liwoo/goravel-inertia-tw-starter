// TypeScript interfaces for Bdsp entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Bdsp interface matching the backend model
export interface Bdsp extends BaseModel {
  ubdsp_number: string;
  name: string;
  postal_address: string;
  physical_address?: string;
  registration_status?: string;
  partners: string[];
  product_types: string[];
  service_list: BdspService[];
}

export interface BdspService {
  name: string;
  cost: number;
  duration: string;
}

// Bdsp creation data (matches BdspCreateRequest)
export interface BdspCreateData {
  name: string;
  postal_address: string;
  physical_address?: string;
  registration_status?: string;
  partners: string[];
  product_types: string[];
  service_list: BdspService[];
}

// Bdsp update data (matches BdspUpdateRequest - all optional)
export interface BdspUpdateData {
  name?: string;
  postal_address?: string;
  physical_address?: string;
  registration_status?: string;
  partners?: string[];
  product_types?: string[];
  service_list?: BdspService[];
}

// Bdsp list response (matches service GetList response)
export interface BdspListResponse extends PaginatedResult<Bdsp> { }

// Bdsp list request (extends base ListRequest with bdsp-specific filters)
export interface BdspListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface BdspFormErrors {
  name?: string;
  postal_address?: string;
  physical_address?: string;
  registration_status?: string;
  partners?: string[];
  product_types?: string[];
  service_list?: BdspService[];
  general?: string;
}

// Bdsp statistics (if provided by backend)
export interface BdspStats {
  totalBdsps: number;
  pendingBdsps: number;
  activeBdsps: number;
  rejectedBdsps: number;
  suspendedBdsps: number;
}
