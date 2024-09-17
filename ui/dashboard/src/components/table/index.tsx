import { useCallback, useEffect, useRef, useState } from 'react';
import { TableProps } from '@types';
import Pagination from 'components/pagination';
import ElementOnEmpty from './element-on-empty';
import TableHeader from './table-header';
import TableHeaderItem from './table-header-item';
import TableRow from './table-row';
import TableRowItem from './table-row-item';
import Flag from './table-row-items/flag';
import Text from './table-row-items/text';
import Title from './table-row-items/title';

const Table = ({ headers, rows = [], elementEmpty }: TableProps) => {
  const [isSelectAllRows, setIsSelectAllRows] = useState(false);
  const [rowsSelected, setRowsSelected] = useState<number[]>([]);

  const initLoadedRef = useRef(true);

  const handleSelectRow = useCallback(
    (rowIndex?: number) => {
      if (typeof rowIndex !== 'number') return;
      if (rowsSelected.includes(rowIndex))
        return setRowsSelected(rowsSelected.filter(item => item !== rowIndex));

      setRowsSelected([...rowsSelected, rowIndex]);
    },
    [rowsSelected]
  );

  const handleToggleSelectAllRows = () => {
    setIsSelectAllRows(!isSelectAllRows);
    initLoadedRef.current = false;
  };

  useEffect(() => {
    if (!initLoadedRef.current) {
      if (!isSelectAllRows) return setRowsSelected([]);
      setRowsSelected(rows.map((_, index) => index));
    }
  }, [isSelectAllRows]);

  return (
    <div>
      <table className="border-separate border-spacing-y-3 w-full mb-6">
        <thead>
          <TableHeader
            data={headers}
            isSelectAllRows={isSelectAllRows}
            handleToggleSelectAllRows={handleToggleSelectAllRows}
            handleSortedData={() => {}}
          />
        </thead>
        <tbody>
          {rows.length > 0 &&
            rows.map((i, index) => (
              <TableRow
                key={index}
                data={i}
                rowIndex={index}
                rowsSelected={rowsSelected}
                tableRows={rows}
                handleSelectRow={handleSelectRow}
              />
            ))}
        </tbody>
      </table>
      {!rows.length && <ElementOnEmpty>{elementEmpty}</ElementOnEmpty>}
      {rows.length > 0 && <Pagination totalItems={50} itemsPerPage={5} />}
    </div>
  );
};

Table.Header = TableHeader;
Table.HeaderItem = TableHeaderItem;
Table.Row = TableRow;
Table.RowItem = TableRowItem;
Table.Text = Text;
Table.ItemTitle = Title;
Table.Flag = Flag;

export default Table;
