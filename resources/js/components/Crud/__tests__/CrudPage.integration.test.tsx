/**
 * Integration tests for CrudPage with minimal mocking
 * Focus on filter functionality behavior
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// Mock only what's absolutely necessary
vi.mock('@inertiajs/react', () => {
  const mockRouter = {
    get: vi.fn(),
    reload: vi.fn(),
    delete: vi.fn(),
  };
  
  return {
    router: mockRouter,
    Head: ({ title }: { title: string }) => React.createElement('title', {}, title),
    usePage: () => ({
      props: {
        auth: { user: { permissions: [] } }
      }
    }),
  };
});

vi.mock('@/hooks/useFilterMetadata', () => ({
  useFilterMetadata: () => ({ metadata: null, loading: false }),
}));

vi.mock('@/contexts/PermissionsContext', () => ({
  usePermissions: () => ({
    canPerformAction: () => true,
  }),
}));

// Import after mocks
import { CrudPage } from '../CrudPage';
import { router } from '@inertiajs/react';

const mockRouter = router as any;

describe('CrudPage Integration - Filter Functionality', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
    // Reset localStorage mock
    (global.localStorage.getItem as any).mockReturnValue(null);
  });

  const createTestData = () => ({
    data: [
      { id: 1, name: 'Item 1', status: 'AVAILABLE' },
      { id: 2, name: 'Item 2', status: 'BORROWED' },
      { id: 3, name: 'Item 3', status: 'AVAILABLE' },
    ],
    currentPage: 1,
    perPage: 10,
    total: 3,
    lastPage: 1,
    from: 1,
    to: 3,
  });

  it('should handle simple filter tab interactions correctly', async () => {
    const simpleFilters = [
      {
        key: 'status-available',
        label: 'Available',
        value: 'AVAILABLE',
        filterParams: { status: 'AVAILABLE' },
        badge: 2,
      },
      {
        key: 'status-borrowed',
        label: 'Borrowed',
        value: 'BORROWED',
        filterParams: { status: 'BORROWED' },
        badge: 1,
      },
    ];

    render(
      <CrudPage
        data={createTestData()}
        filters={{ 
          filters: { status: 'AVAILABLE' },
          search: '',
          sort: undefined,
          direction: undefined 
        }}
        title="Test Items"
        resourceName="items"
        columns={[
          { key: 'name', label: 'Name' },
          { key: 'status', label: 'Status' },
        ]}
        simpleFilters={simpleFilters}
      />
    );

    // Verify Available tab is active initially
    const tabs = screen.getByRole('tablist');
    const availableTab = within(tabs).getByRole('tab', { name: /Available/i });
    expect(availableTab).toHaveAttribute('data-state', 'active');

    // Verify filter badge is shown
    await screen.findByText(/Active filters:/i);
    const filterBadge = screen.getByText('status: AVAILABLE');
    expect(filterBadge).toBeInTheDocument();

    // Click X on filter badge
    const badge = filterBadge.closest('.flex') as HTMLElement;
    const removeButton = within(badge!).getByRole('button', { name: /×/i });
    await user.click(removeButton);

    // Verify router was called without status
    expect(mockRouter.get).toHaveBeenCalledWith(
      '/admin/items',
      expect.not.objectContaining({ status: expect.anything() }),
      expect.any(Object)
    );

    const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
    expect(lastCall[1]).not.toHaveProperty('status');
  });

  it('should clear filters when clicking All tab', async () => {
    const simpleFilters = [
      {
        key: 'status-available',
        label: 'Available',
        value: 'AVAILABLE',
        filterParams: { status: 'AVAILABLE' },
      },
    ];

    render(
      <CrudPage
        data={createTestData()}
        filters={{ 
          filters: { status: 'AVAILABLE' },
          search: '',
          sort: undefined,
          direction: undefined 
        }}
        title="Test Items"
        resourceName="items"
        columns={[{ key: 'name', label: 'Name' }]}
        simpleFilters={simpleFilters}
      />
    );

    // Click All tab
    const tabs = screen.getByRole('tablist');
    const allTab = within(tabs).getByRole('tab', { name: /All/i });
    await user.click(allTab);

    // Verify router was called without status filter
    expect(mockRouter.get).toHaveBeenCalledWith(
      '/admin/items',
      expect.objectContaining({
        page: 1,
        pageSize: 10,
      }),
      expect.any(Object)
    );

    const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
    expect(lastCall[1]).not.toHaveProperty('status');
  });

  it('should handle search while preserving filters', async () => {
    vi.useFakeTimers();

    render(
      <CrudPage
        data={createTestData()}
        filters={{ 
          filters: { status: 'AVAILABLE' },
          search: '',
          sort: undefined,
          direction: undefined 
        }}
        title="Test Items"
        resourceName="items"
        columns={[{ key: 'name', label: 'Name' }]}
      />
    );

    // Type in search
    const searchInput = screen.getByPlaceholderText(/Search items/i);
    await user.type(searchInput, 'test');

    // Fast-forward debounce
    vi.advanceTimersByTime(300);

    // Verify filters are preserved during search
    expect(mockRouter.get).toHaveBeenCalledWith(
      '/admin/items',
      expect.objectContaining({
        search: 'test',
        status: 'AVAILABLE', // Should preserve existing filter
        page: 1,
      }),
      expect.any(Object)
    );

    vi.useRealTimers();
  });

  it('should handle Clear all button correctly', async () => {
    render(
      <CrudPage
        data={createTestData()}
        filters={{ 
          filters: {
            status: 'AVAILABLE',
            category: 'books',
          },
          search: 'test',
          sort: undefined,
          direction: undefined 
        }}
        title="Test Items"
        resourceName="items"
        columns={[{ key: 'name', label: 'Name' }]}
        customFilters={[
          { key: 'category', label: 'Category', type: 'text' }
        ]}
      />
    );

    // Find and click Clear all
    const clearAllButton = screen.getByRole('button', { name: /Clear all/i });
    await user.click(clearAllButton);

    // Verify only search is preserved
    expect(mockRouter.get).toHaveBeenCalledWith(
      '/admin/items',
      expect.objectContaining({
        page: 1,
        pageSize: 10,
        search: 'test', // Search should be preserved
      }),
      expect.any(Object)
    );

    const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
    expect(lastCall[1]).not.toHaveProperty('status');
    expect(lastCall[1]).not.toHaveProperty('category');
    expect(lastCall[1]).toHaveProperty('search', 'test');
  });
});