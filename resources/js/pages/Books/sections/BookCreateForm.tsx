import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { CrudFormProps } from '@/types/crud';
import { BookCreateData, BookStatus } from '@/types/book';
import { Calendar, DollarSign, Hash, Plus, Tag, X, BookOpen, User, Clock, FileText } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface BookCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
}

export const BookCreateForm = forwardRef<any, BookCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const { t } = useTranslation('books');
  const [formData, setFormData] = useState<BookCreateData>({
    title: '',
    author: '',
    isbn: '',
    description: '',
    price: 0,
    status: 'AVAILABLE',
    publishedAt: '',
    tags: [],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [tagInput, setTagInput] = useState('');

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.title.trim()) {
      newErrors.title = t('validation.titleRequired');
    }
    if (!formData.author.trim()) {
      newErrors.author = t('validation.authorRequired');
    }
    if (!formData.isbn.trim()) {
      newErrors.isbn = t('validation.isbnRequired');
    }
    if (formData.price < 0) {
      newErrors.price = t('validation.pricePositive');
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      // Send data as is - backend will handle field name transformation

      const response = await fetch('/api/books', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess(t('toast.created'));
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  const handleAddTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      setFormData({
        ...formData,
        tags: [...(formData.tags || []), tagInput.trim()]
      });
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData({
      ...formData,
      tags: formData.tags?.filter(tag => tag !== tagToRemove) || []
    });
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      {/* Book Information */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.bookInfo')}</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="title">{t('form.title')}</Label>
                <Input
                  id="title"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder={t('form.enterTitle')}
                  className={errors.title ? 'border-destructive' : ''}
                />
                {errors.title && (
                  <p className="text-sm text-destructive">{errors.title}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <User className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="author">{t('form.author')}</Label>
                <Input
                  id="author"
                  value={formData.author}
                  onChange={(e) => setFormData({ ...formData, author: e.target.value })}
                  placeholder={t('form.enterAuthor')}
                  className={errors.author ? 'border-destructive' : ''}
                />
                {errors.author && (
                  <p className="text-sm text-destructive">{errors.author}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Hash className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="isbn">{t('form.isbn')}</Label>
                <Input
                  id="isbn"
                  value={formData.isbn}
                  onChange={(e) => setFormData({ ...formData, isbn: e.target.value })}
                  placeholder={t('form.enterIsbn')}
                  className={errors.isbn ? 'border-destructive' : ''}
                />
                {errors.isbn && (
                  <p className="text-sm text-destructive">{errors.isbn}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="description">{t('form.description')}</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder={t('form.enterDescription')}
                  rows={4}
                  className="resize-none"
                />
              </div>
            </div>
          </div>
        </div>

        <Separator className="my-6" />

        {/* Publication Details */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.publicationDetails')}</h3>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-2">
                  <Label htmlFor="price">{t('form.price')}</Label>
                  <Input
                    id="price"
                    type="number"
                    step="0.01"
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
                    placeholder={t('form.pricePlaceholder')}
                    className={errors.price ? 'border-destructive' : ''}
                  />
                  {errors.price && (
                    <p className="text-sm text-destructive">{errors.price}</p>
                  )}
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-muted">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex-1 space-y-2">
                  <Label htmlFor="publishedAt">{t('form.publishedDate')}</Label>
                  <Input
                    id="publishedAt"
                    type="date"
                    value={formData.publishedAt}
                    onChange={(e) => setFormData({ ...formData, publishedAt: e.target.value })}
                  />
                </div>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Clock className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="status">{t('form.status')}</Label>
                <Select
                  value={formData.status}
                  onValueChange={(value) => setFormData({ ...formData, status: value as BookStatus })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder={t('form.selectStatus')} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="AVAILABLE">{t('status.available')}</SelectItem>
                    <SelectItem value="BORROWED">{t('status.borrowed')}</SelectItem>
                    <SelectItem value="MAINTENANCE">{t('status.maintenance')}</SelectItem>
                    <SelectItem value="RESERVED">{t('status.reserved')}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Tag className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="tags">{t('form.tags')}</Label>
                <div className="flex gap-2">
                  <Input
                    id="tags"
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    onKeyPress={handleKeyPress}
                    placeholder={t('form.addTags')}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={handleAddTag}
                    size="icon"
                  >
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
                {formData.tags && formData.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2 mt-2">
                    {formData.tags.map((tag, index) => (
                      <Badge key={index} variant="secondary" className="gap-1">
                        {tag}
                        <button
                          type="button"
                          onClick={() => handleRemoveTag(tag)}
                          className="ml-1 hover:text-destructive"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
});

BookCreateForm.displayName = 'BookCreateForm';