import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { BdspCreateData, BdspService } from '@/types/bdsp';
import { User, MapPin, Tag, FileText, Plus, X, DollarSign, Clock } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Checkbox } from '@/components/ui/checkbox';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import axios from 'axios';

interface BdspCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
}

export const BdspCreateForm = forwardRef<any, BdspCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<BdspCreateData>({
    name: '',
    postal_address: '',
    physical_address: '',
    registration_status: '',
    partners: [],
    product_types: [],
    service_list: [],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  // Config options from API
  const [partnerOptions, setPartnerOptions] = useState<string[]>([]);
  const [productOptions, setProductOptions] = useState<string[]>([]);
  const [loadingConfigs, setLoadingConfigs] = useState(false);

  // Temporary state for inputs
  const [partnerInput, setPartnerInput] = useState('');
  const [productInput, setProductInput] = useState('');
  const [partnerPopoverOpen, setPartnerPopoverOpen] = useState(false);
  const [productPopoverOpen, setProductPopoverOpen] = useState(false);

  const [newService, setNewService] = useState<BdspService>({
    name: '',
    cost: 0,
    duration: ''
  });

  // Fetch config options on mount
  useEffect(() => {
    const fetchConfigOptions = async () => {
      setLoadingConfigs(true);
      try {
        const [partnersRes, productsRes] = await Promise.all([
          axios.get('/api/configs', {
            params: {
              config_type: 'Partners',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          }),
          axios.get('/api/configs', {
            params: {
              config_type: 'Product Types',
              pageSize: 100,
              sort: 'name',
              direction: 'ASC'
            }
          })
        ]);

        const partnersData = partnersRes.data.data?.data || [];
        const productsData = productsRes.data.data?.data || [];

        setPartnerOptions(partnersData.map((config: any) => config.name));
        setProductOptions(productsData.map((config: any) => config.name));
      } catch (error) {
        console.error('Error fetching config options:', error);
      } finally {
        setLoadingConfigs(false);
      }
    };

    fetchConfigOptions();
  }, []);

  // Partner Handlers
  const addPartner = (value?: string) => {
    const partnerToAdd = value || partnerInput.trim();
    if (partnerToAdd && !formData.partners.includes(partnerToAdd)) {
      setFormData({
        ...formData,
        partners: [...formData.partners, partnerToAdd]
      });
      setPartnerInput('');
    }
  };

  const removePartner = (partner: string) => {
    setFormData({
      ...formData,
      partners: formData.partners.filter(p => p !== partner)
    });
  };

  // Product Type Handlers
  const addProduct = (value?: string) => {
    const productToAdd = value || productInput.trim();
    if (productToAdd && !formData.product_types.includes(productToAdd)) {
      setFormData({
        ...formData,
        product_types: [...formData.product_types, productToAdd]
      });
      setProductInput('');
    }
  };

  const removeProduct = (product: string) => {
    setFormData({
      ...formData,
      product_types: formData.product_types.filter(p => p !== product)
    });
  };

  // Service Repeater Handlers
  const addService = () => {
    if (newService.name.trim() && newService.duration.trim()) {
      setFormData({
        ...formData,
        service_list: [...formData.service_list, { ...newService }]
      });
      setNewService({ name: '', cost: 0, duration: '' });
    }
  };

  const removeService = (index: number) => {
    setFormData({
      ...formData,
      service_list: formData.service_list.filter((_, i) => i !== index)
    });
  };

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    if (!formData.name?.trim()) newErrors.name = 'Name is required';
    if (!formData.postal_address?.trim()) newErrors.postalAddress = 'Postal Address is required';
    if (formData.service_list.length === 0) newErrors.serviceList = 'At least one service is required';

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch('/api/bdsps', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Bdsp created successfully');
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

        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>Enter the core BDSP details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter name"
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="registrationStatus">Registration Status</Label>
                <Select
                  value={formData.registration_status}
                  onValueChange={(value) => setFormData({ ...formData, registration_status: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Pending">Pending</SelectItem>
                    <SelectItem value="Active">Active</SelectItem>
                    <SelectItem value="Rejected">Rejected</SelectItem>
                    <SelectItem value="Inactive">Inactive</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Location */}
        <Card>
          <CardHeader>
            <CardTitle>Location</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="postalAddress">Postal Address *</Label>
                <Textarea
                  id="postalAddress"
                  value={formData.postal_address}
                  onChange={(e) => setFormData({ ...formData, postal_address: e.target.value })}
                  placeholder="Enter postal address"
                  rows={3}
                  className={errors.postalAddress ? 'border-destructive' : ''}
                />
                {errors.postalAddress && <p className="text-sm text-destructive">{errors.postalAddress}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="physicalAddress">Physical Address</Label>
                <Textarea
                  id="physicalAddress"
                  value={formData.physical_address}
                  onChange={(e) => setFormData({ ...formData, physical_address: e.target.value })}
                  placeholder="Enter physical address"
                  rows={3}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Partners & Products */}
        <Card>
          <CardHeader>
            <CardTitle>Partners & Products</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Partners */}
            <div className="space-y-2">
              <Label>Partners</Label>
              <Popover open={partnerPopoverOpen} onOpenChange={setPartnerPopoverOpen}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className="w-full justify-between"
                    type="button"
                  >
                    {formData.partners.length > 0
                      ? `${formData.partners.length} selected`
                      : 'Select partners...'}
                    <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[400px] p-0" align="start">
                  <div className="p-4 space-y-4">
                    <div className="space-y-2">
                      <Input
                        placeholder="Search or add custom..."
                        value={partnerInput}
                        onChange={(e) => setPartnerInput(e.target.value)}
                        onKeyPress={(e) => {
                          if (e.key === 'Enter' && partnerInput.trim()) {
                            e.preventDefault();
                            addPartner();
                          }
                        }}
                      />
                      {partnerInput.trim() && !partnerOptions.includes(partnerInput.trim()) && (
                        <Button
                          variant="outline"
                          size="sm"
                          className="w-full"
                          onClick={() => addPartner()}
                          type="button"
                        >
                          <Plus className="mr-2 h-4 w-4" />
                          Add "{partnerInput}"
                        </Button>
                      )}
                    </div>
                    <Separator />
                    <div className="max-h-[200px] overflow-auto space-y-2">
                      {partnerOptions
                        .filter(option =>
                          option.toLowerCase().includes(partnerInput.toLowerCase())
                        )
                        .map((option) => (
                          <div key={option} className="flex items-center space-x-2">
                            <Checkbox
                              id={`partner-${option}`}
                              checked={formData.partners.includes(option)}
                              onCheckedChange={(checked) => {
                                if (checked) {
                                  addPartner(option);
                                } else {
                                  removePartner(option);
                                }
                              }}
                            />
                            <label
                              htmlFor={`partner-${option}`}
                              className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1"
                            >
                              {option}
                            </label>
                          </div>
                        ))}
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
              <div className="flex flex-wrap gap-2 mt-2">
                {formData.partners.map((partner) => (
                  <Badge key={partner} variant="secondary" className="gap-1">
                    {partner}
                    <X
                      className="h-3 w-3 cursor-pointer hover:text-destructive"
                      onClick={() => removePartner(partner)}
                    />
                  </Badge>
                ))}
              </div>
            </div>

            <Separator />

            {/* Product Types */}
            <div className="space-y-2">
              <Label>Product Types</Label>
              <Popover open={productPopoverOpen} onOpenChange={setProductPopoverOpen}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className="w-full justify-between"
                    type="button"
                  >
                    {formData.product_types.length > 0
                      ? `${formData.product_types.length} selected`
                      : 'Select product types...'}
                    <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[400px] p-0" align="start">
                  <div className="p-4 space-y-4">
                    <div className="space-y-2">
                      <Input
                        placeholder="Search or add custom..."
                        value={productInput}
                        onChange={(e) => setProductInput(e.target.value)}
                        onKeyPress={(e) => {
                          if (e.key === 'Enter' && productInput.trim()) {
                            e.preventDefault();
                            addProduct();
                          }
                        }}
                      />
                      {productInput.trim() && !productOptions.includes(productInput.trim()) && (
                        <Button
                          variant="outline"
                          size="sm"
                          className="w-full"
                          onClick={() => addProduct()}
                          type="button"
                        >
                          <Plus className="mr-2 h-4 w-4" />
                          Add "{productInput}"
                        </Button>
                      )}
                    </div>
                    <Separator />
                    <div className="max-h-[200px] overflow-auto space-y-2">
                      {productOptions
                        .filter(option =>
                          option.toLowerCase().includes(productInput.toLowerCase())
                        )
                        .map((option) => (
                          <div key={option} className="flex items-center space-x-2">
                            <Checkbox
                              id={`product-${option}`}
                              checked={formData.product_types.includes(option)}
                              onCheckedChange={(checked) => {
                                if (checked) {
                                  addProduct(option);
                                } else {
                                  removeProduct(option);
                                }
                              }}
                            />
                            <label
                              htmlFor={`product-${option}`}
                              className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1"
                            >
                              {option}
                            </label>
                          </div>
                        ))}
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
              <div className="flex flex-wrap gap-2 mt-2">
                {formData.product_types.map((product) => (
                  <Badge key={product} variant="secondary" className="gap-1">
                    {product}
                    <X
                      className="h-3 w-3 cursor-pointer hover:text-destructive"
                      onClick={() => removeProduct(product)}
                    />
                  </Badge>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Services */}
        <Card>
          <CardHeader>
            <CardTitle>Services</CardTitle>
            <CardDescription>List the services offered by this BDSP</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end border p-4 rounded-lg bg-muted/20">
              <div className="space-y-2 md:col-span-2">
                <Label>Service Name</Label>
                <Input
                  value={newService.name}
                  onChange={(e) => setNewService({ ...newService, name: e.target.value })}
                  placeholder="e.g. Business Training"
                />
              </div>
              <div className="space-y-2">
                <Label>Cost</Label>
                <Input
                  type="number"
                  value={newService.cost}
                  onChange={(e) => setNewService({ ...newService, cost: parseFloat(e.target.value) || 0 })}
                  placeholder="0.00"
                />
              </div>
              <div className="space-y-2">
                <Label>Duration</Label>
                <Input
                  value={newService.duration}
                  onChange={(e) => setNewService({ ...newService, duration: e.target.value })}
                  placeholder="e.g. 2 weeks"
                />
              </div>
              <Button type="button" onClick={addService} className="md:col-span-4 w-full">
                <Plus className="mr-2 h-4 w-4" /> Add Service
              </Button>
            </div>

            {formData.service_list.length > 0 && (
              <div className="border rounded-lg overflow-hidden">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Service Name</TableHead>
                      <TableHead>Cost</TableHead>
                      <TableHead>Duration</TableHead>
                      <TableHead className="w-[50px]"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {formData.service_list.map((service, index) => (
                      <TableRow key={index}>
                        <TableCell className="font-medium">{service.name}</TableCell>
                        <TableCell>{service.cost.toLocaleString()}</TableCell>
                        <TableCell>{service.duration}</TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => removeService(index)}
                            className="h-8 w-8 text-destructive hover:text-destructive/90"
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
            {errors.serviceList && <p className="text-sm text-destructive">{errors.serviceList}</p>}
          </CardContent>
        </Card>

      </div>
    </form>
  );
});

BdspCreateForm.displayName = 'BdspCreateForm';
