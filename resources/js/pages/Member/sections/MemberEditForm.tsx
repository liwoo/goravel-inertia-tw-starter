import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { Member, MemberUpdateData } from '@/types/member';
import { User, Mail } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface MemberEditFormProps extends CrudEditFormProps<Member> {
  setIsSaving?: (saving: boolean) => void;
}

export const MemberEditForm = forwardRef<any, MemberEditFormProps>(({
  item: member,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<MemberUpdateData>({
    first_name: member.first_name,
    last_name: member.last_name,
    other_names: member.other_names,
    gender: member.gender,
    nationality: member.nationality,
    national_id_number: member.national_id_number,
    date_of_birth: member.date_of_birth ? member.date_of_birth.slice(0, 10) : '',
    email: member.email,
    phone_number: member.phone_number,
    is_intern: member.is_intern,
    is_part_time: member.is_part_time,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    if (!formData.first_name?.trim()) {
      newErrors.first_name = 'First Name is required';
    }
    if (!formData.last_name?.trim()) {
      newErrors.last_name = 'Last Name is required';
    }
    if (!formData.national_id_number?.trim()) {
      newErrors.national_id_number = 'National ID is required';
    }
    if (!formData.phone_number?.trim()) {
      newErrors.phone_number = 'Phone Number is required';
    }


    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`/api/additional_business_members/${member.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Member updated successfully');
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

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Member Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* First Name */}
            <div className="space-y-2">
              <Label htmlFor="first_name">First Name *</Label>
              <Input
                id="first_name"
                value={formData.first_name || ''}
                onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                placeholder="Enter first name"
                className={errors.first_name ? 'border-destructive' : ''}
              />
              {errors.first_name && <p className="text-sm text-destructive">{errors.first_name}</p>}
            </div>

            {/* Last Name */}
            <div className="space-y-2">
              <Label htmlFor="last_name">Last Name *</Label>
              <Input
                id="last_name"
                value={formData.last_name || ''}
                onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                placeholder="Enter last name"
                className={errors.last_name ? 'border-destructive' : ''}
              />
              {errors.last_name && <p className="text-sm text-destructive">{errors.last_name}</p>}
            </div>

            {/* Other Names */}
            <div className="space-y-2">
              <Label htmlFor="other_names">Other Names</Label>
              <Input
                id="other_names"
                value={formData.other_names || ''}
                onChange={(e) => setFormData({ ...formData, other_names: e.target.value })}
                placeholder="Enter other names"
              />
            </div>

            {/* Gender */}
            <div className="space-y-2">
              <Label htmlFor="gender">Gender</Label>
              <select
                id="gender"
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                value={formData.gender || ''}
                onChange={(e) => setFormData({ ...formData, gender: e.target.value })}
              >
                <option value="">Select Gender</option>
                <option value="Male">Male</option>
                <option value="Female">Female</option>
              </select>
            </div>

            {/* Nationality */}
            <div className="space-y-2">
              <Label htmlFor="nationality">Nationality *</Label>
              <Input
                id="nationality"
                value={formData.nationality || ''}
                onChange={(e) => setFormData({ ...formData, nationality: e.target.value })}
                placeholder="Enter nationality"
              />
            </div>

            {/* National ID */}
            <div className="space-y-2">
              <Label htmlFor="national_id_number">National ID Number *</Label>
              <Input
                id="national_id_number"
                value={formData.national_id_number || ''}
                onChange={(e) => setFormData({ ...formData, national_id_number: e.target.value })}
                placeholder="Enter national ID"
                className={errors.national_id_number ? 'border-destructive' : ''}
              />
              {errors.national_id_number && <p className="text-sm text-destructive">{errors.national_id_number}</p>}
            </div>

            {/* Date of Birth */}
            <div className="space-y-2">
              <Label htmlFor="date_of_birth">Date of Birth</Label>
              <Input
                id="date_of_birth"
                type="date"
                value={formData.date_of_birth || ''}
                onChange={(e) => setFormData({ ...formData, date_of_birth: e.target.value })}
              />
            </div>

            {/* Phone Number */}
            <div className="space-y-2">
              <Label htmlFor="phone_number">Phone Number *</Label>
              <Input
                id="phone_number"
                value={formData.phone_number || ''}
                onChange={(e) => setFormData({ ...formData, phone_number: e.target.value })}
                placeholder="Enter phone number"
                className={errors.phone_number ? 'border-destructive' : ''}
              />
              {errors.phone_number && <p className="text-sm text-destructive">{errors.phone_number}</p>}
            </div>

            {/* Email */}
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={formData.email || ''}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                placeholder="Enter email"
                className={errors.email ? 'border-destructive' : ''}
              />
              {errors.email && <p className="text-sm text-destructive">{errors.email}</p>}
            </div>
          </div>

          {/* Checkboxes */}
          <div className="flex gap-6 mt-4">
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_intern"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                checked={formData.is_intern || false}
                onChange={(e) => setFormData({ ...formData, is_intern: e.target.checked })}
              />
              <Label htmlFor="is_intern">Is Intern?</Label>
            </div>
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_part_time"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                checked={formData.is_part_time || false}
                onChange={(e) => setFormData({ ...formData, is_part_time: e.target.checked })}
              />
              <Label htmlFor="is_part_time">Is Part Time?</Label>
            </div>
          </div>
        </div>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">ID</p>
              <p className="font-medium text-foreground">#{member.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">
                {(member.createdAt || member.created_at) ? new Date(member.createdAt || member.created_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">
                {(member.updatedAt || member.updated_at) ? new Date(member.updatedAt || member.updated_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
});

MemberEditForm.displayName = 'MemberEditForm';
