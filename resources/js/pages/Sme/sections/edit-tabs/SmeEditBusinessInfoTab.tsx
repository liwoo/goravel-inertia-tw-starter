import React, { useState, useEffect } from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { X, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Checkbox } from '@/components/ui/checkbox';
import { Separator } from '@/components/ui/separator';
import { DISTRICT_NAMES } from '@/constants/districts';
import { SmeUpdateData } from '@/types/sme';
import {
  validateMalawiPhone,
  validateBusinessRegistration,
  validateTIN,
  formatMalawiPhone,
  formatBusinessRegistration,
  VALIDATION_MESSAGES,
  EXAMPLE_FORMATS
} from '@/lib/malawi-validators';
import axios from 'axios';

interface SmeEditBusinessInfoTabProps {
  businessData: SmeUpdateData;
  onChange: (data: SmeUpdateData) => void;
  errors: Record<string, string>;
}

export const SmeEditBusinessInfoTab: React.FC<SmeEditBusinessInfoTabProps> = ({
  businessData,
  onChange,
  errors
}) => {
  const [categoryOptions, setCategoryOptions] = useState<string[]>([]);
  const [sectorOptions, setSectorOptions] = useState<string[]>([]);
  const [improvementOptions, setImprovementOptions] = useState<string[]>([]);
  const [financingOptions, setFinancingOptions] = useState<string[]>([]);
  const [improvementAspect, setImprovementAspect] = useState('');
  const [financingSource, setFinancingSource] = useState('');
  const [improvementPopoverOpen, setImprovementPopoverOpen] = useState(false);
  const [financingPopoverOpen, setFinancingPopoverOpen] = useState(false);

  useEffect(() => {
    const fetchConfigs = async () => {
      try {
        const [categoryRes, sectorRes, improvementRes, financingRes] = await Promise.all([
          axios.get('/api/configs', { params: { config_type: 'Business Categories', pageSize: 100, sort: 'name', direction: 'ASC' }}),
          axios.get('/api/configs', { params: { config_type: 'Sectors', pageSize: 100, sort: 'name', direction: 'ASC' }}),
          axios.get('/api/configs', { params: { config_type: 'Improvement Aspects', pageSize: 100, sort: 'name', direction: 'ASC' }}),
          axios.get('/api/configs', { params: { config_type: 'Financing', pageSize: 100, sort: 'name', direction: 'ASC' }})
        ]);

        setCategoryOptions(categoryRes.data.data?.data?.map((c: any) => c.name) || []);
        setSectorOptions(sectorRes.data.data?.data?.map((s: any) => s.name) || []);
        setImprovementOptions(improvementRes.data.data?.data?.map((i: any) => i.name) || []);
        setFinancingOptions(financingRes.data.data?.data?.map((f: any) => f.name) || []);
      } catch (error) {
        console.error('Error fetching configs:', error);
      }
    };

    fetchConfigs();
  }, []);

  const addImprovementAspect = (value?: string) => {
    const aspectToAdd = value || improvementAspect.trim();
    if (aspectToAdd && !businessData.businessImprovementAspects.includes(aspectToAdd)) {
      onChange({
        ...businessData,
        businessImprovementAspects: [...businessData.businessImprovementAspects, aspectToAdd]
      });
      setImprovementAspect('');
    }
  };

  const removeImprovementAspect = (aspectToRemove: string) => {
    onChange({
      ...businessData,
      businessImprovementAspects: businessData.businessImprovementAspects.filter(a => a !== aspectToRemove)
    });
  };

  const addFinancingSource = (value?: string) => {
    const sourceToAdd = value || financingSource.trim();
    if (sourceToAdd && !businessData.businessAccessedFinancing.includes(sourceToAdd)) {
      onChange({
        ...businessData,
        businessAccessedFinancing: [...businessData.businessAccessedFinancing, sourceToAdd]
      });
      setFinancingSource('');
    }
  };

  const removeFinancingSource = (sourceToRemove: string) => {
    onChange({
      ...businessData,
      businessAccessedFinancing: businessData.businessAccessedFinancing.filter(s => s !== sourceToRemove)
    });
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Basic Information</CardTitle>
          <CardDescription>Core business details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">Business Name *</Label>
              <Input
                id="name"
                value={businessData.name}
                onChange={(e) => onChange({ ...businessData, name: e.target.value })}
                placeholder="Enter business name"
                className={errors.name ? 'border-destructive' : ''}
              />
              {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="registrationNumber">Registration Number</Label>
              <Input
                id="registrationNumber"
                value={businessData.registrationNumber}
                onChange={(e) => {
                  const formatted = formatBusinessRegistration(e.target.value);
                  onChange({ ...businessData, registrationNumber: formatted });
                }}
                placeholder={EXAMPLE_FORMATS.BUSINESS_REG}
                className={errors.registrationNumber ? 'border-destructive' : ''}
              />
              {errors.registrationNumber && <p className="text-sm text-destructive">{errors.registrationNumber}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="taxIdentificationNumber">Tax ID (TIN)</Label>
              <Input
                id="taxIdentificationNumber"
                value={businessData.taxIdentificationNumber}
                onChange={(e) => onChange({ ...businessData, taxIdentificationNumber: e.target.value })}
                placeholder={EXAMPLE_FORMATS.TIN}
                maxLength={8}
                className={errors.taxIdentificationNumber ? 'border-destructive' : ''}
              />
              {errors.taxIdentificationNumber && <p className="text-sm text-destructive">{errors.taxIdentificationNumber}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="operationalStartDate">Operational Start Date</Label>
              <Input
                id="operationalStartDate"
                type="date"
                value={businessData.operationalStartDate}
                onChange={(e) => onChange({ ...businessData, operationalStartDate: e.target.value })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="businessCategory">Category *</Label>
              <Select
                value={businessData.businessCategory}
                onValueChange={(value) => onChange({ ...businessData, businessCategory: value })}
              >
                <SelectTrigger className={errors.businessCategory ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                  {categoryOptions.map((cat) => (
                    <SelectItem key={cat} value={cat}>{cat}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.businessCategory && <p className="text-sm text-destructive">{errors.businessCategory}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="sector">Sector *</Label>
              <Select
                value={businessData.sector}
                onValueChange={(value) => onChange({ ...businessData, sector: value })}
              >
                <SelectTrigger className={errors.sector ? 'border-destructive' : ''}>
                  <SelectValue placeholder="Select sector" />
                </SelectTrigger>
                <SelectContent>
                  {sectorOptions.map((sec) => (
                    <SelectItem key={sec} value={sec}>{sec}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.sector && <p className="text-sm text-destructive">{errors.sector}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="subSector">Sub Sector</Label>
              <Input
                id="subSector"
                value={businessData.subSector}
                onChange={(e) => onChange({ ...businessData, subSector: e.target.value })}
                placeholder="Enter sub sector"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="contactPhone">Contact Phone *</Label>
              <Input
                id="contactPhone"
                value={businessData.contactPhone}
                onChange={(e) => {
                  const formatted = formatMalawiPhone(e.target.value);
                  onChange({ ...businessData, contactPhone: formatted });
                }}
                placeholder={EXAMPLE_FORMATS.PHONE}
                className={errors.contactPhone ? 'border-destructive' : ''}
              />
              {errors.contactPhone && <p className="text-sm text-destructive">{errors.contactPhone}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="contactEmail">Contact Email *</Label>
              <Input
                id="contactEmail"
                type="email"
                value={businessData.contactEmail}
                onChange={(e) => onChange({ ...businessData, contactEmail: e.target.value })}
                placeholder="Enter email"
                className={errors.contactEmail ? 'border-destructive' : ''}
              />
              {errors.contactEmail && <p className="text-sm text-destructive">{errors.contactEmail}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="website">Website</Label>
              <Input
                id="website"
                value={businessData.website}
                onChange={(e) => onChange({ ...businessData, website: e.target.value })}
                placeholder="https://example.com"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="district">District</Label>
              <Select
                value={businessData.district}
                onValueChange={(value) => onChange({ ...businessData, district: value })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select district" />
                </SelectTrigger>
                <SelectContent>
                  {DISTRICT_NAMES.map((dist) => (
                    <SelectItem key={dist} value={dist}>{dist}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="traditionalAuthority">Traditional Authority</Label>
              <Input
                id="traditionalAuthority"
                value={businessData.traditionalAuthority}
                onChange={(e) => onChange({ ...businessData, traditionalAuthority: e.target.value })}
                placeholder="Enter TA"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="businessDescription">Description</Label>
            <Textarea
              id="businessDescription"
              value={businessData.businessDescription}
              onChange={(e) => onChange({ ...businessData, businessDescription: e.target.value })}
              rows={3}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="physicalAddress">Physical Address</Label>
            <Textarea
              id="physicalAddress"
              value={businessData.physicalAddress}
              onChange={(e) => onChange({ ...businessData, physicalAddress: e.target.value })}
              rows={2}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="postalAddress">Postal Address</Label>
            <Textarea
              id="postalAddress"
              value={businessData.postalAddress}
              onChange={(e) => onChange({ ...businessData, postalAddress: e.target.value })}
              rows={2}
            />
          </div>
        </CardContent>
      </Card>

      {/* Business Development Card */}
      <Card>
        <CardHeader>
          <CardTitle>Business Development</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Business Improvement Aspects *</Label>
            <Popover open={improvementPopoverOpen} onOpenChange={setImprovementPopoverOpen}>
              <PopoverTrigger asChild>
                <Button variant="outline" className="w-full justify-between" type="button">
                  {businessData.businessImprovementAspects.length > 0
                    ? `${businessData.businessImprovementAspects.length} selected`
                    : 'Select improvement aspects...'}
                  <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-[400px] p-0" align="start">
                <div className="p-4 space-y-4">
                  <div className="space-y-2">
                    <Input
                      placeholder="Search or add custom..."
                      value={improvementAspect}
                      onChange={(e) => setImprovementAspect(e.target.value)}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter' && improvementAspect.trim()) {
                          e.preventDefault();
                          addImprovementAspect();
                        }
                      }}
                    />
                    {improvementAspect.trim() && !improvementOptions.includes(improvementAspect.trim()) && (
                      <Button variant="outline" size="sm" className="w-full" onClick={() => addImprovementAspect()} type="button">
                        <Plus className="mr-2 h-4 w-4" />
                        Add "{improvementAspect}"
                      </Button>
                    )}
                  </div>
                  <Separator />
                  <div className="max-h-[200px] overflow-auto space-y-2">
                    {improvementOptions
                      .filter(option => option.toLowerCase().includes(improvementAspect.toLowerCase()))
                      .map((option) => (
                        <div key={option} className="flex items-center space-x-2">
                          <Checkbox
                            id={`improvement-${option}`}
                            checked={businessData.businessImprovementAspects.includes(option)}
                            onCheckedChange={(checked) => {
                              if (checked) {
                                addImprovementAspect(option);
                              } else {
                                removeImprovementAspect(option);
                              }
                            }}
                          />
                          <label htmlFor={`improvement-${option}`} className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1">
                            {option}
                          </label>
                        </div>
                      ))}
                  </div>
                </div>
              </PopoverContent>
            </Popover>
            {errors.businessImprovementAspects && (
              <p className="text-sm text-destructive">{errors.businessImprovementAspects}</p>
            )}
            <div className="flex flex-wrap gap-2 mt-2">
              {businessData.businessImprovementAspects.map((aspect) => (
                <Badge key={aspect} variant="secondary" className="gap-1">
                  {aspect}
                  <X className="h-3 w-3 cursor-pointer" onClick={() => removeImprovementAspect(aspect)} />
                </Badge>
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <Label>Business Accessed Financing *</Label>
            <Popover open={financingPopoverOpen} onOpenChange={setFinancingPopoverOpen}>
              <PopoverTrigger asChild>
                <Button variant="outline" className="w-full justify-between" type="button">
                  {businessData.businessAccessedFinancing.length > 0
                    ? `${businessData.businessAccessedFinancing.length} selected`
                    : 'Select financing sources...'}
                  <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-[400px] p-0" align="start">
                <div className="p-4 space-y-4">
                  <div className="space-y-2">
                    <Input
                      placeholder="Search or add custom..."
                      value={financingSource}
                      onChange={(e) => setFinancingSource(e.target.value)}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter' && financingSource.trim()) {
                          e.preventDefault();
                          addFinancingSource();
                        }
                      }}
                    />
                    {financingSource.trim() && !financingOptions.includes(financingSource.trim()) && (
                      <Button variant="outline" size="sm" className="w-full" onClick={() => addFinancingSource()} type="button">
                        <Plus className="mr-2 h-4 w-4" />
                        Add "{financingSource}"
                      </Button>
                    )}
                  </div>
                  <Separator />
                  <div className="max-h-[200px] overflow-auto space-y-2">
                    {financingOptions
                      .filter(option => option.toLowerCase().includes(financingSource.toLowerCase()))
                      .map((option) => (
                        <div key={option} className="flex items-center space-x-2">
                          <Checkbox
                            id={`financing-${option}`}
                            checked={businessData.businessAccessedFinancing.includes(option)}
                            onCheckedChange={(checked) => {
                              if (checked) {
                                addFinancingSource(option);
                              } else {
                                removeFinancingSource(option);
                              }
                            }}
                          />
                          <label htmlFor={`financing-${option}`} className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1">
                            {option}
                          </label>
                        </div>
                      ))}
                  </div>
                </div>
              </PopoverContent>
            </Popover>
            {errors.businessAccessedFinancing && (
              <p className="text-sm text-destructive">{errors.businessAccessedFinancing}</p>
            )}
            <div className="flex flex-wrap gap-2 mt-2">
              {businessData.businessAccessedFinancing.map((source) => (
                <Badge key={source} variant="secondary" className="gap-1">
                  {source}
                  <X className="h-3 w-3 cursor-pointer" onClick={() => removeFinancingSource(source)} />
                </Badge>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
