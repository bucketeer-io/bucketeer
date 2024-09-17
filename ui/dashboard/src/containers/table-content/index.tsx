import { ReactNode } from 'react';
import EmptyData from 'containers/empty-data';
import { TableHeaders, TableRows } from '@types';
import { cn } from 'utils/style';
import Table from 'components/table';

export type TableContentProps = {
  headers: TableHeaders;
  rows: TableRows;
  emptyTitle: string;
  emptyDescription: string;
  emptyActions?: ReactNode;
  className?: string;
};

const TableContent = ({
  headers,
  rows,
  emptyTitle,
  emptyDescription,
  emptyActions,
  className
}: TableContentProps) => {
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
      />
    </div>
  );
};

export default TableContent;
