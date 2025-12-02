import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Separator } from '@/components/ui/separator';
import {
  User as UserIcon,
  Mail,
  Shield,
  CalendarDays,
  MapPin,
  Pencil,
  X,
  Check,
  Loader2
} from 'lucide-react';
import { toast } from 'sonner';
import axios from '@/lib/axios';
import type {
  ProfileData,
  UpdateProfileRequest,
  ProfileFormErrors,
  ApiErrorResponse
} from '@/types/account';
import type { User } from '@/types/app';

interface ProfileSectionProps {
  user: User | null;
}

export const ProfileSection: React.FC<ProfileSectionProps> = ({ user }) => {
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [formData, setFormData] = useState<UpdateProfileRequest>({
    name: '',
    email: '',
  });
  const [errors, setErrors] = useState<ProfileFormErrors>({});

  // Fetch profile data on mount
  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    setIsLoading(true);
    try {
      const response = await axios.get<{ success: boolean; data: ProfileData }>('/api/account/profile');
      const profileData = response.data.data;
      setProfile(profileData);
      setFormData({
        name: profileData.name,
        email: profileData.email,
      });
    } catch (error) {
      console.error('Failed to fetch profile:', error);
      toast.error('Failed to load profile data');
      // Fall back to user data from shared data if available
      if (user) {
        setFormData({
          name: user.name,
          email: user.email,
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const validateForm = (): boolean => {
    const newErrors: ProfileFormErrors = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    } else if (formData.name.length < 2) {
      newErrors.name = 'Name must be at least 2 characters';
    } else if (formData.name.length > 100) {
      newErrors.name = 'Name must not exceed 100 characters';
    }

    if (!formData.email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSave = async () => {
    if (!validateForm()) {
      return;
    }

    setIsSaving(true);
    try {
      const response = await axios.put<{ success: boolean; data: ProfileData }>('/api/account/profile', formData);
      setProfile(response.data.data);
      setIsEditing(false);
      setErrors({});
      toast.success('Profile updated successfully');
    } catch (error: unknown) {
      console.error('Failed to update profile:', error);

      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;

        if (errorData.errors) {
          const apiErrors: ProfileFormErrors = {};
          if (errorData.errors.name) {
            apiErrors.name = errorData.errors.name[0];
          }
          if (errorData.errors.email) {
            apiErrors.email = errorData.errors.email[0];
          }
          setErrors(apiErrors);
        }

        toast.error(errorData.message || 'Failed to update profile');
      } else {
        toast.error('An unexpected error occurred');
      }
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setErrors({});
    // Reset form data to current profile
    if (profile) {
      setFormData({
        name: profile.name,
        email: profile.email,
      });
    }
  };

  const getInitials = (name: string): string => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const formatDate = (dateString?: string): string => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const formatDateTime = (dateString?: string): string => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <Skeleton className="h-6 w-32" />
              <Skeleton className="h-4 w-48" />
            </div>
            <Skeleton className="h-9 w-20" />
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col sm:flex-row items-center gap-6">
            <Skeleton className="h-24 w-24 rounded-full" />
            <div className="space-y-2 text-center sm:text-left">
              <Skeleton className="h-6 w-40" />
              <Skeleton className="h-4 w-48" />
              <Skeleton className="h-5 w-20" />
            </div>
          </div>
          <Separator />
          <div className="grid gap-4 sm:grid-cols-2">
            <Skeleton className="h-20" />
            <Skeleton className="h-20" />
            <Skeleton className="h-20" />
            <Skeleton className="h-20" />
          </div>
        </CardContent>
      </Card>
    );
  }

  const displayProfile = profile || (user ? {
    id: user.id,
    name: user.name,
    email: user.email,
    role: user.role,
    is_active: user.is_active ?? true,
    created_at: '',
    roles: user.roles || [],
  } : null);

  if (!displayProfile) {
    return (
      <Card>
        <CardContent className="py-8">
          <p className="text-center text-muted-foreground">
            Unable to load profile information.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Profile Information</CardTitle>
            <CardDescription>
              Manage your account details and personal information
            </CardDescription>
          </div>
          {!isEditing ? (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsEditing(true)}
            >
              <Pencil className="h-4 w-4 mr-2" />
              Edit
            </Button>
          ) : (
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleCancel}
                disabled={isSaving}
              >
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button
                size="sm"
                onClick={handleSave}
                disabled={isSaving}
              >
                {isSaving ? (
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                ) : (
                  <Check className="h-4 w-4 mr-2" />
                )}
                Save
              </Button>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Profile Header */}
        <div className="flex flex-col sm:flex-row items-center gap-6">
          <Avatar className="h-24 w-24 ring-4 ring-background shadow-xl">
            <AvatarImage
              src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${displayProfile.email}`}
              alt={displayProfile.name}
            />
            <AvatarFallback className="text-xl bg-gradient-to-br from-teal-500 to-teal-600 text-white">
              {getInitials(displayProfile.name)}
            </AvatarFallback>
          </Avatar>
          <div className="space-y-2 text-center sm:text-left">
            {isEditing ? (
              <div className="space-y-2">
                <Label htmlFor="edit-name">Full Name</Label>
                <Input
                  id="edit-name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter your name"
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>
            ) : (
              <h2 className="text-2xl font-bold text-foreground">{displayProfile.name}</h2>
            )}
            {isEditing ? (
              <div className="space-y-2">
                <Label htmlFor="edit-email">Email Address</Label>
                <Input
                  id="edit-email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  placeholder="Enter your email"
                  className={errors.email ? 'border-destructive' : ''}
                />
                {errors.email && (
                  <p className="text-sm text-destructive">{errors.email}</p>
                )}
              </div>
            ) : (
              <p className="text-muted-foreground">{displayProfile.email}</p>
            )}
            <div className="flex flex-wrap justify-center sm:justify-start gap-2">
              {displayProfile.is_active ? (
                <Badge variant="default" className="bg-green-600">Active</Badge>
              ) : (
                <Badge variant="destructive">Inactive</Badge>
              )}
              {displayProfile.roles?.map((role) => (
                <Badge key={role.id} variant="secondary">
                  {role.name}
                </Badge>
              ))}
            </div>
          </div>
        </div>

        <Separator />

        {/* Profile Details */}
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
            <div className="p-2 rounded-lg bg-background">
              <UserIcon className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Full Name</p>
              <p className="font-medium text-foreground">{displayProfile.name}</p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
            <div className="p-2 rounded-lg bg-background">
              <Mail className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Email Address</p>
              <p className="font-medium text-foreground">{displayProfile.email}</p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
            <div className="p-2 rounded-lg bg-background">
              <Shield className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Role</p>
              <p className="font-medium text-foreground">
                {displayProfile.roles?.map(r => r.name).join(', ') || displayProfile.role || 'N/A'}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
            <div className="p-2 rounded-lg bg-background">
              <CalendarDays className="h-4 w-4 text-muted-foreground" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Member Since</p>
              <p className="font-medium text-foreground">
                {formatDate(profile?.created_at)}
              </p>
            </div>
          </div>

          {profile?.last_login_at && (
            <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
              <div className="p-2 rounded-lg bg-background">
                <CalendarDays className="h-4 w-4 text-muted-foreground" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Last Login</p>
                <p className="font-medium text-foreground">
                  {formatDateTime(profile.last_login_at)}
                </p>
              </div>
            </div>
          )}

          {profile?.last_login_ip && (
            <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
              <div className="p-2 rounded-lg bg-background">
                <MapPin className="h-4 w-4 text-muted-foreground" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Last Login IP</p>
                <p className="font-medium text-foreground font-mono text-sm">
                  {profile.last_login_ip}
                </p>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default ProfileSection;
