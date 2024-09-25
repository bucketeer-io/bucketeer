import { useCallback, useMemo, useRef, useState } from 'react';
import { TableProps } from '@types';
import Pagination from 'components/pagination';
import ElementOnEmpty from './element-on-empty';
import TableRoot from './root';
import TableBody from './table-body';
import TableHeader from './table-header';
import TableHeaderRow from './table-header';
import TableHeaderItem from './table-header-cell';
import TableRow from './table-row';
import TableRowItem from './table-row-item';
import Flag from './table-row-items/flag';
import Text from './table-row-items/text';
import Title from './table-row-items/title';

const Table = <T,>({
  table,
  elementEmpty,
  paginationProps,
  rowsSelected = [],
  sortingState,
  spreadColumn,
  setRowsSelected,
  onSortingTable
}: TableProps<T>) => {
  const [isSelectAllRows, setIsSelectAllRows] = useState(false);
  const initLoadedRef = useRef(true);

  const tableRows = useMemo(
    () => table?.getRowModel()?.rows,
    [table, paginationProps]
  );

  const handleSelectRow = useCallback(
    (rowId?: string) => {
      if (!rowId) return;
      initLoadedRef.current = false;
      const newRows = rowsSelected?.includes(rowId)
        ? rowsSelected.filter(item => item !== rowId)
        : [...rowsSelected, rowId];
      if (setRowsSelected) setRowsSelected(newRows);
      setIsSelectAllRows(newRows.length === tableRows.length ? true : false);
    },
    [rowsSelected, tableRows, setRowsSelected]
  );

  const handleToggleSelectAllRows = useCallback(() => {
    setIsSelectAllRows(!isSelectAllRows);
    if (setRowsSelected) {
      if (isSelectAllRows) {
        return setRowsSelected([]);
      }
      return setRowsSelected(tableRows.map(row => row.id));
    }
  }, [isSelectAllRows, setRowsSelected]);

  return (
    <div className="w-full">
      <TableRoot>
        {tableRows?.length > 0 &&
          table
            ?.getHeaderGroups()
            ?.map(headerGroup => (
              <TableHeaderRow
                sortingState={sortingState}
                key={headerGroup.id}
                data={headerGroup.headers}
                isSelectAllRows={isSelectAllRows}
                handleToggleSelectAllRows={handleToggleSelectAllRows}
                spreadColumn={spreadColumn}
                onSortingTable={onSortingTable}
              />
            ))}
        <TableBody>
          {tableRows?.length > 0 &&
            tableRows.map((row, index) => (
              <TableRow
                key={index}
                row={row}
                rowsSelected={rowsSelected}
                spreadColumn={spreadColumn}
                handleSelectRow={handleSelectRow}
              />
            ))}
        </TableBody>
      </TableRoot>
      {tableRows?.length > 0 &&
        paginationProps &&
        paginationProps?.totalCount > paginationProps?.pageSize && (
          <Pagination paginationProps={paginationProps} />
        )}
      {!tableRows?.length && <ElementOnEmpty>{elementEmpty}</ElementOnEmpty>}
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
