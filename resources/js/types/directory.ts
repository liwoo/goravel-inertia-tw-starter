export interface DirectorySME {
  usme_number: string;
  name: string;
  sector: string;
  sub_sector?: string;
  classification: string;
  district: string;
  region: string;
  contact_phone?: string;
  contact_email?: string;
  website?: string;
  business_description?: string;
  operational_start_date?: string;
}

export interface FilterOption {
  value: string;
  label: string;
  count: number;
}

export interface DistrictGroup {
  region: string;
  districts: FilterOption[];
}

export interface DirectoryFilters {
  sectors: FilterOption[];
  districts: DistrictGroup[];
}

export interface DirectoryPagination {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface DirectoryPageProps {
  smes: DirectorySME[];
  filters: DirectoryFilters;
  pagination: DirectoryPagination;
  userFormalisationScore: number;
  minScoreForContactView: number;
}
