/**
 * Integration Tests for Book CRUD Operations
 * 
 * Tests the complete book lifecycle: Create, Read, Update, Delete
 * Including form validation, API interactions, and UI feedback
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { router } from '@inertiajs/react'
import { BookCreateForm, BookEditForm, BookDetailView } from '@/components/Books/BookForms'
import { CrudPage } from '@/components/Crud/CrudPage'
import { bookColumns } from '@/components/Books/book-columns'
import { Book, BookCreateData, BookUpdateData } from '@/types/book'

// Mock Inertia router
vi.mock('@inertiajs/react', () => ({
  router: {
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    get: vi.fn(),
  },
  Head: ({ children }: { children: string }) => <title>{children}</title>,
}))

// Mock data
const mockBook: Book = {
  id: 1,
  title: 'The Great Gatsby',
  author: 'F. Scott Fitzgerald',
  isbn: '978-0-7432-7356-5',
  description: 'A classic American novel about the Jazz Age',
  price: 12.99,
  status: 'AVAILABLE',
  publishedAt: '1925-04-10',
  tags: ['classic', 'fiction', 'american-literature'],
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
}

const mockBookList = [
  mockBook,
  {
    id: 2,
    title: 'To Kill a Mockingbird',
    author: 'Harper Lee',
    isbn: '978-0-06-112008-4',
    description: 'A gripping tale of racial injustice',
    price: 14.99,
    status: 'BORROWED',
    publishedAt: '1960-07-11',
    tags: ['classic', 'fiction', 'social-justice'],
    createdAt: '2024-01-16T11:00:00Z',
    updatedAt: '2024-01-16T11:00:00Z',
    borrowedBy: 'John Doe',
    borrowedAt: '2024-01-20T09:00:00Z',
    dueDate: '2024-02-20T09:00:00Z',
  },
  {
    id: 3,
    title: '1984',
    author: 'George Orwell',
    isbn: '978-0-452-28423-4',
    description: 'A dystopian social science fiction novel',
    price: 13.99,
    status: 'MAINTENANCE',
    publishedAt: '1949-06-08',
    tags: ['dystopian', 'fiction', 'political'],
    createdAt: '2024-01-17T12:00:00Z',
    updatedAt: '2024-01-17T12:00:00Z',
  },
] as Book[]

const mockPaginatedData = {
  data: mockBookList,
  total: 3,
  perPage: 10,
  currentPage: 1,
  lastPage: 1,
  from: 1,
  to: 3,
}

const mockFilters = {
  search: '',
  sort: '',
  direction: 'asc' as const,
  page: 1,
  perPage: 10,
  filters: {},
}

describe('Book CRUD Operations Integration Tests', () => {
  const user = userEvent.setup()
  
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('Book Creation', () => {
    it('should render create form with all required fields', () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      render(<BookCreateForm {...mockProps} />)

      // Check all form fields are present
      expect(screen.getByLabelText(/title/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/author/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/isbn/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/price/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/status/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/published date/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/tags/i)).toBeInTheDocument()

      // Check action buttons
      expect(screen.getByRole('button', { name: /create book/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
    })

    it('should validate required fields on form submission', async () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      render(<BookCreateForm {...mockProps} />)

      const submitButton = screen.getByRole('button', { name: /create book/i })
      await user.click(submitButton)

      // Check validation errors appear
      await waitFor(() => {
        expect(screen.getByText(/title is required/i)).toBeInTheDocument()
        expect(screen.getByText(/author is required/i)).toBeInTheDocument()
        expect(screen.getByText(/isbn is required/i)).toBeInTheDocument()
      })

      // Ensure form was not submitted
      expect(router.post).not.toHaveBeenCalled()
    })

    it('should submit valid form data', async () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      const mockRouterPost = vi.mocked(router.post)
      mockRouterPost.mockImplementation((url, data, options) => {
        options?.onSuccess?.()
        return Promise.resolve()
      })

      render(<BookCreateForm {...mockProps} />)

      // Fill out form with valid data
      await user.type(screen.getByLabelText(/title/i), 'Test Book')
      await user.type(screen.getByLabelText(/author/i), 'Test Author')
      await user.type(screen.getByLabelText(/isbn/i), '978-1234567890')
      await user.type(screen.getByLabelText(/price/i), '19.99')
      await user.type(screen.getByLabelText(/description/i), 'A test book description')

      // Add tags
      const tagInput = screen.getByPlaceholderText(/add a tag/i)
      await user.type(tagInput, 'test-tag')
      await user.click(screen.getByRole('button', { name: /add/i }))

      // Submit form
      const submitButton = screen.getByRole('button', { name: /create book/i })
      await user.click(submitButton)

      // Verify API call
      await waitFor(() => {
        expect(router.post).toHaveBeenCalledWith('/api/books', expect.objectContaining({
          title: 'Test Book',
          author: 'Test Author',
          isbn: '978-1234567890',
          price: 19.99,
          description: 'A test book description',
          tags: ['test-tag'],
          status: 'AVAILABLE',
        }), expect.any(Object))
      })

      // Verify success callback
      expect(mockProps.onSuccess).toHaveBeenCalled()
    })

    it('should handle API errors gracefully', async () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      const mockRouterPost = vi.mocked(router.post)
      mockRouterPost.mockImplementation((url, data, options) => {
        options?.onError?.({ 
          title: 'Title already exists',
          isbn: 'Invalid ISBN format'
        })
        return Promise.resolve()
      })

      render(<BookCreateForm {...mockProps} />)

      // Fill form with valid data
      await user.type(screen.getByLabelText(/title/i), 'Duplicate Book')
      await user.type(screen.getByLabelText(/author/i), 'Test Author')
      await user.type(screen.getByLabelText(/isbn/i), 'invalid-isbn')

      // Submit form
      await user.click(screen.getByRole('button', { name: /create book/i }))

      // Check error messages appear
      await waitFor(() => {
        expect(screen.getByText(/title already exists/i)).toBeInTheDocument()
        expect(screen.getByText(/invalid isbn format/i)).toBeInTheDocument()
      })

      expect(mockProps.onSuccess).not.toHaveBeenCalled()
    })

    it('should manage tags correctly', async () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      render(<BookCreateForm {...mockProps} />)

      const tagInput = screen.getByPlaceholderText(/add a tag/i)
      const addButton = screen.getByRole('button', { name: /add/i })

      // Add multiple tags
      await user.type(tagInput, 'fiction')
      await user.click(addButton)

      await user.type(tagInput, 'classic')
      await user.click(addButton)

      // Verify tags appear
      expect(screen.getByText('fiction')).toBeInTheDocument()
      expect(screen.getByText('classic')).toBeInTheDocument()

      // Test tag removal
      const removeButtons = screen.getAllByText('×')
      await user.click(removeButtons[0])

      // Verify tag was removed
      expect(screen.queryByText('fiction')).not.toBeInTheDocument()
      expect(screen.getByText('classic')).toBeInTheDocument()
    })
  })

  describe('Book Editing', () => {
    it('should populate form with existing book data', () => {
      const mockProps = {
        item: mockBook,
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      render(<BookEditForm {...mockProps} />)

      // Check form is populated with existing data
      expect(screen.getByDisplayValue(mockBook.title)).toBeInTheDocument()
      expect(screen.getByDisplayValue(mockBook.author)).toBeInTheDocument()
      expect(screen.getByDisplayValue(mockBook.isbn)).toBeInTheDocument()
      expect(screen.getByDisplayValue(mockBook.price.toString())).toBeInTheDocument()
      expect(screen.getByDisplayValue(mockBook.description!)).toBeInTheDocument()

      // Check tags are displayed
      mockBook.tags!.forEach(tag => {
        expect(screen.getByText(tag)).toBeInTheDocument()
      })
    })

    it('should update book with modified data', async () => {
      const mockProps = {
        item: mockBook,
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: false,
      }

      const mockRouterPut = vi.mocked(router.put)
      mockRouterPut.mockImplementation((url, data, options) => {
        options?.onSuccess?.()
        return Promise.resolve()
      })

      render(<BookEditForm {...mockProps} />)

      // Modify form data
      const titleInput = screen.getByDisplayValue(mockBook.title)
      await user.clear(titleInput)
      await user.type(titleInput, 'Updated Title')

      const priceInput = screen.getByDisplayValue(mockBook.price.toString())
      await user.clear(priceInput)
      await user.type(priceInput, '25.99')

      // Submit form
      await user.click(screen.getByRole('button', { name: /update book/i }))

      // Verify API call
      await waitFor(() => {
        expect(router.put).toHaveBeenCalledWith(
          `/api/books/${mockBook.id}`,
          expect.objectContaining({
            title: 'Updated Title',
            price: 25.99,
          }),
          expect.any(Object)
        )
      })

      expect(mockProps.onSuccess).toHaveBeenCalled()
    })
  })

  describe('Book Detail View', () => {
    it('should display all book information correctly', () => {
      const mockProps = {
        item: mockBook,
        onEdit: vi.fn(),
        onClose: vi.fn(),
        canEdit: true,
      }

      render(<BookDetailView {...mockProps} />)

      // Check book information is displayed
      expect(screen.getByText(mockBook.title)).toBeInTheDocument()
      expect(screen.getByText(`by ${mockBook.author}`)).toBeInTheDocument()
      expect(screen.getByText(mockBook.isbn)).toBeInTheDocument()
      expect(screen.getByText(mockBook.description!)).toBeInTheDocument()
      expect(screen.getByText('$12.99')).toBeInTheDocument()

      // Check tags are displayed
      mockBook.tags!.forEach(tag => {
        expect(screen.getByText(tag)).toBeInTheDocument()
      })

      // Check action buttons
      expect(screen.getByRole('button', { name: /edit book/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /close/i })).toBeInTheDocument()
    })

    it('should show borrow button for available books', () => {
      const availableBook = { ...mockBook, status: 'AVAILABLE' as const }
      const mockProps = {
        item: availableBook,
        onEdit: vi.fn(),
        onClose: vi.fn(),
        canEdit: true,
      }

      render(<BookDetailView {...mockProps} />)

      expect(screen.getByRole('button', { name: /borrow book/i })).toBeInTheDocument()
    })

    it('should show return button for borrowed books', () => {
      const borrowedBook = { ...mockBook, status: 'BORROWED' as const }
      const mockProps = {
        item: borrowedBook,
        onEdit: vi.fn(),
        onClose: vi.fn(),
        canEdit: true,
      }

      render(<BookDetailView {...mockProps} />)

      expect(screen.getByRole('button', { name: /return book/i })).toBeInTheDocument()
    })

    it('should handle borrow action', async () => {
      const availableBook = { ...mockBook, status: 'AVAILABLE' as const }
      const mockProps = {
        item: availableBook,
        onEdit: vi.fn(),
        onClose: vi.fn(),
        canEdit: true,
      }

      const mockRouterPost = vi.mocked(router.post)
      mockRouterPost.mockImplementation((url, data, options) => {
        options?.onSuccess?.()
        return Promise.resolve()
      })

      render(<BookDetailView {...mockProps} />)

      const borrowButton = screen.getByRole('button', { name: /borrow book/i })
      await user.click(borrowButton)

      await waitFor(() => {
        expect(router.post).toHaveBeenCalledWith(
          `/api/books/${availableBook.id}/borrow`,
          {},
          expect.any(Object)
        )
      })
    })
  })

  describe('Book List and Table Operations', () => {
    it('should render book list with proper table structure', () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      // Check table headers
      expect(screen.getByText('Title')).toBeInTheDocument()
      expect(screen.getByText('Author')).toBeInTheDocument()
      expect(screen.getByText('ISBN')).toBeInTheDocument()
      expect(screen.getByText('Status')).toBeInTheDocument()
      expect(screen.getByText('Price')).toBeInTheDocument()

      // Check book data is displayed
      expect(screen.getByText(mockBook.title)).toBeInTheDocument()
      expect(screen.getByText(mockBook.author)).toBeInTheDocument()
    })

    it('should handle search functionality', async () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      const mockRouterGet = vi.mocked(router.get)

      render(<CrudPage {...mockProps} />)

      const searchInput = screen.getByPlaceholderText(/search books/i)
      await user.type(searchInput, 'gatsby')

      // Wait for debounced search
      await waitFor(() => {
        expect(router.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            search: 'gatsby',
            page: 1,
          }),
          expect.any(Object)
        )
      }, { timeout: 500 })
    })

    it('should handle sorting', async () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      // Click on title column to sort
      const titleHeader = screen.getByText('Title')
      await user.click(titleHeader)

      await waitFor(() => {
        expect(router.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            sort: 'title',
            direction: 'asc',
          }),
          expect.any(Object)
        )
      })
    })

    it('should handle row selection for bulk operations', async () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
        onBulkAction: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      // Select individual rows
      const checkboxes = screen.getAllByRole('checkbox')
      const firstRowCheckbox = checkboxes[1] // Skip header checkbox

      await user.click(firstRowCheckbox)

      // Check bulk action bar appears
      await waitFor(() => {
        expect(screen.getByText(/1 selected/i)).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
      })
    })

    it('should handle pagination', async () => {
      const mockPaginatedDataMultiPage = {
        ...mockPaginatedData,
        lastPage: 3,
        currentPage: 1,
      }

      const mockProps = {
        data: mockPaginatedDataMultiPage,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      // Check pagination controls exist
      expect(screen.getByText(/page 1 of 3/i)).toBeInTheDocument()

      // Click next page
      const nextButton = screen.getByRole('button', { name: /go to next page/i })
      await user.click(nextButton)

      await waitFor(() => {
        expect(router.get).toHaveBeenCalledWith(
          '/admin/books',
          expect.objectContaining({
            page: 2,
          }),
          expect.any(Object)
        )
      })
    })
  })

  describe('Error Handling and Loading States', () => {
    it('should display loading state during form submission', async () => {
      const mockProps = {
        onSuccess: vi.fn(),
        onCancel: vi.fn(),
        isLoading: true,
      }

      render(<BookCreateForm {...mockProps} />)

      const submitButton = screen.getByRole('button', { name: /creating.../i })
      expect(submitButton).toBeDisabled()
    })

    it('should display loading state in table', () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
        loading: true,
      }

      render(<CrudPage {...mockProps} />)

      expect(screen.getByText(/loading/i)).toBeInTheDocument()
    })

    it('should display empty state when no books exist', () => {
      const emptyData = {
        ...mockPaginatedData,
        data: [],
        total: 0,
      }

      const mockProps = {
        data: emptyData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      expect(screen.getByText(/no results found/i)).toBeInTheDocument()
    })
  })

  describe('Book Status Management', () => {
    it('should display correct status badges', () => {
      const mockProps = {
        data: mockPaginatedData,
        filters: mockFilters,
        title: 'Books',
        resourceName: 'books',
        columns: bookColumns,
        onRefresh: vi.fn(),
      }

      render(<CrudPage {...mockProps} />)

      // Check status badges are displayed
      expect(screen.getByText(/available/i)).toBeInTheDocument()
      expect(screen.getByText(/borrowed/i)).toBeInTheDocument()
      expect(screen.getByText(/maintenance/i)).toBeInTheDocument()
    })

    it('should show due date for borrowed books', () => {
      const borrowedBook = {
        ...mockBook,
        status: 'BORROWED' as const,
        dueDate: '2024-02-20T09:00:00Z',
      }

      const mockProps = {
        item: borrowedBook,
        onEdit: vi.fn(),
        onClose: vi.fn(),
        canEdit: true,
      }

      render(<BookDetailView {...mockProps} />)

      expect(screen.getByText(/due date/i)).toBeInTheDocument()
      expect(screen.getByText(/2\/20\/2024/)).toBeInTheDocument()
    })
  })
})