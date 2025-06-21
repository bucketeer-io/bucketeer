import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, getEditorEnvironmentIds, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { Push } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useFetchTags } from 'pages/members/collection-loader';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { PushActionsType, PushFilters } from '../types';
import { useFetchPushes } from './use-fetch-pushes';

const CollectionLoader = ({
  filters,
  setFilters,
  onAdd,
  onActions,
  onClearFilters
}: {
  filters: PushFilters;
  setFilters: (values: Partial<PushFilters>) => void;
  onAdd: () => void;
  onActions: (item: Push, type: PushActionsType) => void;
  onClearFilters: () => void;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editorEnvironmentIds = getEditorEnvironmentIds(consoleAccount!);

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    organizationId: currentEnvironment.organizationId,
    entityType: 'FEATURE_FLAG'
  });

  const tagList = tagCollection?.tags || [];
  const columns = useColumns({ onActions, tags: tagList });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchPushes({
    ...filters,
    environmentIds: editorEnvironmentIds,
    organizationId: currentEnvironment.organizationId
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

  const pushes = collection?.pushes || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={pushes}
      isFilter={isNotEmpty(filters?.disabled)}
      searchQuery={filters.searchQuery}
      onClear={onClearFilters}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <TableListContent className="min-w-[1000px]">
      <DataTable
        isLoading={isLoading || isLoadingTags}
        data={pushes}
        columns={columns}
        onSortingChange={onSortingChangeHandler}
        emptyCollection={emptyState}
      />
      {!isLoading && (
        <Pagination
          page={filters.page}
          totalCount={totalCount}
          onChange={page => setFilters({ page })}
        />
      )}
    </TableListContent>
  );
};

export default CollectionLoader;
