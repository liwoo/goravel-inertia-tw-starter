import React, { useState, useEffect, useCallback } from 'react';
import { Plus, X, Filter, Trash2, ChevronDown, ChevronUp } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import {
  FilterCondition,
  CompoundFilter,
  FilterMetadata,
  FilterDefinition,
  LogicOperator,
} from '@/types/filters';
import { FilterConditionBuilder } from './FilterConditionBuilder';
import { FilterGroupBuilder } from './FilterGroupBuilder';
import { toast } from 'sonner';

interface DynamicFilterBuilderProps {
  metadata: FilterMetadata | null;
  onApply: (filter: FilterCondition | CompoundFilter | null) => void;
  onClose?: () => void;
  className?: string;
  loading?: boolean;
  initialFilter?: FilterCondition | CompoundFilter | null;
}

export function DynamicFilterBuilder({
  metadata,
  onApply,
  onClose,
  className,
  loading = false,
  initialFilter,
}: DynamicFilterBuilderProps) {
  const [isOpen, setIsOpen] = useState(true);
  
  // Initialize root logic and filters from initialFilter prop
  const initialRootLogic = initialFilter && 'logic' in initialFilter ? initialFilter.logic : 'AND';
  const initialFilters = (() => {
    if (!initialFilter) return [];
    
    // If it's a compound filter with multiple conditions
    if ('logic' in initialFilter) {
      return initialFilter.conditions;
    }
    
    // If it's a single condition
    return [initialFilter];
  })();
  
  const [rootLogic, setRootLogic] = useState<LogicOperator>(initialRootLogic);
  const [filters, setFilters] = useState<(FilterCondition | CompoundFilter)[]>(initialFilters);

  // Update filters when initialFilter changes
  useEffect(() => {
    if (initialFilter) {
      if ('logic' in initialFilter) {
        setRootLogic(initialFilter.logic);
        setFilters(initialFilter.conditions);
      } else {
        setFilters([initialFilter]);
      }
    } else {
      // Don't clear filters if initialFilter becomes null - user might be editing
      // Only clear if explicitly requested via handleClear
    }
  }, [initialFilter]);
  
  // Load metadata on mount
  useEffect(() => {
    if (!metadata && !loading) {
      console.log('No filter metadata available');
    }
  }, [metadata, loading]);

  const handleAddCondition = useCallback(() => {
    if (!metadata) return;
    
    const firstField = metadata.filters[0];
    if (!firstField) return;

    const newCondition: FilterCondition = {
      field: firstField.field,
      operator: firstField.operators[0] || 'equals',
      value: undefined,
      type: firstField.type,
    };
    
    setFilters([...filters, newCondition]);
  }, [filters, metadata]);

  const handleAddGroup = useCallback(() => {
    const newGroup: CompoundFilter = {
      logic: 'AND',
      conditions: [],
    };
    
    setFilters([...filters, newGroup]);
  }, [filters]);

  const handleUpdateCondition = useCallback((index: number, condition: FilterCondition | CompoundFilter) => {
    const newFilters = [...filters];
    newFilters[index] = condition;
    setFilters(newFilters);
  }, [filters]);

  const handleRemoveCondition = useCallback((index: number) => {
    setFilters(filters.filter((_, i) => i !== index));
  }, [filters]);

  const handleClear = useCallback(() => {
    setFilters([]);
    onApply(null);
    toast.success('Filters cleared');
  }, [onApply]);

  const handleApply = useCallback(() => {
    if (filters.length === 0) {
      onApply(null);
      return;
    }

    if (filters.length === 1) {
      onApply(filters[0]);
    } else {
      const compoundFilter: CompoundFilter = {
        logic: rootLogic,
        conditions: filters,
      };
      onApply(compoundFilter);
    }
    
    toast.success('Filters applied');
  }, [filters, rootLogic, onApply]);

  const getActiveFilterCount = useCallback(() => {
    return filters.length;
  }, [filters]);

  if (loading) {
    return (
      <Card className={cn('w-full', className)}>
        <CardContent className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </CardContent>
      </Card>
    );
  }

  if (!metadata) {
    return (
      <Card className={cn('w-full', className)}>
        <CardContent className="text-center py-8">
          <p className="text-muted-foreground">No filter definitions available</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn('w-full', className)}>
      <Collapsible open={isOpen} onOpenChange={setIsOpen}>
        <CardHeader className="p-4 sm:p-6 pb-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Filter className="h-5 w-5 text-muted-foreground" />
              <CardTitle className="text-lg">Advanced Filters</CardTitle>
              {getActiveFilterCount() > 0 && (
                <Badge variant="secondary">
                  {getActiveFilterCount()} active
                </Badge>
              )}
            </div>
            <div className="flex items-center gap-2">
              {onClose && (
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0 md:hidden" onClick={onClose}>
                  <X className="h-4 w-4" />
                </Button>
              )}
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  {isOpen ? (
                    <ChevronUp className="h-4 w-4" />
                  ) : (
                    <ChevronDown className="h-4 w-4" />
                  )}
                </Button>
              </CollapsibleTrigger>
            </div>
          </div>
        </CardHeader>
        
        <CollapsibleContent>
          <CardContent className="p-4 sm:p-6 pt-0 space-y-4">
            {/* Root Logic Selector (when multiple conditions) */}
            {filters.length > 1 && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground">Match</span>
                <div className="flex gap-1">
                  <Button
                    variant={rootLogic === 'AND' ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setRootLogic('AND')}
                    className="h-7"
                  >
                    All (AND)
                  </Button>
                  <Button
                    variant={rootLogic === 'OR' ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setRootLogic('OR')}
                    className="h-7"
                  >
                    Any (OR)
                  </Button>
                </div>
                <span className="text-sm text-muted-foreground">of the following:</span>
              </div>
            )}

            {/* Filter Conditions */}
            <div className="space-y-3">
              {filters.map((filter, index) => (
                <div key={index} className="relative">
                  {index > 0 && (
                    <div className="absolute -top-3 left-8 text-xs text-muted-foreground bg-background px-2">
                      {rootLogic}
                    </div>
                  )}
                  
                  {'logic' in filter ? (
                    <FilterGroupBuilder
                      group={filter}
                      metadata={metadata}
                      onChange={(group) => handleUpdateCondition(index, group)}
                      onRemove={() => handleRemoveCondition(index)}
                    />
                  ) : (
                    <FilterConditionBuilder
                      condition={filter}
                      metadata={metadata}
                      onChange={(condition) => handleUpdateCondition(index, condition)}
                      onRemove={() => handleRemoveCondition(index)}
                    />
                  )}
                </div>
              ))}
              
              {filters.length === 0 && (
                <div className="text-center py-8 border-2 border-dashed rounded-lg">
                  <p className="text-sm text-muted-foreground mb-4">
                    No filters applied. Add a condition to get started.
                  </p>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleAddCondition}
                    className="mr-2"
                  >
                    <Plus className="h-4 w-4 mr-1" />
                    Add Condition
                  </Button>
                </div>
              )}
            </div>

            {/* Action Buttons */}
            {filters.length > 0 && (
              <>
                <Separator />
                <div className="flex flex-col sm:flex-row gap-3 sm:items-center sm:justify-between">
                  <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleAddCondition}
                      className="w-full sm:w-auto"
                    >
                      <Plus className="h-4 w-4 mr-1" />
                      Add Condition
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleAddGroup}
                      className="w-full sm:w-auto"
                    >
                      <Plus className="h-4 w-4 mr-1" />
                      Add Group
                    </Button>
                  </div>
                  
                  <div className="flex gap-2 w-full sm:w-auto">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleClear}
                      className="flex-1 sm:flex-none"
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Clear All
                    </Button>
                    <Button
                      size="sm"
                      onClick={handleApply}
                      className="flex-1 sm:flex-none"
                    >
                      Apply Filters
                    </Button>
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </CollapsibleContent>
      </Collapsible>
    </Card>
  );
}