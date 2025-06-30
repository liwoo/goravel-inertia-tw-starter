"use client"

import * as React from 'react';
import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  Row,
} from '@tanstack/react-table';
import { 
  ChevronUp, 
  ChevronDown, 
  ChevronsUpDown, 
  FileX, 
  MoreVertical,
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  Settings2,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { DataTableProps, CrudAction, CrudColumn } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

// Column header component with improved sorting
function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  className,
}: {
  column: any
  title: string
  className?: string
}) {
  if (!column.getCanSort()) {
    return <div className={cn(className)}>{title}</div>
  }

  return (
    <div className={cn("flex items-center space-x-2", className)}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            className="data-[state=open]:bg-accent -ml-3 h-8"
          >
            <span>{title}</span>
            {column.getIsSorted() === "desc" ? (
              <ChevronDown className="ml-2 h-4 w-4" />
            ) : column.getIsSorted() === "asc" ? (
              <ChevronUp className="ml-2 h-4 w-4" />
            ) : (
              <ChevronsUpDown className="ml-2 h-4 w-4" />
            )}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start">
          <DropdownMenuItem onClick={() => column.toggleSorting(false)}>
            <ChevronUp className="mr-2 h-3 w-3 text-muted-foreground/70" />
            Asc
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => column.toggleSorting(true)}>
            <ChevronDown className="mr-2 h-3 w-3 text-muted-foreground/70" />
            Desc
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}

// Helper function to convert CrudColumn to ColumnDef
const createColumnDef = <T extends { id: number }>(
  column: CrudColumn<T>
): ColumnDef<T> => ({
  accessorKey: column.key,
  id: column.key,
  header: ({ column: tanColumn }) => (
    <DataTableColumnHeader 
      column={tanColumn} 
      title={column.label}
      className={column.className}
    />
  ),
  cell: ({ row }) => {
    const item = row.original;
    return (
      <div className={cn('text-sm', column.className)}>
        {column.render ? 
          column.render(item) : 
          String((item as any)[column.key] || '-')
        }
      </div>
    );
  },
  enableSorting: column.sortable,
  enableHiding: true,
  meta: {
    className: column.className,
    width: column.width,
  },
});

// Helper function to create actions column
const createActionsColumn = <T extends { id: number }>(
  actions: CrudAction<T>[]
): ColumnDef<T> => ({
  id: "actions",
  header: () => <span className="sr-only">Actions</span>,
  cell: ({ row }) => {
    const item = row.original;
    
    const handleActionClick = (action: CrudAction<T>, item: T) => {
      if (action.confirm) {
        const message = action.confirmMessage || `Are you sure you want to ${action.label.toLowerCase()}?`;
        if (confirm(message)) {
          action.onClick(item);
        }
      } else {
        action.onClick(item);
      }
    };

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            className="data-[state=open]:bg-muted hover:bg-muted text-muted-foreground h-8 w-8 p-0"
            size="icon"
            title="Actions menu (Alt+Enter)"
          >
            <MoreVertical className="h-4 w-4" />
            <span className="sr-only">Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-48">
          <DropdownMenuLabel>Actions</DropdownMenuLabel>
          {actions.map((action, actionIndex) => {
            const isDisabled = action.disabled?.(item) || false;
            const isDestructive = action.key === 'delete' || action.className?.includes('destructive');
            
            return (
              <React.Fragment key={action.key}>
                {actionIndex > 0 && actions[actionIndex - 1]?.key !== 'delete' && action.key === 'delete' && (
                  <DropdownMenuSeparator />
                )}
                <DropdownMenuItem
                  onClick={() => !isDisabled && handleActionClick(action, item)}
                  disabled={isDisabled}
                  className={cn(
                    'cursor-pointer',
                    isDestructive && 'text-destructive focus:text-destructive',
                    action.className,
                    isDisabled && 'opacity-50 cursor-not-allowed'
                  )}
                >
                  {action.icon && (
                    <span className="mr-2 h-4 w-4">{action.icon}</span>
                  )}
                  {action.label}
                </DropdownMenuItem>
              </React.Fragment>
            );
          })}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  },
  enableSorting: false,
  enableHiding: false,
});

// Pagination component
function DataTablePagination<TData>({
  table,
}: {
  table: any
}) {
  return (
    <div className="flex items-center justify-between px-2">
      <div className="flex-1 text-sm text-muted-foreground">
        {table.getFilteredSelectedRowModel().rows.length} of{" "}
        {table.getFilteredRowModel().rows.length} row(s) selected.
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select
            value={`${table.getState().pagination.pageSize}`}
            onValueChange={(value) => {
              table.setPageSize(Number(value))
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={table.getState().pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {table.getState().pagination.pageIndex + 1} of{" "}
          {table.getPageCount()}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => table.setPageIndex(0)}
            disabled={!table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to first page</span>
            <ChevronsLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => table.setPageIndex(table.getPageCount() - 1)}
            disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to last page</span>
            <ChevronsRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}

// View options component
function DataTableViewOptions<TData>({
  table,
}: {
  table: any
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="ml-auto hidden h-8 lg:flex"
        >
          <Settings2 className="mr-2 h-4 w-4" />
          View
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[150px]">
        <DropdownMenuLabel>Toggle columns</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {table
          .getAllColumns()
          .filter(
            (column: any) =>
              typeof column.accessorFn !== "undefined" && column.getCanHide()
          )
          .map((column: any) => {
            return (
              <DropdownMenuCheckboxItem
                key={column.id}
                className="capitalize"
                checked={column.getIsVisible()}
                onCheckedChange={(value) => column.toggleVisibility(!!value)}
              >
                {column.id}
              </DropdownMenuCheckboxItem>
            )
          })}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

// Extended interface for modern features
interface ModernDataTableProps<T> extends DataTableProps<T> {
  searchKey?: string
  searchPlaceholder?: string
  enableSearch?: boolean
  enableColumnToggle?: boolean
  enablePagination?: boolean
  onSearch?: (value: string) => void
}

export function CrudDataTable<T extends { id: number }>({
  data,
  columns,
  actions,
  sortField,
  sortDirection,
  onSort,
  selectedIds,
  onSelectionChange,
  enableSelection = false,
  loading = false,
  emptyMessage = 'No results found',
  className,
  searchKey,
  searchPlaceholder = "Search...",
  enableSearch = false,
  enableColumnToggle = false,
  enablePagination = false,
  onSearch,
}: ModernDataTableProps<T>) {
  // State management
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = React.useState({})
  const [globalFilter, setGlobalFilter] = React.useState("")

  // Update sorting when external props change
  React.useEffect(() => {
    if (sortField && sortDirection) {
      setSorting([{ id: sortField, desc: sortDirection === 'desc' }])
    } else {
      setSorting([])
    }
  }, [sortField, sortDirection])

  // Build columns for TanStack Table
  const tableColumns = React.useMemo(() => {
    const cols: ColumnDef<T>[] = []

    // Selection column
    if (enableSelection) {
      cols.push({
        id: "select",
        header: ({ table }) => (
          <div className="flex items-center justify-center">
            <Checkbox
              checked={
                table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && "indeterminate")
              }
              onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
              aria-label="Select all"
            />
          </div>
        ),
        cell: ({ row }) => (
          <div className="flex items-center justify-center">
            <Checkbox
              checked={row.getIsSelected()}
              onCheckedChange={(value) => row.toggleSelected(!!value)}
              aria-label="Select row"
            />
          </div>
        ),
        enableSorting: false,
        enableHiding: false,
        size: 50,
      })
    }

    // Data columns
    columns.forEach(column => {
      cols.push(createColumnDef(column))
    })

    // Actions column - always show if actions are provided
    if (actions && actions.length > 0) {
      cols.push(createActionsColumn(actions))
    }

    return cols
  }, [columns, actions, enableSelection])

  const table = useReactTable({
    data,
    columns: tableColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    // Removed getPaginationRowModel() - we use server-side pagination
    // Note: getSortedRowModel removed - using server-side sorting instead
    getFilteredRowModel: getFilteredRowModel(), // Keep for local search functionality
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    onGlobalFilterChange: setGlobalFilter,
    globalFilterFn: "includesString",
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      globalFilter,
    },
    enableRowSelection: enableSelection,
  })

  // Sync row selection with external state
  React.useEffect(() => {
    const selectedRows = table.getFilteredSelectedRowModel().rows
    const newSelectedIds = selectedRows.map(row => row.original.id)
    
    // Only update if there's a difference to avoid infinite loops
    if (JSON.stringify(newSelectedIds.sort()) !== JSON.stringify(selectedIds.sort())) {
      onSelectionChange(newSelectedIds)
    }
  }, [rowSelection, table, onSelectionChange, selectedIds])

  // Handle external search
  React.useEffect(() => {
    if (onSearch && globalFilter !== "") {
      onSearch(globalFilter)
    }
  }, [globalFilter, onSearch])

  return (
    <div className={cn("w-full space-y-4", className)}>
      {/* Toolbar */}
      {(enableSearch || enableColumnToggle) && (
        <div className="flex items-center justify-between">
          {enableSearch && searchKey && (
            <Input
              placeholder={searchPlaceholder}
              value={globalFilter ?? ""}
              onChange={(event) => setGlobalFilter(String(event.target.value))}
              className="max-w-sm"
            />
          )}
          {enableColumnToggle && <DataTableViewOptions table={table} />}
        </div>
      )}

      {/* Table */}
      <div className="relative">
        {loading && (
          <div className="absolute inset-0 bg-background/50 backdrop-blur-sm flex items-center justify-center z-50 rounded-md">
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
              <span className="text-sm text-muted-foreground">Loading...</span>
            </div>
          </div>
        )}
        
        <div className="rounded-md border overflow-hidden">
          <div className="overflow-x-auto max-h-[600px] overflow-y-auto">
            <Table className="w-full min-w-full">
            <TableHeader className="bg-muted/50 sticky top-0 z-10">
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHead 
                        key={header.id}
                        className={cn(
                          'px-2 py-3 text-left whitespace-nowrap bg-muted/50 border-b',
                          header.id === 'select' && 'w-12 min-w-[3rem]',
                          header.id === 'actions' && 'w-16 min-w-[4rem] sticky right-0 z-20 bg-muted/50 border-l',
                          // Add min-width for other columns to prevent them from being too narrow
                          header.id !== 'select' && header.id !== 'actions' && 'min-w-[8rem]'
                        )}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                      </TableHead>
                    )
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                    className="hover:bg-muted/50"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell 
                        key={cell.id}
                        className={cn(
                          'px-2 py-3 whitespace-nowrap',
                          cell.column.id === 'actions' && 'text-right sticky right-0 z-20 bg-background border-l',
                          cell.column.id === 'select' && 'w-12',
                          // Add consistent padding and prevent text wrapping
                        )}
                      >
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={tableColumns.length}
                    className="h-32 text-center"
                  >
                    <div className="flex flex-col items-center justify-center space-y-2">
                      <FileX className="h-8 w-8 text-muted-foreground/50" />
                      <p className="text-sm text-muted-foreground">{emptyMessage}</p>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
            </Table>
          </div>
        </div>
      </div>

      {/* Pagination */}
      {enablePagination && <DataTablePagination table={table} />}
    </div>
  )
}