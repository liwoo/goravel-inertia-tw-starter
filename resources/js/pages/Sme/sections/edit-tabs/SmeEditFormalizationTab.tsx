import React from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { BusinessFormalisation } from '@/types/business_formalisation';

// Partial type for form state that doesn't require id/smeId
export type FormalizationFormState = Partial<BusinessFormalisation> & {
  hasBankAccount: boolean;
  hasTaxClarification: boolean;
  isRegisteredForVat: boolean;
  isMemberOfAssociation: boolean;
  isAffiliated: boolean;
  hasExportLicense: boolean;
  hasAccessedBds: boolean;
  annualTurnover: number;
  estimatedValueOfAssets: number;
  formalisationScore: number;
};

interface SmeEditFormalizationTabProps {
  formalization: FormalizationFormState | BusinessFormalisation | undefined;
  onChange: (data: FormalizationFormState) => void;
  errors: Record<string, string>;
}

export const SmeEditFormalizationTab: React.FC<SmeEditFormalizationTabProps> = ({
  formalization,
  onChange,
  errors
}) => {
  // Initialize empty formalization if it doesn't exist
  const currentFormalization: FormalizationFormState = formalization || {
    hasBankAccount: false,
    hasTaxClarification: false,
    isRegisteredForVat: false,
    isMemberOfAssociation: false,
    isAffiliated: false,
    hasExportLicense: false,
    hasAccessedBds: false,
    annualTurnover: 0,
    estimatedValueOfAssets: 0,
    formalisationScore: 0,
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Business Formalization Status</CardTitle>
          <CardDescription>
            Indicate the formal compliance and registration status of the business
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
                  checked={currentFormalization.hasBankAccount}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, hasBankAccount: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-2">
                <Label htmlFor="hasTaxClarification">Has Tax Clarification Certificate</Label>
                <Switch
                  id="hasTaxClarification"
                  checked={currentFormalization.hasTaxClarification}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, hasTaxClarification: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-2">
                <Label htmlFor="isRegisteredForVat">Registered for VAT</Label>
                <Switch
                  id="isRegisteredForVat"
                  checked={currentFormalization.isRegisteredForVat}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, isRegisteredForVat: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-2">
                <Label htmlFor="hasExportLicense">Has Export License</Label>
                <Switch
                  id="hasExportLicense"
                  checked={currentFormalization.hasExportLicense}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, hasExportLicense: checked })
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
                  checked={currentFormalization.isMemberOfAssociation}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, isMemberOfAssociation: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-2">
                <Label htmlFor="isAffiliated">Is Affiliated</Label>
                <Switch
                  id="isAffiliated"
                  checked={currentFormalization.isAffiliated}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, isAffiliated: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between space-x-2">
                <Label htmlFor="hasAccessedBds">Has Accessed Business Development Services</Label>
                <Switch
                  id="hasAccessedBds"
                  checked={currentFormalization.hasAccessedBds}
                  onCheckedChange={(checked) =>
                    onChange({ ...currentFormalization, hasAccessedBds: checked })
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
                  type="number"
                  value={currentFormalization.annualTurnover}
                  onChange={(e) =>
                    onChange({ ...currentFormalization, annualTurnover: parseFloat(e.target.value) || 0 })
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
                  type="number"
                  value={currentFormalization.estimatedValueOfAssets}
                  onChange={(e) =>
                    onChange({ ...currentFormalization, estimatedValueOfAssets: parseFloat(e.target.value) || 0 })
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
  );
};
