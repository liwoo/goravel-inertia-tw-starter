import React, { useCallback, useMemo } from 'react';
import { X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Calendar } from '@/components/ui/calendar';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import { CalendarIcon } from 'lucide-react';
import {
  FilterCondition,
  FilterDefinition,
  FilterMetadata,
  FilterOperator,
  getOperatorLabel,
  operatorRequiresValue,
  operatorRequiresRange,
  operatorRequiresArray,
  operatorRequiresNumber,
} from '@/types/filters';

interface FilterConditionBuilderProps {
  condition: FilterCondition;
  metadata: FilterMetadata;
  onChange: (condition: FilterCondition) => void;
  onRemove: () => void;
}

export function FilterConditionBuilder({
  condition,
  metadata,
  onChange,
  onRemove,
}: FilterConditionBuilderProps) {
  const fieldDefinition = useMemo(() => {
    return metadata.filters.find(f => f.field === condition.field);
  }, [metadata.filters, condition.field]);

  const handleFieldChange = useCallback((field: string) => {
    const newFieldDef = metadata.filters.find(f => f.field === field);
    if (!newFieldDef) return;

    onChange({
      ...condition,
      field,
      type: newFieldDef.type,
      operator: newFieldDef.operators[0] || 'equals',
      value: undefined,
    });
  }, [metadata.filters, condition, onChange]);

  const handleOperatorChange = useCallback((operator: FilterOperator) => {
    onChange({
      ...condition,
      operator,
      // Clear value if operator doesn't require it
      value: operatorRequiresValue(operator) ? condition.value : undefined,
    });
  }, [condition, onChange]);

  const handleValueChange = useCallback((value: any) => {
    onChange({
      ...condition,
      value,
    });
  }, [condition, onChange]);

  const renderValueInput = () => {
    if (!operatorRequiresValue(condition.operator)) {
      return null;
    }

    if (!fieldDefinition) {
      return <Input placeholder="Value" disabled />;
    }

    // Handle range operators (between, not_between)
    if (operatorRequiresRange(condition.operator)) {
      const values = Array.isArray(condition.value) ? condition.value : [undefined, undefined];
      
      if (fieldDefinition.type === 'number') {
        return (
          <div className="flex gap-2 items-center">
            <Input
              type="number"
              placeholder="From"
              value={values[0] || ''}
              onChange={(e) => handleValueChange([e.target.value, values[1]])}
              className="flex-1"
            />
            <span className="text-sm text-muted-foreground">to</span>
            <Input
              type="number"
              placeholder="To"
              value={values[1] || ''}
              onChange={(e) => handleValueChange([values[0], e.target.value])}
              className="flex-1"
            />
          </div>
        );
      }

      if (fieldDefinition.type === 'date') {
        return (
          <div className="flex gap-2 items-center">
            <DatePicker
              value={values[0]}
              onChange={(date) => handleValueChange([date, values[1]])}
              placeholder="From date"
            />
            <span className="text-sm text-muted-foreground">to</span>
            <DatePicker
              value={values[1]}
              onChange={(date) => handleValueChange([values[0], date])}
              placeholder="To date"
            />
          </div>
        );
      }
    }

    // Handle array operators (in, not_in, contains_any, contains_all)
    if (operatorRequiresArray(condition.operator)) {
      if (fieldDefinition.type === 'enum' && fieldDefinition.enum_values) {
        return (
          <MultiSelect
            options={fieldDefinition.enum_values}
            value={condition.value || []}
            onChange={handleValueChange}
            placeholder="Select values"
          />
        );
      }
      
      // For other types, use comma-separated input
      return (
        <Input
          placeholder="Enter values separated by commas"
          value={Array.isArray(condition.value) ? condition.value.join(', ') : ''}
          onChange={(e) => handleValueChange(e.target.value.split(',').map(v => v.trim()))}
        />
      );
    }

    // Handle number input for last_n_days, next_n_days
    if (operatorRequiresNumber(condition.operator)) {
      return (
        <Input
          type="number"
          placeholder="Number of days"
          value={condition.value || ''}
          onChange={(e) => handleValueChange(e.target.value)}
        />
      );
    }

    // Handle specific field types
    switch (fieldDefinition.type) {
      case 'string':
        return (
          <Input
            type="text"
            placeholder="Enter value"
            value={condition.value || ''}
            onChange={(e) => handleValueChange(e.target.value)}
          />
        );

      case 'number':
        return (
          <Input
            type="number"
            placeholder="Enter number"
            value={condition.value || ''}
            onChange={(e) => handleValueChange(e.target.value)}
          />
        );

      case 'date':
      case 'datetime':
        return (
          <DatePicker
            value={condition.value}
            onChange={handleValueChange}
            placeholder="Select date"
          />
        );

      case 'boolean':
        return (
          <div className="flex items-center space-x-2">
            <Checkbox
              checked={condition.value === true}
              onCheckedChange={(checked) => handleValueChange(checked)}
            />
            <span className="text-sm">True</span>
          </div>
        );

      case 'enum':
        if (fieldDefinition.enum_values) {
          return (
            <Select
              value={condition.value !== undefined && condition.value !== null ? String(condition.value) : undefined}
              onValueChange={handleValueChange}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select value" />
              </SelectTrigger>
              <SelectContent>
                {fieldDefinition.enum_values.map((value) => {
                  // Skip empty enum values
                  if (!value || value === '') return null;
                  return (
                    <SelectItem key={value} value={value}>
                      {value}
                    </SelectItem>
                  );
                })}
              </SelectContent>
            </Select>
          );
        }
        break;
    }

    return (
      <Input
        placeholder="Enter value"
        value={condition.value || ''}
        onChange={(e) => handleValueChange(e.target.value)}
      />
    );
  };

  return (
    <div className="flex flex-col sm:flex-row gap-2 items-start p-3 border rounded-lg bg-muted/30">
      {/* Field Selector */}
      <Select 
        value={condition.field && condition.field !== '' ? condition.field : undefined} 
        onValueChange={handleFieldChange}
      >
        <SelectTrigger className="w-full sm:w-[180px]">
          <SelectValue placeholder="Select field" />
        </SelectTrigger>
        <SelectContent>
          {metadata.filters.map((filter) => {
            // Skip filters with empty field names
            if (!filter.field || filter.field === '') return null;
            return (
              <SelectItem key={filter.field} value={filter.field}>
                {filter.label}
              </SelectItem>
            );
          })}
        </SelectContent>
      </Select>

      {/* Operator Selector */}
      {fieldDefinition && (
        <Select 
          value={condition.operator && condition.operator !== '' ? condition.operator : undefined} 
          onValueChange={handleOperatorChange}
        >
          <SelectTrigger className="w-full sm:w-[200px]">
            <SelectValue placeholder="Select operator" />
          </SelectTrigger>
          <SelectContent>
            {fieldDefinition.operators.map((operator) => {
              // Skip empty operators
              if (!operator || operator === '') return null;
              return (
                <SelectItem key={operator} value={operator}>
                  {getOperatorLabel(operator)}
                </SelectItem>
              );
            })}
          </SelectContent>
        </Select>
      )}

      {/* Value Input */}
      <div className="flex-1">
        {renderValueInput()}
      </div>

      {/* Remove Button */}
      <Button
        variant="ghost"
        size="sm"
        onClick={onRemove}
        className="h-9 w-9 p-0"
      >
        <X className="h-4 w-4" />
      </Button>
    </div>
  );
}

// Date Picker Component
function DatePicker({
  value,
  onChange,
  placeholder,
}: {
  value: string | undefined;
  onChange: (date: string | undefined) => void;
  placeholder: string;
}) {
  const date = value ? new Date(value) : undefined;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn(
            'justify-start text-left font-normal',
            !value && 'text-muted-foreground'
          )}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {value ? format(date!, 'PPP') : placeholder}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="single"
          selected={date}
          onSelect={(newDate) => {
            onChange(newDate ? format(newDate, 'yyyy-MM-dd') : undefined);
          }}
          initialFocus
        />
      </PopoverContent>
    </Popover>
  );
}

// Multi-Select Component
function MultiSelect({
  options,
  value,
  onChange,
  placeholder,
}: {
  options: string[];
  value: string[];
  onChange: (value: string[]) => void;
  placeholder: string;
}) {
  const [open, setOpen] = React.useState(false);

  const handleToggle = (option: string) => {
    const newValue = value.includes(option)
      ? value.filter(v => v !== option)
      : [...value, option];
    onChange(newValue);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className="justify-between"
        >
          {value.length > 0 ? (
            <span>{value.length} selected</span>
          ) : (
            <span className="text-muted-foreground">{placeholder}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0">
        <div className="max-h-[300px] overflow-auto">
          {options.map((option) => (
            <div
              key={option}
              className="flex items-center space-x-2 p-2 hover:bg-accent cursor-pointer"
              onClick={() => handleToggle(option)}
            >
              <Checkbox checked={value.includes(option)} />
              <span className="text-sm">{option}</span>
            </div>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
}