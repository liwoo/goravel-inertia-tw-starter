import React from 'react';
import { TFunction } from 'i18next';
import { Book, BookStatus } from '@/types/book';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { StatusBadge } from '@/components/Crud/StatusBadge';
import { Badge } from '@/components/ui/badge';
import { BookOpen, Calendar, DollarSign, Hash, User, Tag, CheckCircle, Clock, Wrench } from 'lucide-react';

// Book status configuration with improved theming
function getBookStatusConfig(t: TFunction) {
  return {
    AVAILABLE: {
      label: t('status.available'),
      icon: <CheckCircle className="h-3 w-3" />,
      color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
    },
    BORROWED: {
      label: t('status.borrowed'),
      icon: <Clock className="h-3 w-3" />,
      color: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
    },
    MAINTENANCE: {
      label: t('status.maintenance'),
      icon: <Wrench className="h-3 w-3" />,
      color: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400'
    },
    RESERVED: {
      label: t('status.reserved'),
      icon: <BookOpen className="h-3 w-3" />,
      color: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400'
    },
  };
}

/**
 * Book table columns configuration with improved theming
 */
export function getBookColumns(t: TFunction): CrudColumn<Book>[] {
  const BOOK_STATUS_CONFIG = getBookStatusConfig(t);

  return [
    {
      key: 'title',
      label: t('columns.bookDetails'),
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
      label: t('columns.status'),
      sortable: true,
      className: 'w-32',
      render: (book) => {
        const config = BOOK_STATUS_CONFIG[book.status as keyof typeof BOOK_STATUS_CONFIG];
        if (!config) {
          return <Badge variant="outline">{t('status.unknown')}</Badge>;
        }
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
      label: t('columns.price'),
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
      label: t('columns.published'),
      sortable: true,
      className: 'w-32',
      render: (book) => {
        const publishDate = book.publishedAt || book.published_at;
        return (
          <div className="text-sm">
            {publishDate ? (
              <div className="flex items-center text-muted-foreground">
                <Calendar className="w-3 h-3 mr-1" />
                {new Date(publishDate).toLocaleDateString('en-US', {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric'
                })}
              </div>
            ) : (
              <span className="text-muted-foreground">-</span>
            )}
          </div>
        );
      },
    },
    {
      key: 'tags',
      label: t('columns.tags'),
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
            <span className="text-muted-foreground text-sm">{t('columns.noTags')}</span>
          )}
        </div>
      ),
    },
    {
      key: 'createdAt',
      label: t('columns.added'),
      sortable: true,
      className: 'w-28',
      render: (book) => {
        // Try multiple possible date field locations
        const dateValue = book.createdAt ||
                         (book as any).created_at ||
                         (book as any).CreatedAt ||
                         (book as any).Created_at;

        if (!dateValue) {
          // If no date field found, show a dash
          return <div className="text-sm text-muted-foreground">-</div>;
        }

        // Parse the date - handle various formats
        let date: Date;

        // If it's already a valid date object
        if (dateValue instanceof Date) {
          date = dateValue;
        }
        // If it's a timestamp object with StdTime
        else if (dateValue && typeof dateValue === 'object' && dateValue.StdTime) {
          date = new Date(dateValue.StdTime);
        }
        // If it's a string or number
        else {
          date = new Date(dateValue);
        }

        const isValidDate = !isNaN(date.getTime());

        // Remove debug log for production
        if (!isValidDate && process.env.NODE_ENV === 'development') {
          console.warn('Invalid date format for book:', book.id, 'Date value:', dateValue);
        }

        return (
          <div className="text-sm text-muted-foreground">
            {isValidDate ? date.toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: 'numeric'
            }) : '-'}
          </div>
        );
      },
    },
  ];
}

/**
 * Compact book columns for mobile/smaller screens
 */
export function getBookColumnsMobile(t: TFunction): CrudColumn<Book>[] {
  const BOOK_STATUS_CONFIG = getBookStatusConfig(t);

  return [
    {
      key: 'title',
      label: t('columns.book'),
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
                if (!config) {
                  return <Badge variant="outline" className="text-xs">{t('status.unknown')}</Badge>;
                }
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
}

/**
 * Book filters configuration
 */
export function getBookFilters(t: TFunction): CrudFilter[] {
  return [
    {
      key: 'status',
      label: t('filters.status'),
      type: 'select',
      options: [
        { value: '__all__', label: t('filters.allStatus') },
        { value: 'AVAILABLE', label: t('status.available') },
        { value: 'BORROWED', label: t('status.borrowed') },
        { value: 'MAINTENANCE', label: t('status.maintenance') },
      ],
    },
    {
      key: 'author',
      label: t('filters.author'),
      type: 'text',
      placeholder: t('filters.filterByAuthor'),
    },
    {
      key: 'minPrice',
      label: t('filters.minPrice'),
      type: 'number',
      placeholder: '0.00',
    },
    {
      key: 'maxPrice',
      label: t('filters.maxPrice'),
      type: 'number',
      placeholder: '100.00',
    },
    {
      key: 'publishedAfter',
      label: t('filters.publishedAfter'),
      type: 'date',
    },
    {
      key: 'publishedBefore',
      label: t('filters.publishedBefore'),
      type: 'date',
    },
  ];
}

/**
 * Quick filter buttons for common book queries
 */
export function getBookQuickFilters(t: TFunction) {
  return [
    {
      key: 'all',
      label: t('filters.allBooks'),
      icon: <BookOpen className="h-4 w-4" />,
      filters: {},
    },
    {
      key: 'available',
      label: t('status.available'),
      icon: <CheckCircle className="h-4 w-4 text-green-500" />,
      filters: { status: 'AVAILABLE' },
    },
    {
      key: 'borrowed',
      label: t('status.borrowed'),
      icon: <Clock className="h-4 w-4 text-blue-500" />,
      filters: { status: 'BORROWED' },
    },
    {
      key: 'maintenance',
      label: t('status.maintenance'),
      icon: <Wrench className="h-4 w-4 text-orange-500" />,
      filters: { status: 'MAINTENANCE' },
    },
    {
      key: 'recent',
      label: t('filters.recentlyAdded'),
      icon: <Calendar className="h-4 w-4" />,
      filters: {
        createdAfter: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
      },
    },
  ] as const;
}