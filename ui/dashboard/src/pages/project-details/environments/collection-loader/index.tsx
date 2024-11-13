import { useParams } from 'react-router-dom';
import { SortingState } from '@tanstack/react-table';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import { Environment } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { EnvironmentActionsType, EnvironmentFilters } from '../types';
import { useFetchEnvironments } from './use-fetch-environments';

const CollectionLoader = ({
  onAdd,
  filters,
  setFilters,
  onActions
}: {
  onAdd?: () => void;
  filters: EnvironmentFilters;
  setFilters: (values: Partial<EnvironmentFilters>) => void;
  onActions: (item: Environment, type: EnvironmentActionsType) => void;
}) => {
  const { projectId } = useParams();
  const columns = useColumns({ onActions });
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchEnvironments({ ...filters, projectId });

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

  const environments = collection?.environments || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={environments}
      searchQuery={filters.searchQuery}
      onClear={() => setFilters({ searchQuery: '' })}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <>
      <DataTable
        isLoading={isLoading}
        data={environments}
        columns={columns}
        onSortingChange={onSortingChangeHandler}
        emptyCollection={emptyState}
      />
      {totalCount > LIST_PAGE_SIZE && !isLoading && (
        <Pagination
          page={filters.page}
          totalCount={totalCount}
          onChange={page => setFilters({ page })}
        />
      )}
    </>
  );
};

export default CollectionLoader;
