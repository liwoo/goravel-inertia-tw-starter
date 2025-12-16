import React, { useState, useEffect, useMemo } from 'react';
import axios from 'axios';
import { toast } from 'sonner';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Separator } from '@/components/ui/separator';
import { ScrollArea } from '@/components/ui/scroll-area';
import { snakefiyKeys } from '@/lib/utils';
import { Loader2, Info, AlertCircle } from 'lucide-react';
import {
  validateBusinessRegistration,
  validateTIN,
  formatBusinessRegistration,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS,
} from '@/lib/malawi-validators';

interface FormalisationEditModalProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  smeId: number;
  currentData: {
    registrationNumber?: string;
    taxIdentificationNumber?: string;
    hasBankAccount?: boolean;
    hasTaxClarification?: boolean;
    isRegisteredForVat?: boolean;
    isMemberOfAssociation?: boolean;
    isAffiliated?: boolean;
    hasExportLicense?: boolean;
    hasAccessedBds?: boolean;
    annualTurnover?: number;
    estimatedValueOfAssets?: number;
  };
  hasPendingAmendment?: boolean;
  onSuccess?: () => void;
}

interface FormData {
  registrationNumber?: string;
  taxIdentificationNumber?: string;
  hasBankAccount?: boolean;
  hasTaxClarification?: boolean;
  isRegisteredForVat?: boolean;
  isMemberOfAssociation?: boolean;
  isAffiliated?: boolean;
  hasExportLicense?: boolean;
  hasAccessedBds?: boolean;
  annualTurnover?: string;
  estimatedValueOfAssets?: string;
  changeReason?: string;
}

export const FormalisationEditSheet: React.FC<FormalisationEditModalProps> = ({
  isOpen,
  onOpenChange,
  smeId,
  currentData,
  hasPendingAmendment = false,
  onSuccess,
}) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState<FormData>({});
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Initialize form data when modal opens
  useEffect(() => {
    if (isOpen) {
      setFormData({
        registrationNumber: currentData.registrationNumber || '',
        taxIdentificationNumber: currentData.taxIdentificationNumber || '',
        hasBankAccount: currentData.hasBankAccount ?? false,
        hasTaxClarification: currentData.hasTaxClarification ?? false,
        isRegisteredForVat: currentData.isRegisteredForVat ?? false,
        isMemberOfAssociation: currentData.isMemberOfAssociation ?? false,
        isAffiliated: currentData.isAffiliated ?? false,
        hasExportLicense: currentData.hasExportLicense ?? false,
        hasAccessedBds: currentData.hasAccessedBds ?? false,
        annualTurnover: currentData.annualTurnover?.toString() || '',
        estimatedValueOfAssets: currentData.estimatedValueOfAssets?.toString() || '',
        changeReason: '',
      });
    }
  }, [isOpen, currentData]);

  // Check if any changes have been made
  const hasChanges = useMemo(() => {
    return (
      formData.registrationNumber !== (currentData.registrationNumber || '') ||
      formData.taxIdentificationNumber !== (currentData.taxIdentificationNumber || '') ||
      formData.hasBankAccount !== (currentData.hasBankAccount ?? false) ||
      formData.hasTaxClarification !== (currentData.hasTaxClarification ?? false) ||
      formData.isRegisteredForVat !== (currentData.isRegisteredForVat ?? false) ||
      formData.isMemberOfAssociation !== (currentData.isMemberOfAssociation ?? false) ||
      formData.isAffiliated !== (currentData.isAffiliated ?? false) ||
      formData.hasExportLicense !== (currentData.hasExportLicense ?? false) ||
      formData.hasAccessedBds !== (currentData.hasAccessedBds ?? false) ||
      formData.annualTurnover !== (currentData.annualTurnover?.toString() || '') ||
      formData.estimatedValueOfAssets !== (currentData.estimatedValueOfAssets?.toString() || '')
    );
  }, [formData, currentData]);

  const handleSubmit = async () => {
    if (!hasChanges) {
      toast.error('No changes detected');
      return;
    }

    // Validate form fields
    const newErrors: Record<string, string> = {};

    // Validate Registration Number if provided
    if (formData.registrationNumber && formData.registrationNumber.trim() !== '') {
      if (!validateBusinessRegistration(formData.registrationNumber)) {
        newErrors.registrationNumber = VALIDATION_MESSAGES.BUSINESS_REG;
      }
    }

    // Validate TIN if provided
    if (formData.taxIdentificationNumber && formData.taxIdentificationNumber.trim() !== '') {
      if (!validateTIN(formData.taxIdentificationNumber)) {
        newErrors.taxIdentificationNumber = VALIDATION_MESSAGES.TIN;
      }
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      toast.error('Please fix the validation errors before submitting');
      return;
    }

    setErrors({});
    setIsSubmitting(true);
    try {
      const payload = snakefiyKeys({
        smeId,
        registrationNumber: formData.registrationNumber,
        taxIdentificationNumber: formData.taxIdentificationNumber,
        hasBankAccount: formData.hasBankAccount,
        hasTaxClarification: formData.hasTaxClarification,
        isRegisteredForVat: formData.isRegisteredForVat,
        isMemberOfAssociation: formData.isMemberOfAssociation,
        isAffiliated: formData.isAffiliated,
        hasExportLicense: formData.hasExportLicense,
        hasAccessedBds: formData.hasAccessedBds,
        annualTurnover: formData.annualTurnover ? parseFloat(formData.annualTurnover) : null,
        estimatedValueOfAssets: formData.estimatedValueOfAssets
          ? parseFloat(formData.estimatedValueOfAssets)
          : null,
        changeReason: formData.changeReason || null,
      });

      await axios.post('/api/applications/amendment', payload);

      toast.success('Amendment request submitted successfully');
      onOpenChange(false);

      if (onSuccess) {
        onSuccess();
      }
    } catch (error: any) {
      console.error('Error submitting amendment:', error);
      toast.error(error.response?.data?.message || 'Failed to submit amendment request');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  const formatCurrency = (value: string): string => {
    const cleaned = value.replace(/[^0-9.]/g, '');
    const parts = cleaned.split('.');
    if (parts.length > 2) {
      return parts[0] + '.' + parts.slice(1).join('');
    }
    return cleaned;
  };

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-xl md:max-w-2xl overflow-hidden flex flex-col">
        <SheetHeader className="pb-4">
          <SheetTitle className="text-2xl">Request Formalisation Changes</SheetTitle>
          <SheetDescription>
            Propose changes to your business formalisation details. All changes require admin approval.
          </SheetDescription>
        </SheetHeader>

        <ScrollArea className="flex-1 pr-4">
          {hasPendingAmendment ? (
            <Alert variant="destructive" className="mb-4">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Pending Amendment Request</AlertTitle>
              <AlertDescription>
                You already have a pending formalisation amendment request. Please wait for it to be
                processed before submitting another request.
              </AlertDescription>
            </Alert>
          ) : (
            <Alert className="mb-4">
              <Info className="h-4 w-4" />
              <AlertTitle>Review Process</AlertTitle>
              <AlertDescription>
                Changes will be submitted as an amendment application for review. Your current information
                will remain unchanged until the request is approved by an administrator.
              </AlertDescription>
            </Alert>
          )}

          <div className="space-y-6 pb-4">
          <div>
            <h3 className="text-lg font-semibold mb-4">Registration Details</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="registrationNumber">Registration Number</Label>
                <Input
                  id="registrationNumber"
                  value={formData.registrationNumber || ''}
                  onChange={(e) => {
                    const formatted = formatBusinessRegistration(e.target.value);
                    setFormData({ ...formData, registrationNumber: formatted });
                    // Clear error when user types
                    if (errors.registrationNumber) {
                      setErrors({ ...errors, registrationNumber: '' });
                    }
                  }}
                  placeholder={`e.g., ${EXAMPLE_FORMATS.BUSINESS_REG}`}
                  className={errors.registrationNumber ? 'border-destructive' : ''}
                  disabled={hasPendingAmendment}
                />
                <p className="text-xs text-muted-foreground">
                  Format: BRNR-XXXXXX (e.g., {EXAMPLE_FORMATS.BUSINESS_REG})
                </p>
                {currentData.registrationNumber && (
                  <p className="text-xs text-muted-foreground">
                    Current: {currentData.registrationNumber}
                  </p>
                )}
                {errors.registrationNumber && (
                  <p className="text-sm text-destructive">{errors.registrationNumber}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="taxIdentificationNumber">Tax Identification Number</Label>
                <Input
                  id="taxIdentificationNumber"
                  value={formData.taxIdentificationNumber || ''}
                  onChange={(e) => {
                    // Only allow digits, limit to 8 characters
                    const cleaned = e.target.value.replace(/\D/g, '').slice(0, 8);
                    setFormData({ ...formData, taxIdentificationNumber: cleaned });
                    // Clear error when user types
                    if (errors.taxIdentificationNumber) {
                      setErrors({ ...errors, taxIdentificationNumber: '' });
                    }
                  }}
                  placeholder={`e.g., ${EXAMPLE_FORMATS.TIN}`}
                  maxLength={8}
                  className={errors.taxIdentificationNumber ? 'border-destructive' : ''}
                  disabled={hasPendingAmendment}
                />
                <p className="text-xs text-muted-foreground">
                  8 digits (e.g., {EXAMPLE_FORMATS.TIN})
                </p>
                {currentData.taxIdentificationNumber && (
                  <p className="text-xs text-muted-foreground">
                    Current: {currentData.taxIdentificationNumber}
                  </p>
                )}
                {errors.taxIdentificationNumber && (
                  <p className="text-sm text-destructive">{errors.taxIdentificationNumber}</p>
                )}
              </div>
            </div>
          </div>

          <Separator />

          <div>
            <h3 className="text-lg font-semibold mb-4">Compliance Indicators</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="hasBankAccount" className="text-base">
                    Has Bank Account
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.hasBankAccount ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="hasBankAccount"
                  checked={formData.hasBankAccount ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, hasBankAccount: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="hasTaxClarification" className="text-base">
                    Has Tax Clarification
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.hasTaxClarification ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="hasTaxClarification"
                  checked={formData.hasTaxClarification ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, hasTaxClarification: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="isRegisteredForVat" className="text-base">
                    Registered for VAT
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.isRegisteredForVat ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="isRegisteredForVat"
                  checked={formData.isRegisteredForVat ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, isRegisteredForVat: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="isMemberOfAssociation" className="text-base">
                    Member of Association
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.isMemberOfAssociation ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="isMemberOfAssociation"
                  checked={formData.isMemberOfAssociation ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, isMemberOfAssociation: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="isAffiliated" className="text-base">
                    Is Affiliated
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.isAffiliated ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="isAffiliated"
                  checked={formData.isAffiliated ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, isAffiliated: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="hasExportLicense" className="text-base">
                    Has Export License
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.hasExportLicense ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="hasExportLicense"
                  checked={formData.hasExportLicense ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, hasExportLicense: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-4">
                <div className="flex-1">
                  <Label htmlFor="hasAccessedBds" className="text-base">
                    Has Accessed BDS
                  </Label>
                  <p className="text-xs text-muted-foreground mt-1">
                    Current: {currentData.hasAccessedBds ? 'Yes' : 'No'}
                  </p>
                </div>
                <Switch
                  id="hasAccessedBds"
                  checked={formData.hasAccessedBds ?? false}
                  onCheckedChange={(checked) =>
                    setFormData({ ...formData, hasAccessedBds: checked })
                  }
                />
              </div>
            </div>
          </div>

          <Separator />

          <div>
            <h3 className="text-lg font-semibold mb-4">Financial Information</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="annualTurnover">Annual Turnover (MWK)</Label>
                <Input
                  id="annualTurnover"
                  type="text"
                  inputMode="decimal"
                  value={formData.annualTurnover || ''}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      annualTurnover: formatCurrency(e.target.value),
                    })
                  }
                  placeholder={currentData.annualTurnover?.toString() || 'Enter amount'}
                />
                {currentData.annualTurnover && (
                  <p className="text-xs text-muted-foreground">
                    Current: MWK {currentData.annualTurnover.toLocaleString()}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="estimatedValueOfAssets">Estimated Value of Assets (MWK)</Label>
                <Input
                  id="estimatedValueOfAssets"
                  type="text"
                  inputMode="decimal"
                  value={formData.estimatedValueOfAssets || ''}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      estimatedValueOfAssets: formatCurrency(e.target.value),
                    })
                  }
                  placeholder={currentData.estimatedValueOfAssets?.toString() || 'Enter amount'}
                />
                {currentData.estimatedValueOfAssets && (
                  <p className="text-xs text-muted-foreground">
                    Current: MWK {currentData.estimatedValueOfAssets.toLocaleString()}
                  </p>
                )}
              </div>
            </div>
          </div>

          <Separator />

            <div className="space-y-2">
              <Label htmlFor="changeReason">Reason for Changes (Optional)</Label>
              <Textarea
                id="changeReason"
                value={formData.changeReason || ''}
                onChange={(e) => setFormData({ ...formData, changeReason: e.target.value })}
                placeholder="Provide a brief explanation for the proposed changes"
                rows={4}
                className="resize-none"
              />
              <p className="text-xs text-muted-foreground">
                Help administrators understand the context of your amendment request
              </p>
            </div>
          </div>
        </ScrollArea>

        <SheetFooter className="pt-4 border-t">
          <Button variant="outline" onClick={handleCancel} disabled={isSubmitting}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            disabled={!hasChanges || isSubmitting || hasPendingAmendment}
          >
            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Submit Amendment Request
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};
