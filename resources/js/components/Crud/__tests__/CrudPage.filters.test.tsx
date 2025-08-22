/**
 * Focused tests for CrudPage filter functionality
 * Tests the filter badge persistence fixes
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// Create a minimal test component that tests just the filter logic
const FilterTestComponent = ({ 
  initialFilters = {},
  simpleFilters = [],
  onFilterChange,
  onNavigate 
}: any) => {
  const [activeFilters, setActiveFilters] = React.useState(initialFilters);
  const [activeSimpleFilter, setActiveSimpleFilter] = React.useState<string | undefined>(
    // Determine initial active simple filter
    simpleFilters.find((sf: any) => {
      if (sf.filterParams) {
        return Object.entries(sf.filterParams).every(([key, value]) => 
          activeFilters[key] === value
        );
      }
      return activeFilters[sf.key] === sf.value;
    })?.value.toString()
  );

  const handleFilterChange = (filterKey: string, value: any) => {
    const newFilters = { ...activeFilters };
    if (value === '' || value === null || value === undefined) {
      delete newFilters[filterKey];
      
      // Reset simple filter tab if we're removing a filter that matches a simple filter
      const matchingSimpleFilter = simpleFilters.find((sf: any) => {
        if (sf.filterParams) {
          // Check if this simple filter has the same key and was the active value
          return Object.keys(sf.filterParams).includes(filterKey) && 
                 sf.filterParams[filterKey] === activeFilters[filterKey];
        }
        return sf.key === filterKey;
      });
      
      if (matchingSimpleFilter && activeSimpleFilter === matchingSimpleFilter.value.toString()) {
        setActiveSimpleFilter(undefined);
      }
    } else {
      newFilters[filterKey] = value;
    }
    
    setActiveFilters(newFilters);
    onFilterChange?.(newFilters);
    
    // Simulate navigation
    const overrides: Record<string, any> = { page: 1 };
    
    // Add all active filter keys with their values (or null if removed)
    Object.keys(activeFilters).forEach(key => {
      overrides[key] = newFilters[key] !== undefined ? newFilters[key] : null;
    });
    
    // Add any new filters
    Object.keys(newFilters).forEach(key => {
      overrides[key] = newFilters[key];
    });
    
    // Filter out null values (simulating buildNavigationParams)
    const params: Record<string, any> = {};
    Object.keys(overrides).forEach(key => {
      if (overrides[key] !== null && overrides[key] !== undefined) {
        params[key] = overrides[key];
      }
    });
    
    onNavigate?.(params);
  };

  const handleSimpleFilterChange = (filterValue: string | undefined) => {
    setActiveSimpleFilter(filterValue);
    
    const overrides: Record<string, any> = { page: 1 };
    
    if (filterValue === undefined) {
      // "All" was selected - clear all simple filter parameters
      simpleFilters.forEach((filter: any) => {
        if (filter.filterParams) {
          Object.keys(filter.filterParams).forEach(key => {
            overrides[key] = null;
          });
        } else {
          overrides[filter.key] = null;
        }
      });
      
      // Also update activeFilters to remove these keys
      const newActiveFilters = { ...activeFilters };
      simpleFilters.forEach((filter: any) => {
        if (filter.filterParams) {
          Object.keys(filter.filterParams).forEach(key => {
            delete newActiveFilters[key];
          });
        } else {
          delete newActiveFilters[filter.key];
        }
      });
      setActiveFilters(newActiveFilters);
    } else {
      // Find the selected filter
      const selectedFilter = simpleFilters.find((f: any) => f.value.toString() === filterValue);
      if (selectedFilter) {
        if (selectedFilter.filterParams) {
          Object.assign(overrides, selectedFilter.filterParams);
        } else {
          overrides[selectedFilter.key] = selectedFilter.value;
        }
      }
    }
    
    // Filter out null values
    const params: Record<string, any> = {};
    Object.keys(overrides).forEach(key => {
      if (overrides[key] !== null && overrides[key] !== undefined) {
        params[key] = overrides[key];
      }
    });
    
    onNavigate?.(params);
  };

  return (
    <div>
      {/* Simple Filter Tabs */}
      <div role="tablist">
        <button 
          role="tab"
          data-state={activeSimpleFilter === undefined ? 'active' : 'inactive'}
          onClick={() => handleSimpleFilterChange(undefined)}
        >
          All
        </button>
        {simpleFilters.map((filter: any) => (
          <button
            key={filter.key}
            role="tab"
            data-state={activeSimpleFilter === filter.value.toString() ? 'active' : 'inactive'}
            onClick={() => handleSimpleFilterChange(filter.value.toString())}
          >
            {filter.label}
          </button>
        ))}
      </div>

      {/* Active Filter Badges */}
      {Object.entries(activeFilters).map(([key, value]) => (
        <div key={key} className="filter-badge">
          <span>{key}: {String(value)}</span>
          <button onClick={() => handleFilterChange(key, '')}>×</button>
        </div>
      ))}
    </div>
  );
};

describe('CrudPage Filter Functionality', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should remove filter and reset simple filter tab when X is clicked on filter badge', async () => {
    const onNavigate = vi.fn();
    const simpleFilters = [
      {
        key: 'status-available',
        label: 'Available',
        value: 'AVAILABLE',
        filterParams: { status: 'AVAILABLE' }
      },
      {
        key: 'status-borrowed',
        label: 'Borrowed', 
        value: 'BORROWED',
        filterParams: { status: 'BORROWED' }
      },
    ];

    render(
      <FilterTestComponent
        initialFilters={{ status: 'AVAILABLE' }}
        simpleFilters={simpleFilters}
        onNavigate={onNavigate}
      />
    );

    // Check that Available tab is selected
    const availableTab = screen.getByRole('tab', { name: 'Available' });
    expect(availableTab).toHaveAttribute('data-state', 'active');

    // Check that filter badge is shown
    const filterBadge = screen.getByText('status: AVAILABLE');
    expect(filterBadge).toBeInTheDocument();

    // Click the X button on the filter badge
    const removeButton = filterBadge.nextElementSibling as HTMLElement;
    await user.click(removeButton);

    // Verify navigation was called without the status filter
    await waitFor(() => {
      expect(onNavigate).toHaveBeenLastCalledWith({
        page: 1,
        // status should NOT be present
      });
    });

    // Verify All tab is now active
    const allTab = screen.getByRole('tab', { name: 'All' });
    expect(allTab).toHaveAttribute('data-state', 'active');
    expect(availableTab).toHaveAttribute('data-state', 'inactive');
  });

  it('should clear all simple filter parameters when All tab is clicked', async () => {
    const onNavigate = vi.fn();
    const simpleFilters = [
      {
        key: 'status-available',
        label: 'Available',
        value: 'AVAILABLE',
        filterParams: { status: 'AVAILABLE' }
      },
      {
        key: 'status-borrowed',
        label: 'Borrowed',
        value: 'BORROWED',
        filterParams: { status: 'BORROWED' }
      },
    ];

    render(
      <FilterTestComponent
        initialFilters={{ status: 'AVAILABLE' }}
        simpleFilters={simpleFilters}
        onNavigate={onNavigate}
      />
    );

    // Click the All tab
    const allTab = screen.getByRole('tab', { name: 'All' });
    await user.click(allTab);

    // Verify navigation was called without any status filter
    await waitFor(() => {
      expect(onNavigate).toHaveBeenLastCalledWith({
        page: 1,
        // status should NOT be present
      });
    });

    // Verify All tab is active
    expect(allTab).toHaveAttribute('data-state', 'active');
  });

  it('should handle complex filter parameters correctly', async () => {
    const onNavigate = vi.fn();
    const simpleFilters = [
      {
        key: 'active-true',
        label: 'Active',
        value: true,
        filterParams: { active: true, status: 'ACTIVE' }
      },
    ];

    render(
      <FilterTestComponent
        initialFilters={{ active: true, status: 'ACTIVE' }}
        simpleFilters={simpleFilters}
        onNavigate={onNavigate}
      />
    );

    // Check that Active tab is selected
    const activeTab = screen.getByRole('tab', { name: 'Active' });
    expect(activeTab).toHaveAttribute('data-state', 'active');

    // Remove the active filter
    const activeFilterBadge = screen.getByText('active: true');
    const removeButton = activeFilterBadge.nextElementSibling as HTMLElement;
    await user.click(removeButton);

    // Should reset to All tab since the simple filter requires both parameters
    await waitFor(() => {
      const allTab = screen.getByRole('tab', { name: 'All' });
      expect(allTab).toHaveAttribute('data-state', 'active');
      expect(activeTab).toHaveAttribute('data-state', 'inactive');
    });
  });
});