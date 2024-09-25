import { HTMLAttributes } from 'react';
import { TableHeaderProps } from '@types';
import { TableCommonType } from './root';
import TableHeaderItem from './table-header-cell';

const TableHeaderRow = <T,>({
  data = [],
  isSelectAllRows,
  sortingState,
  handleToggleSelectAllRows,
  spreadColumn,
  onSortingTable,
  ...props
}: TableHeaderProps<T> &
  TableCommonType &
  HTMLAttributes<HTMLTableSectionElement>) => {
  return (
    <thead {...props}>
      <tr>
        {data.map((header, index) => (
          <TableHeaderItem
            header={header}
            key={index}
            isSelectAllRows={isSelectAllRows}
            colIndex={index}
            sortingState={sortingState}
            handleToggleSelectAllRows={handleToggleSelectAllRows}
            spreadColumn={spreadColumn}
            onSortingTable={onSortingTable}
          />
        ))}
      </tr>
    </thead>
  );
};

export default TableHeaderRow;
