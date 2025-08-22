import React, { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { vi } from 'vitest';

// Mock Inertia router
export const mockRouter = {
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

// Mock Inertia hooks
export const mockUsePage = vi.fn(() => ({
  props: {
    auth: {
      user: {
        id: 1,
        name: 'Test User',
        email: 'test@example.com',
        permissions: [],
      },
    },
    errors: {},
  },
  component: 'Test',
  version: '1',
  url: '/',
}));

// Create a wrapper component with all necessary providers
interface AllTheProvidersProps {
  children: React.ReactNode;
}

const AllTheProviders: React.FC<AllTheProvidersProps> = ({ children }) => {
  return <>{children}</>;
};

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options });

export * from '@testing-library/react';
export { customRender as render };

// Helper to create mock data for CrudPage
export function createMockCrudData<T extends { id: number }>(
  items: T[],
  options: {
    currentPage?: number;
    perPage?: number;
    total?: number;
  } = {}
) {
  const currentPage = options.currentPage || 1;
  const perPage = options.perPage || 10;
  const total = options.total || items.length;
  const lastPage = Math.ceil(total / perPage);

  return {
    data: items,
    currentPage,
    perPage,
    total,
    lastPage,
    from: (currentPage - 1) * perPage + 1,
    to: Math.min(currentPage * perPage, total),
  };
}

// Helper to create mock filters
export function createMockFilters(overrides: Record<string, any> = {}) {
  return {
    search: '',
    sort: undefined,
    direction: undefined,
    filters: {},
    ...overrides,
  };
}

// Helper to create simple filter configs
export function createMockSimpleFilter(
  key: string,
  label: string,
  value: any,
  filterParams?: Record<string, any>
) {
  return {
    key,
    label,
    value,
    filterParams: filterParams || { [key]: value },
    badge: 0,
  };
}