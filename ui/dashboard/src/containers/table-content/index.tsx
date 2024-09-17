import { ReactNode } from 'react';
import EmptyData from 'containers/empty-data';
import { TableHeaders, TableRows, TableSignature } from '@types';
import { cn } from 'utils/style';
import Table from 'components/table';

export type TableContentProps<T> = {
  headers: TableHeaders;
  rows: TableRows;
  emptyTitle: string;
  emptyDescription: string;
  emptyActions?: ReactNode;
  className?: string;
  originalData: T[];
  rowsData: T[];
  setRowsData: (data: T[]) => void;
};

const TableContent = <T extends TableSignature>({
  headers,
  rows,
  emptyTitle,
  emptyDescription,
  emptyActions,
  className,
  originalData,
  rowsData,
  setRowsData
}: TableContentProps<T>) => {
  return (
    <div className={cn('grid gap-6 mt-6', className)}>
      <Table
        headers={headers}
        rows={rows}
        elementEmpty={
          <EmptyData
            title={emptyTitle}
            description={emptyDescription}
            emptyActions={emptyActions}
          />
        }
        originalData={originalData}
        rowsData={rowsData}
        setRowsData={setRowsData}
      />
    </div>
  );
};

export default TableContent;
