import React from 'react';
import { X } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  FilterCondition, 
  CompoundFilter, 
  getOperatorLabel,
  FilterMetadata 
} from '@/types/filters';

interface ActiveFilterBadgesProps {
  dynamicFilter?: FilterCondition | CompoundFilter | null;
  customFilters?: Record<string, any>;
  filterMetadata?: FilterMetadata | null;
  onRemoveDynamicFilter?: (index?: number) => void;
  onRemoveCustomFilter?: (key: string) => void;
  onClearAll?: () => void;
}

export function ActiveFilterBadges({
  dynamicFilter,
  customFilters = {},
  filterMetadata,
  onRemoveDynamicFilter,
  onRemoveCustomFilter,
  onClearAll,
}: ActiveFilterBadgesProps) {
  // Extract conditions from dynamic filter
  const getDynamicConditions = (): FilterCondition[] => {
    if (!dynamicFilter) return [];
    
    if ('logic' in dynamicFilter) {
      // It's a compound filter
      return dynamicFilter.conditions.filter(
        (c): c is FilterCondition => !('logic' in c)
      ) as FilterCondition[];
    }
    
    // It's a single condition
    return [dynamicFilter];
  };

  const dynamicConditions = getDynamicConditions();
  
  // Get custom filter entries (excluding special keys)
  const customFilterEntries = Object.entries(customFilters).filter(
    ([key, value]) => value !== undefined && value !== '' && value !== '__all__'
  );

  const totalFilters = dynamicConditions.length + customFilterEntries.length;

  if (totalFilters === 0) return null;

  // Get field label from metadata
  const getFieldLabel = (field: string) => {
    if (!filterMetadata) return field;
    const definition = filterMetadata.filters.find(f => f.field === field);
    return definition?.label || field;
  };

  // Format filter value for display
  const formatValue = (value: any): string => {
    if (value === null || value === undefined) return 'empty';
    if (Array.isArray(value)) {
      if (value.length === 2 && typeof value[0] !== 'object') {
        // Range values
        return `${value[0]} to ${value[1]}`;
      }
      return value.join(', ');
    }
    if (typeof value === 'boolean') return value ? 'Yes' : 'No';
    if (typeof value === 'object') return JSON.stringify(value);
    return String(value);
  };

  return (
    <div className="flex flex-wrap items-center gap-2">
      <span className="text-sm text-muted-foreground">Active filters:</span>
      
      {/* Dynamic filter conditions */}
      {dynamicConditions.map((condition, index) => (
        <Badge 
          key={`dynamic-${index}`} 
          variant="secondary" 
          className="pl-2 pr-1 py-1 flex items-center gap-1"
        >
          <span className="text-xs">
            {getFieldLabel(condition.field)} {getOperatorLabel(condition.operator)} {formatValue(condition.value)}
          </span>
          {onRemoveDynamicFilter && (
            <Button
              variant="ghost"
              size="sm"
              className="h-4 w-4 p-0 ml-1 hover:bg-transparent"
              onClick={() => onRemoveDynamicFilter(index)}
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </Badge>
      ))}
      
      {/* Custom filters */}
      {customFilterEntries.map(([key, value]) => (
        <Badge 
          key={`custom-${key}`} 
          variant="secondary" 
          className="pl-2 pr-1 py-1 flex items-center gap-1"
        >
          <span className="text-xs">
            {key}: {formatValue(value)}
          </span>
          {onRemoveCustomFilter && (
            <Button
              variant="ghost"
              size="sm"
              className="h-4 w-4 p-0 ml-1 hover:bg-transparent"
              onClick={() => onRemoveCustomFilter(key)}
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </Badge>
      ))}
      
      {/* Clear all button */}
      {totalFilters > 1 && onClearAll && (
        <Button
          variant="ghost"
          size="sm"
          className="h-6 px-2 text-xs"
          onClick={onClearAll}
        >
          Clear all
        </Button>
      )}
    </div>
  );
}