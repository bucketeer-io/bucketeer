import Filter from 'containers/filter';
import TableContent, { TableContentProps } from 'containers/table-content';

type Props<T> = TableContentProps<T> & {
  isLoading?: boolean;
  searchValue: string;
  onSearchChange: (value: string) => void;
};

const FilterLayout = <T,>({
  columns,
  data,
  emptyTitle,
  emptyDescription,
  paginationProps,
  isLoading,
  sortingState,
  searchValue,
  onSortingTable,
  onSearchChange
}: Props<T>) => {
  return (
    <div>
      <Filter searchValue={searchValue} onSearchChange={onSearchChange} />
      <TableContent
        isLoading={isLoading}
        columns={columns}
        data={data}
        paginationProps={paginationProps}
        emptyTitle={emptyTitle}
        emptyDescription={emptyDescription}
        sortingState={sortingState}
        onSortingTable={onSortingTable}
      />
    </div>
  );
};

export default FilterLayout;
