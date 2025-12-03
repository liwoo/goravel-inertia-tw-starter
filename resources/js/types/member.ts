import { BaseModel, PaginatedResult, ListRequest } from './crud';

export interface Member extends BaseModel {
  first_name: string;
  last_name: string;
  other_names?: string;
  gender?: string;
  nationality: string;
  national_id_number: string;
  date_of_birth?: string;
  email?: string;
  phone_number: string;
  is_intern: boolean;
  is_part_time: boolean;
  sme_id: number;

  name?: string;
}

export interface MemberCreateData {
  first_name: string;
  last_name: string;
  other_names?: string;
  gender?: string;
  nationality: string;
  national_id_number: string;
  date_of_birth?: string;
  email?: string;
  phone_number: string;
  is_intern: boolean;
  is_part_time: boolean;
  sme_id: number;
}

export interface MemberUpdateData {
  first_name?: string;
  last_name?: string;
  other_names?: string;
  gender?: string;
  nationality?: string;
  national_id_number?: string;
  date_of_birth?: string;
  email?: string;
  phone_number?: string;
  is_intern?: boolean;
  is_part_time?: boolean;
}

export interface MemberListResponse extends PaginatedResult<Member> { }

export interface MemberListRequest extends ListRequest {
}

// Form validation types
export interface MemberFormErrors {
  first_name?: string;
  last_name?: string;
  email?: string;
  phone_number?: string;
  national_id_number?: string;
  general?: string;
}

export interface MemberStats {
  totalmembers: number;
}
