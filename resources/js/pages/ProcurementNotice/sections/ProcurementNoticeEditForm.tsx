import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { ProcurementNoticeUpdateData, ProcurementNotice } from '@/types/procurementnotice';
import { DISTRICT_OPTIONS } from '@/types/district';
import { Plus, X, Check, ChevronLeft, ChevronRight, ChevronsUpDown } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Progress } from '@/components/ui/progress';
import axios from 'axios';
import MarkdownEditor from '@uiw/react-markdown-editor';
import { cn } from '@/lib/utils';

interface ProcurementNoticeEditFormProps extends CrudEditFormProps<ProcurementNotice> {
    setIsSaving?: (saving: boolean) => void;
}

const STEPS = [
    { id: 1, name: 'General Information', description: 'Basic details and settings' },
    { id: 2, name: 'Application Details', description: 'Detailed descriptions' },
];

export const ProcurementNoticeEditForm = forwardRef<any, ProcurementNoticeEditFormProps>(({
    onSuccess,
    onError,
    onCancel,
    isLoading = false,
    setIsSaving,
    item: initialData
}, ref) => {
    const [currentStep, setCurrentStep] = useState(1);
    const [formData, setFormData] = useState<ProcurementNoticeUpdateData>({
        procured_by: initialData.procured_by,
        procurement_type: initialData.procurement_type,
        market_approach: initialData.market_approach,
        invitation: initialData.invitation,
        ref_no: initialData.ref_no,
        open_date: initialData.open_date ? initialData.open_date.slice(0, 10) : '',
        close_date: initialData?.close_date ? initialData.close_date.slice(0, 10) : '',
        partners: initialData?.partners || [],
        qualifying_districts: initialData?.qualifying_districts || [],
        is_published: initialData?.is_published || false,
        organization: initialData?.organization || '',
        classification: initialData?.classification || [],
        interested_smes: initialData?.interested_smes || [],
        details: initialData?.details || '',
        application_details: initialData?.application_details || '',
        minimum_qualifying_score: initialData?.minimum_qualifying_score || 0,
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    // Config options
    const [partnerOptions, setPartnerOptions] = useState<string[]>([]);
    const [classificationOptions, setClassificationOptions] = useState<string[]>([]);
    const [procuredByOptions, setProcuredByOptions] = useState<string[]>([]);
    const [organizationOptions, setOrganizationOptions] = useState<string[]>([]);
    const [procurementTypeOptions, setProcurementTypeOptions] = useState<string[]>([]);

    const [loadingConfigs, setLoadingConfigs] = useState(false);

    // Temporary state for inputs
    const [partnerInput, setPartnerInput] = useState('');
    const [classificationInput, setClassificationInput] = useState('');
    const [procuredByInput, setProcuredByInput] = useState('');
    const [organizationInput, setOrganizationInput] = useState('');
    const [procurementTypeInput, setProcurementTypeInput] = useState('');

    const [partnerPopoverOpen, setPartnerPopoverOpen] = useState(false);
    const [districtPopoverOpen, setDistrictPopoverOpen] = useState(false);
    const [classificationPopoverOpen, setClassificationPopoverOpen] = useState(false);
    const [procuredByPopoverOpen, setProcuredByPopoverOpen] = useState(false);
    const [organizationPopoverOpen, setOrganizationPopoverOpen] = useState(false);
    const [procurementTypePopoverOpen, setProcurementTypePopoverOpen] = useState(false);

    const progress = (currentStep / STEPS.length) * 100;

    // Fetch config options on mount
    useEffect(() => {
        const fetchConfigOptions = async () => {
            setLoadingConfigs(true);
            try {
                const [partnersRes, classificationRes, procuredByRes, organizationRes, procurementTypeRes] = await Promise.all([
                    axios.get('/api/configs', { params: { config_type: 'Development Partners', pageSize: 100, sort: 'name', direction: 'ASC' } }),
                    axios.get('/api/configs', { params: { config_type: 'Procurement Classification', pageSize: 100, sort: 'name', direction: 'ASC' } }),
                    axios.get('/api/configs', { params: { config_type: 'Procured By', pageSize: 100, sort: 'name', direction: 'ASC' } }),
                    axios.get('/api/configs', { params: { config_type: 'Organization', pageSize: 100, sort: 'name', direction: 'ASC' } }),
                    axios.get('/api/configs', { params: { config_type: 'Procurement Type', pageSize: 100, sort: 'name', direction: 'ASC' } }),
                ]);

                const partnersData = partnersRes.data.data?.data || [];
                const classificationData = classificationRes.data.data?.data || [];
                const procuredByData = procuredByRes.data.data?.data || [];
                const organizationData = organizationRes.data.data?.data || [];
                const procurementTypeData = procurementTypeRes.data.data?.data || [];

                setPartnerOptions(partnersData.map((config: any) => config.name));
                setClassificationOptions(classificationData.map((config: any) => config.name));
                setProcuredByOptions(procuredByData.map((config: any) => config.name));
                setOrganizationOptions(organizationData.map((config: any) => config.name));
                setProcurementTypeOptions(procurementTypeData.map((config: any) => config.name));

            } catch (error) {
                console.error('Error fetching config options:', error);
            } finally {
                setLoadingConfigs(false);
            }
        };

        fetchConfigOptions();
    }, []);

    // Generic Array Handlers
    const addItem = (field: keyof ProcurementNoticeUpdateData, value: string) => {
        const currentArray = formData[field] as string[];
        if (value && !currentArray.includes(value)) {
            setFormData({
                ...formData,
                [field]: [...currentArray, value]
            });
        }
    };

    const removeItem = (field: keyof ProcurementNoticeUpdateData, value: string) => {
        const currentArray = formData[field] as string[];
        setFormData({
            ...formData,
            [field]: currentArray.filter(item => item !== value)
        });
    };

    const toggleItem = (field: keyof ProcurementNoticeUpdateData, value: string) => {
        const currentArray = formData[field] as string[];
        if (currentArray.includes(value)) {
            removeItem(field, value);
        } else {
            addItem(field, value);
        }
    };

    const validateStep = (step: number) => {
        const newErrors: Record<string, string> = {};

        if (step === 1) {
            if (!formData.procured_by?.trim()) newErrors.procured_by = 'Procured By is required';
            if (!formData.organization?.trim()) newErrors.organization = 'Organization is required';
            if (!formData.procurement_type?.trim()) newErrors.procurement_type = 'Procurement Type is required';
            if (!formData.market_approach) newErrors.market_approach = 'Market Approach is required';
            if (!formData.invitation) newErrors.invitation = 'Invitation is required';
            if (!formData.open_date) newErrors.open_date = 'Open Date is required';
            if (!formData.close_date) newErrors.close_date = 'Close Date is required';
        } else if (step === 2) {
            if (!formData.details?.trim()) newErrors.details = 'Details are required';
            if (!formData.application_details?.trim()) newErrors.application_details = 'Application Details are required';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleNext = () => {
        if (validateStep(currentStep)) {
            setCurrentStep(prev => Math.min(prev + 1, STEPS.length));
        }
    };

    const handlePrevious = () => {
        setCurrentStep(prev => Math.max(prev - 1, 1));
    };

    const handleSubmit = async () => {
        if (!validateStep(currentStep)) return;

        setIsSaving?.(true);

        try {
            const response = await fetch(`/api/procurement-notices/${initialData.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'X-Requested-With': 'XMLHttpRequest',
                },
                body: JSON.stringify(formData),
            });

            if (response.ok) {
                onSuccess('Procurement Notice updated successfully');
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
        handleSubmit: () => {
            if (currentStep === STEPS.length) {
                handleSubmit();
            } else {
                handleNext();
            }
        }
    }));

    // Helper to render a single-select Popover (Combobox style)
    const renderSingleSelectPopover = (
        label: string,
        value: string,
        onChange: (val: string) => void,
        options: string[],
        inputValue: string,
        setInputValue: (val: string) => void,
        open: boolean,
        setOpen: (val: boolean) => void,
        error?: string
    ) => (
        <div className="space-y-2">
            <Label>{label} *</Label>
            <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={open}
                        className={cn("w-full justify-between", error ? "border-destructive" : "")}
                    >
                        {value || `Select ${label}...`}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[400px] p-0" align="start">
                    <div className="p-4 space-y-4">
                        <Input
                            placeholder={`Search ${label}...`}
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                        />
                        {inputValue.trim() && !options.includes(inputValue.trim()) && (
                            <Button
                                variant="outline"
                                size="sm"
                                className="w-full"
                                onClick={() => {
                                    onChange(inputValue.trim());
                                    setInputValue('');
                                    setOpen(false);
                                }}
                                type="button"
                            >
                                <Plus className="mr-2 h-4 w-4" />
                                Add "{inputValue}"
                            </Button>
                        )}
                        <div className="max-h-[200px] overflow-auto space-y-2">
                            {options
                                .filter(opt => opt.toLowerCase().includes(inputValue.toLowerCase()))
                                .map((option) => (
                                    <div
                                        key={option}
                                        className={cn(
                                            "flex items-center space-x-2 p-2 rounded cursor-pointer hover:bg-accent",
                                            value === option ? "bg-accent" : ""
                                        )}
                                        onClick={() => {
                                            onChange(option);
                                            setOpen(false);
                                        }}
                                    >
                                        <Check
                                            className={cn(
                                                "mr-2 h-4 w-4",
                                                value === option ? "opacity-100" : "opacity-0"
                                            )}
                                        />
                                        <span className="text-sm flex-1">{option}</span>
                                    </div>
                                ))}
                        </div>
                    </div>
                </PopoverContent>
            </Popover>
            {error && <p className="text-sm text-destructive">{error}</p>}
        </div>
    );

    return (
        <div className="space-y-6">
            {/* Progress Bar */}
            <div className="space-y-2">
                <div className="flex justify-between items-center">
                    <div className="space-y-1">
                        <h3 className="text-lg font-semibold">{STEPS[currentStep - 1].name}</h3>
                        <p className="text-sm text-muted-foreground">{STEPS[currentStep - 1].description}</p>
                    </div>
                    <Badge variant="outline" className="text-xs">
                        Step {currentStep} of {STEPS.length}
                    </Badge>
                </div>
                <Progress value={progress} className="h-2" />
            </div>

            {/* Step Indicators */}
            <div className="flex justify-between">
                {STEPS.map((step, index) => (
                    <div key={step.id} className={`flex items-center ${index < STEPS.length - 1 ? 'flex-1' : ''}`}>
                        <div className="flex flex-col items-center">
                            <div
                                className={`w-10 h-10 rounded-full flex items-center justify-center border-2 ${currentStep > step.id
                                    ? 'bg-primary border-primary text-primary-foreground'
                                    : currentStep === step.id
                                        ? 'border-primary text-primary'
                                        : 'border-muted text-muted-foreground'
                                    }`}
                            >
                                {currentStep > step.id ? <Check className="h-5 w-5" /> : <span className="text-sm font-semibold">{step.id}</span>}
                            </div>
                            <span className="text-xs mt-2 text-center hidden sm:block">{step.name}</span>
                        </div>
                        {index < STEPS.length - 1 && (
                            <div className={`flex-1 h-0.5 mx-2 ${currentStep > step.id ? 'bg-primary' : 'bg-muted'}`} />
                        )}
                    </div>
                ))}
            </div>

            <Separator />

            <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
                {/* Step 1: General Information */}
                {currentStep === 1 && (
                    <div className="space-y-6">
                        <Card>
                            <CardHeader>
                                <CardTitle>Basic Information</CardTitle>
                                <CardDescription>Core details of the procurement notice</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    {renderSingleSelectPopover(
                                        "Procured By",
                                        formData.procured_by,
                                        (val) => setFormData({ ...formData, procured_by: val }),
                                        procuredByOptions,
                                        procuredByInput,
                                        setProcuredByInput,
                                        procuredByPopoverOpen,
                                        setProcuredByPopoverOpen,
                                        errors.procured_by
                                    )}

                                    {renderSingleSelectPopover(
                                        "Organization",
                                        formData.organization,
                                        (val) => setFormData({ ...formData, organization: val }),
                                        organizationOptions,
                                        organizationInput,
                                        setOrganizationInput,
                                        organizationPopoverOpen,
                                        setOrganizationPopoverOpen,
                                        errors.organization
                                    )}

                                    {renderSingleSelectPopover(
                                        "Procurement Type",
                                        formData.procurement_type,
                                        (val) => setFormData({ ...formData, procurement_type: val }),
                                        procurementTypeOptions,
                                        procurementTypeInput,
                                        setProcurementTypeInput,
                                        procurementTypePopoverOpen,
                                        setProcurementTypePopoverOpen,
                                        errors.procurement_type
                                    )}

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
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle>Classification & Participants</CardTitle>
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

                                <Separator />

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
                            </CardContent>
                        </Card>
                    </div>
                )}

                {/* Step 2: Details */}
                {currentStep === 2 && (
                    <Card>
                        <CardHeader>
                            <CardTitle>Details</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="details">Details *</Label>
                                <div className="min-h-[400px] border rounded-md p-1">
                                    <MarkdownEditor
                                        value={formData.details}
                                        onChange={(value) => setFormData({ ...formData, details: value })}
                                        height="400px"
                                    />
                                </div>
                                {errors.details && <p className="text-sm text-destructive">{errors.details}</p>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="application_details">Application Details *</Label>
                                <div className="min-h-[400px] border rounded-md p-1">
                                    <MarkdownEditor
                                        value={formData.application_details}
                                        onChange={(value) => setFormData({ ...formData, application_details: value })}
                                        height="400px"
                                    />
                                </div>
                                {errors.application_details && <p className="text-sm text-destructive">{errors.application_details}</p>}
                            </div>
                        </CardContent>
                    </Card>
                )}

                <div className="flex justify-between pt-4">
                    <Button
                        type="button"
                        variant="outline"
                        onClick={handlePrevious}
                        disabled={currentStep === 1}
                    >
                        <ChevronLeft className="mr-2 h-4 w-4" />
                        Previous
                    </Button>

                    {currentStep < STEPS.length ? (
                        <Button
                            type="button"
                            onClick={handleNext}
                            disabled={isLoading}
                        >
                            Next
                            <ChevronRight className="h-4 w-4 ml-2" />
                        </Button>
                    ) : (
                        <Button
                            type="button"
                            onClick={handleSubmit}
                            disabled={isLoading}
                        >
                            {isLoading ? 'Creating...' : 'Create Procurement Notice'}
                        </Button>
                    )}
                </div>
            </form>
        </div>
    );
});

ProcurementNoticeEditForm.displayName = 'ProcurementNoticeEditForm';
