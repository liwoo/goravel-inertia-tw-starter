import { useState, useEffect, useCallback } from 'react';
import { FilterMetadata } from '@/types/filters';
import { toast } from 'sonner';

export function useFilterMetadata(resourceName: string) {
  const [metadata, setMetadata] = useState<FilterMetadata | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchMetadata = useCallback(async () => {
    if (!resourceName) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/${resourceName}/filters`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch filter metadata: ${response.statusText}`);
      }

      const result = await response.json();
      
      if (result.success && result.data) {
        setMetadata(result.data);
      } else {
        throw new Error(result.message || 'Failed to load filter metadata');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load filter metadata';
      setError(errorMessage);
      console.error('Filter metadata error:', err);
      
      // Only show toast for actual errors, not 404s (which mean no filters defined)
      if (!errorMessage.includes('404')) {
        toast.error('Failed to load filter definitions');
      }
    } finally {
      setLoading(false);
    }
  }, [resourceName]);

  useEffect(() => {
    fetchMetadata();
  }, [fetchMetadata]);

  const refetch = useCallback(() => {
    fetchMetadata();
  }, [fetchMetadata]);

  return {
    metadata,
    loading,
    error,
    refetch,
  };
}