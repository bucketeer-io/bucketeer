import { ReactNode } from 'react';
import EmptyData from 'containers/empty-data';
import { SortingType } from 'containers/pages';
import useTable, { ColumnType } from 'hooks/use-table';
import { cn } from 'utils/style';
import { PaginationProps } from 'components/pagination';
import Spinner from 'components/spinner';
import Table from 'components/table';

export type TableContentProps<T> = {
  columns: ColumnType<T>[];
  data?: T[];
  emptyTitle: string;
  emptyDescription: string;
  emptyActions?: ReactNode;
  className?: string;
  paginationProps?: PaginationProps;
  rowsSelected?: string[];
  isLoading?: boolean;
  sortingState?: SortingType;
  setRowsSelected?: (rows: string[]) => void;
  onSortingTable?: (accessorKey: string) => void;
};

const TableContent = <T,>({
  columns,
  data,
  emptyTitle,
  emptyDescription,
  emptyActions,
  className,
  paginationProps,
  rowsSelected,
  isLoading,
  sortingState,
  onSortingTable,
  setRowsSelected
}: TableContentProps<T>) => {
  const { table, spreadColumn } = useTable({
    columns,
    data: data || []
  });

  return (
    <div className={cn('grid gap-6 mt-6', className)}>
      {isLoading ? (
        <div className="pt-20 flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <Table
          table={table}
          paginationProps={paginationProps}
          rowsSelected={rowsSelected}
          sortingState={sortingState}
          elementEmpty={
            <EmptyData
              title={emptyTitle}
              description={emptyDescription}
              emptyActions={emptyActions}
            />
          }
          onSortingTable={onSortingTable}
          spreadColumn={spreadColumn}
          setRowsSelected={setRowsSelected}
        />
      )}
    </div>
  );
};

export default TableContent;
