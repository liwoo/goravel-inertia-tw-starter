import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { ConfigCreateData, ConfigType, CONFIG_TYPES } from '@/types/config';
import { User, FileText } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface ConfigCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
}

export const ConfigCreateForm = forwardRef<any, ConfigCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<ConfigCreateData>({
    name: '',
    code: '',
    configType: ConfigType.Financing,
    description: '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [codeManuallyEdited, setCodeManuallyEdited] = useState(false);

  // Auto-generate code from name
  const generateCode = (name: string): string => {
    const trimmed = name.trim();
    if (!trimmed) return '';

    const words = trimmed.split(/\s+/);

    if (words.length > 1) {
      // Multiple words: take first letter of each word
      return words.map(word => word[0].toUpperCase()).join('');
    }

    // Single word: take first 2 letters
    return trimmed.substring(0, 2).toUpperCase();
  };

  const handleNameChange = (name: string) => {
    setFormData({ ...formData, name });

    // Auto-generate code if it hasn't been manually edited
    if (!codeManuallyEdited) {
      setFormData(prev => ({ ...prev, name, code: generateCode(name) }));
    }
  };

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    if (!formData.name) {
      newErrors.name = 'Name is required';
    }
    if (!formData.configType) {
      newErrors.configType = 'Config Type is required';
    }
    

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      // Convert camelCase to snake_case for backend
      const requestData = {
        name: formData.name,
        code: formData.code,
        config_type: formData.configType,
        description: formData.description,
      };

      const response = await fetch('/api/configs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(requestData),
      });

      if (response.ok) {
        onSuccess('Config created successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Config Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <User className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => handleNameChange(e.target.value)}
                  placeholder="Enter name"
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="code">Code</Label>
                <Input
                  id="code"
                  value={formData.code}
                  onChange={(e) => {
                    setFormData({ ...formData, code: e.target.value });
                    setCodeManuallyEdited(true);
                  }}
                  placeholder="Auto-generated from name"
                  className={errors.code ? 'border-destructive' : ''}
                />
                <p className="text-xs text-muted-foreground">
                  Auto-generates from name. Edit to customize.
                </p>
                {errors.code && (
                  <p className="text-sm text-destructive">{errors.code}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="configType">Config Type *</Label>
                <Select
                  value={formData.configType}
                  onValueChange={(value) => setFormData({ ...formData, configType: value as ConfigType })}
                >
                  <SelectTrigger className={errors.configType ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Select config type" />
                  </SelectTrigger>
                  <SelectContent>
                    {CONFIG_TYPES.map((type) => (
                      <SelectItem key={type} value={type}>
                        {type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.configType && (
                  <p className="text-sm text-destructive">{errors.configType}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Enter description"
                  rows={4}
                  className={errors.description ? 'border-destructive' : ''}
                />
                {errors.description && (
                  <p className="text-sm text-destructive">{errors.description}</p>
                )}
              </div>
            </div>

          </div>
        </div>
      </div>
    </form>
  );
});

ConfigCreateForm.displayName = 'ConfigCreateForm';
