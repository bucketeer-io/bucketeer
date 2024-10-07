import { useState } from 'react';
import {
  type ColumnDef,
  type TableState,
  flexRender,
  getCoreRowModel,
  useReactTable,
  SortingState,
  getSortedRowModel
} from '@tanstack/react-table';
import { IconAngleDown, IconAngleUp, IconSorting } from '@icons';
import Table from 'components/tablev2';

export interface DataTableProps<TData, TValue> {
  data: TData[];
  columns: ColumnDef<TData, TValue>[];
  state?: Partial<TableState>;
  onRowClick?: (data: TData) => void;
}

export const DataTable = <TData, TValue>({
  data,
  columns,
  state,
  onRowClick
}: DataTableProps<TData, TValue>) => {
  const [sorting, setSorting] = useState<SortingState>([]);

  const table = useReactTable({
    data,
    columns,
    state: { ...state, sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualSorting: true
  });

  return (
    <Table.Root>
      <Table.Header>
        {table.getHeaderGroups().map(headerGroup => (
          <Table.Row key={headerGroup.id}>
            {headerGroup.headers.map(header => {
              console.log('test', header.column.columnDef.enableSorting);

              return (
                <Table.Head
                  key={header.id}
                  // align={header.column.columnDef.meta?.align}
                  // data-fit-content={header.column.columnDef.meta?.fitContent}
                  onClick={header.column.getToggleSortingHandler()}
                >
                  {header.isPlaceholder ? null : (
                    <div className="flex items-center gap-3">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                      {
                        {
                          asc: <IconAngleUp />,
                          desc: <IconAngleDown />,
                          false: header.column.columnDef.enableSorting !==
                            false && <IconSorting />
                        }[header.column.getIsSorted() as string]
                      }
                    </div>
                  )}
                </Table.Head>
              );
            })}
          </Table.Row>
        ))}
      </Table.Header>
      <Table.Body>
        {table.getRowModel().rows?.length ? (
          table.getRowModel().rows.map(row => (
            <Table.Row
              key={row.id}
              data-state={row.getIsSelected() && 'selected'}
              data-hoverable={!!onRowClick}
              onClick={() => onRowClick?.(row.original)}
            >
              {row.getVisibleCells().map(cell => (
                <Table.Cell
                  key={cell.id}
                  // align={cell.column.columnDef.meta?.align}
                  // data-fit-content={cell.column.columnDef.meta?.fitContent}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </Table.Cell>
              ))}
            </Table.Row>
          ))
        ) : (
          <Table.Row>
            <Table.Cell colSpan={columns.length} className="h-24 text-center">
              {`No results.`}
            </Table.Cell>
          </Table.Row>
        )}
      </Table.Body>
    </Table.Root>
  );
};
