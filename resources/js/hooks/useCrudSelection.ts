import { useState, useCallback, useMemo } from 'react';
import { UseCrudSelectionReturn } from '@/types/crud';

/**
 * Hook for managing bulk selection in CRUD tables
 */
export function useCrudSelection<T extends { id: number }>(
  data: T[],
  maxSelection?: number
): UseCrudSelectionReturn<T> {
  const [selectedIds, setSelectedIds] = useState<number[]>([]);

  const selectedItems = useMemo(() => {
    return data.filter(item => selectedIds.includes(item.id));
  }, [data, selectedIds]);

  const isAllSelected = useMemo(() => {
    return data.length > 0 && selectedIds.length === data.length;
  }, [data.length, selectedIds.length]);

  const isSomeSelected = useMemo(() => {
    return selectedIds.length > 0 && selectedIds.length < data.length;
  }, [selectedIds.length, data.length]);

  const toggleSelection = useCallback((id: number) => {
    setSelectedIds(prev => {
      const isSelected = prev.includes(id);
      
      if (isSelected) {
        return prev.filter(selectedId => selectedId !== id);
      } else {
        const newSelection = [...prev, id];
        
        // Respect max selection limit
        if (maxSelection && newSelection.length > maxSelection) {
          return newSelection.slice(-maxSelection);
        }
        
        return newSelection;
      }
    });
  }, [maxSelection]);

  const toggleAllSelection = useCallback(() => {
    if (isAllSelected) {
      setSelectedIds([]);
    } else {
      const allIds = data.map(item => item.id);
      
      if (maxSelection && allIds.length > maxSelection) {
        setSelectedIds(allIds.slice(0, maxSelection));
      } else {
        setSelectedIds(allIds);
      }
    }
  }, [data, isAllSelected, maxSelection]);

  const clearSelection = useCallback(() => {
    setSelectedIds([]);
  }, []);

  const setSelection = useCallback((ids: number[]) => {
    if (maxSelection && ids.length > maxSelection) {
      setSelectedIds(ids.slice(0, maxSelection));
    } else {
      setSelectedIds(ids);
    }
  }, [maxSelection]);

  const selectRange = useCallback((startIndex: number, endIndex: number) => {
    const start = Math.min(startIndex, endIndex);
    const end = Math.max(startIndex, endIndex);
    const rangeIds = data.slice(start, end + 1).map(item => item.id);
    
    setSelectedIds(prev => {
      const newSelection = Array.from(new Set([...prev, ...rangeIds]));
      
      if (maxSelection && newSelection.length > maxSelection) {
        return newSelection.slice(0, maxSelection);
      }
      
      return newSelection;
    });
  }, [data, maxSelection]);

  const canSelectMore = useMemo(() => {
    if (!maxSelection) return true;
    return selectedIds.length < maxSelection;
  }, [selectedIds.length, maxSelection]);

  const remainingSelections = useMemo(() => {
    if (!maxSelection) return null;
    return maxSelection - selectedIds.length;
  }, [selectedIds.length, maxSelection]);

  return {
    selectedIds,
    selectedItems,
    isAllSelected,
    isSomeSelected,
    toggleSelection,
    toggleAllSelection,
    clearSelection,
    setSelection,
    selectRange,
    canSelectMore,
    remainingSelections,
    selectionCount: selectedIds.length,
  };
}