import { ReactElement, useCallback, useEffect, useRef, useState } from 'react';
import {
  type ColumnDef,
  SortingState,
  type TableState,
  Updater,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable
} from '@tanstack/react-table';
import { cn } from 'utils/style';
import { IconSorting, IconSortingDown, IconSortingUp } from '@icons';
import Table from 'components/table';
import PageLayout from 'elements/page-layout';

export interface DataTableProps<TData, TValue> {
  data: TData[];
  columns: ColumnDef<TData, TValue>[];
  manualSorting?: boolean;
  state?: Partial<TableState>;
  emptyCollection?: ReactElement;
  isLoading?: boolean;
  onRowClick?: (data: TData) => void;
  onSortingChange?: (v: SortingState) => void;
}

export const DataTable = <TData, TValue>({
  data,
  columns,
  manualSorting = true,
  state,
  emptyCollection,
  isLoading,
  onRowClick,
  onSortingChange
}: DataTableProps<TData, TValue>) => {
  const tableContainerRef = useRef<HTMLDivElement>(null);
  const [isScrolltable, setIsScrollTable] = useState(false);
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
    manualSorting
  });

  useEffect(() => {
    const checkScrollTable = () => {
      if (!tableContainerRef.current) return;
      const isScrollable =
        tableContainerRef.current.scrollWidth >
        tableContainerRef.current.clientWidth;

      setIsScrollTable(isScrollable);
    };
    checkScrollTable();
    window.addEventListener('resize', () => checkScrollTable());
    return () => {
      window.removeEventListener('resize', () => checkScrollTable());
    };
  }, [tableContainerRef.current]);

  return (
    <div
      ref={tableContainerRef}
      className={cn(
        'overflow-x-auto hidden-scroll w-full relative',
        isScrolltable ? 'px-0' : 'px-2'
      )}
    >
      <Table.Root>
        <Table.Header>
          {table.getHeaderGroups().map(headerGroup => (
            <Table.Row key={headerGroup.id}>
              {headerGroup.headers.map((header, index) => (
                <Table.Head
                  key={header.id}
                  onClick={header.column.getToggleSortingHandler()}
                  style={{ width: header.column.columnDef.size }}
                  className={cn({
                    'cursor-pointer select-none':
                      header.column.columnDef.enableSorting !== false,
                    'sticky bg-white':
                      index === 0 || index === headerGroup.headers.length - 1,
                    'left-0': index === 0,
                    'right-0': index === headerGroup.headers.length - 1
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
            table.getRowModel().rows.map(row => {
              return (
                <Table.Row
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                  data-hoverable={!!onRowClick}
                  onClick={() => onRowClick?.(row.original)}
                  className={cn('shadow-card bg-white', {
                    'rounded-lg': !isScrolltable
                  })}
                >
                  {row.getVisibleCells().map((cell, cellIndex) => {
                    const visibleCells = row.getVisibleCells();
                    const lastIndex = visibleCells.length - 1;

                    const isFirst = cellIndex === 0;
                    const isLast = cellIndex === lastIndex;
                    const isSticky = isFirst || isLast;
                    return (
                      <Table.Cell
                        key={cell.id}
                        className={cn(
                          'px-4 py-2 h-[60px] min-h-[60px] border-gray-300',
                          {
                            'first:rounded-l-lg last:rounded-r-lg':
                              !isScrolltable,
                            'sticky bg-white z-[15]': isSticky && isScrolltable,
                            'left-0 shadow-right border-l before:absolute':
                              isFirst && isScrolltable,
                            'right-0 shadow-left border-r':
                              isLast && isScrolltable
                          }
                        )}
                        style={{
                          width: cell.column.columnDef.size,
                          maxWidth: cell.column.columnDef.maxSize
                        }}
                      >
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </Table.Cell>
                    );
                  })}
                </Table.Row>
              );
            })
          ) : (
            <Table.Row>
              <Table.Cell className="pt-32" colSpan={columns.length}>
                {emptyCollection}
              </Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table.Root>
    </div>
  );
};
