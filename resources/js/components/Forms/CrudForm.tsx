import React, { forwardRef, useImperativeHandle } from 'react';
import { useFormHandler, FormConfig, FormActions, FormField } from '@/lib/crud-form-utils';

export interface CrudFormProps<T = any> {
  onSuccess: (message?: string) => void;
  onError?: (errors: any) => void;
  onCancel: () => void;
  setIsSaving?: (saving: boolean) => void;
  config: FormConfig<T>;
  children: (props: {
    register: (name: string) => any;
    errors: Record<string, string>;
    isSubmitting: boolean;
    FormField: typeof FormField;
  }) => React.ReactNode;
  initialData?: T;
}

export interface CrudFormHandle {
  handleSubmit: () => void;
}

/**
 * Generic CRUD form component that handles common form logic
 */
export const CrudForm = forwardRef<CrudFormHandle, CrudFormProps>(
  ({ onSuccess, onError, onCancel, setIsSaving, config, children, initialData }, ref) => {
    const formRef = React.useRef<HTMLFormElement>(null);
    const [formData, setFormData] = React.useState<any>(initialData || {});
    
    // Enhanced config with callbacks
    const enhancedConfig: FormConfig = {
      ...config,
      onSuccess: (response) => {
        config.onSuccess?.(response);
        onSuccess(config.successMessage);
      },
      onError: (errors) => {
        config.onError?.(errors);
        onError?.(errors);
      },
    };
    
    const { handleSubmit, isSubmitting, errors } = useFormHandler(enhancedConfig, setIsSaving);
    
    // Expose submit method to parent
    useImperativeHandle(ref, () => ({
      handleSubmit: () => {
        if (formRef.current) {
          formRef.current.dispatchEvent(
            new Event('submit', { cancelable: true, bubbles: true })
          );
        }
      },
    }));
    
    const onSubmit = async (e: React.FormEvent) => {
      e.preventDefault();
      await handleSubmit(formData);
    };
    
    const register = (name: string) => ({
      name,
      value: formData[name] || '',
      onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        setFormData({ ...formData, [name]: e.target.value });
      },
    });
    
    return (
      <form ref={formRef} onSubmit={onSubmit} className="space-y-4">
        {children({ register, errors, isSubmitting, FormField })}
        <FormActions
          onCancel={onCancel}
          isSubmitting={isSubmitting}
        />
      </form>
    );
  }
);

CrudForm.displayName = 'CrudForm';

/**
 * Example usage for a Create Form:
 * 
 * const UserCreateForm = forwardRef<CrudFormHandle, CrudFormProps>((props, ref) => {
 *   return (
 *     <CrudForm
 *       ref={ref}
 *       {...props}
 *       config={{
 *         endpoint: '/api/users',
 *         method: 'post',
 *         successMessage: 'User created successfully',
 *       }}
 *     >
 *       {({ register, errors, FormField }) => (
 *         <>
 *           <FormField label="Name" name="name" required error={errors.name}>
 *             <input type="text" {...register('name')} className="form-input" />
 *           </FormField>
 *           <FormField label="Email" name="email" required error={errors.email}>
 *             <input type="email" {...register('email')} className="form-input" />
 *           </FormField>
 *         </>
 *       )}
 *     </CrudForm>
 *   );
 * });
 */