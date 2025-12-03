import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { MemberCreateData } from '@/types/member';
import { User, Mail } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface MemberCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
  smeId?: number;
}

export const MemberCreateForm = forwardRef<any, MemberCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving,
  smeId
}, ref) => {
  const [formData, setFormData] = useState<MemberCreateData>({
    first_name: '',
    last_name: '',
    other_names: '',
    gender: '',
    nationality: 'Malawian',
    national_id_number: '',
    date_of_birth: '',
    email: '',
    phone_number: '',
    is_intern: false,
    is_part_time: false,
    sme_id: smeId || 0,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    if (!formData.first_name) {
      newErrors.first_name = 'First Name is required';
    }
    if (!formData.last_name) {
      newErrors.last_name = 'Last Name is required';
    }
    if (!formData.national_id_number) {
      newErrors.national_id_number = 'National ID is required';
    }
    if (!formData.phone_number) {
      newErrors.phone_number = 'Phone Number is required';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch('/api/additional_business_members', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Member created successfully');
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
                value={formData.first_name}
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
                value={formData.last_name}
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
                value={formData.nationality}
                onChange={(e) => setFormData({ ...formData, nationality: e.target.value })}
                placeholder="Enter nationality"
              />
            </div>

            {/* National ID */}
            <div className="space-y-2">
              <Label htmlFor="national_id_number">National ID Number *</Label>
              <Input
                id="national_id_number"
                value={formData.national_id_number}
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
                value={formData.phone_number}
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
                checked={formData.is_intern}
                onChange={(e) => setFormData({ ...formData, is_intern: e.target.checked })}
              />
              <Label htmlFor="is_intern">Is Intern?</Label>
            </div>
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_part_time"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                checked={formData.is_part_time}
                onChange={(e) => setFormData({ ...formData, is_part_time: e.target.checked })}
              />
              <Label htmlFor="is_part_time">Is Part Time?</Label>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
});

MemberCreateForm.displayName = 'MemberCreateForm';
