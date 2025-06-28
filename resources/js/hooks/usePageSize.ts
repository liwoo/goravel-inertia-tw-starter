import { useState, useEffect } from 'react';

const PAGE_SIZE_KEY = 'crud-page-size';

interface PaginationConfig {
  defaultPageSize?: number;
  maxPageSize?: number;
  allowedSizes?: number[];
}

export function usePageSize(config?: PaginationConfig) {
  const defaultPageSize = config?.defaultPageSize || 20;
  const allowedSizes = config?.allowedSizes || [10, 20, 50, 100];
  
  const [pageSize, setPageSize] = useState<number>(() => {
    // Try to get saved page size from localStorage
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(PAGE_SIZE_KEY);
      if (saved) {
        const parsed = parseInt(saved, 10);
        if (!isNaN(parsed) && allowedSizes.includes(parsed)) {
          return parsed;
        }
      }
    }
    return defaultPageSize;
  });

  // Save to localStorage whenever page size changes
  useEffect(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem(PAGE_SIZE_KEY, pageSize.toString());
    }
  }, [pageSize]);

  return { pageSize, setPageSize, allowedSizes, defaultPageSize };
}