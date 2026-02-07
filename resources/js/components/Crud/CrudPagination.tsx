"use client"

import * as React from 'react';
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { PaginationProps } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const PAGE_SIZE_KEY = 'crud-page-size';
const DEFAULT_PAGE_SIZE = 20;

export function CrudPagination({
  currentPage,
  lastPage,
  total,
  perPage,
  onPageChange,
  onPageSizeChange,
  allowedPageSizes = [10, 20, 50, 100],
  showInfo = true,
  className,
}: PaginationProps) {
  const { t } = useTranslation('common');
  const showingFrom = (currentPage - 1) * perPage + 1;
  const showingTo = Math.min(currentPage * perPage, total);

  if (lastPage <= 1) {
    return null;
  }

  return (
    <div className={cn("flex items-center justify-between px-2", className)}>
      <div className="flex-1 text-sm text-muted-foreground">
        {showInfo && (
          <span>
            {t('pagination.showing', {from: showingFrom, to: showingTo, total})}
          </span>
        )}
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">{t('pagination.rowsPerPage')}</p>
          <Select
            value={`${perPage}`}
            onValueChange={(value) => {
              const newPageSize = parseInt(value, 10);
              // Save to localStorage
              localStorage.setItem(PAGE_SIZE_KEY, value);
              // Call the callback if provided
              if (onPageSizeChange) {
                onPageSizeChange(newPageSize);
              }
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={perPage} />
            </SelectTrigger>
            <SelectContent side="top">
              {allowedPageSizes.map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          {t('pagination.pageOf', {current: currentPage, last: lastPage})}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => onPageChange(1)}
            disabled={currentPage === 1}
          >
            <span className="sr-only">{t('pagination.goToFirstPage')}</span>
            <ChevronsLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
          >
            <span className="sr-only">{t('pagination.goToPreviousPage')}</span>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage === lastPage}
          >
            <span className="sr-only">{t('pagination.goToNextPage')}</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => onPageChange(lastPage)}
            disabled={currentPage === lastPage}
          >
            <span className="sr-only">{t('pagination.goToLastPage')}</span>
            <ChevronsRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}