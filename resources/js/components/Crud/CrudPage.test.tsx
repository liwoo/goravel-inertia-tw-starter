import { describe, it, expect, vi, beforeEach } from 'vitest';
import React from 'react';
import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { 
  render, 
  createMockCrudData, 
  createMockFilters,
  createMockSimpleFilter 
} from '@/test/test-utils';

// Mock modules first
vi.mock('@inertiajs/react', () => {
  const mockRouter = {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    reload: vi.fn(),
    replace: vi.fn(),
    visit: vi.fn(),
    on: vi.fn(),
    off: vi.fn(),
  };
  
  return {
    default: {},
    router: mockRouter,
    Head: ({ title }: { title: string }) => React.createElement('title', {}, title),
    usePage: () => ({
      props: {
        auth: {
          user: {
            id: 1,
            name: 'Test User',
            permissions: ['books.create', 'books.update', 'books.delete', 'books.read'],
          },
        },
      },
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
  PermissionsProvider: ({ children }: { children: React.ReactNode }) => React.createElement('div', {}, children),
}));

// Import CrudPage after mocks are set up
import { CrudPage } from './CrudPage';
import { router } from '@inertiajs/react';

// Cast router to get the mocked version
const mockRouter = router as any;

// Sample book data for testing
const mockBooks = [
  { id: 1, title: 'Book 1', author: 'Author 1', status: 'AVAILABLE' },
  { id: 2, title: 'Book 2', author: 'Author 2', status: 'BORROWED' },
  { id: 3, title: 'Book 3', author: 'Author 3', status: 'AVAILABLE' },
];

const mockColumns = [
  { key: 'title', label: 'Title' },
  { key: 'author', label: 'Author' },
  { key: 'status', label: 'Status' },
];

describe('CrudPage', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Filter Badge Functionality', () => {
    it('should remove filter and reset simple filter tab when X is clicked on filter badge', async () => {
      const simpleFilters = [
        createMockSimpleFilter('status-available', 'Available', 'AVAILABLE', { status: 'AVAILABLE' }),
        createMockSimpleFilter('status-borrowed', 'Borrowed', 'BORROWED', { status: 'BORROWED' }),
      ];

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters({ status: 'AVAILABLE' })}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          simpleFilters={simpleFilters}
        />
      );

      // Check that Available tab is selected
      const availableTab = screen.getByRole('tab', { name: /Available/i });
      expect(availableTab).toHaveAttribute('data-state', 'active');

      // Check that filter badge is shown
      const filterBadge = screen.getByText('status: AVAILABLE');
      expect(filterBadge).toBeInTheDocument();

      // Find and click the X button on the filter badge
      const filterBadgeContainer = filterBadge.closest('.flex') as HTMLElement;
      const removeButton = within(filterBadgeContainer!).getByRole('button');
      await user.click(removeButton);

      // Verify router was called without the status filter
      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            page: 1,
            pageSize: 10,
            // status should NOT be present
          }),
          expect.any(Object)
        );
      });

      // Verify the URL params don't include status
      const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
      expect(lastCall[1]).not.toHaveProperty('status');
    });

    it('should clear all simple filter parameters when All tab is clicked', async () => {
      const simpleFilters = [
        createMockSimpleFilter('status-available', 'Available', 'AVAILABLE', { status: 'AVAILABLE' }),
        createMockSimpleFilter('status-borrowed', 'Borrowed', 'BORROWED', { status: 'BORROWED' }),
      ];

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters({ status: 'AVAILABLE' })}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          simpleFilters={simpleFilters}
        />
      );

      // Click the All tab
      const allTab = screen.getByRole('tab', { name: /All/i });
      await user.click(allTab);

      // Verify router was called without any status filter
      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            page: 1,
            pageSize: 10,
            // status should NOT be present
          }),
          expect.any(Object)
        );
      });

      // Verify the URL params don't include status
      const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
      expect(lastCall[1]).not.toHaveProperty('status');
    });

    it('should clear all filters when Clear all button is clicked', async () => {
      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters({ 
            status: 'AVAILABLE',
            author: 'Author 1',
            search: 'test' 
          })}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          customFilters={[
            { key: 'author', label: 'Author', type: 'text' }
          ]}
        />
      );

      // Find and click Clear all button
      const clearAllButton = screen.getByRole('button', { name: /Clear all/i });
      await user.click(clearAllButton);

      // Verify router was called with only preserved params (search)
      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            page: 1,
            pageSize: 10,
            search: 'test', // search should be preserved
            // status and author should NOT be present
          }),
          expect.any(Object)
        );
      });

      const lastCall = mockRouter.get.mock.calls[mockRouter.get.mock.calls.length - 1];
      expect(lastCall[1]).not.toHaveProperty('status');
      expect(lastCall[1]).not.toHaveProperty('author');
      expect(lastCall[1]).toHaveProperty('search', 'test');
    });
  });

  describe('Search Functionality', () => {
    it('should debounce search input and preserve other filters', async () => {
      vi.useFakeTimers();

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters({ status: 'AVAILABLE' })}
          title="Books"
          resourceName="books"
          columns={mockColumns}
        />
      );

      const searchInput = screen.getByPlaceholderText(/Search books/i);
      await user.type(searchInput, 'test');

      // Fast-forward debounce timer
      vi.advanceTimersByTime(300);

      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            search: 'test',
            status: 'AVAILABLE', // existing filter should be preserved
            page: 1,
          }),
          expect.any(Object)
        );
      });

      vi.useRealTimers();
    });
  });

  describe('Pagination', () => {
    it('should change page size and reset to page 1', async () => {
      render(
        <CrudPage
          data={createMockCrudData(mockBooks, { 
            total: 100, 
            currentPage: 3,
            perPage: 10 
          })}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
        />
      );

      // Find page size selector
      const pageSizeButton = screen.getByRole('button', { name: /10 \/ page/i });
      await user.click(pageSizeButton);

      // Select 25 per page
      const option25 = screen.getByRole('menuitem', { name: /25/i });
      await user.click(option25);

      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            page: 1, // should reset to page 1
            pageSize: 25,
          }),
          expect.any(Object)
        );
      });
    });
  });

  describe('Sorting', () => {
    it('should toggle sort direction when clicking sorted column', async () => {
      const mockColumnsWithSort = [
        { key: 'title', label: 'Title', sortable: true },
        { key: 'author', label: 'Author', sortable: true },
      ];

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters({ sort: 'title', direction: 'asc' })}
          title="Books"
          resourceName="books"
          columns={mockColumnsWithSort}
        />
      );

      // Click on the already sorted column header
      const titleHeader = screen.getByRole('button', { name: /Title/i });
      await user.click(titleHeader);

      await waitFor(() => {
        expect(mockRouter.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            sort: 'title',
            direction: 'desc', // should toggle from asc to desc
            page: 1,
          }),
          expect.any(Object)
        );
      });
    });
  });

  describe('Bulk Actions', () => {
    it('should show bulk action bar when items are selected', async () => {
      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          onBulkAction={vi.fn()}
        />
      );

      // Select first checkbox
      const checkboxes = screen.getAllByRole('checkbox');
      await user.click(checkboxes[1]); // First item checkbox (0 is select all)

      // Bulk action bar should appear
      expect(screen.getByText('1 selected')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /Delete/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /Clear/i })).toBeInTheDocument();
    });
  });

  describe('CRUD Operations', () => {
    it('should open create drawer when Add button is clicked', async () => {
      const CreateForm = () => <div>Create Form</div>;

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          createForm={CreateForm as any}
        />
      );

      const addButton = screen.getByRole('button', { name: /Add book/i });
      await user.click(addButton);

      // Drawer should open with create form
      await waitFor(() => {
        expect(screen.getByText('Create Form')).toBeInTheDocument();
      });
    });

    it('should handle keyboard shortcut for create (Cmd+N)', async () => {
      const CreateForm = () => <div>Create Form</div>;

      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          createForm={CreateForm as any}
        />
      );

      // Simulate Cmd+N
      await user.keyboard('{Meta>}n{/Meta}');

      // Drawer should open with create form
      await waitFor(() => {
        expect(screen.getByText('Create Form')).toBeInTheDocument();
      });
    });
  });

  describe('Permission-based Rendering', () => {
    it('should not show create button when canCreate is false', () => {
      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          createForm={(() => <div>Create Form</div>) as any}
          canCreate={false}
        />
      );

      expect(screen.queryByRole('button', { name: /Add book/i })).not.toBeInTheDocument();
    });

    it('should not show delete action when canDelete is false', async () => {
      render(
        <CrudPage
          data={createMockCrudData(mockBooks)}
          filters={createMockFilters()}
          title="Books"
          resourceName="books"
          columns={mockColumns}
          canDelete={false}
        />
      );

      // Open action menu for first item
      const actionButtons = screen.getAllByRole('button', { name: '' });
      const firstActionButton = actionButtons.find(btn => 
        btn.querySelector('svg')?.classList.contains('lucide-more-vertical')
      );
      
      if (firstActionButton) {
        await user.click(firstActionButton);
        
        // Delete option should not be present
        expect(screen.queryByRole('menuitem', { name: /Delete/i })).not.toBeInTheDocument();
      }
    });
  });
});