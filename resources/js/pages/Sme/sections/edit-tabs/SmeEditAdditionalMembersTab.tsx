import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { Plus, Trash2, Pencil } from 'lucide-react';
import { NATIONALITY_OPTIONS } from '@/types/nationalities';
import {
  validateMalawiPhone,
  validateNationalID,
  formatMalawiPhone,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';

interface AdditionalMemberData {
  id?: number;
  firstName: string;
  lastName: string;
  otherNames?: string;
  nationality: string;
  nationalIdNumber: string;
  dateOfBirth?: string;
  gender: string;
  email?: string;
  phoneNumber: string;
  isIntern: boolean;
  isPartTime: boolean;
}

interface SmeEditAdditionalMembersTabProps {
  members: AdditionalMemberData[];
  onChange: (members: AdditionalMemberData[]) => void;
}

const emptyMember: AdditionalMemberData = {
  firstName: '',
  lastName: '',
  otherNames: '',
  nationality: '',
  nationalIdNumber: '',
  dateOfBirth: '',
  gender: '',
  email: '',
  phoneNumber: '',
  isIntern: false,
  isPartTime: false,
};

const GENDER_OPTIONS = [
  { value: 'MALE', label: 'Male' },
  { value: 'FEMALE', label: 'Female' },
];

// Reusable member form component
const MemberForm: React.FC<{
  member: AdditionalMemberData;
  onChange: (member: AdditionalMemberData) => void;
  idPrefix: string;
}> = ({ member, onChange, idPrefix }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-firstName`}>First Name *</Label>
        <Input
          id={`${idPrefix}-firstName`}
          value={member.firstName}
          onChange={(e) => onChange({ ...member, firstName: e.target.value })}
          placeholder="Enter first name"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-lastName`}>Last Name *</Label>
        <Input
          id={`${idPrefix}-lastName`}
          value={member.lastName}
          onChange={(e) => onChange({ ...member, lastName: e.target.value })}
          placeholder="Enter last name"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-otherNames`}>Other Names</Label>
        <Input
          id={`${idPrefix}-otherNames`}
          value={member.otherNames || ''}
          onChange={(e) => onChange({ ...member, otherNames: e.target.value })}
          placeholder="Enter other names"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-nationality`}>Nationality *</Label>
        <Select
          value={member.nationality}
          onValueChange={(value) => onChange({ ...member, nationality: value })}
        >
          <SelectTrigger>
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
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-nationalIdNumber`}>National ID *</Label>
        <Input
          id={`${idPrefix}-nationalIdNumber`}
          value={member.nationalIdNumber}
          onChange={(e) => onChange({ ...member, nationalIdNumber: e.target.value.toUpperCase() })}
          placeholder={EXAMPLE_FORMATS.NATIONAL_ID}
          maxLength={10}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-dateOfBirth`}>Date of Birth</Label>
        <Input
          id={`${idPrefix}-dateOfBirth`}
          type="date"
          value={member.dateOfBirth || ''}
          onChange={(e) => onChange({ ...member, dateOfBirth: e.target.value })}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-gender`}>Gender *</Label>
        <Select
          value={member.gender}
          onValueChange={(value) => onChange({ ...member, gender: value })}
        >
          <SelectTrigger>
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
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-email`}>Email</Label>
        <Input
          id={`${idPrefix}-email`}
          type="email"
          value={member.email || ''}
          onChange={(e) => onChange({ ...member, email: e.target.value })}
          placeholder="Enter email"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${idPrefix}-phoneNumber`}>Phone Number *</Label>
        <Input
          id={`${idPrefix}-phoneNumber`}
          value={member.phoneNumber}
          onChange={(e) => {
            const formatted = formatMalawiPhone(e.target.value);
            onChange({ ...member, phoneNumber: formatted });
          }}
          placeholder={EXAMPLE_FORMATS.PHONE}
        />
      </div>

      <div className="space-y-2">
        <Label className="text-sm text-muted-foreground">Employment Type</Label>
        <div className="flex flex-col gap-3 pt-1">
          <div className="flex items-center space-x-2">
            <Checkbox
              id={`${idPrefix}-isIntern`}
              checked={member.isIntern}
              onCheckedChange={(checked) => onChange({ ...member, isIntern: checked as boolean })}
            />
            <Label htmlFor={`${idPrefix}-isIntern`} className="cursor-pointer font-normal">Is Intern</Label>
          </div>
          <div className="flex items-center space-x-2">
            <Checkbox
              id={`${idPrefix}-isPartTime`}
              checked={member.isPartTime}
              onCheckedChange={(checked) => onChange({ ...member, isPartTime: checked as boolean })}
            />
            <Label htmlFor={`${idPrefix}-isPartTime`} className="cursor-pointer font-normal">Is Part Time</Label>
          </div>
        </div>
      </div>
    </div>
  );
};

export const SmeEditAdditionalMembersTab: React.FC<SmeEditAdditionalMembersTabProps> = ({
  members,
  onChange
}) => {
  const [activeTab, setActiveTab] = useState('new-member');
  const [newMember, setNewMember] = useState<AdditionalMemberData>(emptyMember);

  const addMember = () => {
    // Validate required fields
    if (!newMember.firstName.trim() || !newMember.lastName.trim() ||
        !newMember.nationalIdNumber.trim() || !newMember.phoneNumber.trim() ||
        !newMember.gender) {
      return;
    }

    // Validate national ID format
    if (!validateNationalID(newMember.nationalIdNumber)) {
      alert(VALIDATION_MESSAGES.NATIONAL_ID);
      return;
    }

    // Validate phone format
    if (!validateMalawiPhone(newMember.phoneNumber)) {
      alert(VALIDATION_MESSAGES.PHONE);
      return;
    }

    onChange([...members, newMember]);
    setNewMember(emptyMember);

    // Switch to the newly added member's tab
    setActiveTab(`member-${members.length}`);
  };

  const updateMember = (index: number, updatedMember: AdditionalMemberData) => {
    const newMembers = [...members];
    newMembers[index] = updatedMember;
    onChange(newMembers);
  };

  const removeMember = (index: number) => {
    const newMembers = members.filter((_, i) => i !== index);
    onChange(newMembers);
    // Switch to new-member tab after deletion
    setActiveTab('new-member');
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Additional Business Members</CardTitle>
          <CardDescription>
            Add or edit team members (optional). You can add multiple members.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
            <TabsList
              className="w-full"
              style={{
                display: 'grid',
                gridTemplateColumns: `repeat(${members.length + 1}, 1fr)`
              }}
            >
              {members.map((member, index) => (
                <TabsTrigger key={`member-${index}`} value={`member-${index}`}>
                  <span className="truncate max-w-[120px]">
                    {member.firstName} {member.lastName}
                  </span>
                </TabsTrigger>
              ))}
              <TabsTrigger value="new-member">
                <Plus className="h-4 w-4 mr-1" />
                Add New
              </TabsTrigger>
            </TabsList>

            {/* Tabs for existing members - now editable */}
            {members.map((member, index) => (
              <TabsContent key={`member-${index}`} value={`member-${index}`} className="space-y-4">
                <p className="text-xs text-muted-foreground">
                  Changes are tracked automatically. Click "Save" to persist updates.
                </p>
                <Card>
                  <CardHeader className="pb-4">
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg flex items-center gap-2">
                          <Pencil className="h-4 w-4 text-muted-foreground" />
                          Edit: {member.firstName} {member.lastName}
                        </CardTitle>
                        <CardDescription className="mt-1">
                          Update team member details below
                        </CardDescription>
                        <div className="flex gap-2 mt-2">
                          {member.isIntern && <Badge variant="secondary">Intern</Badge>}
                          {member.isPartTime && <Badge variant="outline">Part Time</Badge>}
                        </div>
                      </div>
                      <Button
                        type="button"
                        variant="destructive"
                        size="sm"
                        onClick={() => removeMember(index)}
                      >
                        <Trash2 className="h-4 w-4 mr-1" />
                        Remove
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <MemberForm
                      member={member}
                      onChange={(updatedMember) => updateMember(index, updatedMember)}
                      idPrefix={`member-${index}`}
                    />
                  </CardContent>
                </Card>
              </TabsContent>
            ))}

            {/* Tab for adding new member */}
            <TabsContent value="new-member" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">New Team Member</CardTitle>
                  <CardDescription>Fill in the details to add a new team member</CardDescription>
                </CardHeader>
                <CardContent>
                  <MemberForm
                    member={newMember}
                    onChange={setNewMember}
                    idPrefix="new-member"
                  />

                  <Button
                    type="button"
                    onClick={addMember}
                    disabled={!newMember.firstName || !newMember.lastName || !newMember.nationalIdNumber || !newMember.phoneNumber || !newMember.gender}
                    className="w-full mt-6"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Add Team Member
                  </Button>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>

          {members.length === 0 && (
            <p className="text-sm text-muted-foreground text-center py-4">
              No team members added yet. Click "Add New" to add your first team member.
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
