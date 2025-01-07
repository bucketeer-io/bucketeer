import { SortingState } from '@tanstack/react-table';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { collection } from '../page-loader';
import { UserSegments, UserSegmentsActionsType, UserSegmentsFilters } from '../types';

const CollectionLoader = ({
  onAdd,
  filters,
  setFilters,
  onActionHandler
}: {
  onAdd?: () => void;
  filters: UserSegmentsFilters;
  setFilters: (values: Partial<UserSegmentsFilters>) => void;
  organizationIds?: string[];
  onActionHandler: (value: UserSegments, type: UserSegmentsActionsType) => void;
}) => {
  const columns = useColumns({ onActionHandler });
  const isLoading = false;
  const isError = false;
  
  const onSortingChangeHandler = (sorting: SortingState) => {
    const updateOrderBy =
      sorting.length > 0
        ? sortingListFields[sorting[0].id]
        : sortingListFields.default;

    setFilters({
      orderBy: updateOrderBy,
      orderDirection: sorting[0]?.desc ? 'DESC' : 'ASC'
    });
  };

  const userSegments = collection?.userSegments || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={userSegments}
      searchQuery={filters.searchQuery as string}
      onClear={() => setFilters({ searchQuery: '' })}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={() => {}} />
  ) : (
    <>
      <DataTable
        isLoading={isLoading}
        data={userSegments}
        columns={columns}
        onSortingChange={onSortingChangeHandler}
        emptyCollection={emptyState}
      />
      {totalCount > LIST_PAGE_SIZE && !isLoading && (
        <Pagination
          page={filters.page as number}
          totalCount={totalCount}
          onChange={page => setFilters({ page })}
        />
      )}
    </>
  );
};

export default CollectionLoader;
