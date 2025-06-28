import { useState, useCallback, useMemo, useEffect } from 'react';
import { router } from '@inertiajs/react';
import { UseCrudFiltersReturn } from '@/types/crud';

/**
 * Hook for managing CRUD filters with Inertia.js integration
 */
export function useCrudFilters(
  resourceName: string, 
  initialFilters: Record<string, any> = {},
  options: {
    preserveState?: boolean;
    preserveScroll?: boolean;
    replace?: boolean;
  } = {}
): UseCrudFiltersReturn {
  const [filters, setFilters] = useState(initialFilters);
  const [pendingFilters, setPendingFilters] = useState(initialFilters);

  const updateFilter = useCallback((key: string, value: any) => {
    const newFilters = { ...pendingFilters };
    
    if (value === '' || value === null || value === undefined) {
      delete newFilters[key];
    } else {
      newFilters[key] = value;
    }
    
    setPendingFilters(newFilters);
  }, [pendingFilters]);

  const applyFilters = useCallback(() => {
    setFilters(pendingFilters);
    
    const queryParams: Record<string, any> = {
      page: 1, // Reset to first page when filtering
    };

    if (Object.keys(pendingFilters).length > 0) {
      queryParams.filters = pendingFilters;
    }

    router.get(`/${resourceName}`, queryParams, {
      preserveState: options.preserveState ?? true,
      preserveScroll: options.preserveScroll ?? true,
      replace: options.replace ?? false,
    });
  }, [resourceName, pendingFilters, options]);

  const clearFilters = useCallback(() => {
    const emptyFilters = {};
    setFilters(emptyFilters);
    setPendingFilters(emptyFilters);
    
    router.get(`/${resourceName}`, { page: 1 }, {
      preserveState: options.preserveState ?? true,
      preserveScroll: options.preserveScroll ?? true,
      replace: options.replace ?? false,
    });
  }, [resourceName, options]);

  const hasActiveFilters = useMemo(() => {
    return Object.keys(filters).length > 0;
  }, [filters]);

  const hasPendingChanges = useMemo(() => {
    return JSON.stringify(filters) !== JSON.stringify(pendingFilters);
  }, [filters, pendingFilters]);

  return {
    filters: pendingFilters,
    updateFilter,
    clearFilters,
    hasActiveFilters,
    applyFilters,
    hasPendingChanges,
    appliedFilters: filters,
  };
}

/**
 * Hook for managing search with debouncing
 */
export function useCrudSearch(
  resourceName: string,
  initialSearch: string = '',
  debounceMs: number = 300,
  options: {
    preserveState?: boolean;
    preserveScroll?: boolean;
    replace?: boolean;
  } = {}
) {
  const [searchTerm, setSearchTerm] = useState(initialSearch);
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState(initialSearch);

  // Debounce the search term
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearchTerm(searchTerm);
    }, debounceMs);

    return () => clearTimeout(timer);
  }, [searchTerm, debounceMs]);

  // Trigger search when debounced term changes
  useEffect(() => {
    if (debouncedSearchTerm !== initialSearch) {
      const queryParams: Record<string, any> = {
        page: 1, // Reset to first page when searching
      };

      if (debouncedSearchTerm.trim()) {
        queryParams.search = debouncedSearchTerm;
      }

      router.get(`/${resourceName}`, queryParams, {
        preserveState: options.preserveState ?? true,
        preserveScroll: options.preserveScroll ?? true,
        replace: options.replace ?? false,
      });
    }
  }, [debouncedSearchTerm, resourceName, initialSearch, options]);

  const handleSearch = useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
  }, []);

  return {
    searchTerm,
    handleSearch,
    clearSearch,
    isSearching: searchTerm !== debouncedSearchTerm,
  };
}