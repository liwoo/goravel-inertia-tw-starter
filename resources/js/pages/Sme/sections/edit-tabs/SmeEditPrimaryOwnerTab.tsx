import React, { useState, useEffect } from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { DISTRICT_NAMES } from '@/constants/districts';
import { NATIONALITY_OPTIONS } from '@/types/nationalities';
import { GENDER_OPTIONS } from '@/types/gender';
import { EDUCATION_OPTIONS } from '@/types/education';
import { MALAWIAN_STATUS_OPTIONS } from '@/types/malawian-status';
import {
  formatMalawiPhone,
  validateMalawiPhone,
  validateNationalID,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';

interface PrimaryOwnerData {
  id?: number;
  firstName: string;
  lastName: string;
  otherNames?: string;
  nationality: string;
  nationalIdNumber: string;
  dateOfBirth?: string;
  gender: string;
  educationLevel: string;
  malawianStatus: string;
  hasSpecialNeeds: boolean;
  phoneNumber: string;
  landlineNumber?: string;
  email?: string;
  physicalAddress?: string;
  postalAddress?: string;
  district?: string;
  traditionalAuthority?: string;
  altContactName?: string;
  altContactRelationship?: string;
  altContactPhone?: string;
}

interface SmeEditPrimaryOwnerTabProps {
  primaryOwner: PrimaryOwnerData | undefined;
  onChange: (data: PrimaryOwnerData) => void;
  errors: Record<string, string>;
}

export const SmeEditPrimaryOwnerTab: React.FC<SmeEditPrimaryOwnerTabProps> = ({
  primaryOwner,
  onChange,
  errors
}) => {
  if (!primaryOwner) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Primary Business Owner</CardTitle>
          <CardDescription>No primary owner data available</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">Unable to load primary owner information.</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Personal Information</CardTitle>
          <CardDescription>Primary business owner details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="firstName">First Name *</Label>
              <Input
                id="firstName"
                value={primaryOwner.firstName}
                onChange={(e) => onChange({ ...primaryOwner, firstName: e.target.value })}
                placeholder="Enter first name"
                className={errors.firstName ? 'border-destructive' : ''}
              />
              {errors.firstName && <p className="text-sm text-destructive">{errors.firstName}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="lastName">Last Name *</Label>
              <Input
                id="lastName"
                value={primaryOwner.lastName}
                onChange={(e) => onChange({ ...primaryOwner, lastName: e.target.value })}
                placeholder="Enter last name"
                className={errors.lastName ? 'border-destructive' : ''}
              />
              {errors.lastName && <p className="text-sm text-destructive">{errors.lastName}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="otherNames">Other Names</Label>
              <Input
                id="otherNames"
                value={primaryOwner.otherNames}
                onChange={(e) => onChange({ ...primaryOwner, otherNames: e.target.value })}
                placeholder="Enter other names"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="nationalIdNumber">National ID *</Label>
              <Input
                id="nationalIdNumber"
                value={primaryOwner.nationalIdNumber}
                onChange={(e) => onChange({ ...primaryOwner, nationalIdNumber: e.target.value.toUpperCase() })}
                placeholder={EXAMPLE_FORMATS.NATIONAL_ID}
                maxLength={8}
                className={errors.nationalIdNumber ? 'border-destructive' : ''}
              />
              {errors.nationalIdNumber && <p className="text-sm text-destructive">{errors.nationalIdNumber}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="nationality">Nationality *</Label>
              <Select
                value={primaryOwner.nationality}
                onValueChange={(value) => onChange({ ...primaryOwner, nationality: value })}
              >
                <SelectTrigger className={errors.nationality ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select nationality" />
                </SelectTrigger>
                <SelectContent>
                  {NATIONALITY_OPTIONS.map((nationality) => (
                    <SelectItem key={nationality.value} value={nationality.value}>
                      {nationality.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.nationality && <p className="text-sm text-destructive">{errors.nationality}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="dateOfBirth">Date of Birth *</Label>
              <Input
                id="dateOfBirth"
                type="date"
                value={primaryOwner.dateOfBirth}
                onChange={(e) => onChange({ ...primaryOwner, dateOfBirth: e.target.value })}
                className={errors.dateOfBirth ? 'border-destructive' : ''}
              />
              {errors.dateOfBirth && <p className="text-sm text-destructive">{errors.dateOfBirth}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="gender">Gender *</Label>
              <Select
                value={primaryOwner.gender}
                onValueChange={(value) => onChange({ ...primaryOwner, gender: value })}
              >
                <SelectTrigger className={errors.gender ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select gender" />
                </SelectTrigger>
                <SelectContent>
                  {GENDER_OPTIONS.map((gender) => (
                    <SelectItem key={gender.value} value={gender.value}>
                      {gender.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.gender && <p className="text-sm text-destructive">{errors.gender}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="educationLevel">Education Level *</Label>
              <Select
                value={primaryOwner.educationLevel}
                onValueChange={(value) => onChange({ ...primaryOwner, educationLevel: value })}
              >
                <SelectTrigger className={errors.educationLevel ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select education level" />
                </SelectTrigger>
                <SelectContent>
                  {EDUCATION_OPTIONS.map((education) => (
                    <SelectItem key={education.value} value={education.value}>
                      {education.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.educationLevel && <p className="text-sm text-destructive">{errors.educationLevel}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="malawianStatus">Malawian Status *</Label>
              <Select
                value={primaryOwner.malawianStatus}
                onValueChange={(value) => onChange({ ...primaryOwner, malawianStatus: value })}
              >
                <SelectTrigger className={errors.malawianStatus ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  {MALAWIAN_STATUS_OPTIONS.map((status) => (
                    <SelectItem key={status.value} value={status.value}>
                      {status.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.malawianStatus && <p className="text-sm text-destructive">{errors.malawianStatus}</p>}
            </div>

            <div className="flex items-center space-x-2 pt-7">
              <Checkbox
                id="hasSpecialNeeds"
                checked={primaryOwner.hasSpecialNeeds}
                onCheckedChange={(checked) => onChange({ ...primaryOwner, hasSpecialNeeds: checked as boolean })}
              />
              <Label htmlFor="hasSpecialNeeds" className="cursor-pointer">Has Special Needs</Label>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Contact Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="phoneNumber">Phone Number *</Label>
              <Input
                id="phoneNumber"
                value={primaryOwner.phoneNumber}
                onChange={(e) => {
                  const formatted = formatMalawiPhone(e.target.value);
                  onChange({ ...primaryOwner, phoneNumber: formatted });
                }}
                placeholder={EXAMPLE_FORMATS.PHONE}
                className={errors.phoneNumber ? 'border-destructive' : ''}
              />
              {errors.phoneNumber && <p className="text-sm text-destructive">{errors.phoneNumber}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="landlineNumber">Landline</Label>
              <Input
                id="landlineNumber"
                value={primaryOwner.landlineNumber}
                onChange={(e) => {
                  const formatted = formatMalawiPhone(e.target.value);
                  onChange({ ...primaryOwner, landlineNumber: formatted });
                }}
                placeholder={EXAMPLE_FORMATS.PHONE}
                className={errors.landlineNumber ? 'border-destructive' : ''}
              />
              {errors.landlineNumber && <p className="text-sm text-destructive">{errors.landlineNumber}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={primaryOwner.email}
                onChange={(e) => onChange({ ...primaryOwner, email: e.target.value })}
                placeholder="Enter email"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="ownerDistrict">District</Label>
              <Select
                value={primaryOwner.district}
                onValueChange={(value) => onChange({ ...primaryOwner, district: value })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select district" />
                </SelectTrigger>
                <SelectContent>
                  {DISTRICT_NAMES.map((district) => (
                    <SelectItem key={district} value={district}>
                      {district}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="ownerTraditionalAuthority">Traditional Authority</Label>
              <Input
                id="ownerTraditionalAuthority"
                value={primaryOwner.traditionalAuthority}
                onChange={(e) => onChange({ ...primaryOwner, traditionalAuthority: e.target.value })}
                placeholder="Enter TA"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="ownerPhysicalAddress">Physical Address</Label>
            <Textarea
              id="ownerPhysicalAddress"
              value={primaryOwner.physicalAddress}
              onChange={(e) => onChange({ ...primaryOwner, physicalAddress: e.target.value })}
              placeholder="Enter physical address"
              rows={2}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="ownerPostalAddress">Postal Address</Label>
            <Textarea
              id="ownerPostalAddress"
              value={primaryOwner.postalAddress}
              onChange={(e) => onChange({ ...primaryOwner, postalAddress: e.target.value })}
              placeholder="Enter postal address"
              rows={2}
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Alternative Contact</CardTitle>
          <CardDescription>Emergency contact person (optional)</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="altContactName">Contact Name</Label>
              <Input
                id="altContactName"
                value={primaryOwner.altContactName}
                onChange={(e) => onChange({ ...primaryOwner, altContactName: e.target.value })}
                placeholder="Enter contact name"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="altContactRelationship">Relationship</Label>
              <Input
                id="altContactRelationship"
                value={primaryOwner.altContactRelationship}
                onChange={(e) => onChange({ ...primaryOwner, altContactRelationship: e.target.value })}
                placeholder="Enter relationship"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="altContactPhone">Contact Phone</Label>
              <Input
                id="altContactPhone"
                value={primaryOwner.altContactPhone}
                onChange={(e) => {
                  const formatted = formatMalawiPhone(e.target.value);
                  onChange({ ...primaryOwner, altContactPhone: formatted });
                }}
                placeholder={EXAMPLE_FORMATS.PHONE}
                className={errors.altContactPhone ? 'border-destructive' : ''}
              />
              {errors.altContactPhone && <p className="text-sm text-destructive">{errors.altContactPhone}</p>}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
