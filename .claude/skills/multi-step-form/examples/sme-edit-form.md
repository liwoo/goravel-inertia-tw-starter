# SME Edit Form - Reference Example

This is the real-world example of the tab-based edit pattern, from `resources/js/pages/Sme/sections/SmeEditForm.tsx` and `SmeEditFormSimple.tsx`.

## Architecture: Wrapper + Simple Pattern

The edit form is split into two components:

1. **`SmeEditForm`** (wrapper) - Fetches all related data in parallel, shows a loading spinner, then renders the simple form.
2. **`SmeEditFormSimple`** (form) - Receives pre-loaded data as props, renders tab-based UI, handles saves.

This separation keeps data fetching concerns out of the form rendering logic.

## Wrapper Component

```typescript
export const SmeEditForm = forwardRef<any, SmeEditFormProps>(({
  item: sme, onSuccess, onError, onCancel, isLoading, setIsSaving
}, ref) => {
  const [primaryOwner, setPrimaryOwner] = useState<PrimaryOwnerData | undefined>();
  const [formalization, setFormalization] = useState<BusinessFormalisation | undefined>();
  const [additionalMembers, setAdditionalMembers] = useState<AdditionalMemberData[]>([]);
  const [loadingRelatedData, setLoadingRelatedData] = useState(true);

  useEffect(() => {
    const fetchRelatedData = async () => {
      setLoadingRelatedData(true);
      try {
        const [ownerRes, formalizationRes, membersRes] = await Promise.all([
          axios.get(`/api/smes/${sme.id}/primary_business_owner`)
            .catch(() => ({ data: { data: null } })),
          axios.get(`/api/smes/${sme.id}/business_formalisation`)
            .catch(() => ({ data: { data: null } })),
          axios.get(`/api/smes/${sme.id}/additional_business_members`)
            .catch(() => ({ data: { data: [] } })),
        ]);

        // Convert snake_case to camelCase
        const ownerData = ownerRes.data.data;
        if (ownerData) {
          setPrimaryOwner({
            id: ownerData.id,
            firstName: ownerData.first_name || '',
            lastName: ownerData.last_name || '',
            dateOfBirth: ownerData.date_of_birth ? ownerData.date_of_birth.slice(0, 10) : undefined,
            hasSpecialNeeds: ownerData.has_special_needs || false,
            // ... all other fields
          });
        }

        // Members array conversion
        setAdditionalMembers(membersRes.data.data.map((m: any) => ({
          id: m.id,
          firstName: m.first_name || '',
          lastName: m.last_name || '',
          dateOfBirth: m.date_of_birth ? m.date_of_birth.slice(0, 10) : undefined,
          isIntern: m.is_intern || false,
          isPartTime: m.is_part_time || false,
          // ...
        })));
      } catch (error) {
        console.error('Error fetching related data:', error);
      } finally {
        setLoadingRelatedData(false);
      }
    };

    if (sme?.id) fetchRelatedData();
  }, [sme?.id]);

  if (loadingRelatedData) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        <span className="ml-2 text-muted-foreground">Loading data...</span>
      </div>
    );
  }

  return (
    <SmeEditFormSimple
      ref={ref} item={sme}
      primaryOwner={primaryOwner} formalization={formalization}
      additionalMembers={additionalMembers}
      onSuccess={onSuccess} onError={onError} onCancel={onCancel}
      isLoading={isLoading} setIsSaving={setIsSaving}
    />
  );
});
```

Key patterns:
- **Parallel fetch** with `Promise.all` for all related entities
- **Graceful failure** with `.catch(() => ({ data: { data: null } }))` per request
- **Manual snake-to-camel** conversion (no automatic utility for this direction)
- **Date slicing**: `date_of_birth.slice(0, 10)` to strip timestamps

## Tab-Based Form

```typescript
export const SmeEditFormSimple = forwardRef<any, Props>(({
  item: sme, primaryOwner: initialPrimaryOwner, formalization: initialFormalization,
  additionalMembers: initialAdditionalMembers = [], onSuccess, onError, setIsSaving
}, ref) => {
  const [activeTab, setActiveTab] = useState('business-info');
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Each section gets its own state, initialized from props
  const [businessData, setBusinessData] = useState<SmeUpdateData>({
    name: sme.name,
    registrationNumber: sme.registration_number || sme.registrationNumber || undefined,
    // Handle both snake_case and camelCase from backend
    businessCategory: sme.business_category || sme.businessCategory || '',
    businessImprovementAspects: Array.isArray(sme.business_improvement_aspects || sme.businessImprovementAspects)
      ? (sme.business_improvement_aspects || sme.businessImprovementAspects)
      : String(sme.business_improvement_aspects || '').split(',').filter(Boolean),
    // ...
  });

  const [primaryOwner, setPrimaryOwner] = useState(initialPrimaryOwner);
  const [formalization, setFormalization] = useState(initialFormalization);
  const [additionalMembers, setAdditionalMembers] = useState(initialAdditionalMembers);
  const [originalMemberIds] = useState(initialAdditionalMembers.filter(m => m.id).map(m => m.id!));
```

Note the defensive handling of `businessImprovementAspects` which may arrive as an array or comma-separated string.

## Tab-Based useImperativeHandle

```typescript
useImperativeHandle(ref, () => ({
  handleSubmit: () => {
    switch (activeTab) {
      case 'business-info': return handleSaveBusinessInfo();
      case 'primary-owner': return handleSavePrimaryOwner();
      case 'formalization': return handleSaveFormalization();
      case 'additional-members': return handleSaveAdditionalMembers();
      default: return handleSaveBusinessInfo();
    }
  }
}));
```

## Tab UI Layout

```tsx
<Tabs value={activeTab} onValueChange={setActiveTab}>
  <TabsList className="w-full">
    <TabsTrigger value="business-info" className="flex-1" title="Business Info">
      <Building2 className="h-4 w-4 sm:mr-2" />
      <span className="hidden sm:inline">Business Info</span>
    </TabsTrigger>
    <TabsTrigger value="primary-owner" className="flex-1" title="Primary Owner">
      <User className="h-4 w-4 sm:mr-2" />
      <span className="hidden sm:inline">Primary Owner</span>
    </TabsTrigger>
    <TabsTrigger value="formalization" className="flex-1" title="Formalization">
      <FileText className="h-4 w-4 sm:mr-2" />
      <span className="hidden sm:inline">Formalization</span>
    </TabsTrigger>
    <TabsTrigger value="additional-members" className="flex-1" title="Team Members">
      <Users className="h-4 w-4 sm:mr-2" />
      <span className="hidden sm:inline">Team Members</span>
    </TabsTrigger>
  </TabsList>

  <TabsContent value="business-info">
    <SmeEditBusinessInfoTab businessData={businessData} onChange={setBusinessData} errors={errors} />
  </TabsContent>
  {/* ... other tabs */}
</Tabs>
```

Responsive pattern: Icon always visible, label hidden on small screens with `hidden sm:inline`.

## Formalization Upsert Pattern

Backend may auto-create a formalization record when an SME is created. The edit form handles both scenarios:

```typescript
const handleSaveFormalization = async () => {
  const isUpdate = formalization.id !== undefined;
  const url = isUpdate
    ? `/api/business-formalisations/${formalization.id}`
    : `/api/business-formalisations/upsert`;
  const method = isUpdate ? 'PUT' : 'POST';

  const payload = isUpdate
    ? snakefiyKeys(formalization)
    : snakefiyKeys({ ...formalization, sme_id: sme.id });

  const response = await fetch(url, { method, headers: {...}, body: JSON.stringify(payload) });

  if (response.ok) {
    const data = await response.json();
    if (data.data?.id) {
      setFormalization({ ...formalization, id: data.data.id }); // Capture ID from upsert
    }
  }
};
```

## Additional Members CRUD

Full create/update/delete cycle:

```typescript
const handleSaveAdditionalMembers = async () => {
  const currentIds = additionalMembers.filter(m => m.id).map(m => m.id!);
  const deletedIds = originalMemberIds.filter(id => !currentIds.includes(id));

  // 1. Delete removed members
  for (const id of deletedIds) {
    await fetch(`/api/additional_business_members/${id}`, { method: 'DELETE', headers });
  }

  // 2. Create new / update existing
  for (const member of additionalMembers) {
    const payload = snakefiyKeys({ ...member, smeId: sme.id });

    if (member.id) {
      await fetch(`/api/additional_business_members/${member.id}`, {
        method: 'PUT', headers, body: JSON.stringify(payload),
      });
    } else {
      const res = await fetch('/api/additional_business_members', {
        method: 'POST', headers, body: JSON.stringify(payload),
      });
      if (res.ok) {
        const data = await res.json();
        if (data.data?.id) member.id = data.data.id; // Track new ID
      }
    }
  }

  setAdditionalMembers([...additionalMembers]); // Trigger re-render with new IDs
};
```

## Edit Tab Component Pattern

Each tab is a simple component receiving data, onChange, and errors:

```typescript
interface SmeEditBusinessInfoTabProps {
  businessData: SmeUpdateData;
  onChange: (data: SmeUpdateData) => void;
  errors: Record<string, string>;
}

export function SmeEditBusinessInfoTab({ businessData, onChange, errors }: SmeEditBusinessInfoTabProps) {
  return (
    <Card>
      <CardHeader><CardTitle>Business Information</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Business Name *</Label>
            <Input
              value={businessData.name}
              onChange={(e) => onChange({ ...businessData, name: e.target.value })}
            />
            {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
          </div>
          {/* ... more fields */}
        </div>
      </CardContent>
    </Card>
  );
}
```

## Source Files

- Edit wrapper: `resources/js/pages/Sme/sections/SmeEditForm.tsx`
- Edit form: `resources/js/pages/Sme/sections/SmeEditFormSimple.tsx`
- Tab components: `resources/js/pages/Sme/sections/edit-tabs/SmeEdit*Tab.tsx`
- Types: `resources/js/types/sme.ts`, `resources/js/types/business_formalisation.ts`
