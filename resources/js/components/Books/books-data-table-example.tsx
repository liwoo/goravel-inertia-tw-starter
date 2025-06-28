"use client"

import { EnhancedDataTable } from "@/components/Crud/enhanced-data-table"
import { bookColumns } from "./book-columns"
import { Book } from "@/types/book"

interface BooksDataTableProps {
  data: Book[]
  loading?: boolean
  onRowSelectionChange?: (selectedBooks: Book[]) => void
}

export function BooksDataTable({
  data,
  loading = false,
  onRowSelectionChange,
}: BooksDataTableProps) {
  return (
    <EnhancedDataTable
      columns={bookColumns}
      data={data}
      searchKey="title"
      searchPlaceholder="Search books by title..."
      loading={loading}
      emptyMessage="No books found."
      onRowSelectionChange={onRowSelectionChange}
      enableSearch={true}
      enableColumnToggle={true}
      enablePagination={true}
    />
  )
}

// Example usage component
export function BooksPageExample() {
  const [books, setBooks] = React.useState<Book[]>([])
  const [loading, setLoading] = React.useState(false)
  const [selectedBooks, setSelectedBooks] = React.useState<Book[]>([])

  // Example data
  const exampleBooks: Book[] = [
    {
      id: 1,
      title: "The Great Gatsby",
      author: "F. Scott Fitzgerald",
      isbn: "978-0-7432-7356-5",
      description: "A classic American novel",
      price: 12.99,
      status: "AVAILABLE",
      publishedAt: "1925-04-10",
      tags: ["classic", "fiction", "american-literature"],
      createdAt: "2024-01-15T10:00:00Z",
      updatedAt: "2024-01-15T10:00:00Z",
    },
    {
      id: 2,
      title: "To Kill a Mockingbird",
      author: "Harper Lee",
      isbn: "978-0-06-112008-4",
      description: "A gripping tale of racial injustice",
      price: 14.99,
      status: "BORROWED",
      publishedAt: "1960-07-11",
      tags: ["classic", "fiction", "social-justice"],
      createdAt: "2024-01-16T11:00:00Z",
      updatedAt: "2024-01-16T11:00:00Z",
      borrowedBy: "John Doe",
      borrowedAt: "2024-01-20T09:00:00Z",
      dueDate: "2024-02-20T09:00:00Z",
    },
    {
      id: 3,
      title: "1984",
      author: "George Orwell",
      isbn: "978-0-452-28423-4",
      description: "A dystopian social science fiction novel",
      price: 13.99,
      status: "MAINTENANCE",
      publishedAt: "1949-06-08",
      tags: ["dystopian", "fiction", "political"],
      createdAt: "2024-01-17T12:00:00Z",
      updatedAt: "2024-01-17T12:00:00Z",
    },
  ]

  React.useEffect(() => {
    // Simulate loading data
    setLoading(true)
    setTimeout(() => {
      setBooks(exampleBooks)
      setLoading(false)
    }, 1000)
  }, [])

  const handleRowSelectionChange = (selectedRows: Book[]) => {
    setSelectedBooks(selectedRows)
    console.log("Selected books:", selectedRows)
  }

  return (
    <div className="container mx-auto py-10 space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Books</h1>
          <p className="text-muted-foreground">
            Manage your book collection
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {selectedBooks.length > 0 && (
            <p className="text-sm text-muted-foreground">
              {selectedBooks.length} book(s) selected
            </p>
          )}
        </div>
      </div>

      <BooksDataTable
        data={books}
        loading={loading}
        onRowSelectionChange={handleRowSelectionChange}
      />
    </div>
  )
}