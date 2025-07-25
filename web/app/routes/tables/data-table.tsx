import type {
  ColumnDef,
  OnChangeFn,
  PaginationState,
  SortingState,
  RowSelectionState,
} from "@tanstack/react-table";

import {
  flexRender,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronsLeftIcon,
  ChevronsRightIcon,
} from "lucide-react";
import React, {
  useState,
  useEffect,
  useRef,
  useMemo,
  useCallback,
  Fragment,
} from "react";
import { Button } from "~/components/ui/button";
import { Label } from "~/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { cn } from "~/lib/utils";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  emptyMessage?: string;
  pagination: PaginationState;
  totalRows: number;
  onPaginationChange: OnChangeFn<PaginationState>;
  onRowClick?: (row: TData) => void;
  tableMeta?: any;
  onRowSelectionChange?: (selectedRows: TData[]) => void;
  rowSelection?: RowSelectionState;
  onRowSelectionStateChange?: OnChangeFn<RowSelectionState>;
}

export const DataTable = React.memo(
  <TData, TValue>({
    columns,
    data,
    emptyMessage = "No results.",
    pagination,
    totalRows,
    onPaginationChange,
    onRowClick,
    tableMeta,
    onRowSelectionChange,
    rowSelection: controlledRowSelection,
    onRowSelectionStateChange,
  }: DataTableProps<TData, TValue>) => {
    const tableRef = useRef<HTMLTableElement>(null);

    const [sorting, setSorting] = useState<SortingState>([]);
    const [internalRowSelection, setInternalRowSelection] =
      useState<RowSelectionState>({});

    // Use controlled selection if provided, otherwise use internal state
    const rowSelection = controlledRowSelection ?? internalRowSelection;
    const setRowSelection =
      onRowSelectionStateChange ?? setInternalRowSelection;

    // Notify parent of selection changes
    useEffect(() => {
      if (onRowSelectionChange) {
        const selectedRows = Object.keys(rowSelection)
          .filter((key) => rowSelection[key])
          .map((index) => data[parseInt(index)])
          .filter(Boolean);
        onRowSelectionChange(selectedRows);
      }
    }, [rowSelection, data, onRowSelectionChange]);

    // Memoize row click handler to prevent unnecessary re-renders
    const handleRowClick = useCallback(
      (row: TData) => {
        if (onRowClick) {
          onRowClick(row);
        }
      },
      [onRowClick]
    );

    const table = useReactTable({
      data: Array.isArray(data) ? data : [],
      columns,
      getCoreRowModel: getCoreRowModel(),
      onSortingChange: setSorting,
      getSortedRowModel: getSortedRowModel(),
      getPaginationRowModel: getPaginationRowModel(),
      manualPagination: true,
      rowCount: totalRows,
      onPaginationChange: onPaginationChange,
      onRowSelectionChange: setRowSelection,
      state: {
        sorting,
        pagination,
        rowSelection,
      },
      meta: tableMeta,
    });

    return (
      <div className="flex flex-col h-full">
        <div className="flex-grow overflow-auto min-h-0 rounded-lg">
          <Table ref={tableRef}>
            <TableHeader className="sticky top-0 z-10 [&_tr:first-child_th:first-child]:rounded-tl-lg [&_tr:first-child_th:last-child]:rounded-tr-lg">
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHead
                        key={header.id}
                        className={cn({
                          "sticky-column":
                            header.column.columnDef.meta &&
                            (header.column.columnDef.meta as any).isSticky,
                        })}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                      </TableHead>
                    );
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows &&
              table.getRowModel().rows.length > 0 ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                    className={cn(
                      "data-table-row",
                      onRowClick && "cursor-pointer",
                      (row.original as any)?.id ===
                        (table.options.meta as any)?.editingRowId &&
                        "bg-muted/50"
                    )}
                    onClick={() => handleRowClick(row.original)}
                    tabIndex={onRowClick ? 0 : undefined}
                    onKeyDown={(e) => {
                      if (onRowClick && (e.key === "Enter" || e.key === " ")) {
                        e.preventDefault();
                        handleRowClick(row.original);
                      }
                    }}
                  >
                    {row.getVisibleCells().map((cell) => {
                      return (
                        <TableCell
                          key={cell.id}
                          className={cn({
                            "sticky-column":
                              cell.column.columnDef.meta &&
                              (cell.column.columnDef.meta as any).isSticky,
                          })}
                        >
                          {flexRender(
                            cell.column.columnDef.cell,
                            cell.getContext()
                          )}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    {emptyMessage}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
        <div className="flex items-center justify-between px-4 py-3 border-t bg-background sticky bottom-0 z-10 rounded-b-lg">
          <div className="hidden flex-1 text-sm text-muted-foreground lg:flex">
            {table.getFilteredSelectedRowModel().rows.length} of{" "}
            {table.getFilteredRowModel().rows.length} row(s) selected.
          </div>
          <div className="flex  items-center justify-end gap-6 lg:gap-8 lg:w-fit">
            <div className="hidden items-center gap-2 lg:flex">
              <Label htmlFor="rows-per-page" className="text-sm font-medium">
                Rows per page
              </Label>
              <Select
                value={`${table.getState().pagination.pageSize}`}
                onValueChange={(value) => {
                  table.setPageSize(Number(value));
                }}
              >
                <SelectTrigger className="w-20" id="rows-per-page">
                  <SelectValue
                    placeholder={table.getState().pagination.pageSize}
                  />
                </SelectTrigger>
                <SelectContent side="top">
                  {[50, 100, 150, 200, 250].map((pageSize) => (
                    <SelectItem key={pageSize} value={`${pageSize}`}>
                      {pageSize}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="flex w-fit items-center justify-center text-sm font-medium">
              Page {table.getState().pagination.pageIndex + 1} of{" "}
              {Math.ceil(totalRows / pagination.pageSize) || 1}
              <span className="ml-1 text-muted-foreground">
                ({pagination.pageSize * pagination.pageIndex + 1}-
                {Math.min(
                  pagination.pageSize * (pagination.pageIndex + 1),
                  totalRows
                )}{" "}
                of {totalRows})
              </span>
            </div>
            <div className="ml-auto flex items-center gap-2 lg:ml-0">
              <Button
                variant="outline"
                className="hidden h-8 w-8 p-0 lg:flex"
                onClick={() => table.setPageIndex(0)}
                disabled={!table.getCanPreviousPage()}
              >
                <span className="sr-only">Go to first page</span>
                <ChevronsLeftIcon />
              </Button>
              <Button
                variant="outline"
                className="size-8"
                size="icon"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
              >
                <span className="sr-only">Go to previous page</span>
                <ChevronLeftIcon />
              </Button>
              <Button
                variant="outline"
                className="size-8"
                size="icon"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
              >
                <span className="sr-only">Go to next page</span>
                <ChevronRightIcon />
              </Button>
              <Button
                variant="outline"
                className="hidden size-8 lg:flex"
                size="icon"
                onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                disabled={!table.getCanNextPage()}
              >
                <span className="sr-only">Go to last page</span>
                <ChevronsRightIcon />
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }
) as <TData, TValue>(
  props: DataTableProps<TData, TValue>
) => React.ReactElement;
