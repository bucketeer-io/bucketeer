import { ReactElement, useCallback, useState } from 'react';
import {
  type ColumnDef,
  type TableState,
  flexRender,
  getCoreRowModel,
  useReactTable,
  SortingState,
  getSortedRowModel,
  Updater
} from '@tanstack/react-table';
import { cn } from 'utils/style';
import { IconSorting, IconSortingDown, IconSortingUp } from '@icons';
import Table from 'components/table';
import PageLayout from 'elements/page-layout';

export interface DataTableProps<TData, TValue> {
  data: TData[];
  columns: ColumnDef<TData, TValue>[];
  state?: Partial<TableState>;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
  rowClassName?: string;
  onRowClick?: (data: TData) => void;
  onSortingChange?: (v: SortingState) => void;
}

export const DataTable = <TData, TValue>({
  data,
  columns,
  state,
  emptyCollection,
  isLoading,
  rowClassName,
  onRowClick,
  onSortingChange
}: DataTableProps<TData, TValue>) => {
  const [sorting, setSorting] = useState<SortingState>([]);

  const onSortingChangeHandler = useCallback(
    (updater: Updater<SortingState>) => {
      const newSorting =
        typeof updater === 'function' ? updater(sorting) : updater;

      setSorting(newSorting);
      if (data.length > 0) onSortingChange?.(newSorting);
    },
    [sorting, onSortingChange]
  );

  const table = useReactTable({
    data,
    columns,
    state: { ...state, sorting },
    onSortingChange: onSortingChangeHandler,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualSorting: true
  });

  return (
    <Table.Root>
      <Table.Header>
        {table.getHeaderGroups().map(headerGroup => (
          <Table.Row key={headerGroup.id}>
            {headerGroup.headers.map(header => (
              <Table.Head
                key={header.id}
                onClick={header.column.getToggleSortingHandler()}
                style={{ width: header.column.columnDef.size }}
                className={cn({
                  'cursor-pointer select-none':
                    header.column.columnDef.enableSorting !== false
                })}
              >
                {header.isPlaceholder ? null : (
                  <div className="flex items-center gap-3">
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    {header.column.columnDef.enableSorting !== false &&
                      {
                        asc: <IconSortingUp />,
                        desc: <IconSortingDown />,
                        false: <IconSorting />
                      }[header.column.getIsSorted() as string]}
                  </div>
                )}
              </Table.Head>
            ))}
          </Table.Row>
        ))}
      </Table.Header>
      <Table.Body>
        {isLoading ? (
          <Table.Row>
            <Table.Cell colSpan={columns.length}>
              <PageLayout.LoadingState className="py-10" />
            </Table.Cell>
          </Table.Row>
        ) : table.getRowModel().rows?.length ? (
          table.getRowModel().rows.map(row => (
            <Table.Row
              key={row.id}
              data-state={row.getIsSelected() && 'selected'}
              data-hoverable={!!onRowClick}
              onClick={() => onRowClick?.(row.original)}
              className={cn('shadow-card rounded-lg bg-white', rowClassName)}
            >
              {row.getVisibleCells().map(cell => (
                <Table.Cell
                  key={cell.id}
                  style={{ width: cell.column.columnDef.size }}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </Table.Cell>
              ))}
            </Table.Row>
          ))
        ) : (
          <Table.Row>
            <Table.Cell className="pt-32" colSpan={columns.length}>
              {emptyCollection}
            </Table.Cell>
          </Table.Row>
        )}
      </Table.Body>
    </Table.Root>
  );
};
