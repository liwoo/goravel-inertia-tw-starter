---
name: multi-step-form
description: Guide for building multi-step forms with sequential wizard creation and tab-based editing. Use when creating complex forms that span multiple sections, require per-step validation, or involve related entities (e.g., parent + children records). Based on the SME create/edit pattern.
---

# Multi-Step Form Pattern

This skill defines how to build multi-step forms in this codebase. Follow these patterns when creating forms that collect data across multiple sections or involve related entities.

## Architecture Overview

Multi-step forms use **two distinct strategies**:

| Mode | Pattern | Navigation | Save Behavior |
|------|---------|------------|---------------|
| **Create** | Sequential wizard with `Progress` bar | Next/Previous buttons, per-step validation gates | Single final submission creates all entities |
| **Edit** | Tab-based interface with `Tabs` | Click any tab freely | Each tab saves independently |

## Key Principles

1. **Create = wizard, Edit = tabs** - Creation walks users through steps in order. Editing lets users jump to any section.
2. **Per-step validation before advancing** - Never let the user proceed to the next step with invalid data.
3. **Single error state** - Use one `errors: Record<string, string>` for the entire form, cleared on step change or successful validation.
4. **forwardRef + useImperativeHandle** - Expose `handleSubmit` to the parent `CrudPage` component. On create forms, this either advances to next step or submits on the final step.
5. **Sequential API calls on create** - Create the parent entity first, then child entities using the parent's returned ID.
6. **Independent saves on edit** - Each tab has its own save handler. Only the active tab saves when the user clicks Save.
7. **Double submission guard** - Use an `isSubmitting` flag to prevent duplicate submissions.
8. **Snake/camel case conversion** - Frontend uses camelCase, backend uses snake_case. Convert with `snakefiyKeys()` on submit, manual mapping on load.

## File Structure

For an entity called `{Entity}`:

```
resources/js/pages/{Entity}/sections/
  {Entity}CreateForm.tsx          # Sequential wizard (forwardRef)
  {Entity}EditForm.tsx            # Wrapper that fetches related data
  {Entity}EditFormSimple.tsx      # Tab-based form (forwardRef)
  edit-tabs/
    {Entity}Edit{Section}Tab.tsx  # One component per tab
```

## Form State Management

Each logical section gets its own `useState`:

```typescript
// Step 1
const [businessData, setBusinessData] = useState<EntityCreateData>({...defaults});
// Step 2
const [ownerData, setOwnerData] = useState<OwnerData>({...defaults});
// Step 3
const [detailsData, setDetailsData] = useState<DetailsData>({...defaults});
// Shared
const [errors, setErrors] = useState<Record<string, string>>({});
const [currentStep, setCurrentStep] = useState(1);
const [isSubmitting, setIsSubmitting] = useState(false);
```

Update with spread: `setBusinessData({ ...businessData, fieldName: value })`.

For array fields (multi-select), use add/remove helpers:

```typescript
const addItem = (value: string) => {
  if (value && !data.items.includes(value)) {
    setData({ ...data, items: [...data.items, value] });
  }
};
const removeItem = (value: string) => {
  setData({ ...data, items: data.items.filter(v => v !== value) });
};
```

## Step Definitions

Define steps as a constant array:

```typescript
const STEPS = [
  { id: 1, name: 'Business Information', description: 'Basic details' },
  { id: 2, name: 'Primary Owner', description: 'Owner information' },
  { id: 3, name: 'Formalization', description: 'Status details' },
  { id: 4, name: 'Additional Members', description: 'Optional team members' },
  { id: 5, name: 'Review', description: 'Review and submit' },
];
```

## Validation Pattern

Validate per-step before allowing navigation:

```typescript
const handleNext = () => {
  let isValid = true;

  if (currentStep === 1) isValid = validateStep1();
  else if (currentStep === 2) isValid = validateStep2();
  // ...

  if (isValid && currentStep < STEPS.length) {
    setCurrentStep(currentStep + 1);
    setErrors({});
  }
};
```

Each validation function builds an errors object and returns a boolean:

```typescript
const validateStep1 = (): boolean => {
  const newErrors: Record<string, string> = {};

  // Required fields
  if (!data.name?.trim()) newErrors.name = 'Name is required';

  // Format validation (only if value present)
  if (data.phone && !validateMalawiPhone(data.phone)) {
    newErrors.phone = VALIDATION_MESSAGES.PHONE;
  }

  // Array fields
  if (!data.items || data.items.length === 0) {
    newErrors.items = 'At least one item is required';
  }

  setErrors(newErrors);
  return Object.keys(newErrors).length === 0;
};
```

Use validators from `@/lib/malawi-validators` for phone, national ID, TIN, and business registration fields.

## useImperativeHandle Pattern

Expose form submission to parent CrudPage:

```typescript
// Create form: advances step or submits on final step
useImperativeHandle(ref, () => ({
  handleSubmit: () => {
    if (currentStep === STEPS.length) {
      handleSubmit();
    } else {
      handleNext();
    }
  }
}));

// Edit form: saves the active tab
useImperativeHandle(ref, () => ({
  handleSubmit: () => {
    switch (activeTab) {
      case 'section-a': return handleSaveSectionA();
      case 'section-b': return handleSaveSectionB();
    }
  }
}));
```

## Create Form Submission

Submit all entities sequentially, using the parent ID for children:

```typescript
const handleSubmit = async () => {
  if (isSubmitting) return; // Double-submit guard
  setIsSubmitting(true);

  try {
    // 1. Create parent entity
    const parentRes = await fetch('/api/entities', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Accept': 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
      body: JSON.stringify(snakefiyKeys(parentData)),
    });
    const parentJson = await parentRes.json();
    const parentId = parentJson.data.id;

    // 2. Create child entities using parent ID
    await fetch('/api/child-entities', {
      method: 'POST',
      body: JSON.stringify({ ...snakefiyKeys(childData), parent_id: parentId }),
      headers: { /* same headers */ },
    });

    // 3. Upsert for entities that may be auto-created by backend
    await fetch('/api/auto-entities/upsert', {
      method: 'POST',
      body: JSON.stringify({ ...snakefiyKeys(autoData), parent_id: parentId }),
      headers: { /* same headers */ },
    });

    // 4. Create array children (loop)
    for (const member of members) {
      await fetch('/api/members', {
        method: 'POST',
        body: JSON.stringify({ ...snakefiyKeys(member), parent_id: parentId }),
        headers: { /* same headers */ },
      });
    }

    onSuccess('Entity created successfully');
    onEntityCreated?.({ id: parentId, ...parentJson.data });
  } catch (error) {
    onError?.(error);
  } finally {
    setIsSubmitting(false);
  }
};
```

## Edit Form: Data Loading Wrapper

The edit form uses a wrapper component to fetch related data in parallel:

```typescript
export const EntityEditForm = forwardRef<any, Props>(({ item, ...props }, ref) => {
  const [relatedA, setRelatedA] = useState<TypeA | undefined>();
  const [relatedB, setRelatedB] = useState<TypeB | undefined>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRelated = async () => {
      setLoading(true);
      const [aRes, bRes] = await Promise.all([
        axios.get(`/api/entities/${item.id}/related_a`).catch(() => ({ data: { data: null } })),
        axios.get(`/api/entities/${item.id}/related_b`).catch(() => ({ data: { data: [] } })),
      ]);
      setRelatedA(convertSnakeToCamel(aRes.data.data));
      setRelatedB(convertSnakeToCamel(bRes.data.data));
      setLoading(false);
    };
    if (item?.id) fetchRelated();
  }, [item?.id]);

  if (loading) return <Loader2 className="h-8 w-8 animate-spin" />;

  return <EntityEditFormSimple ref={ref} item={item} relatedA={relatedA} relatedB={relatedB} {...props} />;
});
```

## Edit Form: Tab-Based Save

Each tab saves independently with its own validation and API call:

```typescript
const handleSaveSection = async () => {
  const newErrors: Record<string, string> = {};
  // validate...
  if (Object.keys(newErrors).length > 0) { setErrors(newErrors); return; }

  setErrors({});
  setIsSaving?.(true);
  try {
    const response = await fetch(`/api/sections/${section.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', 'Accept': 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
      body: JSON.stringify(snakefiyKeys(sectionData)),
    });
    if (response.ok) onSuccess('Section updated');
    else onError?.(await response.json().catch(() => ({})));
  } catch (error) { onError?.(error); }
  finally { setIsSaving?.(false); }
};
```

## Edit Form: Additional Members (Array Children)

Track original IDs to detect deletions:

```typescript
const [members, setMembers] = useState<MemberData[]>(initialMembers);
const [originalMemberIds] = useState<number[]>(initialMembers.filter(m => m.id).map(m => m.id!));

const handleSaveMembers = async () => {
  const currentIds = members.filter(m => m.id).map(m => m.id!);
  const deletedIds = originalMemberIds.filter(id => !currentIds.includes(id));

  // Delete removed
  for (const id of deletedIds) {
    await fetch(`/api/members/${id}`, { method: 'DELETE', headers: {...} });
  }
  // Create or update
  for (const member of members) {
    if (member.id) {
      await fetch(`/api/members/${member.id}`, { method: 'PUT', body: JSON.stringify(snakefiyKeys({...member, parentId})), headers: {...} });
    } else {
      const res = await fetch('/api/members', { method: 'POST', body: JSON.stringify(snakefiyKeys({...member, parentId})), headers: {...} });
      if (res.ok) { const data = await res.json(); member.id = data.data?.id; }
    }
  }
  setMembers([...members]);
  onSuccess('Members saved');
};
```

## UI Components Used

| Purpose | Component | Import |
|---------|-----------|--------|
| Text input | `Input` | `@/components/ui/input` |
| Select dropdown | `Select, SelectContent, SelectItem, SelectTrigger, SelectValue` | `@/components/ui/select` |
| Multi-line text | `Textarea` | `@/components/ui/textarea` |
| Boolean toggle | `Switch` | `@/components/ui/switch` |
| Checkbox | `Checkbox` | `@/components/ui/checkbox` |
| Multi-select | `Popover` + `Checkbox` + `Badge` | `@/components/ui/popover`, `@/components/ui/badge` |
| Progress bar | `Progress` | `@/components/ui/progress` |
| Section cards | `Card, CardContent, CardHeader, CardTitle` | `@/components/ui/card` |
| Tab navigation | `Tabs, TabsContent, TabsList, TabsTrigger` | `@/components/ui/tabs` |
| Labels | `Label` | `@/components/ui/label` |

## Data Conversion

Frontend camelCase to backend snake_case: use `snakefiyKeys()` from `@/lib/utils`.

Backend snake_case to frontend camelCase: manually map in `useEffect` when loading:

```typescript
setPrimaryOwner({
  firstName: data.first_name || '',
  dateOfBirth: data.date_of_birth ? data.date_of_birth.slice(0, 10) : undefined,
  hasSpecialNeeds: data.has_special_needs || false,
});
```

Date fields from backend include timestamps (`2024-01-15T00:00:00Z`) - slice to `YYYY-MM-DD` for date inputs.

## Config-Driven Options

Fetch dynamic select options from the configs API:

```typescript
useEffect(() => {
  const fetchConfigs = async () => {
    const [catRes, secRes] = await Promise.all([
      axios.get('/api/configs?config_type=business_category'),
      axios.get('/api/configs?config_type=sector'),
    ]);
    setCategoryOptions(catRes.data.data.map((c: any) => c.value));
    setSectorOptions(secRes.data.data.map((s: any) => s.value));
  };
  fetchConfigs();
}, []);
```

## Checklist

When building a new multi-step form, ensure:

- [ ] Create form uses `forwardRef` + `useImperativeHandle`
- [ ] Edit form uses wrapper pattern to fetch related data
- [ ] Per-step validation with clear error messages
- [ ] `isSubmitting` guard against double submission
- [ ] `snakefiyKeys()` on all outbound payloads
- [ ] Manual snake-to-camel mapping on inbound data
- [ ] Date fields sliced to `YYYY-MM-DD`
- [ ] Malawi validators applied to phone, national ID, TIN, registration fields
- [ ] Array children tracked with original IDs for deletion detection
- [ ] Each edit tab saves independently
- [ ] `setIsSaving` callback to notify parent of save state
- [ ] Error display below each invalid field
- [ ] Required field indicators (`*`)
- [ ] Progress bar and step indicators on create wizard

For detailed code examples, see [examples/sme-create-form.md](examples/sme-create-form.md) and [examples/sme-edit-form.md](examples/sme-edit-form.md).
