import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import { Experiment } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { ExperimentActionsType, ExperimentFilters } from '../types';
import { useFetchExperiments } from './use-fetch-experiment';

const CollectionLoader = ({
  filters,
  setFilters,
  onAdd,
  onActions
}: {
  filters: ExperimentFilters;
  setFilters: (values: Partial<ExperimentFilters>) => void;
  onAdd: () => void;
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
}) => {
  const columns = useColumns({ onActions });
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchExperiments({
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

  const experiments = collection?.experiments || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={experiments}
      searchQuery={filters.searchQuery}
      isFilter={filters.isFilter}
      description="No experiments match your search filters. Try changing your filters."
      buttonText="Clear Filters"
      buttonVariant={'secondary'}
      onClear={() =>
        setFilters({
          searchQuery: '',
          isFilter: undefined,
          status: 'ACTIVE',
          statuses: ['WAITING', 'RUNNING']
        })
      }
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <>
      <DataTable
        isLoading={isLoading}
        data={experiments}
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
