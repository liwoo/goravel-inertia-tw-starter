# SME Create Form - Reference Example

This is the real-world example of the multi-step wizard pattern, from `resources/js/pages/Sme/sections/SmeCreateForm.tsx`.

## Component Signature

```typescript
interface SmeCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
  initialData?: InitialSmeData;
  onSmeCreated?: (sme: { id: number; name: string; usme_number: string }) => void;
}

export const SmeCreateForm = forwardRef<any, SmeCreateFormProps>(({
  onSuccess, onError, onCancel, isLoading, setIsSaving, initialData, onSmeCreated
}, ref) => {
  // ...
});
SmeCreateForm.displayName = 'SmeCreateForm';
```

## Steps

```typescript
const STEPS = [
  { id: 1, name: 'Business Information', description: 'Basic MSME details' },
  { id: 2, name: 'Primary Owner', description: 'Business owner information' },
  { id: 3, name: 'Formalization', description: 'Business formalization status' },
  { id: 4, name: 'Additional Members', description: 'Optional team members' },
  { id: 5, name: 'Review', description: 'Review and submit' },
];
```

## State Declarations

```typescript
const [currentStep, setCurrentStep] = useState(1);
const [errors, setErrors] = useState<Record<string, string>>({});
const [isSubmitting, setIsSubmitting] = useState(false);

// Step 1: Business Information
const [businessData, setBusinessData] = useState<SmeCreateData>({
  name: '', registrationNumber: undefined, taxIdentificationNumber: undefined,
  operationalStartDate: '', businessCategory: '', sector: '', subSector: undefined,
  businessDescription: undefined, contactPhone: '', contactEmail: '',
  physicalAddress: undefined, postalAddress: undefined, website: undefined,
  district: undefined, traditionalAuthority: undefined,
  businessImprovementAspects: [], businessAccessedFinancing: [],
});

// Step 2: Primary Owner
const [primaryOwner, setPrimaryOwner] = useState<PrimaryOwnerData>({
  firstName: '', lastName: '', otherNames: undefined, nationality: '',
  nationalIdNumber: '', dateOfBirth: undefined, gender: '', educationLevel: '',
  malawianStatus: '', hasSpecialNeeds: false, phoneNumber: '',
  landlineNumber: undefined, email: undefined, physicalAddress: undefined,
  postalAddress: undefined, district: undefined, traditionalAuthority: undefined,
  altContactName: undefined, altContactRelationship: undefined, altContactPhone: undefined,
});

// Step 3: Formalization
const [formalizationData, setFormalizationData] = useState<FormalisationFormData>(DEFAULT_FORMALISATION_DATA);

// Step 4: Additional Members
const [additionalMembers, setAdditionalMembers] = useState<AdditionalMemberData[]>([]);

// Config options fetched from API
const [categoryOptions, setCategoryOptions] = useState<string[]>([]);
const [sectorOptions, setSectorOptions] = useState<string[]>([]);
```

## Pre-fill from Initial Data

```typescript
useEffect(() => {
  if (initialData) {
    setBusinessData(prev => ({
      ...prev,
      name: initialData.name || prev.name,
      contactPhone: initialData.contactPhone || prev.contactPhone,
      // ... map all initialData fields
    }));
    setPrimaryOwner(prev => ({
      ...prev,
      firstName: initialData.ownerFirstName || prev.firstName,
      // ... map all owner fields
    }));
  }
}, [initialData]);
```

## Validation Functions

### Step 1 - Business Information

```typescript
const validateStep1 = (): boolean => {
  const newErrors: Record<string, string> = {};

  if (!businessData.name?.trim()) newErrors.name = 'Business Name is required';
  if (!businessData.businessCategory) newErrors.businessCategory = 'Category is required';
  if (!businessData.sector) newErrors.sector = 'Sector is required';

  if (!businessData.contactPhone) {
    newErrors.contactPhone = 'Phone is required';
  } else if (!validateMalawiPhone(businessData.contactPhone)) {
    newErrors.contactPhone = VALIDATION_MESSAGES.PHONE;
  }

  if (!businessData.contactEmail) {
    newErrors.contactEmail = 'Email is required';
  }

  // Optional but validated if present
  if (businessData.registrationNumber && !validateBusinessRegistration(businessData.registrationNumber)) {
    newErrors.registrationNumber = VALIDATION_MESSAGES.BUSINESS_REG;
  }
  if (businessData.taxIdentificationNumber && !validateTIN(businessData.taxIdentificationNumber)) {
    newErrors.taxIdentificationNumber = VALIDATION_MESSAGES.TIN;
  }

  // Array fields - at least one required
  if (!businessData.businessImprovementAspects?.length) {
    newErrors.businessImprovementAspects = 'At least one improvement aspect is required';
  }
  if (!businessData.businessAccessedFinancing?.length) {
    newErrors.businessAccessedFinancing = 'At least one financing source is required';
  }

  setErrors(newErrors);
  return Object.keys(newErrors).length === 0;
};
```

### Step 2 - Primary Owner

```typescript
const validateStep2 = (): boolean => {
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
  if (!primaryOwner.malawianStatus) newErrors.malawianStatus = 'Status is required';

  if (!primaryOwner.phoneNumber) {
    newErrors.phoneNumber = 'Phone is required';
  } else if (!validateMalawiPhone(primaryOwner.phoneNumber)) {
    newErrors.phoneNumber = VALIDATION_MESSAGES.PHONE;
  }

  // Optional phone fields validated if present
  if (primaryOwner.landlineNumber && !validateMalawiPhone(primaryOwner.landlineNumber)) {
    newErrors.landlineNumber = VALIDATION_MESSAGES.PHONE;
  }
  if (primaryOwner.altContactPhone && !validateMalawiPhone(primaryOwner.altContactPhone)) {
    newErrors.altContactPhone = VALIDATION_MESSAGES.PHONE;
  }

  setErrors(newErrors);
  return Object.keys(newErrors).length === 0;
};
```

## Navigation with useImperativeHandle

```typescript
useImperativeHandle(ref, () => ({
  handleSubmit: () => {
    if (currentStep === STEPS.length) {
      handleSubmit();    // Final step -> submit all
    } else {
      handleNext();      // Intermediate step -> validate and advance
    }
  }
}));

const handleNext = () => {
  let isValid = true;
  if (currentStep === 1) isValid = validateStep1();
  else if (currentStep === 2) isValid = validateStep2();
  else if (currentStep === 3) isValid = validateStep3();
  // Step 4 (additional members) has no required validation

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
```

## Final Submission

```typescript
const handleSubmit = async () => {
  if (isSubmitting) return;
  setIsSubmitting(true);

  const headers = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
  };

  try {
    // 1. Create SME
    const smeResponse = await fetch('/api/smes', {
      method: 'POST', headers,
      body: JSON.stringify(snakefiyKeys(businessData)),
    });
    const smeData = await smeResponse.json();
    const smeId = smeData.data.id;

    // 2. Create Primary Owner (linked to SME)
    await fetch('/api/primary-business-owners', {
      method: 'POST', headers,
      body: JSON.stringify({ ...snakefiyKeys(primaryOwner), sme_id: smeId }),
    });

    // 3. Upsert Formalization (backend may auto-create one)
    await fetch('/api/business-formalisations/upsert', {
      method: 'POST', headers,
      body: JSON.stringify({
        sme_id: smeId,
        ...snakefiyKeys(formalizationData),
      }),
    });

    // 4. Create Additional Members (loop)
    for (const member of additionalMembers) {
      await fetch('/api/additional_business_members', {
        method: 'POST', headers,
        body: JSON.stringify({ ...snakefiyKeys(member), sme_id: smeId }),
      });
    }

    onSuccess('MSME created successfully');
    onSmeCreated?.({ id: smeId, name: smeData.data.name, usme_number: smeData.data.usme_number });
  } catch (error) {
    onError?.(error);
  } finally {
    setIsSubmitting(false);
  }
};
```

## Progress Bar UI

```tsx
{/* Progress bar */}
<div className="mb-4">
  <Progress value={(currentStep / STEPS.length) * 100} className="h-2" />
</div>

{/* Step indicators */}
<div className="flex items-center justify-between mb-6">
  {STEPS.map((step) => (
    <div key={step.id} className="flex flex-col items-center">
      <div className={cn(
        "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium",
        step.id < currentStep && "bg-primary text-primary-foreground",
        step.id === currentStep && "bg-primary text-primary-foreground ring-2 ring-primary/30",
        step.id > currentStep && "bg-muted text-muted-foreground"
      )}>
        {step.id < currentStep ? <Check className="h-4 w-4" /> : step.id}
      </div>
      <span className="text-xs mt-1 hidden sm:block">{step.name}</span>
    </div>
  ))}
</div>

{/* Navigation buttons */}
<div className="flex justify-between mt-6">
  <Button variant="outline" onClick={handlePrevious} disabled={currentStep === 1}>
    <ChevronLeft className="h-4 w-4 mr-1" /> Previous
  </Button>
  <Button onClick={() => ref.current?.handleSubmit()}>
    {currentStep === STEPS.length ? 'Submit' : 'Next'}
    {currentStep < STEPS.length && <ChevronRight className="h-4 w-4 ml-1" />}
  </Button>
</div>
```

## Multi-Select Pattern (Popover + Badges)

Used for `businessImprovementAspects` and `businessAccessedFinancing`:

```tsx
const [improvementAspect, setImprovementAspect] = useState('');
const [improvementOptions, setImprovementOptions] = useState<string[]>([]);

const addImprovementAspect = (value?: string) => {
  const aspectToAdd = value || improvementAspect.trim();
  if (aspectToAdd && !businessData.businessImprovementAspects.includes(aspectToAdd)) {
    setBusinessData({
      ...businessData,
      businessImprovementAspects: [...businessData.businessImprovementAspects, aspectToAdd],
    });
    setImprovementAspect('');
  }
};

const removeImprovementAspect = (aspect: string) => {
  setBusinessData({
    ...businessData,
    businessImprovementAspects: businessData.businessImprovementAspects.filter(a => a !== aspect),
  });
};

// In JSX:
<Popover>
  <PopoverTrigger asChild>
    <Button variant="outline" className="w-full justify-between">
      Select aspects...
    </Button>
  </PopoverTrigger>
  <PopoverContent className="w-full p-2">
    {improvementOptions
      .filter(opt => opt.toLowerCase().includes(improvementAspect.toLowerCase()))
      .map(option => (
        <div key={option} className="flex items-center gap-2 p-2 hover:bg-accent rounded cursor-pointer"
          onClick={() => addImprovementAspect(option)}>
          <Checkbox checked={businessData.businessImprovementAspects.includes(option)} />
          <span>{option}</span>
        </div>
      ))}
  </PopoverContent>
</Popover>

{/* Show selected as badges */}
<div className="flex flex-wrap gap-1 mt-2">
  {businessData.businessImprovementAspects.map(aspect => (
    <Badge key={aspect} variant="secondary">
      {aspect}
      <X className="h-3 w-3 ml-1 cursor-pointer" onClick={() => removeImprovementAspect(aspect)} />
    </Badge>
  ))}
</div>
```

## Source Files

- Create form: `resources/js/pages/Sme/sections/SmeCreateForm.tsx`
- Types: `resources/js/types/sme.ts`
- Validators: `resources/js/lib/malawi-validators.ts`
- Form utilities: `resources/js/lib/crud-form-utils.tsx`
