import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { useSearchParams } from 'utils/search-params';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
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
  setFilters: (
    values: Partial<ExperimentFilters>,
    isChangeParams?: boolean
  ) => void;
  onAdd: () => void;
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
}) => {
  const { t } = useTranslation(['message']);
  const columns = useColumns({ onActions });
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions, onChangSearchParams } = useSearchParams();
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
      isFilter={
        filters.isFilter ||
        (!!searchOptions?.statuses?.length && !filters?.filterByTab)
      }
      description={t('message:empty:experiment-match')}
      onClear={() => {
        setFilters(
          {
            searchQuery: '',
            isFilter: undefined,
            status: 'ACTIVE',
            statuses: ['WAITING', 'RUNNING']
          },
          false
        );
        onChangSearchParams({});
      }}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <TableListContent>
      <DataTable
        isLoading={isLoading}
        data={experiments}
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
