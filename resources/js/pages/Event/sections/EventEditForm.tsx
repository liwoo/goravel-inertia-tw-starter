import React, { useState, forwardRef, useImperativeHandle, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { Event, EventUpdateData } from '@/types/event';
import { Hash, BookOpen, FileText, Calendar, MapPin, Users, Plus, X } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { DISTRICT_OPTIONS } from '@/types/district';
import axios from 'axios';

interface EventEditFormProps extends CrudEditFormProps<Event> {
  setIsSaving?: (saving: boolean) => void;
}

export const EventEditForm = forwardRef<any, EventEditFormProps>(({
  item: event,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<EventUpdateData>({
    title: event.title,
    description: event.description,
    date: event.date ? event.date.slice(0, 10) : '',
    end_date: event.end_date ? event.end_date.slice(0, 10) : '',
    venue: event.venue,
    partners: event.partners || [],
    district: event.district,
    notes: event.notes || '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  // Config options
  const [partnerOptions, setPartnerOptions] = useState<string[]>([]);
  const [loadingConfigs, setLoadingConfigs] = useState(false);

  // Temporary state for inputs
  const [partnerInput, setPartnerInput] = useState('');
  const [partnerPopoverOpen, setPartnerPopoverOpen] = useState(false);

  // Fetch config options on mount
  useEffect(() => {
    const fetchConfigOptions = async () => {
      setLoadingConfigs(true);
      try {
        const partnersRes = await axios.get('/api/configs', {
          params: {
            config_type: 'Development Partners',
            pageSize: 100,
            sort: 'name',
            direction: 'ASC'
          }
        });

        const partnersData = partnersRes.data.data?.data || [];

        setPartnerOptions(partnersData.map((config: any) => config.name));
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
    if (partnerToAdd && !formData.partners?.includes(partnerToAdd)) {
      setFormData({
        ...formData,
        partners: [...(formData.partners || []), partnerToAdd]
      });
      setPartnerInput('');
    }
  };

  const removePartner = (partner: string) => {
    setFormData({
      ...formData,
      partners: formData.partners?.filter(p => p !== partner) || []
    });
  };


  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    if (!formData.title?.trim()) newErrors.title = 'Title is required';
    if (!formData.description?.trim()) newErrors.description = 'Description is required';
    if (!formData.date) newErrors.date = 'Date is required';
    if (!formData.venue?.trim()) newErrors.venue = 'Venue is required';
    if (!formData.district?.trim()) newErrors.district = 'District is required';

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`/api/events/${event.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Event updated successfully');
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
  console.log(formData)
  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">

        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle>Event Details</CardTitle>
            <CardDescription>Update the core event information</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="title">Title *</Label>
                <Input
                  id="title"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder="Enter title"
                  className={errors.title ? 'border-destructive' : ''}
                />
                {errors.title && <p className="text-sm text-destructive">{errors.title}</p>}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="date">Date *</Label>
                <Input
                  id="date"
                  type="date"
                  value={formData.date}
                  onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                  className={errors.date ? 'border-destructive' : ''}
                />
                {errors.date && <p className="text-sm text-destructive">{errors.date}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="end_date">End Date</Label>
                <Input
                  id="end_date"
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  className={errors.end_date ? 'border-destructive' : ''}
                />
                {errors.end_date && <p className="text-sm text-destructive">{errors.end_date}</p>}
              </div>

              <div className="space-y-2 md:col-span-2">
                <Label htmlFor="description">Description *</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Enter description"
                  rows={3}
                  className={errors.description ? 'border-destructive' : ''}
                />
                {errors.description && <p className="text-sm text-destructive">{errors.description}</p>}
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
                <Label htmlFor="venue">Venue *</Label>
                <Input
                  id="venue"
                  value={formData.venue}
                  onChange={(e) => setFormData({ ...formData, venue: e.target.value })}
                  placeholder="Enter venue"
                  className={errors.venue ? 'border-destructive' : ''}
                />
                {errors.venue && <p className="text-sm text-destructive">{errors.venue}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="district">District *</Label>
                <Select
                  value={formData.district}
                  onValueChange={(value) => setFormData({ ...formData, district: value })}
                >
                  <SelectTrigger className={errors.district ? 'border-destructive' : ''}>
                    <SelectValue placeholder="Select district" />
                  </SelectTrigger>
                  <SelectContent>
                    {DISTRICT_OPTIONS.map((option) => (
                      <SelectItem key={option.value} value={option.value}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.district && <p className="text-sm text-destructive">{errors.district}</p>}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Participants */}
        <Card>
          <CardHeader>
            <CardTitle>Participants</CardTitle>
            <CardDescription>Manage event partners</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
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
                    {(formData.partners?.length || 0) > 0
                      ? `${formData.partners?.length} selected`
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
                              checked={formData.partners?.includes(option)}
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
                {formData.partners?.map((partner) => (
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
          </CardContent>
        </Card>

        {/* Additional Info */}
        <Card>
          <CardHeader>
            <CardTitle>Additional Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Label htmlFor="notes">Notes</Label>
              <Textarea
                id="notes"
                value={formData.notes || ''}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Enter notes"
                rows={3}
              />
            </div>
          </CardContent>
        </Card>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">ID</p>
              <p className="font-medium text-foreground">#{event.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">
                {(event.createdAt || event.created_at) ? new Date(event.createdAt || event.created_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">
                {(event.updatedAt || event.updated_at) ? new Date(event.updatedAt || event.updated_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
          </div>
        </div>

      </div>
    </form>
  );
});

EventEditForm.displayName = 'EventEditForm';
