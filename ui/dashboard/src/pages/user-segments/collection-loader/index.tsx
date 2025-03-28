import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import { UserSegment } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { UserSegmentsActionsType, UserSegmentsFilters } from '../types';
import { useFetchSegments } from './use-fetch-segment';

const CollectionLoader = ({
  segmentUploading,
  onAdd,
  filters,
  setFilters,
  onActionHandler
}: {
  segmentUploading: UserSegment | null;
  onAdd?: () => void;
  filters: UserSegmentsFilters;
  setFilters: (values: Partial<UserSegmentsFilters>) => void;
  organizationIds?: string[];
  onActionHandler: (value: UserSegment, type: UserSegmentsActionsType) => void;
}) => {
  const columns = useColumns({ segmentUploading, onActionHandler });
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchSegments({
    ...filters,
    environmentId: currenEnvironment.id
  });

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

  const userSegments = collection?.segments || [];
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
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <TableListContent>
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
    </TableListContent>
  );
};

export default CollectionLoader;
