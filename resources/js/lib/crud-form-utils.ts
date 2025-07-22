import { router } from '@inertiajs/react';
import { toast } from 'sonner';

export interface FormConfig<T = any> {
  endpoint: string;
  method?: 'post' | 'put' | 'patch';
  successMessage?: string;
  errorMessage?: string;
  transformData?: (data: T) => any;
  onSuccess?: (response?: any) => void;
  onError?: (errors: any) => void;
  headers?: Record<string, string>;
}

export interface FormHandlerResult {
  handleSubmit: (data: any) => Promise<void>;
  isSubmitting: boolean;
  errors: Record<string, string>;
}

/**
 * Generic form submission handler for CRUD operations
 */
export function useFormHandler<T = any>(
  config: FormConfig<T>,
  setIsSaving?: (saving: boolean) => void
): FormHandlerResult {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [errors, setErrors] = React.useState<Record<string, string>>({});

  const handleSubmit = async (data: T) => {
    setIsSubmitting(true);
    setIsSaving?.(true);
    setErrors({});

    try {
      // Transform data if transformer provided
      const transformedData = config.transformData ? config.transformData(data) : data;

      // Get CSRF token
      const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');

      // Prepare headers
      const headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        ...(csrfToken && { 'X-CSRF-TOKEN': csrfToken }),
        ...config.headers,
      };

      // Make request
      const response = await fetch(config.endpoint, {
        method: config.method || 'post',
        headers,
        body: JSON.stringify(transformedData),
      });

      const responseData = await response.json().catch(() => null);

      if (!response.ok) {
        throw responseData || { message: `Request failed with status ${response.status}` };
      }

      // Success
      toast.success(config.successMessage || 'Operation completed successfully');
      config.onSuccess?.(responseData);
    } catch (error: any) {
      console.error('Form submission error:', error);
      
      // Extract validation errors
      if (error?.errors && typeof error.errors === 'object') {
        setErrors(error.errors);
      }
      
      // Show error message
      const errorMessage = error?.message || config.errorMessage || 'Operation failed';
      toast.error(errorMessage);
      
      config.onError?.(error);
    } finally {
      setIsSubmitting(false);
      setIsSaving?.(false);
    }
  };

  return {
    handleSubmit,
    isSubmitting,
    errors,
  };
}

/**
 * Common form field wrapper with error handling
 */
export interface FormFieldProps {
  label: string;
  name: string;
  required?: boolean;
  error?: string;
  children: React.ReactNode;
  hint?: string;
}

export function FormField({ label, name, required, error, children, hint }: FormFieldProps) {
  return (
    <div className="space-y-2">
      <label htmlFor={name} className="text-sm font-medium">
        {label}
        {required && <span className="text-destructive ml-1">*</span>}
      </label>
      {children}
      {hint && !error && (
        <p className="text-sm text-muted-foreground">{hint}</p>
      )}
      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}
    </div>
  );
}

/**
 * Common form actions (Save/Cancel buttons)
 */
export interface FormActionsProps {
  onCancel: () => void;
  isSubmitting?: boolean;
  submitLabel?: string;
  cancelLabel?: string;
}

export function FormActions({ 
  onCancel, 
  isSubmitting, 
  submitLabel = 'Save',
  cancelLabel = 'Cancel' 
}: FormActionsProps) {
  return (
    <div className="flex justify-end gap-2 pt-4 border-t">
      <button
        type="button"
        onClick={onCancel}
        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
        disabled={isSubmitting}
      >
        {cancelLabel}
      </button>
      <button
        type="submit"
        className="px-4 py-2 text-sm font-medium text-white bg-primary rounded-md hover:bg-primary/90 disabled:opacity-50"
        disabled={isSubmitting}
      >
        {isSubmitting ? 'Saving...' : submitLabel}
      </button>
    </div>
  );
}

// Add React import
import React from 'react';