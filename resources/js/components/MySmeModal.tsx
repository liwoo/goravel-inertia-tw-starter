import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { DISTRICT_NAMES } from '@/constants/districts';
import { NATIONALITY_OPTIONS } from '@/types/nationalities';
import { GENDER_OPTIONS } from '@/types/gender';
import { EDUCATION_OPTIONS } from '@/types/education';
import { MALAWIAN_STATUS_OPTIONS } from '@/types/malawian-status';
import {
  formatMalawiPhone,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';
import { snakefiyKeys, camelifyKeys } from '@/lib/utils';
import { Loader2 } from 'lucide-react';

interface MySmeModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  smeId: number;
}

interface SmeData {
  id?: number;
  name: string;
  sector: string;
  subSector?: string;
  contactPhone: string;
  contactEmail: string;
  website?: string;
  district?: string;
  physicalAddress?: string;
  businessDescription?: string;
}

interface PrimaryOwnerData {
  id?: number;
  firstName: string;
  lastName: string;
  nationalIdNumber: string;
  dateOfBirth?: string;
  gender: string;
  phoneNumber: string;
  email?: string;
}

export const MySmeModal: React.FC<MySmeModalProps> = ({ open, onOpenChange, smeId }) => {
  const [activeTab, setActiveTab] = useState('business');
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Business Details State
  const [smeData, setSmeData] = useState<SmeData>({
    name: '',
    sector: '',
    contactPhone: '',
    contactEmail: ''
  });
  const [smeErrors, setSmeErrors] = useState<Record<string, string>>({});

  // Primary Owner State
  const [ownerData, setOwnerData] = useState<PrimaryOwnerData | null>(null);
  const [ownerErrors, setOwnerErrors] = useState<Record<string, string>>({});

  // Configuration options
  const [sectorOptions, setSectorOptions] = useState<string[]>([]);

  // Fetch SME and Owner Data
  useEffect(() => {
    if (open && smeId) {
      fetchData();
      fetchConfigs();
    }
  }, [open, smeId]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [smeRes, ownerRes] = await Promise.all([
        axios.get(`/api/smes/${smeId}`),
        axios.get(`/api/smes/${smeId}/primary_business_owner`)
      ]);

      const smeDataRaw = smeRes.data.data;
      setSmeData({
        id: smeDataRaw.id,
        name: smeDataRaw.name || '',
        sector: smeDataRaw.sector || '',
        subSector: smeDataRaw.sub_sector || smeDataRaw.subSector || '',
        contactPhone: smeDataRaw.contact_phone || smeDataRaw.contactPhone || '',
        contactEmail: smeDataRaw.contact_email || smeDataRaw.contactEmail || '',
        website: smeDataRaw.website || '',
        district: smeDataRaw.district || '',
        physicalAddress: smeDataRaw.physical_address || smeDataRaw.physicalAddress || '',
        businessDescription: smeDataRaw.business_description || smeDataRaw.businessDescription || ''
      });

      const ownerDataRaw = ownerRes.data.data;
      if (ownerDataRaw) {
        setOwnerData({
          id: ownerDataRaw.id,
          firstName: ownerDataRaw.first_name || ownerDataRaw.firstName || '',
          lastName: ownerDataRaw.last_name || ownerDataRaw.lastName || '',
          nationalIdNumber: ownerDataRaw.national_id_number || ownerDataRaw.nationalIdNumber || '',
          dateOfBirth: ownerDataRaw.date_of_birth || ownerDataRaw.dateOfBirth || '',
          gender: ownerDataRaw.gender || '',
          phoneNumber: ownerDataRaw.phone_number || ownerDataRaw.phoneNumber || '',
          email: ownerDataRaw.email || ''
        });
      }
    } catch (error) {
      console.error('Error fetching data:', error);
      toast.error('Failed to load SME details');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchConfigs = async () => {
    try {
      const sectorRes = await axios.get('/api/configs', {
        params: { config_type: 'Sectors', pageSize: 100, sort: 'name', direction: 'ASC' }
      });
      setSectorOptions(sectorRes.data.data?.data?.map((s: any) => s.name) || []);
    } catch (error) {
      console.error('Error fetching configs:', error);
    }
  };

  const validateBusinessData = (): boolean => {
    const errors: Record<string, string> = {};

    if (!smeData.name?.trim()) {
      errors.name = 'Business name is required';
    }
    if (!smeData.sector?.trim()) {
      errors.sector = 'Sector is required';
    }
    if (!smeData.contactPhone?.trim()) {
      errors.contactPhone = 'Contact phone is required';
    }
    if (!smeData.contactEmail?.trim()) {
      errors.contactEmail = 'Contact email is required';
    }

    setSmeErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const validateOwnerData = (): boolean => {
    if (!ownerData) return false;

    const errors: Record<string, string> = {};

    if (!ownerData.firstName?.trim()) {
      errors.firstName = 'First name is required';
    }
    if (!ownerData.lastName?.trim()) {
      errors.lastName = 'Last name is required';
    }
    if (!ownerData.nationalIdNumber?.trim()) {
      errors.nationalIdNumber = 'National ID is required';
    }
    if (!ownerData.gender?.trim()) {
      errors.gender = 'Gender is required';
    }
    if (!ownerData.phoneNumber?.trim()) {
      errors.phoneNumber = 'Phone number is required';
    }

    setOwnerErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSaveBusiness = async () => {
    if (!validateBusinessData()) {
      toast.error('Please fix the validation errors');
      return;
    }

    setIsSaving(true);
    try {
      const payload = snakefiyKeys({
        name: smeData.name,
        sector: smeData.sector,
        subSector: smeData.subSector,
        contactPhone: smeData.contactPhone,
        contactEmail: smeData.contactEmail,
        website: smeData.website,
        district: smeData.district,
        physicalAddress: smeData.physicalAddress,
        businessDescription: smeData.businessDescription
      });

      await axios.put(`/api/smes/${smeId}`, payload);
      toast.success('Business details updated successfully');
      setSmeErrors({});
    } catch (error: any) {
      console.error('Error updating SME:', error);
      toast.error(error.response?.data?.message || 'Failed to update business details');
    } finally {
      setIsSaving(false);
    }
  };

  const handleSaveOwner = async () => {
    if (!validateOwnerData() || !ownerData?.id) {
      toast.error('Please fix the validation errors');
      return;
    }

    setIsSaving(true);
    try {
      const payload = snakefiyKeys({
        firstName: ownerData.firstName,
        lastName: ownerData.lastName,
        nationalIdNumber: ownerData.nationalIdNumber,
        dateOfBirth: ownerData.dateOfBirth,
        gender: ownerData.gender,
        phoneNumber: ownerData.phoneNumber,
        email: ownerData.email
      });

      await axios.put(`/api/primary-business-owners/${ownerData.id}`, payload);
      toast.success('Owner details updated successfully');
      setOwnerErrors({});
    } catch (error: any) {
      console.error('Error updating owner:', error);
      toast.error(error.response?.data?.message || 'Failed to update owner details');
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    onOpenChange(false);
    setSmeErrors({});
    setOwnerErrors({});
    setActiveTab('business');
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-2xl">My SME Details</DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <>
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="business">Business Details</TabsTrigger>
                <TabsTrigger value="owner">Primary Owner</TabsTrigger>
              </TabsList>

              <TabsContent value="business" className="space-y-4 mt-4">
                <Card>
                  <CardHeader>
                    <CardTitle>Business Information</CardTitle>
                    <CardDescription>Update your business details</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="name">Business Name *</Label>
                        <Input
                          id="name"
                          value={smeData.name}
                          onChange={(e) => setSmeData({ ...smeData, name: e.target.value })}
                          placeholder="Enter business name"
                          className={smeErrors.name ? 'border-destructive' : ''}
                        />
                        {smeErrors.name && <p className="text-sm text-destructive">{smeErrors.name}</p>}
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="sector">Sector *</Label>
                        <Select
                          value={smeData.sector}
                          onValueChange={(value) => setSmeData({ ...smeData, sector: value })}
                        >
                          <SelectTrigger className={smeErrors.sector ? 'border-destructive' : ''}>
                            <SelectValue placeholder="Select sector" />
                          </SelectTrigger>
                          <SelectContent>
                            {sectorOptions.map((sec) => (
                              <SelectItem key={sec} value={sec}>{sec}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {smeErrors.sector && <p className="text-sm text-destructive">{smeErrors.sector}</p>}
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="subSector">Sub Sector</Label>
                        <Input
                          id="subSector"
                          value={smeData.subSector}
                          onChange={(e) => setSmeData({ ...smeData, subSector: e.target.value })}
                          placeholder="Enter sub sector"
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="contactPhone">Contact Phone *</Label>
                        <Input
                          id="contactPhone"
                          value={smeData.contactPhone}
                          onChange={(e) => {
                            const formatted = formatMalawiPhone(e.target.value);
                            setSmeData({ ...smeData, contactPhone: formatted });
                          }}
                          placeholder={EXAMPLE_FORMATS.PHONE}
                          className={smeErrors.contactPhone ? 'border-destructive' : ''}
                        />
                        {smeErrors.contactPhone && <p className="text-sm text-destructive">{smeErrors.contactPhone}</p>}
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="contactEmail">Contact Email *</Label>
                        <Input
                          id="contactEmail"
                          type="email"
                          value={smeData.contactEmail}
                          onChange={(e) => setSmeData({ ...smeData, contactEmail: e.target.value })}
                          placeholder="Enter email"
                          className={smeErrors.contactEmail ? 'border-destructive' : ''}
                        />
                        {smeErrors.contactEmail && <p className="text-sm text-destructive">{smeErrors.contactEmail}</p>}
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="website">Website</Label>
                        <Input
                          id="website"
                          value={smeData.website}
                          onChange={(e) => setSmeData({ ...smeData, website: e.target.value })}
                          placeholder="https://example.com"
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="district">District</Label>
                        <Select
                          value={smeData.district}
                          onValueChange={(value) => setSmeData({ ...smeData, district: value })}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select district" />
                          </SelectTrigger>
                          <SelectContent>
                            {DISTRICT_NAMES.map((dist) => (
                              <SelectItem key={dist} value={dist}>{dist}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="physicalAddress">Physical Address</Label>
                      <Textarea
                        id="physicalAddress"
                        value={smeData.physicalAddress}
                        onChange={(e) => setSmeData({ ...smeData, physicalAddress: e.target.value })}
                        rows={2}
                        placeholder="Enter physical address"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="businessDescription">Business Description</Label>
                      <Textarea
                        id="businessDescription"
                        value={smeData.businessDescription}
                        onChange={(e) => setSmeData({ ...smeData, businessDescription: e.target.value })}
                        rows={3}
                        placeholder="Describe your business"
                      />
                    </div>
                  </CardContent>
                </Card>

                <DialogFooter>
                  <Button variant="outline" onClick={handleCancel} disabled={isSaving}>
                    Cancel
                  </Button>
                  <Button onClick={handleSaveBusiness} disabled={isSaving}>
                    {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    Save Changes
                  </Button>
                </DialogFooter>
              </TabsContent>

              <TabsContent value="owner" className="space-y-4 mt-4">
                {!ownerData ? (
                  <Card>
                    <CardHeader>
                      <CardTitle>Primary Owner</CardTitle>
                      <CardDescription>No owner data available</CardDescription>
                    </CardHeader>
                    <CardContent>
                      <p className="text-muted-foreground">Unable to load primary owner information.</p>
                    </CardContent>
                  </Card>
                ) : (
                  <>
                    <Card>
                      <CardHeader>
                        <CardTitle>Personal Information</CardTitle>
                        <CardDescription>Update primary owner details</CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div className="space-y-2">
                            <Label htmlFor="firstName">First Name *</Label>
                            <Input
                              id="firstName"
                              value={ownerData.firstName}
                              onChange={(e) => setOwnerData({ ...ownerData, firstName: e.target.value })}
                              placeholder="Enter first name"
                              className={ownerErrors.firstName ? 'border-destructive' : ''}
                            />
                            {ownerErrors.firstName && <p className="text-sm text-destructive">{ownerErrors.firstName}</p>}
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="lastName">Last Name *</Label>
                            <Input
                              id="lastName"
                              value={ownerData.lastName}
                              onChange={(e) => setOwnerData({ ...ownerData, lastName: e.target.value })}
                              placeholder="Enter last name"
                              className={ownerErrors.lastName ? 'border-destructive' : ''}
                            />
                            {ownerErrors.lastName && <p className="text-sm text-destructive">{ownerErrors.lastName}</p>}
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="nationalIdNumber">National ID *</Label>
                            <Input
                              id="nationalIdNumber"
                              value={ownerData.nationalIdNumber}
                              onChange={(e) => setOwnerData({ ...ownerData, nationalIdNumber: e.target.value.toUpperCase() })}
                              placeholder={EXAMPLE_FORMATS.NATIONAL_ID}
                              maxLength={8}
                              className={ownerErrors.nationalIdNumber ? 'border-destructive' : ''}
                            />
                            {ownerErrors.nationalIdNumber && <p className="text-sm text-destructive">{ownerErrors.nationalIdNumber}</p>}
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="dateOfBirth">Date of Birth</Label>
                            <Input
                              id="dateOfBirth"
                              type="date"
                              value={ownerData.dateOfBirth}
                              onChange={(e) => setOwnerData({ ...ownerData, dateOfBirth: e.target.value })}
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="gender">Gender *</Label>
                            <Select
                              value={ownerData.gender}
                              onValueChange={(value) => setOwnerData({ ...ownerData, gender: value })}
                            >
                              <SelectTrigger className={ownerErrors.gender ? 'border-destructive' : ''}>
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
                            {ownerErrors.gender && <p className="text-sm text-destructive">{ownerErrors.gender}</p>}
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="phoneNumber">Phone Number *</Label>
                            <Input
                              id="phoneNumber"
                              value={ownerData.phoneNumber}
                              onChange={(e) => {
                                const formatted = formatMalawiPhone(e.target.value);
                                setOwnerData({ ...ownerData, phoneNumber: formatted });
                              }}
                              placeholder={EXAMPLE_FORMATS.PHONE}
                              className={ownerErrors.phoneNumber ? 'border-destructive' : ''}
                            />
                            {ownerErrors.phoneNumber && <p className="text-sm text-destructive">{ownerErrors.phoneNumber}</p>}
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input
                              id="email"
                              type="email"
                              value={ownerData.email}
                              onChange={(e) => setOwnerData({ ...ownerData, email: e.target.value })}
                              placeholder="Enter email"
                            />
                          </div>
                        </div>
                      </CardContent>
                    </Card>

                    <DialogFooter>
                      <Button variant="outline" onClick={handleCancel} disabled={isSaving}>
                        Cancel
                      </Button>
                      <Button onClick={handleSaveOwner} disabled={isSaving}>
                        {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Save Changes
                      </Button>
                    </DialogFooter>
                  </>
                )}
              </TabsContent>
            </Tabs>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
};
