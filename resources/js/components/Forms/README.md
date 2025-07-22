# CRUD Form Components

This directory contains reusable form components and utilities for CRUD operations.

## CrudForm Component

A generic form wrapper that handles:
- Form submission with loading states
- Error handling and display
- CSRF token management
- Success/error callbacks

### Usage Example

```tsx
import { CrudForm, CrudFormHandle } from '@/components/Forms/CrudForm';

const UserCreateForm = forwardRef<CrudFormHandle, CrudFormProps>((props, ref) => {
  const [roles, setRoles] = useState([]);
  
  return (
    <CrudForm
      ref={ref}
      {...props}
      config={{
        endpoint: '/api/users',
        method: 'post',
        successMessage: 'User created successfully',
        transformData: (data) => ({
          ...data,
          role_id: parseInt(data.role_id),
        }),
      }}
    >
      {({ register, errors, FormField }) => (
        <>
          <FormField label="Name" name="name" required error={errors.name}>
            <Input {...register('name')} />
          </FormField>
          
          <FormField label="Email" name="email" required error={errors.email}>
            <Input type="email" {...register('email')} />
          </FormField>
          
          <FormField label="Role" name="role_id" error={errors.role_id}>
            <Select {...register('role_id')}>
              <option value="">Select a role</option>
              {roles.map(role => (
                <option key={role.id} value={role.id}>{role.name}</option>
              ))}
            </Select>
          </FormField>
        </>
      )}
    </CrudForm>
  );
});
```

## Form Utilities

### useFormHandler

A custom hook that manages form submission:

```tsx
const { handleSubmit, isSubmitting, errors } = useFormHandler({
  endpoint: '/api/resource',
  method: 'post',
  transformData: (data) => ({ ...data, processed: true }),
  onSuccess: (response) => console.log('Success!', response),
  onError: (errors) => console.error('Failed!', errors),
});
```

### FormField Component

A field wrapper that handles labels, errors, and hints:

```tsx
<FormField 
  label="Username" 
  name="username" 
  required 
  error={errors.username}
  hint="Must be unique"
>
  <Input {...register('username')} />
</FormField>
```

### FormActions Component

Standard form buttons (Save/Cancel):

```tsx
<FormActions
  onCancel={handleCancel}
  isSubmitting={isLoading}
  submitLabel="Create User"
  cancelLabel="Go Back"
/>
```

## Benefits

1. **Consistency**: All forms follow the same patterns
2. **Less Boilerplate**: Common logic is abstracted
3. **Type Safety**: Full TypeScript support
4. **Error Handling**: Automatic error extraction and display
5. **Loading States**: Built-in loading indicators
6. **CSRF Protection**: Automatic token handling