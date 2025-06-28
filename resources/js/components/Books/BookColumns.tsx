import React from 'react';
import { Book, BookStatus, BOOK_STATUS_CONFIG } from '@/types/book';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { StatusBadge } from '@/components/Crud/StatusBadge';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { BookOpen, Calendar, DollarSign, Hash, User, Tag } from 'lucide-react';

/**
 * Book table columns configuration
 */
export const bookColumns: CrudColumn<Book>[] = [
  {
    key: 'title',
    label: 'Book Details',
    sortable: true,
    className: 'min-w-[250px]',
    render: (book) => (
      <div className="space-y-1">
        <div className="font-medium text-gray-900 flex items-center">
          📖
          {book.title}
        </div>
        <div className="text-sm text-gray-600 flex items-center">
          <User className="w-3 h-3 mr-1 text-gray-400" />
          by {book.author}
        </div>
        <div className="text-xs text-gray-500 flex items-center">
          <Hash className="w-3 h-3 mr-1 text-gray-400" />
          ISBN: {book.isbn}
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
      const config = BOOK_STATUS_CONFIG[book.status];
      return (
        <div className="flex items-center space-x-2">
          <StatusBadge 
            status={book.status}
            variant={
              book.status === 'AVAILABLE' ? 'success' :
              book.status === 'BORROWED' ? 'info' :
              'warning'
            }
          />
          {book.status === 'BORROWED' && book.dueDate && (
            <div className="text-xs text-gray-500">
              Due: {new Date(book.dueDate).toLocaleDateString()}
            </div>
          )}
        </div>
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
        <div className="font-medium text-gray-900 flex items-center justify-end">
          <DollarSign className="w-3 h-3 mr-1 text-green-500" />
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
          <div className="flex items-center text-gray-600">
            <Calendar className="w-3 h-3 mr-1 text-gray-400" />
            {new Date(book.publishedAt).toLocaleDateString()}
          </div>
        ) : (
          <span className="text-gray-400">-</span>
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
          book.tags.slice(0, 3).map((tag, index) => (
            <Badge 
              key={index} 
              variant="outline" 
              className="text-xs flex items-center"
            >
              <Tag className="w-2 h-2 mr-1" />
              {tag}
            </Badge>
          ))
        ) : (
          <span className="text-gray-400 text-sm">No tags</span>
        )}
        {book.tags && book.tags.length > 3 && (
          <Badge variant="outline" className="text-xs">
            +{book.tags.length - 3} more
          </Badge>
        )}
      </div>
    ),
  },
  {
    key: 'description',
    label: 'Description',
    className: 'min-w-[200px] max-w-[300px]',
    render: (book) => (
      <div className="text-sm text-gray-600">
        {book.description ? (
          <div className="line-clamp-2" title={book.description}>
            {book.description.length > 100 
              ? `${book.description.substring(0, 100)}...` 
              : book.description
            }
          </div>
        ) : (
          <span className="text-gray-400 italic">No description</span>
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
      <div className="text-sm text-gray-500">
        {new Date(book.createdAt).toLocaleDateString()}
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
      <div className="space-y-2">
        <div className="font-medium text-gray-900">{book.title}</div>
        <div className="text-sm text-gray-600">by {book.author}</div>
        <div className="flex items-center justify-between">
          <StatusBadge 
            status={book.status}
            variant={
              book.status === 'AVAILABLE' ? 'success' :
              book.status === 'BORROWED' ? 'info' :
              'warning'
            }
          />
          <div className="text-sm font-medium text-gray-900">
            {new Intl.NumberFormat('en-US', {
              style: 'currency',
              currency: 'USD',
            }).format(book.price)}
          </div>
        </div>
        {book.tags && book.tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {book.tags.slice(0, 2).map((tag, index) => (
              <Badge key={index} variant="outline" className="text-xs">
                {tag}
              </Badge>
            ))}
            {book.tags.length > 2 && (
              <Badge variant="outline" className="text-xs">
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
  {
    key: 'tags',
    label: 'Tags',
    type: 'text',
    placeholder: 'Filter by tags (comma-separated)',
  },
  {
    key: 'isAvailable',
    label: 'Available Only',
    type: 'boolean',
  },
];

/**
 * Advanced book filters for power users
 */
export const bookAdvancedFilters: CrudFilter[] = [
  ...bookFilters,
  {
    key: 'isbn',
    label: 'ISBN',
    type: 'text',
    placeholder: 'Search by ISBN',
  },
  {
    key: 'createdAfter',
    label: 'Added After',
    type: 'date',
  },
  {
    key: 'createdBefore',
    label: 'Added Before',
    type: 'date',
  },
  {
    key: 'borrowedBy',
    label: 'Borrowed By',
    type: 'text',
    placeholder: 'Filter by borrower',
  },
];

/**
 * Quick filter buttons for common book queries
 */
export const bookQuickFilters = [
  {
    key: 'all',
    label: 'All Books',
    icon: <BookOpen className="w-4 h-4" />,
    filters: {},
  },
  {
    key: 'available',
    label: 'Available',
    icon: <span className="w-4 h-4 bg-green-500 rounded-full inline-block" />,
    filters: { status: 'AVAILABLE' },
  },
  {
    key: 'borrowed',
    label: 'Borrowed',
    icon: <span className="w-4 h-4 bg-blue-500 rounded-full inline-block" />,
    filters: { status: 'BORROWED' },
  },
  {
    key: 'maintenance',
    label: 'Maintenance',
    icon: <span className="w-4 h-4 bg-orange-500 rounded-full inline-block" />,
    filters: { status: 'MAINTENANCE' },
  },
  {
    key: 'recent',
    label: 'Recently Added',
    icon: <Calendar className="w-4 h-4" />,
    filters: { 
      createdAfter: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    },
  },
] as const;