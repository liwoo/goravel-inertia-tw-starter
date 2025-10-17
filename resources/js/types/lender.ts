// TypeScript interfaces for Lender entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Lender interface matching the backend model
export interface Lender extends BaseModel {
  name: string;
  email: string;
  phone?: string;
  address?: string;
  gender?: string;
}

// Lender creation data (matches LenderCreateRequest)
export interface LenderCreateData {
  name: string;
  email: string;
  phone?: string;
  address?: string;
  gender?: string;
}

// Lender update data (matches LenderUpdateRequest - all optional)
export interface LenderUpdateData {
  name?: string;
  email?: string;
  phone?: string;
  address?: string;
  gender?: string;
}

// Lender list response (matches service GetList response)
export interface LenderListResponse extends PaginatedResult<Lender> {}

// Lender list request (extends base ListRequest with lender-specific filters)
export interface LenderListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface LenderFormErrors {
  name?: string;
  email?: string;
  phone?: string;
  address?: string;
  gender?: string;
  general?: string;
}

// Lender statistics (if provided by backend)
export interface LenderStats {
  totallenders: number;
  // Add your custom stats here
}
