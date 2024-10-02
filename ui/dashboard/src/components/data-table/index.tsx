import {
  type ColumnDef,
  type TableState,
  flexRender,
  getCoreRowModel,
  useReactTable
} from '@tanstack/react-table';
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
  const table = useReactTable({
    data,
    columns,
    state,
    getCoreRowModel: getCoreRowModel()
  });

  return (
    <Table.Root>
      <Table.Header>
        {table.getHeaderGroups().map(headerGroup => (
          <Table.Row key={headerGroup.id}>
            {headerGroup.headers.map(header => {
              return (
                <Table.Head
                  key={header.id}
                  //   align={header.column.columnDef.meta?.align}
                  //   data-fit-content={header.column.columnDef.meta?.fitContent}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
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
                  //   align={cell.column.columnDef.meta?.align}
                  //   data-fit-content={cell.column.columnDef.meta?.fitContent}
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
