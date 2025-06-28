import React from 'react';
import { X, Filter } from 'lucide-react';
import { cn } from '@/lib/utils';
import { FilterPanelProps, CrudFilter } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';

export function FilterPanel({ 
  filters, 
  values, 
  onChange, 
  onClear,
  className 
}: FilterPanelProps) {
  const hasActiveFilters = Object.keys(values).some(key => 
    values[key] !== undefined && values[key] !== '' && values[key] !== null
  );

  const activeFilterCount = Object.keys(values).filter(key => 
    values[key] !== undefined && values[key] !== '' && values[key] !== null
  ).length;

  const renderFilterInput = (filter: CrudFilter) => {
    const value = values[filter.key];

    switch (filter.type) {
      case 'select':
        return (
          <Select
            value={value || ''}
            onValueChange={(newValue) => onChange(filter.key, newValue)}
          >
            <SelectTrigger>
              <SelectValue placeholder={filter.placeholder || 'Select option'} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">All</SelectItem>
              {filter.options?.map((option) => (
                <SelectItem key={option.value} value={String(option.value)}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        );

      case 'boolean':
        return (
          <div className="flex items-center space-x-2">
            <Checkbox
              id={filter.key}
              checked={value === true || value === 'true'}
              onCheckedChange={(checked) => onChange(filter.key, checked)}
            />
            <Label htmlFor={filter.key} className="text-sm font-normal">
              {filter.placeholder || filter.label}
            </Label>
          </div>
        );

      case 'date':
        return (
          <Input
            type="date"
            value={value || ''}
            onChange={(e) => onChange(filter.key, e.target.value)}
            placeholder={filter.placeholder}
          />
        );

      case 'number':
        return (
          <Input
            type="number"
            value={value || ''}
            onChange={(e) => onChange(filter.key, e.target.value)}
            placeholder={filter.placeholder}
          />
        );

      default: // text
        return (
          <Input
            type="text"
            value={value || ''}
            onChange={(e) => onChange(filter.key, e.target.value)}
            placeholder={filter.placeholder}
          />
        );
    }
  };

  return (
    <div className={cn(
      'bg-white border border-gray-200 rounded-lg p-4 space-y-4',
      className
    )}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Filter className="h-4 w-4 text-gray-500" />
          <h3 className="text-sm font-medium text-gray-900">Filters</h3>
          {activeFilterCount > 0 && (
            <Badge variant="secondary" className="text-xs">
              {activeFilterCount}
            </Badge>
          )}
        </div>
        
        {hasActiveFilters && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onClear}
            className="text-sm text-blue-600 hover:text-blue-800 h-auto p-0"
          >
            Clear all
          </Button>
        )}
      </div>

      {/* Active Filters Summary */}
      {hasActiveFilters && (
        <div className="flex flex-wrap gap-2">
          {Object.entries(values).map(([key, value]) => {
            if (!value || value === '') return null;
            
            const filter = filters.find(f => f.key === key);
            if (!filter) return null;

            let displayValue = String(value);
            if (filter.type === 'select' && filter.options) {
              const option = filter.options.find(opt => String(opt.value) === String(value));
              displayValue = option?.label || displayValue;
            }

            return (
              <Badge
                key={key}
                variant="outline"
                className="text-xs flex items-center gap-1"
              >
                <span className="font-medium">{filter.label}:</span>
                <span>{displayValue}</span>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => onChange(key, '')}
                  className="h-auto p-0 ml-1 hover:bg-transparent"
                >
                  <X className="h-3 w-3" />
                  <span className="sr-only">Remove {filter.label} filter</span>
                </Button>
              </Badge>
            );
          })}
        </div>
      )}

      {/* Filter Controls */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {filters.map((filter) => (
          <div key={filter.key} className="space-y-2">
            <Label htmlFor={filter.key} className="text-sm font-medium text-gray-700">
              {filter.label}
            </Label>
            {renderFilterInput(filter)}
          </div>
        ))}
      </div>

      {/* Helper Text */}
      {filters.length === 0 && (
        <div className="text-center py-4">
          <p className="text-sm text-gray-500">No filters available</p>
        </div>
      )}
    </div>
  );
}