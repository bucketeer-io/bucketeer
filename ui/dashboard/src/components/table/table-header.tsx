import { TableHeaderProps } from '@types';
import TableHeaderItem from './table-header-item';

const TableHeader = ({
  data = [],
  isSelectAllRows,
  sortedObj,
  handleToggleSelectAllRows,
  handleSortedData
}: TableHeaderProps) => {
  return (
    <tr>
      {data.map((i, index) => (
        <TableHeaderItem
          {...i}
          key={index}
          isSelectAllRows={isSelectAllRows}
          colIndex={index}
          sortedObj={sortedObj}
          handleToggleSelectAllRows={handleToggleSelectAllRows}
          handleSortedData={handleSortedData}
        />
      ))}
    </tr>
  );
};

export default TableHeader;
