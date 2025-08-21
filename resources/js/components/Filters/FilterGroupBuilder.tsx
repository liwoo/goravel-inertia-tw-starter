import React, { useCallback } from 'react';
import { Plus, X, Folder } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  CompoundFilter,
  FilterCondition,
  FilterMetadata,
  LogicOperator,
} from '@/types/filters';
import { FilterConditionBuilder } from './FilterConditionBuilder';

interface FilterGroupBuilderProps {
  group: CompoundFilter;
  metadata: FilterMetadata;
  onChange: (group: CompoundFilter) => void;
  onRemove: () => void;
  depth?: number;
}

export function FilterGroupBuilder({
  group,
  metadata,
  onChange,
  onRemove,
  depth = 0,
}: FilterGroupBuilderProps) {
  const handleLogicChange = useCallback((logic: LogicOperator) => {
    onChange({
      ...group,
      logic,
    });
  }, [group, onChange]);

  const handleAddCondition = useCallback(() => {
    const firstField = metadata.filters[0];
    if (!firstField) return;

    const newCondition: FilterCondition = {
      field: firstField.field,
      operator: firstField.operators[0] || 'equals',
      value: undefined,
      type: firstField.type,
    };

    onChange({
      ...group,
      conditions: [...group.conditions, newCondition],
    });
  }, [group, metadata, onChange]);

  const handleAddGroup = useCallback(() => {
    if (depth >= 2) {
      // Limit nesting depth
      return;
    }

    const newGroup: CompoundFilter = {
      logic: 'AND',
      conditions: [],
    };

    onChange({
      ...group,
      conditions: [...group.conditions, newGroup],
    });
  }, [group, depth, onChange]);

  const handleUpdateCondition = useCallback((index: number, condition: FilterCondition | CompoundFilter) => {
    const newConditions = [...group.conditions];
    newConditions[index] = condition;
    onChange({
      ...group,
      conditions: newConditions,
    });
  }, [group, onChange]);

  const handleRemoveCondition = useCallback((index: number) => {
    onChange({
      ...group,
      conditions: group.conditions.filter((_, i) => i !== index),
    });
  }, [group, onChange]);

  return (
    <Card className="p-4 border-l-4 border-l-primary/50">
      <div className="space-y-3">
        {/* Group Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Folder className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">Group</span>
            <div className="flex gap-1">
              <Button
                variant={group.logic === 'AND' ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleLogicChange('AND')}
                className="h-6 text-xs"
              >
                AND
              </Button>
              <Button
                variant={group.logic === 'OR' ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleLogicChange('OR')}
                className="h-6 text-xs"
              >
                OR
              </Button>
            </div>
            {group.conditions.length > 0 && (
              <Badge variant="secondary" className="text-xs">
                {group.conditions.length} condition{group.conditions.length > 1 ? 's' : ''}
              </Badge>
            )}
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={onRemove}
            className="h-8 w-8 p-0"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Group Conditions */}
        <div className="space-y-2 ml-4">
          {group.conditions.map((condition, index) => (
            <div key={index} className="relative">
              {index > 0 && (
                <div className="absolute -top-2 left-4 text-xs text-muted-foreground bg-background px-2">
                  {group.logic}
                </div>
              )}
              
              {'logic' in condition ? (
                <FilterGroupBuilder
                  group={condition}
                  metadata={metadata}
                  onChange={(updatedGroup) => handleUpdateCondition(index, updatedGroup)}
                  onRemove={() => handleRemoveCondition(index)}
                  depth={depth + 1}
                />
              ) : (
                <FilterConditionBuilder
                  condition={condition}
                  metadata={metadata}
                  onChange={(updatedCondition) => handleUpdateCondition(index, updatedCondition)}
                  onRemove={() => handleRemoveCondition(index)}
                />
              )}
            </div>
          ))}

          {group.conditions.length === 0 && (
            <div className="text-center py-4 border-2 border-dashed rounded-lg">
              <p className="text-xs text-muted-foreground mb-2">
                Empty group. Add conditions to this group.
              </p>
            </div>
          )}
        </div>

        {/* Add Buttons */}
        <div className="flex gap-2 ml-4">
          <Button
            variant="outline"
            size="sm"
            onClick={handleAddCondition}
            className="text-xs"
          >
            <Plus className="h-3 w-3 mr-1" />
            Add Condition
          </Button>
          {depth < 2 && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleAddGroup}
              className="text-xs"
            >
              <Plus className="h-3 w-3 mr-1" />
              Add Nested Group
            </Button>
          )}
        </div>
      </div>
    </Card>
  );
}