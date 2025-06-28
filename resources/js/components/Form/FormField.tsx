import React from 'react';
import { cn } from '@/lib/utils';
import { FormFieldProps } from '@/types/crud';
import { Label } from '@/components/ui/label';

export function FormField({ 
  children, 
  label, 
  error, 
  required, 
  className,
  description,
  hint 
}: FormFieldProps) {
  return (
    <div className={cn('space-y-2', className)}>
      {label && (
        <Label className="text-sm font-medium text-gray-900">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </Label>
      )}
      
      {description && (
        <p className="text-sm text-gray-600">{description}</p>
      )}
      
      <div className="relative">
        {children}
      </div>
      
      {hint && !error && (
        <p className="text-xs text-gray-500">{hint}</p>
      )}
      
      {error && (
        <p className="text-sm text-red-600 flex items-center">
          <span className="mr-1">⚠</span>
          {error}
        </p>
      )}
    </div>
  );
}