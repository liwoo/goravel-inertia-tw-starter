import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Calendar,
  DollarSign,
  Hash,
  Tag,
  User as UserIcon,
  BookOpen,
  Edit,
  Trash2,
  Clock,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Book } from '@/types/book';

export function BookDetailView({
  item: book,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<Book>) {
  const { t } = useTranslation('books');

  const formatDate = (date: string | Date | null | undefined) => {
    if (!date) return t('form.notSpecified');
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { color: string; icon: React.ReactNode; label: string }> = {
      'AVAILABLE': { color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400', icon: <CheckCircle className="h-3 w-3" />, label: t('status.available') },
      'BORROWED': { color: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400', icon: <Clock className="h-3 w-3" />, label: t('status.borrowed') },
      'MAINTENANCE': { color: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400', icon: <XCircle className="h-3 w-3" />, label: t('status.maintenance') },
      'RESERVED': { color: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400', icon: <BookOpen className="h-3 w-3" />, label: t('status.reserved') },
    };

    const config = statusConfig[status] || { color: 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400', icon: null, label: t('status.unknown') };

    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        {config.icon}
        {config.label}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">

      {/* Book Information Section */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.bookInfo')}</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.title').replace(' *', '')}</p>
                <p className="font-medium text-foreground">{book.title}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.author').replace(' *', '')}</p>
                <p className="font-medium text-foreground">{book.author}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Hash className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.isbn').replace(' *', '')}</p>
                <p className="font-medium text-foreground">{book.isbn}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.price')}</p>
                <p className="font-medium text-foreground">{formatCurrency(book.price)}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.publishedDate')}</p>
                <p className="font-medium text-foreground">{formatDate(book.publishedAt)}</p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Clock className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">{t('form.status')}</p>
                <div className="mt-1">
                  {getStatusBadge(book.status)}
                </div>
              </div>
            </div>

            {book.tags && book.tags.length > 0 && (
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Tag className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm text-muted-foreground">{t('form.tags')}</p>
                  <div className="flex flex-wrap gap-2 mt-1">
                    {book.tags.map((tag, index) => (
                      <Badge key={index} variant="secondary">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

      {book.description && (
        <>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.description')}</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              {book.description}
            </p>
          </div>
        </>
      )}

      <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.metadata')}</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">{t('form.created')}</p>
              <p className="font-medium text-sm text-foreground">{formatDate(book.createdAt || book.created_at)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t('form.lastUpdated')}</p>
              <p className="font-medium text-sm text-foreground">{formatDate(book.updatedAt || book.updated_at)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t('form.bookId')}</p>
              <p className="font-medium text-sm text-foreground">#{book.id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t('form.inStock')}</p>
              <p className="font-medium text-sm text-foreground">
                {book.status === 'AVAILABLE' ? t('form.yes') : t('form.no')}
              </p>
            </div>
          </div>
        </div>
      </div>

    </div>
  );
}