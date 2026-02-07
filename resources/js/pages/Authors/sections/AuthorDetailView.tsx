import React from 'react';
import { useTranslation } from 'react-i18next';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Author } from '@/types/author';
import { getAuthorStatusConfig } from './authorStatusConfig';

export function AuthorDetailView({
  item: author,
  onEdit,
  onClose,
  canEdit,
}: CrudDetailViewProps<Author>) {
  const { t } = useTranslation('authors');

  const statusConfig = getAuthorStatusConfig(t);

  const formatDate = (date: string | Date | null | undefined) => {
    if (!date) return t('form.notSpecified');
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getStatusBadge = (status: string) => {
    const config = statusConfig[status as keyof typeof statusConfig];
    const fallback = {
      color:
        'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400',
      icon: null,
      label: t('status.unknown'),
    };
    const cfg = config || fallback;

    return (
      <Badge className={`${cfg.color} flex items-center gap-1`}>
        {cfg.icon}
        {cfg.label}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      {/* Author Information Section */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-foreground">
          {t('form.authorInfo')}
        </h3>
        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.firstName').replace(' *', '')}
              </p>
              <p className="font-medium text-foreground">{author.firstName}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.lastName').replace(' *', '')}
              </p>
              <p className="font-medium text-foreground">{author.lastName}</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.email')}
              </p>
              <p className="font-medium text-foreground">
                {author.email || t('form.notSpecified')}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.website')}
              </p>
              <p className="font-medium text-foreground">
                {author.website ? (
                  <a
                    href={author.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline"
                  >
                    {author.website}
                  </a>
                ) : (
                  t('form.notSpecified')
                )}
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.birthDate')}
              </p>
              <p className="font-medium text-foreground">
                {formatDate(author.birthDate)}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.nationality')}
              </p>
              <p className="font-medium text-foreground">
                {author.nationality || t('form.notSpecified')}
              </p>
            </div>
          </div>

          <div>
            <p className="text-sm text-muted-foreground">
              {t('form.status')}
            </p>
            <div className="mt-1">{getStatusBadge(author.status)}</div>
          </div>

          {author.photoUrl && (
            <div>
              <p className="text-sm text-muted-foreground">
                {t('form.photoUrl')}
              </p>
              <a
                href={author.photoUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-primary hover:underline"
              >
                {author.photoUrl}
              </a>
            </div>
          )}
        </div>
      </div>

      {author.bio && (
        <>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold mb-4 text-foreground">
              {t('form.bio')}
            </h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              {author.bio}
            </p>
          </div>
        </>
      )}

      <Separator />

      {/* Metadata Section */}
      <div>
        <h3 className="text-lg font-semibold mb-4 text-foreground">
          {t('form.metadata')}
        </h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">
              {t('form.authorId')}
            </p>
            <p className="font-medium text-sm text-foreground">
              #{author.id}
            </p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">
              {t('form.created')}
            </p>
            <p className="font-medium text-sm text-foreground">
              {formatDate(author.createdAt || author.created_at)}
            </p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">
              {t('form.lastUpdated')}
            </p>
            <p className="font-medium text-sm text-foreground">
              {formatDate(author.updatedAt || author.updated_at)}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
