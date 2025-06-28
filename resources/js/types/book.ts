// TypeScript interfaces for Book entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core Book interface matching the backend model
export interface Book extends BaseModel {
  title: string;
  author: string;
  isbn: string;
  description?: string;
  price: number;
  status: BookStatus;
  publishedAt?: string;
  tags?: string[];
  // Additional computed fields that might come from the backend
  isAvailable?: boolean;
  borrowedBy?: string;
  borrowedAt?: string;
  dueDate?: string;
}

// Book status enum matching backend validation
export type BookStatus = 'AVAILABLE' | 'BORROWED' | 'MAINTENANCE';

// Book creation data (matches BookCreateRequest)
export interface BookCreateData {
  title: string;
  author: string;
  isbn: string;
  description?: string;
  price: number;
  status?: BookStatus;
  publishedAt?: string;
  tags?: string[];
}

// Book update data (matches BookUpdateRequest - all optional)
export interface BookUpdateData {
  title?: string;
  author?: string;
  isbn?: string;
  description?: string;
  price?: number;
  status?: BookStatus;
  publishedAt?: string;
  tags?: string[];
}

// Book list response (matches service GetList response)
export interface BookListResponse extends PaginatedResult<Book> {
  // Additional book-specific metadata if needed
  statusCounts?: {
    available: number;
    borrowed: number;
    maintenance: number;
  };
  totalValue?: number;
}

// Book list request (extends base ListRequest with book-specific filters)
export interface BookListRequest extends ListRequest {
  // Basic filters
  status?: BookStatus;
  author?: string;
  isbn?: string;
  
  // Advanced filters (matches GetListAdvanced)
  minPrice?: number;
  maxPrice?: number;
  publishedAfter?: string;
  publishedBefore?: string;
  tags?: string[];
  
  // Special filters
  isAvailable?: boolean;
  borrowedBy?: string;
}

// Book operations data
export interface BookBorrowData {
  borrowerId?: string;
  dueDate?: string;
  notes?: string;
}

export interface BookReturnData {
  returnDate?: string;
  condition?: 'GOOD' | 'FAIR' | 'DAMAGED';
  notes?: string;
}

// Form validation types
export interface BookFormErrors {
  title?: string;
  author?: string;
  isbn?: string;
  description?: string;
  price?: string;
  status?: string;
  publishedAt?: string;
  tags?: string;
  general?: string;
}

// Book statistics (if provided by backend)
export interface BookStats {
  totalBooks: number;
  availableBooks: number;
  borrowedBooks: number;
  maintenanceBooks: number;
  totalValue: number;
  averagePrice: number;
  topAuthors: Array<{
    name: string;
    count: number;
  }>;
  recentlyAdded: Book[];
  popularBooks: Array<{
    book: Book;
    borrowCount: number;
  }>;
}

// Advanced search/filter options
export interface BookAdvancedFilters {
  // Text filters
  title?: string;
  author?: string;
  isbn?: string;
  description?: string;
  
  // Status filter
  status?: BookStatus[];
  
  // Price range
  priceRange?: {
    min?: number;
    max?: number;
  };
  
  // Date range
  publishedDate?: {
    from?: string;
    to?: string;
  };
  
  // Date range for created
  createdDate?: {
    from?: string;
    to?: string;
  };
  
  // Tags
  tags?: string[];
  tagMatch?: 'any' | 'all'; // Match any tag or all tags
  
  // Availability
  onlyAvailable?: boolean;
  onlyBorrowed?: boolean;
  onlyMaintenance?: boolean;
}

// Book import/export types
export interface BookImportData {
  file: File;
  format: 'csv' | 'json' | 'excel';
  skipErrors?: boolean;
  updateExisting?: boolean;
}

export interface BookExportOptions {
  format: 'csv' | 'json' | 'excel' | 'pdf';
  fields?: string[];
  filters?: BookListRequest;
  includeStats?: boolean;
}

// Bulk operations
export interface BookBulkOperation {
  action: 'delete' | 'updateStatus' | 'updatePrice' | 'addTags' | 'removeTags' | 'export';
  bookIds: number[];
  data?: {
    status?: BookStatus;
    price?: number;
    tags?: string[];
    reason?: string;
  };
}

// Book borrowing history
export interface BookBorrowHistory {
  id: number;
  bookId: number;
  book?: Book;
  borrowerId: string;
  borrowerName: string;
  borrowedAt: string;
  dueDate?: string;
  returnedAt?: string;
  status: 'BORROWED' | 'RETURNED' | 'OVERDUE';
  notes?: string;
}

// Response types for specific endpoints
export interface BookByISBNResponse {
  book: Book | null;
  found: boolean;
}

export interface BooksByAuthorResponse extends PaginatedResult<Book> {
  author: string;
  totalByAuthor: number;
}

export interface BookAvailableResponse extends PaginatedResult<Book> {
  totalAvailable: number;
}

// Form props for Book components
export interface BookFormProps {
  book?: Book;
  onSubmit: (data: BookCreateData | BookUpdateData) => Promise<void>;
  onCancel: () => void;
  isLoading?: boolean;
  errors?: BookFormErrors;
  mode: 'create' | 'edit';
}

// Book detail view props
export interface BookDetailProps {
  book: Book;
  onEdit?: () => void;
  onDelete?: () => void;
  onBorrow?: (data: BookBorrowData) => Promise<void>;
  onReturn?: (data: BookReturnData) => Promise<void>;
  onClose: () => void;
  canEdit?: boolean;
  canDelete?: boolean;
  canBorrow?: boolean;
  canReturn?: boolean;
  borrowHistory?: BookBorrowHistory[];
}

// Book list props (extends CRUD props)
export interface BookListProps {
  data: BookListResponse;
  filters: BookListRequest;
  stats?: BookStats;
  canCreate?: boolean;
  canEdit?: boolean;
  canDelete?: boolean;
  canBorrow?: boolean;
  canManageLibrary?: boolean;
  onRefresh?: () => void;
  onImport?: (data: BookImportData) => Promise<void>;
  onExport?: (options: BookExportOptions) => Promise<void>;
  onBulkOperation?: (operation: BookBulkOperation) => Promise<void>;
}

// API response types (matching backend responses)
export interface BookApiResponse<T = Book> {
  data: T;
  message?: string;
  status: 'success' | 'error';
}

export interface BookListApiResponse {
  data: Book[];
  pagination: {
    current_page: number;
    last_page: number;
    per_page: number;
    total: number;
    from?: number;
    to?: number;
  };
  filters: BookListRequest;
  stats?: BookStats;
}

// Error types
export interface BookError {
  code: string;
  message: string;
  field?: string;
  details?: Record<string, any>;
}

// Book validation rules (for frontend validation)
export interface BookValidationRules {
  title: {
    required: boolean;
    maxLength: number;
  };
  author: {
    required: boolean;
    maxLength: number;
  };
  isbn: {
    required: boolean;
    pattern: RegExp;
    unique?: boolean;
  };
  price: {
    required: boolean;
    min: number;
    type: 'number';
  };
  status: {
    values: BookStatus[];
  };
  description: {
    maxLength: number;
  };
  tags: {
    maxItems: number;
    maxItemLength: number;
  };
}

// Default validation rules matching backend
export const BOOK_VALIDATION_RULES: BookValidationRules = {
  title: {
    required: true,
    maxLength: 255,
  },
  author: {
    required: true,
    maxLength: 100,
  },
  isbn: {
    required: true,
    pattern: /^[\d-]{10,17}$/,
    unique: true,
  },
  price: {
    required: true,
    min: 0,
    type: 'number',
  },
  status: {
    values: ['AVAILABLE', 'BORROWED', 'MAINTENANCE'],
  },
  description: {
    maxLength: 1000,
  },
  tags: {
    maxItems: 10,
    maxItemLength: 50,
  },
};

// Book status display configuration
export const BOOK_STATUS_CONFIG = {
  AVAILABLE: {
    label: 'Available',
    color: 'green',
    icon: '✓',
    description: 'Book is available for borrowing',
  },
  BORROWED: {
    label: 'Borrowed',
    color: 'blue',
    icon: '📖',
    description: 'Book is currently borrowed',
  },
  MAINTENANCE: {
    label: 'Maintenance',
    color: 'orange',
    icon: '🔧',
    description: 'Book is under maintenance',
  },
} as const;