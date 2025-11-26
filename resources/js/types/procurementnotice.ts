// TypeScript interfaces for ProcurementNotice entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core ProcurementNotice interface matching the backend model
export interface ProcurementNotice extends BaseModel {
  procured_by: string;
  procurement_type: string;
  market_approach: string;
  invitation: string;
  ref_no: string;
  open_date: string;
  close_date: string;
  partners: string[];
  qualifying_districts: string[];
  is_published: boolean;
  organization: string;
  classification: string[];
  interested_smes: string[];
  details: string;
  application_details: string;
  minimum_qualifying_score: number;
}

// ProcurementNotice creation data (matches ProcurementNoticeCreateRequest)
export interface ProcurementNoticeCreateData {
  procured_by: string;
  procurement_type: string;
  market_approach: string;
  invitation: string;
  ref_no: string;
  open_date: string;
  close_date: string;
  partners: string[];
  qualifying_districts: string[];
  is_published: boolean;
  organization: string;
  classification: string[];
  interested_smes: string[];
  details: string;
  application_details: string;
  minimum_qualifying_score: number;
}

// ProcurementNotice update data (matches ProcurementNoticeUpdateRequest - all optional)
export interface ProcurementNoticeUpdateData extends Partial<ProcurementNoticeCreateData> { }

// ProcurementNotice list response (matches service GetList response)
export interface ProcurementNoticeListResponse extends PaginatedResult<ProcurementNotice> { }

// ProcurementNotice list request (extends base ListRequest with procurementnotice-specific filters)
export interface ProcurementNoticeListRequest extends ListRequest {
  procured_by?: string;
  procurement_type?: string;
  market_approach?: string;
  invitation?: string;
  ref_no?: string;
  organization?: string;
  is_published?: boolean;
}

// Form validation types
export interface ProcurementNoticeFormErrors {
  procured_by?: string;
  procurement_type?: string;
  market_approach?: string;
  invitation?: string;
  ref_no?: string;
  open_date?: string;
  close_date?: string;
  partners?: string;
  qualifying_districts?: string;
  is_published?: string;
  organization?: string;
  classification?: string;
  interested_smes?: string;
  details?: string;
  application_details?: string;
  minimum_qualifying_score?: string;
  general?: string;
}

// ProcurementNotice statistics (if provided by backend)
export interface ProcurementNoticeStats {
  total_procurement_notices: number;
  published_notices: number;
  // Add your custom stats here
}
