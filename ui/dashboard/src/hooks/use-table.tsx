import { ReactNode } from 'react';
import {
  Cell,
  Column,
  ColumnDef,
  getCoreRowModel,
  Header,
  PaginationState,
  useReactTable
} from '@tanstack/react-table';
import { OrderBy, TableRowItemProps, TableRowItemType } from '@types';
import { PopoverOption, PopoverValue } from 'components/popover';

export type ColumnType<T> = Omit<
  ColumnDef<T>,
  'size' | 'minSize' | 'header' | 'cell' | 'sortingFn' | 'enableSorting'
> &
  TableRowItemProps<T> & {
    cell?: Cell<T, unknown>;
    header?: string;
    headerCellType?: 'title' | 'checkbox' | 'empty';
    size?: number | string;
    minSize?: number | string;
    cellType?: TableRowItemType;
    options?: PopoverOption<PopoverValue>[];
    accessorKey?: string;
    sorting?: boolean;
    sortingKey?: OrderBy;
    renderFunc?: (row?: T) => ReactNode;
  };

export type SpreadColumn<T> = {
  column: Column<T, unknown>;
  columnDef: ColumnType<T>;
};

export type ColumnData<T> = Header<T, unknown> | Cell<T, unknown>;

type UseTableProps<T> = {
  data: T[];
  columns: ColumnType<T>[];
  paginationProps?: PaginationState;
};

const useTable = <T,>({ data, columns }: UseTableProps<T>) => {
  const table = useReactTable({
    data,
    columns: columns as ColumnDef<T>[],
    getCoreRowModel: getCoreRowModel()
  });

  const spreadColumn = <T,>(data: ColumnData<T>): SpreadColumn<T> => {
    const { column } = data;
    const { columnDef } = column;

    return { column, columnDef: columnDef as ColumnType<T> };
  };

  return { table, spreadColumn };
};

export default useTable;
