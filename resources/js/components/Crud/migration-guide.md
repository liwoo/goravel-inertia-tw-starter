# Migration Guide: CrudDataTable to shadcn DataTable

This guide shows how to migrate from the custom `CrudDataTable` component to the new standard shadcn `DataTable` implementation.

## Overview

The new implementation provides:
- **Standard shadcn/ui patterns** - Follows official shadcn documentation
- **TanStack Table integration** - Full-featured table with sorting, filtering, pagination
- **Better TypeScript support** - Strong typing with generics
- **Enhanced accessibility** - Better ARIA labels and keyboard navigation
- **Improved performance** - Optimized rendering and state management

## Components Available

### Core Components
- `DataTable` - Basic table component
- `EnhancedDataTable` - Feature-rich table with search, pagination, column toggle
- `DataTableColumnHeader` - Sortable column headers with dropdowns
- `DataTablePagination` - Advanced pagination controls
- `DataTableViewOptions` - Column visibility toggle

## Migration Steps

### 1. Column Definitions

**Before (CrudDataTable):**
```tsx
const columns: CrudColumn<Book>[] = [
  {
    key: 'title',
    label: 'Title',
    sortable: true,
    render: (book) => <span className="font-medium">{book.title}</span>
  },
  {
    key: 'status',
    label: 'Status',
    render: (book) => <StatusBadge status={book.status} />
  }
]
```

**After (shadcn DataTable):**
```tsx
const bookColumns: ColumnDef<Book>[] = [
  {
    accessorKey: "title",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Title" />
    ),
    cell: ({ row }) => {
      const title = row.getValue("title") as string
      return <span className="font-medium">{title}</span>
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status") as BookStatus
      return <StatusBadge status={status} />
    },
  },
]
```

### 2. Table Usage

**Before:**
```tsx
<CrudDataTable
  data={books}
  columns={columns}
  actions={actions}
  sortField={sortField}
  sortDirection={sortDirection}
  onSort={handleSort}
  selectedIds={selectedIds}
  onSelectionChange={setSelectedIds}
  enableSelection={true}
  loading={loading}
/>
```

**After:**
```tsx
<EnhancedDataTable
  columns={bookColumns}
  data={books}
  searchKey="title"
  searchPlaceholder="Search books..."
  loading={loading}
  onRowSelectionChange={handleRowSelectionChange}
  enableSearch={true}
  enableColumnToggle={true}
  enablePagination={true}
/>
```

### 3. Row Selection

**Before:**
```tsx
// selectedIds: number[]
// onSelectionChange: (ids: number[]) => void
const handleSelectionChange = (ids: number[]) => {
  setSelectedIds(ids)
}
```

**After:**
```tsx
// selectedRows: Book[]
// onRowSelectionChange: (rows: Book[]) => void
const handleRowSelectionChange = (selectedRows: Book[]) => {
  setSelectedBooks(selectedRows)
  // Get IDs if needed: selectedRows.map(row => row.id)
}
```

### 4. Row Actions

**Before:**
```tsx
const actions: CrudAction<Book>[] = [
  {
    key: 'edit',
    label: 'Edit',
    icon: <Edit className="h-4 w-4" />,
    onClick: (book) => handleEdit(book)
  },
  {
    key: 'delete',
    label: 'Delete',
    icon: <Trash className="h-4 w-4" />,
    onClick: (book) => handleDelete(book),
    confirm: true,
    confirmMessage: 'Are you sure you want to delete this book?'
  }
]
```

**After:**
```tsx
// Include in column definition
{
  id: "actions",
  header: "Actions",
  cell: ({ row }) => {
    const book = row.original
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={() => handleEdit(book)}>
            <Edit className="mr-2 h-4 w-4" />
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem 
            onClick={() => handleDelete(book)}
            className="text-destructive"
          >
            <Trash className="mr-2 h-4 w-4" />
            Delete
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    )
  },
}
```

### 5. Selection Column

Add selection column to enable row selection:

```tsx
{
  id: "select",
  header: ({ table }) => (
    <Checkbox
      checked={
        table.getIsAllPageRowsSelected() ||
        (table.getIsSomePageRowsSelected() && "indeterminate")
      }
      onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
      aria-label="Select all"
    />
  ),
  cell: ({ row }) => (
    <Checkbox
      checked={row.getIsSelected()}
      onCheckedChange={(value) => row.toggleSelected(!!value)}
      aria-label="Select row"
    />
  ),
  enableSorting: false,
  enableHiding: false,
},
```

## Key Differences

### State Management
- **Before**: Manual state management for sorting, selection
- **After**: TanStack Table handles all state internally

### Type Safety
- **Before**: Generic `CrudColumn<T>` with string keys
- **After**: Full TypeScript integration with `ColumnDef<T, TValue>`

### Customization
- **Before**: Limited to predefined props
- **After**: Full access to TanStack Table API and customization

### Performance
- **Before**: Re-renders entire table on state changes
- **After**: Optimized rendering with React Table's built-in optimizations

## Benefits of Migration

1. **Standard Patterns**: Follows shadcn/ui conventions
2. **Better Performance**: TanStack Table optimizations
3. **Enhanced Features**: Built-in filtering, sorting, pagination
4. **Type Safety**: Full TypeScript support
5. **Accessibility**: Better ARIA support and keyboard navigation
6. **Maintainability**: Standard patterns, easier to maintain
7. **Extensibility**: Easy to add new features and customizations

## Complete Example

See `books-data-table-example.tsx` for a complete working example with:
- Column definitions
- Row selection
- Search functionality
- Loading states
- Pagination
- Column visibility toggle