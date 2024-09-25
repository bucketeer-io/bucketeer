import { KeyboardEvent } from 'react';
import Filter from 'containers/filter';
import TableContent, { TableContentProps } from 'containers/table-content';

type Props<T> = TableContentProps<T> & {
  isLoading?: boolean;
  searchValue: string;
  onChangeSearchValue: (value: string) => void;
  onKeyDown?: (e: KeyboardEvent<HTMLInputElement>) => void;
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
  onChangeSearchValue,
  onKeyDown
}: Props<T>) => {
  return (
    <div>
      <Filter
        searchValue={searchValue}
        onChangeSearchValue={onChangeSearchValue}
        onKeyDown={onKeyDown}
      />
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
