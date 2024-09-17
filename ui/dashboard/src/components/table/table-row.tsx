import { TableRowProps } from '@types';
import Card from 'components/card';
import TableRowItem from './table-row-item';

const TableRow = ({
  data = [],
  rowIndex,
  rowsSelected,
  tableRows,
  handleSelectRow
}: TableRowProps) => {
  return (
    <Card tag="tr" className="relative">
      {data.map((i, index) => (
        <TableRowItem
          {...i}
          key={index}
          rowIndex={rowIndex}
          rowsSelected={rowsSelected}
          handleSelectRow={handleSelectRow}
          tableRows={tableRows}
        />
      ))}
    </Card>
  );
};

export default TableRow;
