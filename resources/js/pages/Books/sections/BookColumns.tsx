import React from 'react';
import { Book, BookStatus } from '@/types/book';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { StatusBadge } from '@/components/Crud/StatusBadge';
import { Badge } from '@/components/ui/badge';
import { BookOpen, Calendar, DollarSign, Hash, User, Tag, CheckCircle, Clock, Wrench } from 'lucide-react';

// Book status configuration with improved theming
const BOOK_STATUS_CONFIG = {
  AVAILABLE: { 
    label: 'Available', 
    icon: <CheckCircle className="h-3 w-3" />,
    color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
  },
  BORROWED: { 
    label: 'Borrowed', 
    icon: <Clock className="h-3 w-3" />,
    color: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
  },
  MAINTENANCE: { 
    label: 'Maintenance', 
    icon: <Wrench className="h-3 w-3" />,
    color: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400'
  },
};

/**
 * Book table columns configuration with improved theming
 */
export const bookColumns: CrudColumn<Book>[] = [
  {
    key: 'title',
    label: 'Book Details',
    sortable: true,
    className: 'min-w-[300px]',
    render: (book) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <BookOpen className="h-5 w-5 text-muted-foreground" />
        </div>
        <div className="space-y-1">
          <div className="font-medium text-foreground">{book.title}</div>
          <div className="text-sm text-muted-foreground flex items-center">
            <User className="w-3 h-3 mr-1" />
            {book.author}
          </div>
          <div className="text-xs text-muted-foreground flex items-center">
            <Hash className="w-3 h-3 mr-1" />
             {book.isbn}
          </div>
        </div>
      </div>
    ),
  },
  {
    key: 'status',
    label: 'Status',
    sortable: true,
    className: 'w-32',
    render: (book) => {
      const config = BOOK_STATUS_CONFIG[book.status as keyof typeof BOOK_STATUS_CONFIG];
      return (
        <Badge className={`${config.color} flex items-center gap-1`}>
          {config.icon}
          {config.label}
        </Badge>
      );
    },
  },
  {
    key: 'price',
    label: 'Price',
    sortable: true,
    className: 'w-24 text-right',
    render: (book) => (
      <div className="text-right">
        <div className="font-medium text-foreground flex items-center justify-end">
          <DollarSign className="w-3 h-3 mr-1 text-green-500 dark:text-green-400" />
          {new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
          }).format(book.price)}
        </div>
      </div>
    ),
  },
  {
    key: 'publishedAt',
    label: 'Published',
    sortable: true,
    className: 'w-32',
    render: (book) => (
      <div className="text-sm">
        {book.publishedAt ? (
          <div className="flex items-center text-muted-foreground">
            <Calendar className="w-3 h-3 mr-1" />
            {new Date(book.publishedAt).toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: 'numeric'
            })}
          </div>
        ) : (
          <span className="text-muted-foreground">-</span>
        )}
      </div>
    ),
  },
  {
    key: 'tags',
    label: 'Tags',
    className: 'w-48',
    render: (book) => (
      <div className="flex flex-wrap gap-1">
        {book.tags && book.tags.length > 0 ? (
          <>
            {book.tags.slice(0, 3).map((tag, index) => (
              <Badge 
                key={index} 
                variant="secondary" 
                className="text-xs bg-secondary/50 dark:bg-secondary/30"
              >
                {tag}
              </Badge>
            ))}
            {book.tags.length > 3 && (
              <Badge variant="secondary" className="text-xs bg-secondary/50 dark:bg-secondary/30">
                +{book.tags.length - 3}
              </Badge>
            )}
          </>
        ) : (
          <span className="text-muted-foreground text-sm">No tags</span>
        )}
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Added',
    sortable: true,
    className: 'w-28',
    render: (book) => (
      <div className="text-sm text-muted-foreground">
        {new Date(book.createdAt).toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric'
        })}
      </div>
    ),
  },
];

/**
 * Compact book columns for mobile/smaller screens
 */
export const bookColumnsMobile: CrudColumn<Book>[] = [
  {
    key: 'title',
    label: 'Book',
    sortable: true,
    render: (book) => (
      <div className="space-y-3">
        <div className="flex items-start gap-3">
          <div className="p-2 rounded-lg bg-muted">
            <BookOpen className="h-5 w-5 text-muted-foreground" />
          </div>
          <div className="flex-1 space-y-1">
            <div className="font-medium text-foreground">{book.title}</div>
            <div className="text-sm text-muted-foreground">by {book.author}</div>
          </div>
        </div>
        <div className="flex items-center justify-between pl-12">
          <div>
            {(() => {
              const config = BOOK_STATUS_CONFIG[book.status as keyof typeof BOOK_STATUS_CONFIG];
              return (
                <Badge className={`${config.color} flex items-center gap-1 text-xs`}>
                  {config.icon}
                  {config.label}
                </Badge>
              );
            })()}
          </div>
          <div className="text-sm font-medium text-foreground">
            {new Intl.NumberFormat('en-US', {
              style: 'currency',
              currency: 'USD',
            }).format(book.price)}
          </div>
        </div>
        {book.tags && book.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 pl-12">
            {book.tags.slice(0, 2).map((tag, index) => (
              <Badge key={index} variant="secondary" className="text-xs bg-secondary/50 dark:bg-secondary/30">
                {tag}
              </Badge>
            ))}
            {book.tags.length > 2 && (
              <Badge variant="secondary" className="text-xs bg-secondary/50 dark:bg-secondary/30">
                +{book.tags.length - 2}
              </Badge>
            )}
          </div>
        )}
      </div>
    ),
  },
];

/**
 * Book filters configuration
 */
export const bookFilters: CrudFilter[] = [
  {
    key: 'status',
    label: 'Status',
    type: 'select',
    options: [
      { value: '', label: 'All Status' },
      { value: 'AVAILABLE', label: 'Available' },
      { value: 'BORROWED', label: 'Borrowed' },
      { value: 'MAINTENANCE', label: 'Maintenance' },
    ],
  },
  {
    key: 'author',
    label: 'Author',
    type: 'text',
    placeholder: 'Filter by author name',
  },
  {
    key: 'minPrice',
    label: 'Min Price',
    type: 'number',
    placeholder: '0.00',
  },
  {
    key: 'maxPrice',
    label: 'Max Price',
    type: 'number',
    placeholder: '100.00',
  },
  {
    key: 'publishedAfter',
    label: 'Published After',
    type: 'date',
  },
  {
    key: 'publishedBefore',
    label: 'Published Before',
    type: 'date',
  },
];

/**
 * Quick filter buttons for common book queries
 */
export const bookQuickFilters = [
  {
    key: 'all',
    label: 'All Books',
    icon: <BookOpen className="h-4 w-4" />,
    filters: {},
  },
  {
    key: 'available',
    label: 'Available',
    icon: <CheckCircle className="h-4 w-4 text-green-500" />,
    filters: { status: 'AVAILABLE' },
  },
  {
    key: 'borrowed',
    label: 'Borrowed',
    icon: <Clock className="h-4 w-4 text-blue-500" />,
    filters: { status: 'BORROWED' },
  },
  {
    key: 'maintenance',
    label: 'Maintenance',
    icon: <Wrench className="h-4 w-4 text-orange-500" />,
    filters: { status: 'MAINTENANCE' },
  },
  {
    key: 'recent',
    label: 'Recently Added',
    icon: <Calendar className="h-4 w-4" />,
    filters: { 
      createdAfter: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    },
  },
] as const;