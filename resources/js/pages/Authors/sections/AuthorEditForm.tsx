import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { useTranslation } from 'react-i18next';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { CrudEditFormProps } from '@/types/crud';
import { Author, AuthorUpdateData, AuthorStatus } from '@/types/author';
import { Separator } from '@/components/ui/separator';
import { NATIONALITIES } from '@/config/options';

interface AuthorEditFormProps extends CrudEditFormProps<Author> {
  setIsSaving?: (saving: boolean) => void;
}

export const AuthorEditForm = forwardRef<any, AuthorEditFormProps>(
  (
    {
      item: author,
      onSuccess,
      onError,
      onCancel,
      isLoading = false,
      setIsSaving,
    },
    ref,
  ) => {
    const { t } = useTranslation('authors');
    const [formData, setFormData] = useState<AuthorUpdateData>({
      firstName: author.firstName,
      lastName: author.lastName,
      bio: author.bio || '',
      email: author.email || '',
      website: author.website || '',
      birthDate: author.birthDate
        ? new Date(author.birthDate).toISOString().split('T')[0]
        : '',
      nationality: author.nationality || '',
      photoUrl: author.photoUrl || '',
      status: author.status,
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    const handleSubmit = async () => {
      const newErrors: Record<string, string> = {};
      if (!formData.firstName?.trim()) {
        newErrors.firstName = t('validation.firstNameRequired');
      }
      if (!formData.lastName?.trim()) {
        newErrors.lastName = t('validation.lastNameRequired');
      }
      if (
        formData.email &&
        formData.email.trim() &&
        !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)
      ) {
        newErrors.email = t('validation.emailInvalid');
      }

      if (Object.keys(newErrors).length > 0) {
        setErrors(newErrors);
        return;
      }

      setErrors({});
      setIsSaving?.(true);

      try {
        const apiData = {
          first_name: formData.firstName,
          last_name: formData.lastName,
          bio: formData.bio || null,
          email: formData.email || null,
          website: formData.website || null,
          birth_date: formData.birthDate || null,
          nationality: formData.nationality || null,
          photo_url: formData.photoUrl || null,
          status: formData.status,
        };

        const response = await fetch(`/api/authors/${author.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
          body: JSON.stringify(apiData),
        });

        if (response.ok) {
          onSuccess(t('toast.updated'));
        } else {
          const errorData = await response.json().catch(() => ({}));
          if (errorData.errors) {
            const serverErrors: Record<string, string> = {};
            Object.entries(errorData.errors).forEach(([key, messages]) => {
              const camelKey = key.replace(/_([a-z])/g, (_, c) =>
                c.toUpperCase(),
              );
              serverErrors[camelKey] = Array.isArray(messages)
                ? messages[0]
                : String(messages);
            });
            setErrors(serverErrors);
          }
          onError?.(errorData);
        }
      } catch (error) {
        onError?.(error);
      } finally {
        setIsSaving?.(false);
      }
    };

    useImperativeHandle(ref, () => ({
      handleSubmit,
    }));

    return (
      <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
        {/* Author Information */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">
            {t('form.authorInfo')}
          </h3>
          <div className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">{t('form.firstName')}</Label>
                <Input
                  id="firstName"
                  value={formData.firstName}
                  onChange={(e) =>
                    setFormData({ ...formData, firstName: e.target.value })
                  }
                  placeholder={t('form.enterFirstName')}
                  className={errors.firstName ? 'border-destructive' : ''}
                />
                {errors.firstName && (
                  <p className="text-sm text-destructive">
                    {errors.firstName}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName">{t('form.lastName')}</Label>
                <Input
                  id="lastName"
                  value={formData.lastName}
                  onChange={(e) =>
                    setFormData({ ...formData, lastName: e.target.value })
                  }
                  placeholder={t('form.enterLastName')}
                  className={errors.lastName ? 'border-destructive' : ''}
                />
                {errors.lastName && (
                  <p className="text-sm text-destructive">
                    {errors.lastName}
                  </p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="bio">{t('form.bio')}</Label>
              <Textarea
                id="bio"
                value={formData.bio}
                onChange={(e) =>
                  setFormData({ ...formData, bio: e.target.value })
                }
                placeholder={t('form.enterBio')}
                rows={4}
                className="resize-none"
              />
            </div>
          </div>
        </div>

        <Separator />

        {/* Contact Details */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">
            {t('form.contactDetails')}
          </h3>
          <div className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="email">{t('form.email')}</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  placeholder={t('form.enterEmail')}
                  className={errors.email ? 'border-destructive' : ''}
                />
                {errors.email && (
                  <p className="text-sm text-destructive">{errors.email}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="website">{t('form.website')}</Label>
                <Input
                  id="website"
                  type="url"
                  value={formData.website}
                  onChange={(e) =>
                    setFormData({ ...formData, website: e.target.value })
                  }
                  placeholder={t('form.enterWebsite')}
                />
              </div>
            </div>
          </div>
        </div>

        <Separator />

        {/* Additional Information */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">
            {t('form.additionalInfo')}
          </h3>
          <div className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="birthDate">{t('form.birthDate')}</Label>
                <Input
                  id="birthDate"
                  type="date"
                  value={formData.birthDate}
                  onChange={(e) =>
                    setFormData({ ...formData, birthDate: e.target.value })
                  }
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="nationality">{t('form.nationality')}</Label>
                <Select
                  value={formData.nationality || undefined}
                  onValueChange={(value) =>
                    setFormData({ ...formData, nationality: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue
                      placeholder={t('form.selectNationality')}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {NATIONALITIES.map((n) => (
                      <SelectItem key={n.value} value={n.value}>
                        {n.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="photoUrl">{t('form.photoUrl')}</Label>
              <Input
                id="photoUrl"
                type="url"
                value={formData.photoUrl}
                onChange={(e) =>
                  setFormData({ ...formData, photoUrl: e.target.value })
                }
                placeholder={t('form.enterPhotoUrl')}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="status">{t('form.status')}</Label>
              <Select
                value={formData.status}
                onValueChange={(value) =>
                  setFormData({ ...formData, status: value as AuthorStatus })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder={t('form.selectStatus')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ACTIVE">{t('status.active')}</SelectItem>
                  <SelectItem value="INACTIVE">
                    {t('status.inactive')}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        {/* Metadata */}
        <Separator />
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">
            {t('form.metadata')}
          </h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">{t('form.authorId')}</p>
              <p className="font-medium">#{author.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">{t('form.created')}</p>
              <p className="font-medium">
                {(() => {
                  const created = author.createdAt || author.created_at;
                  return created
                    ? new Date(created).toLocaleDateString()
                    : '-';
                })()}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">{t('form.lastUpdated')}</p>
              <p className="font-medium">
                {(() => {
                  const updated = author.updatedAt || author.updated_at;
                  return updated
                    ? new Date(updated).toLocaleDateString()
                    : '-';
                })()}
              </p>
            </div>
          </div>
        </div>
      </form>
    );
  },
);

AuthorEditForm.displayName = 'AuthorEditForm';
