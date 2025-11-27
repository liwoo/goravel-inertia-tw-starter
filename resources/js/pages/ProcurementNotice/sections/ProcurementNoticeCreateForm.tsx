import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { ProcurementNoticeCreateData } from '@/types/procurementnotice';
import { DISTRICT_OPTIONS } from '@/types/district';
import { Plus, X, Check } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import axios from 'axios';

interface ProcurementNoticeCreateFormProps extends CrudFormProps {
    setIsSaving?: (saving: boolean) => void;
    initialData?: any;
}

export const ProcurementNoticeCreateForm = forwardRef<any, ProcurementNoticeCreateFormProps>(({
    onSuccess,
    onError,
    onCancel,
    isLoading = false,
    setIsSaving,
    initialData
}, ref) => {
    const [formData, setFormData] = useState<ProcurementNoticeCreateData>({
        procured_by: '',
        procurement_type: '',
        market_approach: '',
        invitation: '',
        ref_no: '',
        open_date: '',
        close_date: '',
        is_published: false,
        organization: '',
        details: '',
        application_details: '',
        minimum_qualifying_score: 0,
        ...initialData,
        // Ensure arrays are initialized
        partners: initialData?.partners || [],
        qualifying_districts: initialData?.qualifying_districts || [],
        classification: initialData?.classification || [],
        interested_smes: initialData?.interested_smes || []
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    // Config options
    const [partnerOptions, setPartnerOptions] = useState<string[]>([]);
    const [classificationOptions, setClassificationOptions] = useState<string[]>([]);
    const [smeOptions, setSmeOptions] = useState<{ id: number, name: string }[]>([]);
    const [loadingConfigs, setLoadingConfigs] = useState(false);

    // Temporary state for inputs
    const [partnerInput, setPartnerInput] = useState('');
    const [classificationInput, setClassificationInput] = useState('');
    const [smeInput, setSmeInput] = useState('');

    const [partnerPopoverOpen, setPartnerPopoverOpen] = useState(false);
    const [districtPopoverOpen, setDistrictPopoverOpen] = useState(false);
    const [classificationPopoverOpen, setClassificationPopoverOpen] = useState(false);
    const [smePopoverOpen, setSmePopoverOpen] = useState(false);

    // Fetch config options on mount
    useEffect(() => {
        const fetchConfigOptions = async () => {
            setLoadingConfigs(true);
            try {
                const [partnersRes, classificationRes, smesRes] = await Promise.all([
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
                            config_type: 'Classification',
                            pageSize: 100,
                            sort: 'name',
                            direction: 'ASC'
                        }
                    }),
                    axios.get('/api/smes', {
                        params: {
                            pageSize: 100,
                            sort: 'business_name',
                            direction: 'ASC'
                        }
                    })
                ]);

                const partnersData = partnersRes.data.data?.data || [];
                const classificationData = classificationRes.data.data?.data || [];
                const smesData = smesRes.data.data?.data || [];

                setPartnerOptions(partnersData.map((config: any) => config.name));
                setClassificationOptions(classificationData.map((config: any) => config.name));
                setSmeOptions(smesData.map((sme: any) => ({ id: sme.id, name: sme.name })));
            } catch (error) {
                console.error('Error fetching config options:', error);
            } finally {
                setLoadingConfigs(false);
            }
        };

        fetchConfigOptions();
    }, []);

    // Generic Array Handlers
    const addItem = (field: keyof ProcurementNoticeCreateData, value: string) => {
        const currentArray = formData[field] as string[];
        if (value && !currentArray.includes(value)) {
            setFormData({
                ...formData,
                [field]: [...currentArray, value]
            });
        }
    };

    const removeItem = (field: keyof ProcurementNoticeCreateData, value: string) => {
        const currentArray = formData[field] as string[];
        setFormData({
            ...formData,
            [field]: currentArray.filter(item => item !== value)
        });
    };

    const toggleItem = (field: keyof ProcurementNoticeCreateData, value: string) => {
        const currentArray = formData[field] as string[];
        if (currentArray.includes(value)) {
            removeItem(field, value);
        } else {
            addItem(field, value);
        }
    };

    const handleSubmit = async () => {
        // Basic validation
        const newErrors: Record<string, string> = {};

        if (!formData.procured_by?.trim()) newErrors.procured_by = 'Procured By is required';
        if (!formData.procurement_type?.trim()) newErrors.procurement_type = 'Procurement Type is required';
        if (!formData.market_approach) newErrors.market_approach = 'Market Approach is required';
        if (!formData.invitation) newErrors.invitation = 'Invitation is required';
        if (!formData.ref_no?.trim()) newErrors.ref_no = 'Reference Number is required';
        if (!formData.open_date) newErrors.open_date = 'Open Date is required';
        if (!formData.close_date) newErrors.close_date = 'Close Date is required';
        if (!formData.organization?.trim()) newErrors.organization = 'Organization is required';

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        setErrors({});
        setIsSaving?.(true);

        try {
            const response = await fetch('/api/procurement-notices', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'X-Requested-With': 'XMLHttpRequest',
                },
                body: JSON.stringify(formData),
            });

            if (response.ok) {
                onSuccess('Procurement Notice created successfully');
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
                        <CardDescription>Core details of the procurement notice</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="procured_by">Procured By *</Label>
                                <Input
                                    id="procured_by"
                                    value={formData.procured_by}
                                    onChange={(e) => setFormData({ ...formData, procured_by: e.target.value })}
                                    placeholder="e.g. Ministry of Health"
                                    className={errors.procured_by ? 'border-destructive' : ''}
                                />
                                {errors.procured_by && <p className="text-sm text-destructive">{errors.procured_by}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="organization">Organization *</Label>
                                <Input
                                    id="organization"
                                    value={formData.organization}
                                    onChange={(e) => setFormData({ ...formData, organization: e.target.value })}
                                    placeholder="Enter organization"
                                    className={errors.organization ? 'border-destructive' : ''}
                                />
                                {errors.organization && <p className="text-sm text-destructive">{errors.organization}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="procurement_type">Procurement Type *</Label>
                                <Input
                                    id="procurement_type"
                                    value={formData.procurement_type}
                                    onChange={(e) => setFormData({ ...formData, procurement_type: e.target.value })}
                                    placeholder="e.g. Goods, Works, Services"
                                    className={errors.procurement_type ? 'border-destructive' : ''}
                                />
                                {errors.procurement_type && <p className="text-sm text-destructive">{errors.procurement_type}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="ref_no">Reference Number *</Label>
                                <Input
                                    id="ref_no"
                                    value={formData.ref_no}
                                    onChange={(e) => setFormData({ ...formData, ref_no: e.target.value })}
                                    placeholder="Enter reference number"
                                    className={errors.ref_no ? 'border-destructive' : ''}
                                />
                                {errors.ref_no && <p className="text-sm text-destructive">{errors.ref_no}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="market_approach">Market Approach *</Label>
                                <Select
                                    value={formData.market_approach}
                                    onValueChange={(value) => setFormData({ ...formData, market_approach: value })}
                                >
                                    <SelectTrigger className={errors.market_approach ? 'border-destructive' : ''}>
                                        <SelectValue placeholder="Select approach" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="National">National</SelectItem>
                                        <SelectItem value="International">International</SelectItem>
                                    </SelectContent>
                                </Select>
                                {errors.market_approach && <p className="text-sm text-destructive">{errors.market_approach}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="invitation">Invitation Type *</Label>
                                <Select
                                    value={formData.invitation}
                                    onValueChange={(value) => setFormData({ ...formData, invitation: value })}
                                >
                                    <SelectTrigger className={errors.invitation ? 'border-destructive' : ''}>
                                        <SelectValue placeholder="Select invitation type" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="open">Open</SelectItem>
                                        <SelectItem value="limited">Limited</SelectItem>
                                        <SelectItem value="single-source">Single Source</SelectItem>
                                    </SelectContent>
                                </Select>
                                {errors.invitation && <p className="text-sm text-destructive">{errors.invitation}</p>}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Dates & Settings */}
                <Card>
                    <CardHeader>
                        <CardTitle>Dates & Settings</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="open_date">Open Date *</Label>
                                <Input
                                    id="open_date"
                                    type="date"
                                    value={formData.open_date}
                                    onChange={(e) => setFormData({ ...formData, open_date: e.target.value })}
                                    className={errors.open_date ? 'border-destructive' : ''}
                                />
                                {errors.open_date && <p className="text-sm text-destructive">{errors.open_date}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="close_date">Close Date *</Label>
                                <Input
                                    id="close_date"
                                    type="date"
                                    value={formData.close_date}
                                    onChange={(e) => setFormData({ ...formData, close_date: e.target.value })}
                                    className={errors.close_date ? 'border-destructive' : ''}
                                />
                                {errors.close_date && <p className="text-sm text-destructive">{errors.close_date}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="minimum_qualifying_score">Min. Qualifying Score</Label>
                                <Input
                                    id="minimum_qualifying_score"
                                    type="number"
                                    min="0"
                                    max="100"
                                    value={formData.minimum_qualifying_score}
                                    onChange={(e) => setFormData({ ...formData, minimum_qualifying_score: parseInt(e.target.value) || 0 })}
                                />
                            </div>

                            <div className="flex items-center space-x-2 pt-8">
                                <Switch
                                    id="is_published"
                                    checked={formData.is_published}
                                    onCheckedChange={(checked) => setFormData({ ...formData, is_published: checked })}
                                />
                                <Label htmlFor="is_published">Published</Label>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Classification & Districts */}
                <Card>
                    <CardHeader>
                        <CardTitle>Classification & Location</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        {/* Classification */}
                        <div className="space-y-2">
                            <Label>Classification</Label>
                            <Popover open={classificationPopoverOpen} onOpenChange={setClassificationPopoverOpen}>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" className="w-full justify-between">
                                        {formData.classification.length > 0
                                            ? `${formData.classification.length} selected`
                                            : 'Select classification...'}
                                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                    <div className="p-4 space-y-4">
                                        <Input
                                            placeholder="Search..."
                                            value={classificationInput}
                                            onChange={(e) => setClassificationInput(e.target.value)}
                                        />
                                        <div className="max-h-[200px] overflow-auto space-y-2">
                                            {classificationOptions
                                                .filter(opt => opt.toLowerCase().includes(classificationInput.toLowerCase()))
                                                .map((option) => (
                                                    <div key={option} className="flex items-center space-x-2">
                                                        <Checkbox
                                                            id={`class-${option}`}
                                                            checked={formData.classification.includes(option)}
                                                            onCheckedChange={() => toggleItem('classification', option)}
                                                        />
                                                        <label htmlFor={`class-${option}`} className="text-sm cursor-pointer flex-1">
                                                            {option}
                                                        </label>
                                                    </div>
                                                ))}
                                        </div>
                                    </div>
                                </PopoverContent>
                            </Popover>
                            <div className="flex flex-wrap gap-2 mt-2">
                                {formData.classification.map((item) => (
                                    <Badge key={item} variant="secondary" className="gap-1">
                                        {item}
                                        <X className="h-3 w-3 cursor-pointer" onClick={() => removeItem('classification', item)} />
                                    </Badge>
                                ))}
                            </div>
                        </div>

                        <Separator />

                        {/* Qualifying Districts */}
                        <div className="space-y-2">
                            <Label>Qualifying Districts</Label>
                            <Popover open={districtPopoverOpen} onOpenChange={setDistrictPopoverOpen}>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" className="w-full justify-between">
                                        {formData.qualifying_districts.length > 0
                                            ? `${formData.qualifying_districts.length} selected`
                                            : 'Select districts...'}
                                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                    <div className="p-4 space-y-4">
                                        <div className="max-h-[300px] overflow-auto space-y-2">
                                            {DISTRICT_OPTIONS.map((option) => (
                                                <div key={option.value} className="flex items-center space-x-2">
                                                    <Checkbox
                                                        id={`district-${option.value}`}
                                                        checked={formData.qualifying_districts.includes(option.value)}
                                                        onCheckedChange={() => toggleItem('qualifying_districts', option.value)}
                                                    />
                                                    <label htmlFor={`district-${option.value}`} className="text-sm cursor-pointer flex-1">
                                                        {option.label}
                                                    </label>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </PopoverContent>
                            </Popover>
                            <div className="flex flex-wrap gap-2 mt-2">
                                {formData.qualifying_districts.map((item) => (
                                    <Badge key={item} variant="secondary" className="gap-1">
                                        {item}
                                        <X className="h-3 w-3 cursor-pointer" onClick={() => removeItem('qualifying_districts', item)} />
                                    </Badge>
                                ))}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Partners & SMEs */}
                <Card>
                    <CardHeader>
                        <CardTitle>Participants</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        {/* Partners */}
                        <div className="space-y-2">
                            <Label>Partners</Label>
                            <Popover open={partnerPopoverOpen} onOpenChange={setPartnerPopoverOpen}>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" className="w-full justify-between">
                                        {formData.partners.length > 0
                                            ? `${formData.partners.length} selected`
                                            : 'Select partners...'}
                                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                    <div className="p-4 space-y-4">
                                        <Input
                                            placeholder="Search or add custom..."
                                            value={partnerInput}
                                            onChange={(e) => setPartnerInput(e.target.value)}
                                            onKeyPress={(e) => {
                                                if (e.key === 'Enter' && partnerInput.trim()) {
                                                    e.preventDefault();
                                                    addItem('partners', partnerInput.trim());
                                                    setPartnerInput('');
                                                }
                                            }}
                                        />
                                        {partnerInput.trim() && !partnerOptions.includes(partnerInput.trim()) && (
                                            <Button
                                                variant="outline"
                                                size="sm"
                                                className="w-full"
                                                onClick={() => {
                                                    addItem('partners', partnerInput.trim());
                                                    setPartnerInput('');
                                                }}
                                                type="button"
                                            >
                                                <Plus className="mr-2 h-4 w-4" />
                                                Add "{partnerInput}"
                                            </Button>
                                        )}
                                        <div className="max-h-[200px] overflow-auto space-y-2">
                                            {partnerOptions
                                                .filter(opt => opt.toLowerCase().includes(partnerInput.toLowerCase()))
                                                .map((option) => (
                                                    <div key={option} className="flex items-center space-x-2">
                                                        <Checkbox
                                                            id={`partner-${option}`}
                                                            checked={formData.partners.includes(option)}
                                                            onCheckedChange={() => toggleItem('partners', option)}
                                                        />
                                                        <label htmlFor={`partner-${option}`} className="text-sm cursor-pointer flex-1">
                                                            {option}
                                                        </label>
                                                    </div>
                                                ))}
                                        </div>
                                    </div>
                                </PopoverContent>
                            </Popover>
                            <div className="flex flex-wrap gap-2 mt-2">
                                {formData.partners.map((item) => (
                                    <Badge key={item} variant="secondary" className="gap-1">
                                        {item}
                                        <X className="h-3 w-3 cursor-pointer" onClick={() => removeItem('partners', item)} />
                                    </Badge>
                                ))}
                            </div>
                        </div>

                        <Separator />

                        {/* Interested SMEs */}
                        <div className="space-y-2">
                            <Label>Interested SMEs</Label>
                            <Popover open={smePopoverOpen} onOpenChange={setSmePopoverOpen}>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" className="w-full justify-between">
                                        {formData.interested_smes.length > 0
                                            ? `${formData.interested_smes.length} selected`
                                            : 'Select SMEs...'}
                                        <Plus className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                    <div className="p-4 space-y-4">
                                        <Input
                                            placeholder="Search SMEs..."
                                            value={smeInput}
                                            onChange={(e) => setSmeInput(e.target.value)}
                                        />
                                        <div className="max-h-[200px] overflow-auto space-y-2">
                                            {smeOptions
                                                .filter(sme => sme.name.toLowerCase().includes(smeInput.toLowerCase()))
                                                .map((sme) => (
                                                    <div key={sme.id} className="flex items-center space-x-2">
                                                        <Checkbox
                                                            id={`sme-${sme.id}`}
                                                            checked={formData.interested_smes.includes(sme.id.toString())}
                                                            onCheckedChange={() => toggleItem('interested_smes', sme.id.toString())}
                                                        />
                                                        <label htmlFor={`sme-${sme.id}`} className="text-sm cursor-pointer flex-1">
                                                            {sme.name}
                                                        </label>
                                                    </div>
                                                ))}
                                        </div>
                                    </div>
                                </PopoverContent>
                            </Popover>
                            <div className="flex flex-wrap gap-2 mt-2">
                                {formData.interested_smes.map((smeId) => {
                                    const sme = smeOptions.find(s => s.id.toString() === smeId);
                                    return sme ? (
                                        <Badge key={smeId} variant="secondary" className="gap-1">
                                            {sme.name}
                                            <X className="h-3 w-3 cursor-pointer" onClick={() => removeItem('interested_smes', smeId)} />
                                        </Badge>
                                    ) : null;
                                })}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Details */}
                <Card>
                    <CardHeader>
                        <CardTitle>Details</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="details">Details *</Label>
                            <Textarea
                                id="details"
                                value={formData.details}
                                onChange={(e) => setFormData({ ...formData, details: e.target.value })}
                                placeholder="Enter details"
                                rows={5}
                                className={errors.details ? 'border-destructive' : ''}
                            />
                            {errors.details && <p className="text-sm text-destructive">{errors.details}</p>}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="application_details">Application Details *</Label>
                            <Textarea
                                id="application_details"
                                value={formData.application_details}
                                onChange={(e) => setFormData({ ...formData, application_details: e.target.value })}
                                placeholder="Enter application details"
                                rows={5}
                                className={errors.application_details ? 'border-destructive' : ''}
                            />
                            {errors.application_details && <p className="text-sm text-destructive">{errors.application_details}</p>}
                        </div>
                    </CardContent>
                </Card>

            </div>
        </form>
    );
});

ProcurementNoticeCreateForm.displayName = 'ProcurementNoticeCreateForm';
