import React from 'react';
import { TFunction } from 'i18next';
import { Author } from '@/types/author';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { User, CheckCircle, XCircle, Calendar } from 'lucide-react';
import { getAuthorStatusConfig } from './authorStatusConfig';
import { NATIONALITIES } from '@/config/options';

/**
 * Author table columns configuration
 * Desktop: 4 columns — Author (composite), Nationality, Status, Added
 */
export function getAuthorColumns(t: TFunction): CrudColumn<Author>[] {
  const AUTHOR_STATUS_CONFIG = getAuthorStatusConfig(t);

  return [
    {
      key: 'firstName',
      label: t('columns.authorDetails'),
      sortable: true,
      className: 'min-w-[280px]',
      render: (author) => (
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-muted">
            <User className="h-4 w-4 text-muted-foreground" />
          </div>
          <div>
            <div className="font-medium text-foreground">
              {author.firstName} {author.lastName}
            </div>
            {author.email && (
              <div className="text-sm text-muted-foreground">
                {author.email}
              </div>
            )}
          </div>
        </div>
      ),
    },
    {
      key: 'nationality',
      label: t('columns.nationality'),
      sortable: true,
      className: 'w-36',
      render: (author) => (
        <div className="text-sm text-foreground">
          {author.nationality || (
            <span className="text-muted-foreground">
              {t('columns.noNationality')}
            </span>
          )}
        </div>
      ),
    },
    {
      key: 'status',
      label: t('columns.status'),
      sortable: true,
      className: 'w-32',
      render: (author) => {
        const config =
          AUTHOR_STATUS_CONFIG[
            author.status as keyof typeof AUTHOR_STATUS_CONFIG
          ];
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
      key: 'createdAt',
      label: t('columns.added'),
      sortable: true,
      className: 'w-28',
      render: (author) => {
        const dateValue =
          author.createdAt ||
          (author as any).created_at ||
          (author as any).CreatedAt;

        if (!dateValue) {
          return <div className="text-sm text-muted-foreground">-</div>;
        }

        let date: Date;

        if (dateValue instanceof Date) {
          date = dateValue;
        } else if (
          dateValue &&
          typeof dateValue === 'object' &&
          dateValue.StdTime
        ) {
          date = new Date(dateValue.StdTime);
        } else {
          date = new Date(dateValue);
        }

        const isValidDate = !isNaN(date.getTime());

        return (
          <div className="text-sm text-muted-foreground">
            {isValidDate
              ? date.toLocaleDateString('en-US', {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric',
                })
              : '-'}
          </div>
        );
      },
    },
  ];
}

/**
 * Compact author columns for mobile/smaller screens
 */
export function getAuthorColumnsMobile(t: TFunction): CrudColumn<Author>[] {
  const AUTHOR_STATUS_CONFIG = getAuthorStatusConfig(t);

  return [
    {
      key: 'firstName',
      label: t('columns.author'),
      sortable: true,
      render: (author) => (
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-muted">
              <User className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <div className="font-medium text-foreground">
                {author.firstName} {author.lastName}
              </div>
              {author.email && (
                <div className="text-sm text-muted-foreground">
                  {author.email}
                </div>
              )}
            </div>
          </div>
          {(() => {
            const config =
              AUTHOR_STATUS_CONFIG[
                author.status as keyof typeof AUTHOR_STATUS_CONFIG
              ];
            if (!config) {
              return (
                <Badge variant="outline" className="text-xs">
                  {t('status.unknown')}
                </Badge>
              );
            }
            return (
              <Badge
                className={`${config.color} flex items-center gap-1 text-xs`}
              >
                {config.icon}
                {config.label}
              </Badge>
            );
          })()}
        </div>
      ),
    },
  ];
}

/**
 * Author filters configuration
 */
export function getAuthorFilters(t: TFunction): CrudFilter[] {
  return [
    {
      key: 'status',
      label: t('filters.status'),
      type: 'select',
      options: [
        { value: '__all__', label: t('filters.allStatus') },
        { value: 'ACTIVE', label: t('status.active') },
        { value: 'INACTIVE', label: t('status.inactive') },
      ],
    },
    {
      key: 'nationality',
      label: t('filters.nationality'),
      type: 'select',
      options: [
        { value: '__all__', label: t('filters.allNationalities') },
        ...NATIONALITIES.map((n) => ({ value: n.value, label: n.label })),
      ],
    },
  ];
}

/**
 * Quick filter buttons for common author queries
 */
export function getAuthorQuickFilters(t: TFunction) {
  return [
    {
      key: 'all',
      label: t('filters.allAuthors'),
      icon: <User className="h-4 w-4" />,
      filters: {},
    },
    {
      key: 'active',
      label: t('status.active'),
      icon: <CheckCircle className="h-4 w-4 text-green-500" />,
      filters: { status: 'ACTIVE' },
    },
    {
      key: 'inactive',
      label: t('status.inactive'),
      icon: <XCircle className="h-4 w-4 text-gray-500" />,
      filters: { status: 'INACTIVE' },
    },
    {
      key: 'recent',
      label: t('filters.recentlyAdded'),
      icon: <Calendar className="h-4 w-4" />,
      filters: {
        createdAfter: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
          .toISOString()
          .split('T')[0],
      },
    },
  ] as const;
}
