import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { CrudFormProps } from '@/types/crud';
import { SmeCreateData } from '@/types/sme';
import {
  MapPin, FileText, User, FolderOpen, Phone, Mail, Building2, Globe,
  ChevronRight, ChevronLeft, Plus, Trash2, Calendar, X, Check
} from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { DISTRICTS, DISTRICT_NAMES } from '@/constants/districts';
import { NATIONALITY_OPTIONS } from '@/types/nationalities';
import { GENDER_OPTIONS } from '@/types/gender';
import { EDUCATION_OPTIONS } from '@/types/education';
import { MALAWIAN_STATUS_OPTIONS } from '@/types/malawian-status';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import axios from 'axios';
import {snakefiyKeys} from "@/lib/utils";
import {
  validateMalawiPhone,
  validateBusinessRegistration,
  validateTIN,
  validateNationalID,
  formatMalawiPhone,
  formatBusinessRegistration,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';
import {
  FormalisationFormData,
  DEFAULT_FORMALISATION_DATA
} from '@/types/business_formalisation';
import { Switch } from '@/components/ui/switch';

export interface InitialSmeData {
  // Business Information (Step 1)
  name?: string;
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  contactPhone?: string;
  contactEmail?: string;
  physicalAddress?: string;
  postalAddress?: string;
  district?: string;
  traditionalAuthority?: string;
  // Primary Owner (Step 2)
  ownerFirstName?: string;
  ownerLastName?: string;
  ownerOtherNames?: string;
  ownerNationality?: string;
  ownerNationalIdNumber?: string;
  ownerDateOfBirth?: string;
  ownerGender?: string;
  ownerEducationLevel?: string;
  ownerMalawianStatus?: string;
  ownerHasSpecialNeeds?: boolean;
  ownerPhoneNumber?: string;
  ownerLandlineNumber?: string;
  ownerEmail?: string;
  ownerPhysicalAddress?: string;
  ownerPostalAddress?: string;
  ownerDistrict?: string;
  ownerTraditionalAuthority?: string;
  ownerAltContactName?: string;
  ownerAltContactRelationship?: string;
  ownerAltContactPhone?: string;
}

interface SmeCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
  initialData?: InitialSmeData;
  // Callback for when SME is created - provides the SME data object
  // Use this instead of onSuccess when you need the SME data (e.g., for auto-approval workflow)
  onSmeCreated?: (sme: { id: number; name: string; usme_number: string }) => void;
}

interface PrimaryOwnerData {
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
  // region is inferred from district
  district?: string;
  traditionalAuthority?: string;
  altContactName?: string;
  altContactRelationship?: string;
  altContactPhone?: string;
}

interface AdditionalMemberData {
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

const STEPS = [
  { id: 1, name: 'Business Information', description: 'Basic SME details' },
  { id: 2, name: 'Primary Owner', description: 'Business owner information' },
  { id: 3, name: 'Formalization', description: 'Business formalization status' },
  { id: 4, name: 'Additional Members', description: 'Optional team members' },
  { id: 5, name: 'Review', description: 'Review and submit' },
];

export const SmeCreateForm = forwardRef<any, SmeCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving,
  initialData,
  onSmeCreated
}, ref) => {
  const [currentStep, setCurrentStep] = useState(1);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false); // Guard against double submission

  // Step 1: Business Information
  const [businessData, setBusinessData] = useState<SmeCreateData>({
    // usmeNumber is auto-generated - no longer needed in form
    name: '',
    registrationNumber: undefined,
    taxIdentificationNumber: undefined,
    operationalStartDate: '',
    businessCategory: '',
    sector: '',
    subSector: undefined,
    businessDescription: undefined,
    contactPhone: '',
    contactEmail: '',
    physicalAddress: undefined,
    postalAddress: undefined,
    website: undefined,
    // region is inferred from district
    district: undefined,
    traditionalAuthority: undefined,
    businessImprovementAspects: [],
    businessAccessedFinancing: [],
  });

  // Step 2: Primary Owner
  const [primaryOwner, setPrimaryOwner] = useState<PrimaryOwnerData>({
    firstName: '',
    lastName: '',
    otherNames: undefined,
    nationality: '',
    nationalIdNumber: '',
    dateOfBirth: undefined,
    gender: '',
    educationLevel: '',
    malawianStatus: '',
    hasSpecialNeeds: false,
    phoneNumber: '',
    landlineNumber: undefined,
    email: undefined,
    physicalAddress: undefined,
    postalAddress: undefined,
    // region is inferred from district
    district: undefined,
    traditionalAuthority: undefined,
    altContactName: undefined,
    altContactRelationship: undefined,
    altContactPhone: undefined,
  });

  // Step 3: Formalization
  const [formalizationData, setFormalizationData] = useState<FormalisationFormData>(DEFAULT_FORMALISATION_DATA);

  // Step 4: Additional Members
  const [additionalMembers, setAdditionalMembers] = useState<AdditionalMemberData[]>([]);
  const [activeTab, setActiveTab] = useState<string>('new-member');
  const [newMember, setNewMember] = useState<AdditionalMemberData>({
    firstName: '',
    lastName: '',
    otherNames: undefined,
    nationality: '',
    nationalIdNumber: '',
    dateOfBirth: undefined,
    email: undefined,
    phoneNumber: '',
    isIntern: false,
    isPartTime: false,
  });

  // Multi-input helpers for arrays
  const [improvementAspect, setImprovementAspect] = useState('');
  const [financingSource, setFinancingSource] = useState('');

  // Config options from API
  const [improvementOptions, setImprovementOptions] = useState<string[]>([]);
  const [financingOptions, setFinancingOptions] = useState<string[]>([]);
  const [sectorOptions, setSectorOptions] = useState<string[]>([]);
  const [categoryOptions, setCategoryOptions] = useState<string[]>([]);
  // District options use static DISTRICT_NAMES constant instead of API
  const [improvementPopoverOpen, setImprovementPopoverOpen] = useState(false);
  const [financingPopoverOpen, setFinancingPopoverOpen] = useState(false);
  const [loadingConfigs, setLoadingConfigs] = useState(false);

  const progress = (currentStep / STEPS.length) * 100;

  // Apply initial data when provided (e.g., from application)
  useEffect(() => {
    if (initialData) {
      // Prefill business data (Step 1)
      setBusinessData(prev => ({
        ...prev,
        name: initialData.name || prev.name,
        registrationNumber: initialData.registrationNumber || prev.registrationNumber,
        taxIdentificationNumber: initialData.taxIdentificationNumber || prev.taxIdentificationNumber,
        contactPhone: initialData.contactPhone || prev.contactPhone,
        contactEmail: initialData.contactEmail || prev.contactEmail,
        physicalAddress: initialData.physicalAddress || prev.physicalAddress,
        postalAddress: initialData.postalAddress || prev.postalAddress,
        district: initialData.district || prev.district,
        traditionalAuthority: initialData.traditionalAuthority || prev.traditionalAuthority,
      }));

      // Prefill primary owner data (Step 2)
      setPrimaryOwner(prev => ({
        ...prev,
        firstName: initialData.ownerFirstName || prev.firstName,
        lastName: initialData.ownerLastName || prev.lastName,
        otherNames: initialData.ownerOtherNames || prev.otherNames,
        nationality: initialData.ownerNationality || prev.nationality,
        nationalIdNumber: initialData.ownerNationalIdNumber || prev.nationalIdNumber,
        dateOfBirth: initialData.ownerDateOfBirth || prev.dateOfBirth,
        gender: initialData.ownerGender || prev.gender,
        educationLevel: initialData.ownerEducationLevel || prev.educationLevel,
        malawianStatus: initialData.ownerMalawianStatus || prev.malawianStatus,
        hasSpecialNeeds: initialData.ownerHasSpecialNeeds ?? prev.hasSpecialNeeds,
        phoneNumber: initialData.ownerPhoneNumber || prev.phoneNumber,
        landlineNumber: initialData.ownerLandlineNumber || prev.landlineNumber,
        email: initialData.ownerEmail || prev.email,
        physicalAddress: initialData.ownerPhysicalAddress || prev.physicalAddress,
        postalAddress: initialData.ownerPostalAddress || prev.postalAddress,
        district: initialData.ownerDistrict || prev.district,
        traditionalAuthority: initialData.ownerTraditionalAuthority || prev.traditionalAuthority,
        altContactName: initialData.ownerAltContactName || prev.altContactName,
        altContactRelationship: initialData.ownerAltContactRelationship || prev.altContactRelationship,
        altContactPhone: initialData.ownerAltContactPhone || prev.altContactPhone,
      }));
    }
  }, [initialData]);

  // Fetch config options on mount
  useEffect(() => {
    const fetchConfigOptions = async () => {
      setLoadingConfigs(true);
      try {
        // Fetch all config types in parallel (districts use static constant)
        const [improvementRes, financingRes, sectorRes, categoryRes] = await Promise.all([
          axios.get('/api/configs', {
            params: {
              config_type: 'Improvement Aspects',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          }),
          axios.get('/api/configs', {
            params: {
              config_type: 'Financing',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          }),
          axios.get('/api/configs', {
            params: {
              config_type: 'Sectors',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          }),
          axios.get('/api/configs', {
            params: {
              config_type: 'Business Categories',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          })
        ]);

        // Extract config names from the paginated responses
        const improvementData = improvementRes.data.data?.data || [];
        const financingData = financingRes.data.data?.data || [];
        const sectorData = sectorRes.data.data?.data || [];
        const categoryData = categoryRes.data.data?.data || [];
        // Districts use static DISTRICT_NAMES constant

        const improvements = improvementData.map((config: any) => config.name);
        const financing = financingData.map((config: any) => config.name);
        const sectors = sectorData.map((config: any) => config.name);
        const categories = categoryData.map((config: any) => config.name);

        setImprovementOptions(improvements);
        setFinancingOptions(financing);
        setSectorOptions(sectors);
        setCategoryOptions(categories);
        // Districts use static DISTRICT_NAMES constant
      } catch (error) {
        console.error('Error fetching config options:', error);
      } finally {
        setLoadingConfigs(false);
      }
    };

    fetchConfigOptions();
  }, []);

  const validateStep1 = () => {
    const newErrors: Record<string, string> = {};
    // usmeNumber is auto-generated - no validation needed
    if (!businessData.name) newErrors.name = 'Business Name is required';
    if (!businessData.businessCategory) newErrors.businessCategory = 'Business Category is required';
    if (!businessData.sector) newErrors.sector = 'Sector is required';
    if (!businessData.contactPhone) {
      newErrors.contactPhone = 'Contact Phone is required';
    } else if (!validateMalawiPhone(businessData.contactPhone)) {
      newErrors.contactPhone = VALIDATION_MESSAGES.PHONE;
    }
    if (!businessData.contactEmail) newErrors.contactEmail = 'Contact Email is required';

    // Validate business registration if provided
    if (businessData.registrationNumber && !validateBusinessRegistration(businessData.registrationNumber)) {
      newErrors.registrationNumber = VALIDATION_MESSAGES.BUSINESS_REG;
    }

    // Validate TIN if provided
    if (businessData.taxIdentificationNumber && !validateTIN(businessData.taxIdentificationNumber)) {
      newErrors.taxIdentificationNumber = VALIDATION_MESSAGES.TIN;
    }

    if (businessData.businessImprovementAspects.length === 0) {
      newErrors.businessImprovementAspects = 'At least one improvement aspect is required';
    }
    if (businessData.businessAccessedFinancing.length === 0) {
      newErrors.businessAccessedFinancing = 'At least one financing source is required';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const validateStep2 = () => {
    const newErrors: Record<string, string> = {};
    if (!primaryOwner.firstName) newErrors.firstName = 'First Name is required';
    if (!primaryOwner.lastName) newErrors.lastName = 'Last Name is required';
    if (!primaryOwner.nationality) newErrors.nationality = 'Nationality is required';
    if (!primaryOwner.nationalIdNumber) {
      newErrors.nationalIdNumber = 'National ID is required';
    } else if (!validateNationalID(primaryOwner.nationalIdNumber)) {
      newErrors.nationalIdNumber = VALIDATION_MESSAGES.NATIONAL_ID;
    }
    if (!primaryOwner.dateOfBirth) newErrors.dateOfBirth = 'Date of Birth is required';
    if (!primaryOwner.gender) newErrors.gender = 'Gender is required';
    if (!primaryOwner.educationLevel) newErrors.educationLevel = 'Education Level is required';
    if (!primaryOwner.malawianStatus) newErrors.malawianStatus = 'Malawian Status is required';
    if (!primaryOwner.phoneNumber) {
      newErrors.phoneNumber = 'Phone Number is required';
    } else if (!validateMalawiPhone(primaryOwner.phoneNumber)) {
      newErrors.phoneNumber = VALIDATION_MESSAGES.PHONE;
    }

    // Validate landline if provided
    if (primaryOwner.landlineNumber && !validateMalawiPhone(primaryOwner.landlineNumber)) {
      newErrors.landlineNumber = VALIDATION_MESSAGES.PHONE;
    }

    // Validate alternative contact phone if provided
    if (primaryOwner.altContactPhone && !validateMalawiPhone(primaryOwner.altContactPhone)) {
      newErrors.altContactPhone = VALIDATION_MESSAGES.PHONE;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const validateStep3 = () => {
    const newErrors: Record<string, string> = {};

    // Validate annual turnover if provided
    if (formalizationData.annualTurnover && isNaN(parseFloat(formalizationData.annualTurnover))) {
      newErrors.annualTurnover = 'Annual turnover must be a valid number';
    }

    // Validate estimated value of assets if provided
    if (formalizationData.estimatedValueOfAssets && isNaN(parseFloat(formalizationData.estimatedValueOfAssets))) {
      newErrors.estimatedValueOfAssets = 'Estimated value of assets must be a valid number';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const addImprovementAspect = (value?: string) => {
    const aspectToAdd = value || improvementAspect.trim();
    if (aspectToAdd && !businessData.businessImprovementAspects.includes(aspectToAdd)) {
      setBusinessData({
        ...businessData,
        businessImprovementAspects: [...businessData.businessImprovementAspects, aspectToAdd]
      });
      setImprovementAspect(''); // Clear search input
      // Keep popover open for multi-select
    }
  };

  const removeImprovementAspect = (aspectToRemove: string) => {
    setBusinessData({
      ...businessData,
      businessImprovementAspects: businessData.businessImprovementAspects.filter(a => a !== aspectToRemove)
    });
  };

  const addFinancingSource = (value?: string) => {
    const sourceToAdd = value || financingSource.trim();
    if (sourceToAdd && !businessData.businessAccessedFinancing.includes(sourceToAdd)) {
      setBusinessData({
        ...businessData,
        businessAccessedFinancing: [...businessData.businessAccessedFinancing, sourceToAdd]
      });
      setFinancingSource(''); // Clear search input
      // Keep popover open for multi-select
    }
  };

  const removeFinancingSource = (sourceToRemove: string) => {
    setBusinessData({
      ...businessData,
      businessAccessedFinancing: businessData.businessAccessedFinancing.filter(s => s !== sourceToRemove)
    });
  };

  const removeSmeOnRelationshipError = (smeId: number) => {
    // Implement logic to remove SME if there is a relationship error
    // This could involve calling an API endpoint to delete the SME by ID
    axios.delete(`/api/smes/${smeId}`)
      .then(() => {
        console.log(`SME with ID ${smeId} removed due to relationship error.`);
      })
      .catch((error) => {
        console.error(`Failed to remove SME with ID ${smeId}:`, error);
      });
  }

  const addMember = () => {
    if (newMember.firstName && newMember.lastName && newMember.nationalIdNumber && newMember.phoneNumber) {
      const newMembersList = [...additionalMembers, newMember];
      setAdditionalMembers(newMembersList);
      // Switch to the newly added member's tab
      setActiveTab(`member-${newMembersList.length - 1}`);
      setNewMember({
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
      });
    }
  };

  const removeMember = (index: number) => {
    setAdditionalMembers(additionalMembers.filter((_, i) => i !== index));
    // Switch to new member tab after removing a member
    setActiveTab('new-member');
  };

  const handleNext = () => {
    let isValid = true;

    if (currentStep === 1) {
      isValid = validateStep1();
    } else if (currentStep === 2) {
      isValid = validateStep2();
    } else if (currentStep === 3) {
      isValid = validateStep3();
    }

    if (isValid && currentStep < STEPS.length) {
      setCurrentStep(currentStep + 1);
      setErrors({});
    }
  };

  const handlePrevious = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
      setErrors({});
    }
  };

  const handleSubmit = async () => {
    // Prevent double submission
    if (isSubmitting) {
      console.log('Submission already in progress, ignoring duplicate call');
      return;
    }

    setIsSubmitting(true);
    setIsSaving?.(true);
    setErrors({});

    try {
      // Step 1: Create SME
      const smeResponse = await fetch('/api/smes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(snakefiyKeys(businessData)),
      });

      if (!smeResponse.ok) {
        const errorData = await smeResponse.json().catch(() => ({}));
        onError?.(errorData);
        return;
      }

      const smeData = await smeResponse.json();
      const smeId = smeData.data.id;


      // Step 2: Create Primary Business Owner
      const primaryOwnerResponse = await fetch('/api/primary-business-owners', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({
          ...snakefiyKeys(primaryOwner),
          "sme_id": smeId,
        }),
      });

      if (!primaryOwnerResponse.ok) {
        const errorData = await primaryOwnerResponse.json().catch(() => ({}));
        onError?.(errorData);
        removeSmeOnRelationshipError(smeId);
        return;
      }

      // Step 3: Create or Update Business Formalization (upsert to prevent duplicates)
      // Note: A BusinessFormalisation may have already been auto-created by the PrimaryBusinessOwner's
      // AfterCreate hook which triggers CalculateFormalisationScore. Using upsert ensures we update
      // any existing record rather than creating a duplicate.
      const formalizationPayload = {
        sme_id: smeId,
        has_bank_account: formalizationData.hasBankAccount,
        has_tax_clarification: formalizationData.hasTaxClarification,
        is_registered_for_vat: formalizationData.isRegisteredForVat,
        is_member_of_association: formalizationData.isMemberOfAssociation,
        is_affiliated: formalizationData.isAffiliated,
        has_export_license: formalizationData.hasExportLicense,
        has_accessed_bds: formalizationData.hasAccessedBds,
        annual_turnover: formalizationData.annualTurnover ? parseFloat(formalizationData.annualTurnover) : 0,
        estimated_value_of_assets: formalizationData.estimatedValueOfAssets ? parseFloat(formalizationData.estimatedValueOfAssets) : 0,
      };

      const formalizationResponse = await fetch('/api/business-formalisations/upsert', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formalizationPayload),
      });

      if (!formalizationResponse.ok) {
        const errorData = await formalizationResponse.json().catch(() => ({}));
        console.error('Failed to save formalization data:', errorData);
        // Continue with the process even if formalization fails (it's optional data)
      }

      // Step 4: Create Additional Business Members (if any)
      if (additionalMembers.length > 0) {
        const memberPromises = additionalMembers.map(member =>
          fetch('/api/additional_business_members', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'X-Requested-With': 'XMLHttpRequest',
            },
            body: JSON.stringify({
              ...snakefiyKeys(member),
              "sme_id": smeId,
            }),
          })
        );

        await Promise.all(memberPromises);
      }

      // Call onSuccess with a string message (for CrudPage toast display)
      onSuccess?.('SME created successfully');

      // Call onSmeCreated with the SME data object (for approval workflow)
      if (onSmeCreated) {
        onSmeCreated({
          id: smeData.data.id,
          name: smeData.data.name,
          usme_number: smeData.data.usme_number,
        });
      }
    } catch (error: any) {
      onError?.(error.message || 'Failed to create SME');
    } finally {
      setIsSubmitting(false);
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit: () => {
      if (currentStep === STEPS.length) {
        handleSubmit();
      } else {
        handleNext();
      }
    }
  }));

  return (
    <div className="space-y-6">
      {/* Progress Bar */}
      <div className="space-y-2">
        <div className="flex justify-between items-center">
          <div className="space-y-1">
            <h3 className="text-lg font-semibold">{STEPS[currentStep - 1].name}</h3>
            <p className="text-sm text-muted-foreground">{STEPS[currentStep - 1].description}</p>
          </div>
          <Badge variant="outline" className="text-xs">
            Step {currentStep} of {STEPS.length}
          </Badge>
        </div>
        <Progress value={progress} className="h-2" />
      </div>

      {/* Step Indicators */}
      <div className="flex justify-between">
        {STEPS.map((step, index) => (
          <div
            key={step.id}
            className={`flex items-center ${index < STEPS.length - 1 ? 'flex-1' : ''}`}
          >
            <div className="flex flex-col items-center">
              <div
                className={`w-10 h-10 rounded-full flex items-center justify-center border-2 ${
                  currentStep > step.id
                    ? 'bg-primary border-primary text-primary-foreground'
                    : currentStep === step.id
                    ? 'border-primary text-primary'
                    : 'border-muted text-muted-foreground'
                }`}
              >
                {currentStep > step.id ? (
                  <Check className="h-5 w-5" />
                ) : (
                  <span className="text-sm font-semibold">{step.id}</span>
                )}
              </div>
              <span className="text-xs mt-2 text-center hidden sm:block">{step.name}</span>
            </div>
            {index < STEPS.length - 1 && (
              <div className={`flex-1 h-0.5 mx-2 ${currentStep > step.id ? 'bg-primary' : 'bg-muted'}`} />
            )}
          </div>
        ))}
      </div>

      <Separator />

      <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
        {/* Step 1: Business Information */}
        {currentStep === 1 && (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Basic Information</CardTitle>
                <CardDescription>Enter the core business details</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* USME Number is auto-generated - removed from form */}

                  <div className="space-y-2">
                    <Label htmlFor="name">Business Name *</Label>
                    <Input
                      id="name"
                      value={businessData.name}
                      onChange={(e) => setBusinessData({ ...businessData, name: e.target.value })}
                      placeholder="Enter business name"
                      className={errors.name ? 'border-destructive' : ''}
                    />
                    {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="registrationNumber">Business Registration Number</Label>
                    <Input
                      id="registrationNumber"
                      value={businessData.registrationNumber}
                      onChange={(e) => {
                        const formatted = formatBusinessRegistration(e.target.value);
                        setBusinessData({ ...businessData, registrationNumber: formatted });
                      }}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.BUSINESS_REG}`}
                      className={errors.registrationNumber ? 'border-destructive' : ''}
                    />
                    {errors.registrationNumber && <p className="text-sm text-destructive">{errors.registrationNumber}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="taxIdentificationNumber">Tax ID Number (TIN)</Label>
                    <Input
                      id="taxIdentificationNumber"
                      value={businessData.taxIdentificationNumber}
                      onChange={(e) => setBusinessData({ ...businessData, taxIdentificationNumber: e.target.value })}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.TIN}`}
                      maxLength={8}
                      className={errors.taxIdentificationNumber ? 'border-destructive' : ''}
                    />
                    {errors.taxIdentificationNumber && <p className="text-sm text-destructive">{errors.taxIdentificationNumber}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="operationalStartDate">Operational Start Date</Label>
                    <Input
                      id="operationalStartDate"
                      type="date"
                      value={businessData.operationalStartDate}
                      onChange={(e) => setBusinessData({ ...businessData, operationalStartDate: e.target.value })}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="businessCategory">Business Category *</Label>
                    <Select
                      value={businessData.businessCategory}
                      onValueChange={(value) => {
                        console.log('value selected:', value);
                        setBusinessData({...businessData, businessCategory: value});
                      }}
                    >
                      <SelectTrigger className={errors.businessCategory ? 'border-destructive' : ''}>
                        <SelectValue placeholder="Select business category" />
                      </SelectTrigger>
                      <SelectContent>
                        {categoryOptions.map((category) => (
                          <SelectItem key={category} value={category}>
                            {category}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.businessCategory && <p className="text-sm text-destructive">{errors.businessCategory}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="sector">Sector *</Label>
                    <Select
                      value={businessData.sector}
                      onValueChange={(value) => setBusinessData({ ...businessData, sector: value })}
                    >
                      <SelectTrigger className={errors.sector ? 'border-destructive' : ''}>
                        <SelectValue placeholder="Select sector" />
                      </SelectTrigger>
                      <SelectContent>
                        {sectorOptions.map((sector) => (
                          <SelectItem key={sector} value={sector}>
                            {sector}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.sector && <p className="text-sm text-destructive">{errors.sector}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="subSector">Sub Sector</Label>
                    <Input
                      id="subSector"
                      value={businessData.subSector}
                      onChange={(e) => setBusinessData({ ...businessData, subSector: e.target.value })}
                      placeholder="Enter sub sector"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="businessDescription">Business Description</Label>
                  <Textarea
                    id="businessDescription"
                    value={businessData.businessDescription}
                    onChange={(e) => setBusinessData({ ...businessData, businessDescription: e.target.value })}
                    placeholder="Describe the business"
                    rows={3}
                  />
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
                    <Label htmlFor="contactPhone">Contact Phone *</Label>
                    <Input
                      id="contactPhone"
                      value={businessData.contactPhone}
                      onChange={(e) => {
                        const formatted = formatMalawiPhone(e.target.value);
                        setBusinessData({ ...businessData, contactPhone: formatted });
                      }}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.PHONE}`}
                      className={errors.contactPhone ? 'border-destructive' : ''}
                    />
                    {errors.contactPhone && <p className="text-sm text-destructive">{errors.contactPhone}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="contactEmail">Contact Email *</Label>
                    <Input
                      id="contactEmail"
                      type="email"
                      value={businessData.contactEmail}
                      onChange={(e) => setBusinessData({ ...businessData, contactEmail: e.target.value })}
                      placeholder="Enter email"
                      className={errors.contactEmail ? 'border-destructive' : ''}
                    />
                    {errors.contactEmail && <p className="text-sm text-destructive">{errors.contactEmail}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="website">Website</Label>
                    <Input
                      id="website"
                      value={businessData.website}
                      onChange={(e) => setBusinessData({ ...businessData, website: e.target.value })}
                      placeholder="https://example.com"
                    />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Location</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* Region is inferred from district - removed from form */}

                  <div className="space-y-2">
                    <Label htmlFor="district">District</Label>
                    <Select
                      value={businessData.district}
                      onValueChange={(value) => setBusinessData({ ...businessData, district: value })}
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
                    <Label htmlFor="traditionalAuthority">Traditional Authority</Label>
                    <Input
                      id="traditionalAuthority"
                      value={businessData.traditionalAuthority}
                      onChange={(e) => setBusinessData({ ...businessData, traditionalAuthority: e.target.value })}
                      placeholder="Enter TA"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="physicalAddress">Physical Address</Label>
                  <Textarea
                    id="physicalAddress"
                    value={businessData.physicalAddress}
                    onChange={(e) => setBusinessData({ ...businessData, physicalAddress: e.target.value })}
                    placeholder="Enter physical address"
                    rows={2}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="postalAddress">Postal Address</Label>
                  <Textarea
                    id="postalAddress"
                    value={businessData.postalAddress}
                    onChange={(e) => setBusinessData({ ...businessData, postalAddress: e.target.value })}
                    placeholder="Enter postal address"
                    rows={2}
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Business Development</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Business Improvement Aspects *</Label>
                  <Popover open={improvementPopoverOpen} onOpenChange={setImprovementPopoverOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        className="w-full justify-between"
                        type="button"
                      >
                        {businessData.businessImprovementAspects.length > 0
                          ? `${businessData.businessImprovementAspects.length} selected`
                          : 'Select improvement aspects...'}
                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[400px] p-0" align="start">
                      <div className="p-4 space-y-4">
                        <div className="space-y-2">
                          <Input
                            placeholder="Search or add custom..."
                            value={improvementAspect}
                            onChange={(e) => setImprovementAspect(e.target.value)}
                            onKeyPress={(e) => {
                              if (e.key === 'Enter' && improvementAspect.trim()) {
                                e.preventDefault();
                                addImprovementAspect();
                              }
                            }}
                          />
                          {improvementAspect.trim() && !improvementOptions.includes(improvementAspect.trim()) && (
                            <Button
                              variant="outline"
                              size="sm"
                              className="w-full"
                              onClick={() => addImprovementAspect()}
                              type="button"
                            >
                              <Plus className="mr-2 h-4 w-4" />
                              Add "{improvementAspect}"
                            </Button>
                          )}
                        </div>
                        <Separator />
                        <div className="max-h-[200px] overflow-auto space-y-2">
                          {improvementOptions
                            .filter(option =>
                              option.toLowerCase().includes(improvementAspect.toLowerCase())
                            )
                            .map((option) => (
                              <div key={option} className="flex items-center space-x-2">
                                <Checkbox
                                  id={`improvement-${option}`}
                                  checked={businessData.businessImprovementAspects.includes(option)}
                                  onCheckedChange={(checked) => {
                                    if (checked) {
                                      addImprovementAspect(option);
                                    } else {
                                      removeImprovementAspect(option);
                                    }
                                  }}
                                />
                                <label
                                  htmlFor={`improvement-${option}`}
                                  className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1"
                                >
                                  {option}
                                </label>
                              </div>
                            ))}
                        </div>
                      </div>
                    </PopoverContent>
                  </Popover>
                  {errors.businessImprovementAspects && (
                    <p className="text-sm text-destructive">{errors.businessImprovementAspects}</p>
                  )}
                  <div className="flex flex-wrap gap-2 mt-2">
                    {businessData.businessImprovementAspects.map((aspect) => (
                      <Badge key={aspect} variant="secondary" className="gap-1">
                        {aspect}
                        <X
                          className="h-3 w-3 cursor-pointer"
                          onClick={() => removeImprovementAspect(aspect)}
                        />
                      </Badge>
                    ))}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Business Accessed Financing *</Label>
                  <Popover open={financingPopoverOpen} onOpenChange={setFinancingPopoverOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        className="w-full justify-between"
                        type="button"
                      >
                        {businessData.businessAccessedFinancing.length > 0
                          ? `${businessData.businessAccessedFinancing.length} selected`
                          : 'Select financing sources...'}
                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[400px] p-0" align="start">
                      <div className="p-4 space-y-4">
                        <div className="space-y-2">
                          <Input
                            placeholder="Search or add custom..."
                            value={financingSource}
                            onChange={(e) => setFinancingSource(e.target.value)}
                            onKeyPress={(e) => {
                              if (e.key === 'Enter' && financingSource.trim()) {
                                e.preventDefault();
                                addFinancingSource();
                              }
                            }}
                          />
                          {financingSource.trim() && !financingOptions.includes(financingSource.trim()) && (
                            <Button
                              variant="outline"
                              size="sm"
                              className="w-full"
                              onClick={() => addFinancingSource()}
                              type="button"
                            >
                              <Plus className="mr-2 h-4 w-4" />
                              Add "{financingSource}"
                            </Button>
                          )}
                        </div>
                        <Separator />
                        <div className="max-h-[200px] overflow-auto space-y-2">
                          {financingOptions
                            .filter(option =>
                              option.toLowerCase().includes(financingSource.toLowerCase())
                            )
                            .map((option) => (
                              <div key={option} className="flex items-center space-x-2">
                                <Checkbox
                                  id={`financing-${option}`}
                                  checked={businessData.businessAccessedFinancing.includes(option)}
                                  onCheckedChange={(checked) => {
                                    if (checked) {
                                      addFinancingSource(option);
                                    } else {
                                      removeFinancingSource(option);
                                    }
                                  }}
                                />
                                <label
                                  htmlFor={`financing-${option}`}
                                  className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1"
                                >
                                  {option}
                                </label>
                              </div>
                            ))}
                        </div>
                      </div>
                    </PopoverContent>
                  </Popover>
                  {errors.businessAccessedFinancing && (
                    <p className="text-sm text-destructive">{errors.businessAccessedFinancing}</p>
                  )}
                  <div className="flex flex-wrap gap-2 mt-2">
                    {businessData.businessAccessedFinancing.map((source) => (
                      <Badge key={source} variant="secondary" className="gap-1">
                        {source}
                        <X
                          className="h-3 w-3 cursor-pointer"
                          onClick={() => removeFinancingSource(source)}
                        />
                      </Badge>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Step 2: Primary Owner */}
        {currentStep === 2 && (
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, firstName: e.target.value })}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, lastName: e.target.value })}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, otherNames: e.target.value })}
                      placeholder="Enter other names"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="nationalIdNumber">National ID *</Label>
                    <Input
                      id="nationalIdNumber"
                      value={primaryOwner.nationalIdNumber}
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, nationalIdNumber: e.target.value.toUpperCase() })}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.NATIONAL_ID}`}
                      maxLength={8}
                      className={errors.nationalIdNumber ? 'border-destructive' : ''}
                    />
                    {errors.nationalIdNumber && <p className="text-sm text-destructive">{errors.nationalIdNumber}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="nationality">Nationality *</Label>
                    <Select
                      value={primaryOwner.nationality}
                      onValueChange={(value) => setPrimaryOwner({ ...primaryOwner, nationality: value })}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, dateOfBirth: e.target.value })}
                      className={errors.dateOfBirth ? 'border-destructive' : ''}
                      required
                    />
                    {errors.dateOfBirth && <p className="text-sm text-destructive">{errors.dateOfBirth}</p>}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="gender">Gender *</Label>
                    <Select
                      value={primaryOwner.gender}
                      onValueChange={(value) => setPrimaryOwner({ ...primaryOwner, gender: value })}
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
                      onValueChange={(value) => setPrimaryOwner({ ...primaryOwner, educationLevel: value })}
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
                      onValueChange={(value) => setPrimaryOwner({ ...primaryOwner, malawianStatus: value })}
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

                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="hasSpecialNeeds"
                      checked={primaryOwner.hasSpecialNeeds}
                      onCheckedChange={(checked) => setPrimaryOwner({ ...primaryOwner, hasSpecialNeeds: checked as boolean })}
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
                        setPrimaryOwner({ ...primaryOwner, phoneNumber: formatted });
                      }}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.PHONE}`}
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
                        setPrimaryOwner({ ...primaryOwner, landlineNumber: formatted });
                      }}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.PHONE}`}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, email: e.target.value })}
                      placeholder="Enter email"
                    />
                  </div>

                  {/* Region is inferred from district - removed from form */}

                  <div className="space-y-2">
                    <Label htmlFor="ownerDistrict">District</Label>
                    <Select
                      value={primaryOwner.district}
                      onValueChange={(value) => setPrimaryOwner({ ...primaryOwner, district: value })}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, traditionalAuthority: e.target.value })}
                      placeholder="Enter TA"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="ownerPhysicalAddress">Physical Address</Label>
                  <Textarea
                    id="ownerPhysicalAddress"
                    value={primaryOwner.physicalAddress}
                    onChange={(e) => setPrimaryOwner({ ...primaryOwner, physicalAddress: e.target.value })}
                    placeholder="Enter physical address"
                    rows={2}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="ownerPostalAddress">Postal Address</Label>
                  <Textarea
                    id="ownerPostalAddress"
                    value={primaryOwner.postalAddress}
                    onChange={(e) => setPrimaryOwner({ ...primaryOwner, postalAddress: e.target.value })}
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
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, altContactName: e.target.value })}
                      placeholder="Enter contact name"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="altContactRelationship">Relationship</Label>
                    <Input
                      id="altContactRelationship"
                      value={primaryOwner.altContactRelationship}
                      onChange={(e) => setPrimaryOwner({ ...primaryOwner, altContactRelationship: e.target.value })}
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
                        setPrimaryOwner({ ...primaryOwner, altContactPhone: formatted });
                      }}
                      placeholder={`e.g., ${EXAMPLE_FORMATS.PHONE}`}
                      className={errors.altContactPhone ? 'border-destructive' : ''}
                    />
                    {errors.altContactPhone && <p className="text-sm text-destructive">{errors.altContactPhone}</p>}
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Step 3: Formalization */}
        {currentStep === 3 && (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Business Formalization Status</CardTitle>
                <CardDescription>
                  Indicate the formal compliance and registration status of your business
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Registration & Compliance */}
                <div className="space-y-4">
                  <h3 className="font-semibold">Registration & Compliance</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="hasBankAccount">Has Bank Account</Label>
                      <Switch
                        id="hasBankAccount"
                        checked={formalizationData.hasBankAccount}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, hasBankAccount: checked })
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="hasTaxClarification">Has Tax Clarification Certificate</Label>
                      <Switch
                        id="hasTaxClarification"
                        checked={formalizationData.hasTaxClarification}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, hasTaxClarification: checked })
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="isRegisteredForVat">Registered for VAT</Label>
                      <Switch
                        id="isRegisteredForVat"
                        checked={formalizationData.isRegisteredForVat}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, isRegisteredForVat: checked })
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="hasExportLicense">Has Export License</Label>
                      <Switch
                        id="hasExportLicense"
                        checked={formalizationData.hasExportLicense}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, hasExportLicense: checked })
                        }
                      />
                    </div>
                  </div>
                </div>

                <Separator />

                {/* Business Associations */}
                <div className="space-y-4">
                  <h3 className="font-semibold">Business Associations & Support</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="isMemberOfAssociation">Member of Business Association</Label>
                      <Switch
                        id="isMemberOfAssociation"
                        checked={formalizationData.isMemberOfAssociation}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, isMemberOfAssociation: checked })
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="isAffiliated">Is Affiliated</Label>
                      <Switch
                        id="isAffiliated"
                        checked={formalizationData.isAffiliated}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, isAffiliated: checked })
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between space-x-2">
                      <Label htmlFor="hasAccessedBds">Has Accessed Business Development Services</Label>
                      <Switch
                        id="hasAccessedBds"
                        checked={formalizationData.hasAccessedBds}
                        onCheckedChange={(checked) =>
                          setFormalizationData({ ...formalizationData, hasAccessedBds: checked })
                        }
                      />
                    </div>
                  </div>
                </div>

                <Separator />

                {/* Financial Information */}
                <div className="space-y-4">
                  <h3 className="font-semibold">Financial Information</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="annualTurnover">Annual Turnover (MWK)</Label>
                      <Input
                        id="annualTurnover"
                        type="text"
                        value={formalizationData.annualTurnover}
                        onChange={(e) =>
                          setFormalizationData({ ...formalizationData, annualTurnover: e.target.value })
                        }
                        placeholder="e.g., 5000000"
                        className={errors.annualTurnover ? 'border-destructive' : ''}
                      />
                      {errors.annualTurnover && <p className="text-sm text-destructive">{errors.annualTurnover}</p>}
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="estimatedValueOfAssets">Estimated Value of Assets (MWK)</Label>
                      <Input
                        id="estimatedValueOfAssets"
                        type="text"
                        value={formalizationData.estimatedValueOfAssets}
                        onChange={(e) =>
                          setFormalizationData({ ...formalizationData, estimatedValueOfAssets: e.target.value })
                        }
                        placeholder="e.g., 10000000"
                        className={errors.estimatedValueOfAssets ? 'border-destructive' : ''}
                      />
                      {errors.estimatedValueOfAssets && <p className="text-sm text-destructive">{errors.estimatedValueOfAssets}</p>}
                    </div>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    Leave financial fields empty if you prefer not to disclose this information
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Step 4: Additional Members */}
        {currentStep === 4 && (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Additional Business Members</CardTitle>
                <CardDescription>
                  Add team members (optional). You can skip this step if there are no additional members.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                  <TabsList className="grid w-full" style={{ gridTemplateColumns: `repeat(${additionalMembers.length + 1}, 1fr)` }}>
                    {additionalMembers.map((member, index) => (
                      <TabsTrigger key={`member-${index}`} value={`member-${index}`} title={`${member.firstName} ${member.lastName}`}>
                        <User className="h-4 w-4 sm:mr-1" />
                        <span className="hidden sm:inline truncate max-w-[80px]">
                          {member.firstName}
                        </span>
                      </TabsTrigger>
                    ))}
                    <TabsTrigger value="new-member" title="Add New Member">
                      <Plus className="h-4 w-4 sm:mr-1" />
                      <span className="hidden sm:inline">Add New</span>
                    </TabsTrigger>
                  </TabsList>

                  {/* Tabs for existing members */}
                  {additionalMembers.map((member, index) => (
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
                            <Label htmlFor="memberNationality">Nationality</Label>
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
                              onChange={(e) => setNewMember({ ...newMember, nationalIdNumber: e.target.value })}
                              placeholder="Enter national ID"
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
                              onChange={(e) => setNewMember({ ...newMember, phoneNumber: e.target.value })}
                              placeholder="Enter phone"
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
              </CardContent>
            </Card>
          </div>
        )}

        {/* Step 5: Review */}
        {currentStep === 5 && (
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Review Your Information</CardTitle>
                <CardDescription>Please review all information before submitting</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Business Summary */}
                <div>
                  <h3 className="font-semibold mb-2 flex items-center gap-2">
                    <Building2 className="h-5 w-5" />
                    Business Information
                  </h3>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    {/* USME Number is auto-generated and will be assigned after creation */}
                    <div className="text-muted-foreground">Business Name:</div>
                    <div className="font-medium">{businessData.name}</div>
                    <div className="text-muted-foreground">Category:</div>
                    <div className="font-medium">{businessData.businessCategory}</div>
                    <div className="text-muted-foreground">Sector:</div>
                    <div className="font-medium">{businessData.sector}</div>
                    <div className="text-muted-foreground">Phone:</div>
                    <div className="font-medium">{businessData.contactPhone}</div>
                    <div className="text-muted-foreground">Email:</div>
                    <div className="font-medium">{businessData.contactEmail}</div>
                  </div>
                </div>

                <Separator />

                {/* Primary Owner Summary */}
                <div>
                  <h3 className="font-semibold mb-2 flex items-center gap-2">
                    <User className="h-5 w-5" />
                    Primary Business Owner
                  </h3>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div className="text-muted-foreground">Name:</div>
                    <div className="font-medium">{primaryOwner.firstName} {primaryOwner.lastName}</div>
                    <div className="text-muted-foreground">National ID:</div>
                    <div className="font-medium">{primaryOwner.nationalIdNumber}</div>
                    <div className="text-muted-foreground">Phone:</div>
                    <div className="font-medium">{primaryOwner.phoneNumber}</div>
                    {primaryOwner.email && (
                      <>
                        <div className="text-muted-foreground">Email:</div>
                        <div className="font-medium">{primaryOwner.email}</div>
                      </>
                    )}
                  </div>
                </div>

                <Separator />

                {/* Formalization Summary */}
                <div>
                  <h3 className="font-semibold mb-2 flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    Formalization Status
                  </h3>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div className="text-muted-foreground">Bank Account:</div>
                    <div className="font-medium">{formalizationData.hasBankAccount ? 'Yes' : 'No'}</div>
                    <div className="text-muted-foreground">Tax Certificate:</div>
                    <div className="font-medium">{formalizationData.hasTaxClarification ? 'Yes' : 'No'}</div>
                    <div className="text-muted-foreground">VAT Registration:</div>
                    <div className="font-medium">{formalizationData.isRegisteredForVat ? 'Yes' : 'No'}</div>
                    <div className="text-muted-foreground">Business Association:</div>
                    <div className="font-medium">{formalizationData.isMemberOfAssociation ? 'Member' : 'Non-Member'}</div>
                    {formalizationData.annualTurnover && (
                      <>
                        <div className="text-muted-foreground">Annual Turnover:</div>
                        <div className="font-medium">MWK {parseFloat(formalizationData.annualTurnover).toLocaleString()}</div>
                      </>
                    )}
                    {formalizationData.estimatedValueOfAssets && (
                      <>
                        <div className="text-muted-foreground">Asset Value:</div>
                        <div className="font-medium">MWK {parseFloat(formalizationData.estimatedValueOfAssets).toLocaleString()}</div>
                      </>
                    )}
                  </div>
                </div>

                {additionalMembers.length > 0 && (
                  <>
                    <Separator />
                    <div>
                      <h3 className="font-semibold mb-2 flex items-center gap-2">
                        <User className="h-5 w-5" />
                        Additional Members ({additionalMembers.length})
                      </h3>
                      <div className="space-y-2">
                        {additionalMembers.map((member, index) => (
                          <div key={index} className="flex items-center justify-between p-2 bg-muted rounded">
                            <span className="text-sm font-medium">
                              {member.firstName} {member.lastName}
                            </span>
                            <div className="flex gap-1">
                              {member.isIntern && <Badge variant="secondary" className="text-xs">Intern</Badge>}
                              {member.isPartTime && <Badge variant="outline" className="text-xs">Part Time</Badge>}
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>
          </div>
        )}

        {/* Navigation Buttons */}
        <div className="flex justify-between pt-4 border-t">
          <Button
            type="button"
            variant="outline"
            onClick={handlePrevious}
            disabled={currentStep === 1 || isLoading}
          >
            <ChevronLeft className="h-4 w-4 mr-2" />
            Previous
          </Button>

          {currentStep < STEPS.length ? (
            <Button
              type="button"
              onClick={handleNext}
              disabled={isLoading}
            >
              Next
              <ChevronRight className="h-4 w-4 ml-2" />
            </Button>
          ) : (
            <Button
              type="button"
              onClick={handleSubmit}
              disabled={isLoading}
            >
              {isLoading ? 'Creating...' : 'Create SME'}
            </Button>
          )}
        </div>
      </form>
    </div>
  );
});

SmeCreateForm.displayName = 'SmeCreateForm';
