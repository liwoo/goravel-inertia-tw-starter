// TypeScript interfaces for Author entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Author interface matching the backend model
export interface Author extends BaseModel {
  firstName: string;
  lastName: string;
  bio?: string;
  email?: string;
  website?: string;
  birthDate?: string;
  nationality?: string;
  photoUrl?: string;
  status: AuthorStatus;
}

// Author status enum matching backend validation
export type AuthorStatus = 'ACTIVE' | 'INACTIVE';

// Author creation data (matches AuthorCreateRequest)
export interface AuthorCreateData {
  firstName: string;
  lastName: string;
  bio?: string;
  email?: string;
  website?: string;
  birthDate?: string;
  nationality?: string;
  photoUrl?: string;
  status?: AuthorStatus;
}

// Author update data (matches AuthorUpdateRequest - all optional)
export interface AuthorUpdateData {
  firstName?: string;
  lastName?: string;
  bio?: string;
  email?: string;
  website?: string;
  birthDate?: string;
  nationality?: string;
  photoUrl?: string;
  status?: AuthorStatus;
}

// Author list response (matches service GetList response)
export interface AuthorListResponse extends PaginatedResult<Author> {
  statusCounts?: {
    active: number;
    inactive: number;
  };
}

// Author list request (extends base ListRequest with author-specific filters)
export interface AuthorListRequest extends ListRequest {
  status?: AuthorStatus;
  nationality?: string;
}

// Form validation types
export interface AuthorFormErrors {
  firstName?: string;
  lastName?: string;
  bio?: string;
  email?: string;
  website?: string;
  birthDate?: string;
  nationality?: string;
  photoUrl?: string;
  status?: string;
  general?: string;
}

// Author statistics (matches GetAuthorStatistics backend response)
export interface AuthorStats {
  totalAuthors: number;
  activeAuthors: number;
  inactiveAuthors: number;
}

// Author validation rules (for frontend validation)
export interface AuthorValidationRules {
  firstName: {
    required: boolean;
    maxLength: number;
  };
  lastName: {
    required: boolean;
    maxLength: number;
  };
  email: {
    maxLength: number;
  };
  website: {
    maxLength: number;
  };
  nationality: {
    maxLength: number;
  };
  photoUrl: {
    maxLength: number;
  };
  status: {
    values: AuthorStatus[];
  };
  bio: {
    maxLength: number;
  };
}

// Default validation rules matching backend
export const AUTHOR_VALIDATION_RULES: AuthorValidationRules = {
  firstName: {
    required: true,
    maxLength: 100,
  },
  lastName: {
    required: true,
    maxLength: 100,
  },
  email: {
    maxLength: 255,
  },
  website: {
    maxLength: 255,
  },
  nationality: {
    maxLength: 100,
  },
  photoUrl: {
    maxLength: 500,
  },
  status: {
    values: ['ACTIVE', 'INACTIVE'],
  },
  bio: {
    maxLength: 2000,
  },
};

// Author status display configuration
export const AUTHOR_STATUS_CONFIG = {
  ACTIVE: {
    label: 'Active',
    color: 'green',
    description: 'Author is active and visible',
  },
  INACTIVE: {
    label: 'Inactive',
    color: 'gray',
    description: 'Author is inactive and hidden',
  },
} as const;
