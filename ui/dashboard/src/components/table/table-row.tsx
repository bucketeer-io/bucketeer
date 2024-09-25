import { TableRowProps } from '@types';
import Card from 'components/card';
import TableRowItem from './table-row-item';

const TableRow = <T,>({
  row,
  rowsSelected,
  handleSelectRow,
  spreadColumn
}: TableRowProps<T>) => {
  return (
    <Card tag="tr" className="relative">
      {row.getVisibleCells().map(cell => {
        const { columnDef } = spreadColumn(cell);

        return (
          <TableRowItem
            {...columnDef}
            cell={cell}
            key={cell.id}
            rowId={row.id}
            rowsSelected={rowsSelected}
            handleSelectRow={handleSelectRow}
            spreadColumn={spreadColumn}
          />
        );
      })}
    </Card>
  );
};

export default TableRow;
