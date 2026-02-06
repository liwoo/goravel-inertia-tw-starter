import React, { useState, useEffect, useMemo } from 'react';
import axios from 'axios';
import { toast } from 'sonner';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from '@/components/ui/sheet';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { DISTRICT_NAMES } from '@/constants/districts';
import { SECTOR_NAMES } from '@/constants/sectors';
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

// Helper to parse date from various formats (carbon.DateTime returns "2003-07-16 16:22:02")
const parseDateForInput = (dateValue: any): string => {
  if (!dateValue) return '';

  // If it's already in YYYY-MM-DD format, return as-is
  if (typeof dateValue === 'string') {
    // Handle carbon.DateTime format "2003-07-16 16:22:02"
    if (dateValue.includes(' ')) {
      return dateValue.split(' ')[0];
    }
    // Handle ISO format "2003-07-16T16:22:02Z"
    if (dateValue.includes('T')) {
      return dateValue.split('T')[0];
    }
    // Already in YYYY-MM-DD format
    if (/^\d{4}-\d{2}-\d{2}$/.test(dateValue)) {
      return dateValue;
    }
    // Try to parse and format
    const parsed = new Date(dateValue);
    if (!isNaN(parsed.getTime())) {
      return parsed.toISOString().split('T')[0];
    }
  }

  // Handle date object
  if (dateValue instanceof Date && !isNaN(dateValue.getTime())) {
    return dateValue.toISOString().split('T')[0];
  }

  return '';
};

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

  // Fetch SME and Owner Data when modal opens
  useEffect(() => {
    if (open && smeId) {
      console.log('MySmeModal: Fetching data for smeId:', smeId);
      fetchData();
    }
  }, [open, smeId]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      console.log('MySmeModal: Making API calls for my-sme');
      // Use /api/my-sme endpoints which bypass admin permissions
      const [smeRes, ownerRes] = await Promise.all([
        axios.get('/api/my-sme'),
        axios.get('/api/my-sme/primary-owner')
      ]);

      console.log('MySmeModal: MSME Response:', smeRes.data);
      console.log('MySmeModal: Owner Response:', ownerRes.data);

      // Handle nested data structure - API returns { success, data, message }
      const smeDataRaw = smeRes.data?.data || smeRes.data;
      if (smeDataRaw) {
        setSmeData({
          id: smeDataRaw.id || smeDataRaw.ID,
          name: smeDataRaw.name || smeDataRaw.Name || '',
          sector: smeDataRaw.sector || smeDataRaw.Sector || '',
          subSector: smeDataRaw.sub_sector || smeDataRaw.subSector || smeDataRaw.SubSector || '',
          contactPhone: smeDataRaw.contact_phone || smeDataRaw.contactPhone || smeDataRaw.ContactPhone || '',
          contactEmail: smeDataRaw.contact_email || smeDataRaw.contactEmail || smeDataRaw.ContactEmail || '',
          website: smeDataRaw.website || smeDataRaw.Website || '',
          district: smeDataRaw.district || smeDataRaw.District || '',
          physicalAddress: smeDataRaw.physical_address || smeDataRaw.physicalAddress || smeDataRaw.PhysicalAddress || '',
          businessDescription: smeDataRaw.business_description || smeDataRaw.businessDescription || smeDataRaw.BusinessDescription || ''
        });
        console.log('MySmeModal: MSME Data set:', smeDataRaw);
      }

      const ownerDataRaw = ownerRes.data?.data || ownerRes.data;
      if (ownerDataRaw) {
        const rawDateOfBirth = ownerDataRaw.date_of_birth || ownerDataRaw.dateOfBirth || ownerDataRaw.DateOfBirth;
        setOwnerData({
          id: ownerDataRaw.id || ownerDataRaw.ID,
          firstName: ownerDataRaw.first_name || ownerDataRaw.firstName || ownerDataRaw.FirstName || '',
          lastName: ownerDataRaw.last_name || ownerDataRaw.lastName || ownerDataRaw.LastName || '',
          nationalIdNumber: ownerDataRaw.national_id_number || ownerDataRaw.nationalIdNumber || ownerDataRaw.NationalIdNumber || '',
          dateOfBirth: parseDateForInput(rawDateOfBirth),
          gender: ownerDataRaw.gender || ownerDataRaw.Gender || '',
          phoneNumber: ownerDataRaw.phone_number || ownerDataRaw.phoneNumber || ownerDataRaw.PhoneNumber || '',
          email: ownerDataRaw.email || ownerDataRaw.Email || ''
        });
        console.log('MySmeModal: Owner Data set:', ownerDataRaw, 'Parsed dateOfBirth:', parseDateForInput(rawDateOfBirth));
      }
    } catch (error) {
      console.error('MySmeModal: Error fetching data:', error);
      toast.error('Failed to load MSME details');
    } finally {
      setIsLoading(false);
    }
  };

  const validateBusinessData = (): boolean => {
    const errors: Record<string, string> = {};

    // Note: name and contactEmail are disabled fields, not validated
    if (!smeData.sector?.trim()) {
      errors.sector = 'Sector is required';
    }
    if (!smeData.contactPhone?.trim()) {
      errors.contactPhone = 'Contact phone is required';
    }

    setSmeErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const validateOwnerData = (): boolean => {
    if (!ownerData) return false;

    const errors: Record<string, string> = {};

    // Note: firstName, lastName, nationalIdNumber, and email are disabled fields, not validated
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
      // Note: name and contactEmail are not sent as they cannot be changed
      const payload = snakefiyKeys({
        sector: smeData.sector,
        subSector: smeData.subSector,
        contactPhone: smeData.contactPhone,
        website: smeData.website,
        district: smeData.district,
        physicalAddress: smeData.physicalAddress,
        businessDescription: smeData.businessDescription
      });

      // Use /api/my-sme endpoint which bypasses admin permissions
      await axios.put('/api/my-sme', payload);
      toast.success('Business details updated successfully');
      setSmeErrors({});
    } catch (error: any) {
      console.error('Error updating MSME:', error);
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
      // Note: firstName, lastName, nationalIdNumber, and email are not sent as they cannot be changed
      const payload = snakefiyKeys({
        dateOfBirth: ownerData.dateOfBirth,
        gender: ownerData.gender,
        phoneNumber: ownerData.phoneNumber
      });

      // Use /api/my-sme/primary-owner endpoint which bypasses admin permissions
      await axios.put('/api/my-sme/primary-owner', payload);
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

  // Ensure current sector is included in options for display
  const effectiveSectorOptions = useMemo(() => {
    const options: string[] = [...SECTOR_NAMES];
    // Include current sector if not in standard list (legacy data)
    if (smeData.sector && !options.includes(smeData.sector)) {
      options.unshift(smeData.sector);
    }
    return options;
  }, [smeData.sector]);

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-xl md:max-w-2xl overflow-hidden flex flex-col">
        <SheetHeader className="flex-shrink-0 pb-4">
          <SheetTitle className="text-2xl">My MSME Details</SheetTitle>
        </SheetHeader>

        {isLoading ? (
          <div className="flex items-center justify-center flex-1">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <div className="flex-1 overflow-hidden flex flex-col">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col overflow-hidden">
              <TabsList className="grid w-full grid-cols-2 flex-shrink-0">
                <TabsTrigger value="business">Business Details</TabsTrigger>
                <TabsTrigger value="owner">Primary Owner</TabsTrigger>
              </TabsList>

              <TabsContent value="business" className="flex-1 overflow-y-auto space-y-4 mt-4 pr-2">
                <Card>
                  <CardHeader>
                    <CardTitle>Business Information</CardTitle>
                    <CardDescription>Update your business details</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="name">Business Name</Label>
                        <Input
                          id="name"
                          value={smeData.name}
                          disabled
                          className="bg-muted cursor-not-allowed"
                        />
                        <p className="text-xs text-muted-foreground">Business name cannot be changed</p>
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
                            {effectiveSectorOptions.map((sec) => (
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
                        <Label htmlFor="contactEmail">Contact Email</Label>
                        <Input
                          id="contactEmail"
                          type="email"
                          value={smeData.contactEmail}
                          disabled
                          className="bg-muted cursor-not-allowed"
                        />
                        <p className="text-xs text-muted-foreground">Email cannot be changed</p>
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

                <SheetFooter className="pt-4 border-t">
                  <Button variant="outline" onClick={handleCancel} disabled={isSaving}>
                    Cancel
                  </Button>
                  <Button onClick={handleSaveBusiness} disabled={isSaving}>
                    {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    Save Changes
                  </Button>
                </SheetFooter>
              </TabsContent>

              <TabsContent value="owner" className="flex-1 overflow-y-auto space-y-4 mt-4 pr-2">
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
                            <Label htmlFor="firstName">First Name</Label>
                            <Input
                              id="firstName"
                              value={ownerData.firstName}
                              disabled
                              className="bg-muted cursor-not-allowed"
                            />
                            <p className="text-xs text-muted-foreground">Name cannot be changed</p>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="lastName">Last Name</Label>
                            <Input
                              id="lastName"
                              value={ownerData.lastName}
                              disabled
                              className="bg-muted cursor-not-allowed"
                            />
                            <p className="text-xs text-muted-foreground">Name cannot be changed</p>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="nationalIdNumber">National ID</Label>
                            <Input
                              id="nationalIdNumber"
                              value={ownerData.nationalIdNumber}
                              disabled
                              className="bg-muted cursor-not-allowed"
                            />
                            <p className="text-xs text-muted-foreground">National ID cannot be changed</p>
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
                              disabled
                              className="bg-muted cursor-not-allowed"
                            />
                            <p className="text-xs text-muted-foreground">Email cannot be changed</p>
                          </div>
                        </div>
                      </CardContent>
                    </Card>

                    <SheetFooter className="pt-4 border-t">
                      <Button variant="outline" onClick={handleCancel} disabled={isSaving}>
                        Cancel
                      </Button>
                      <Button onClick={handleSaveOwner} disabled={isSaving}>
                        {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Save Changes
                      </Button>
                    </SheetFooter>
                  </>
                )}
              </TabsContent>
            </Tabs>
          </div>
        )}
      </SheetContent>
    </Sheet>
  );
};
