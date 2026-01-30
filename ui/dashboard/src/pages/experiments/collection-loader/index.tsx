import { memo, useEffect } from 'react';
import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { useScreen } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { useSearchParams } from 'utils/search-params';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { CardCollection } from '../collection-layout/card-collection';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { ExperimentActionsType, ExperimentFilters } from '../types';
import { useFetchExperiments } from './use-fetch-experiment';

const CollectionLoader = memo(
  ({
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
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { searchOptions, onChangSearchParams } = useSearchParams();
    const { fromMobileScreen } = useScreen();
    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchExperiments({
      ...filters,
      environmentId: currentEnvironment.id
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
    useEffect(() => {
      if (!collection?.experiments?.length) return;

      const hasWaiting = collection.experiments.some(
        exp => exp.status === 'WAITING'
      );

      if (hasWaiting) {
        const intervalId = setInterval(() => {
          refetch();
        }, 60 * 1000);
        return () => {
          clearInterval(intervalId);
        };
      }
    }, [collection, refetch]);

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
              maintainer: undefined,
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
        {fromMobileScreen ? (
          <DataTable
            isLoading={isLoading}
            data={experiments}
            columns={columns}
            onSortingChange={onSortingChangeHandler}
            emptyCollection={emptyState}
          />
        ) : (
          <CardCollection
            data={experiments}
            isLoading={isLoading}
            emptyCollection={emptyState}
            onActions={onActions}
          />
        )}

        {!isLoading && (
          <Pagination
            page={filters.page}
            totalCount={totalCount}
            onChange={page => setFilters({ page })}
          />
        )}
      </TableListContent>
    );
  }
);

export default CollectionLoader;
