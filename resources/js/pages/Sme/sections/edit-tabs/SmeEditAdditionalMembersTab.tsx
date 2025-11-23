import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { Plus, Trash2 } from 'lucide-react';
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
  email: '',
  phoneNumber: '',
  isIntern: false,
  isPartTime: false,
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
        !newMember.nationalIdNumber.trim() || !newMember.phoneNumber.trim()) {
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
            Add team members (optional). You can add multiple members.
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

            {/* Tabs for existing members */}
            {members.map((member, index) => (
              <TabsContent key={`member-${index}`} value={`member-${index}`} className="space-y-4">
                <Card>
                  <CardHeader className="pb-4">
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg">{member.firstName} {member.lastName}</CardTitle>
                        <div className="flex gap-2 mt-2">
                          {member.isIntern && <Badge variant="secondary">Intern</Badge>}
                          {member.isPartTime && <Badge variant="outline">Part Time</Badge>}
                        </div>
                      </div>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => removeMember(index)}
                        title="Remove member"
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <p className="text-muted-foreground">National ID</p>
                        <p className="font-medium">{member.nationalIdNumber}</p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">Phone</p>
                        <p className="font-medium">{member.phoneNumber}</p>
                      </div>
                      {member.email && (
                        <div>
                          <p className="text-muted-foreground">Email</p>
                          <p className="font-medium">{member.email}</p>
                        </div>
                      )}
                      {member.nationality && (
                        <div>
                          <p className="text-muted-foreground">Nationality</p>
                          <p className="font-medium">{member.nationality}</p>
                        </div>
                      )}
                      {member.dateOfBirth && (
                        <div>
                          <p className="text-muted-foreground">Date of Birth</p>
                          <p className="font-medium">{new Date(member.dateOfBirth).toLocaleDateString()}</p>
                        </div>
                      )}
                      {member.otherNames && (
                        <div>
                          <p className="text-muted-foreground">Other Names</p>
                          <p className="font-medium">{member.otherNames}</p>
                        </div>
                      )}
                    </div>
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
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="memberFirstName">First Name *</Label>
                      <Input
                        id="memberFirstName"
                        value={newMember.firstName}
                        onChange={(e) => setNewMember({ ...newMember, firstName: e.target.value })}
                        placeholder="Enter first name"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberLastName">Last Name *</Label>
                      <Input
                        id="memberLastName"
                        value={newMember.lastName}
                        onChange={(e) => setNewMember({ ...newMember, lastName: e.target.value })}
                        placeholder="Enter last name"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberOtherNames">Other Names</Label>
                      <Input
                        id="memberOtherNames"
                        value={newMember.otherNames}
                        onChange={(e) => setNewMember({ ...newMember, otherNames: e.target.value })}
                        placeholder="Enter other names"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberNationality">Nationality *</Label>
                      <Select
                        value={newMember.nationality}
                        onValueChange={(value) => setNewMember({ ...newMember, nationality: value })}
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
                      <Label htmlFor="memberNationalIdNumber">National ID *</Label>
                      <Input
                        id="memberNationalIdNumber"
                        value={newMember.nationalIdNumber}
                        onChange={(e) => setNewMember({ ...newMember, nationalIdNumber: e.target.value.toUpperCase() })}
                        placeholder={EXAMPLE_FORMATS.NATIONAL_ID}
                        maxLength={10}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberDateOfBirth">Date of Birth</Label>
                      <Input
                        id="memberDateOfBirth"
                        type="date"
                        value={newMember.dateOfBirth}
                        onChange={(e) => setNewMember({ ...newMember, dateOfBirth: e.target.value })}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberEmail">Email</Label>
                      <Input
                        id="memberEmail"
                        type="email"
                        value={newMember.email}
                        onChange={(e) => setNewMember({ ...newMember, email: e.target.value })}
                        placeholder="Enter email"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="memberPhoneNumber">Phone Number *</Label>
                      <Input
                        id="memberPhoneNumber"
                        value={newMember.phoneNumber}
                        onChange={(e) => {
                          const formatted = formatMalawiPhone(e.target.value);
                          setNewMember({ ...newMember, phoneNumber: formatted });
                        }}
                        placeholder={EXAMPLE_FORMATS.PHONE}
                      />
                    </div>

                    <div className="flex items-center space-x-2">
                      <Checkbox
                        id="memberIsIntern"
                        checked={newMember.isIntern}
                        onCheckedChange={(checked) => setNewMember({ ...newMember, isIntern: checked as boolean })}
                      />
                      <Label htmlFor="memberIsIntern" className="cursor-pointer">Is Intern</Label>
                    </div>

                    <div className="flex items-center space-x-2">
                      <Checkbox
                        id="memberIsPartTime"
                        checked={newMember.isPartTime}
                        onCheckedChange={(checked) => setNewMember({ ...newMember, isPartTime: checked as boolean })}
                      />
                      <Label htmlFor="memberIsPartTime" className="cursor-pointer">Is Part Time</Label>
                    </div>
                  </div>

                  <Button
                    type="button"
                    onClick={addMember}
                    disabled={!newMember.firstName || !newMember.lastName || !newMember.nationalIdNumber || !newMember.phoneNumber}
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
