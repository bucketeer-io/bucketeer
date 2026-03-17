import { memo } from 'react';
import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { useScreen } from 'hooks';
import { APIKey } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { CardCollection } from '../collection-layout/card-collection';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { APIKeyActionsType, APIKeysFilters } from '../types';
import { useFetchAPIKeys } from './use-fetch-apikey';

const CollectionLoader = memo(
  ({
    filters,
    setFilters,
    onAdd,
    onActions,
    onClearFilters
  }: {
    filters: APIKeysFilters;
    setFilters: (values: Partial<APIKeysFilters>) => void;
    onAdd: () => void;
    onActions: (item: APIKey, type: APIKeyActionsType) => void;
    onClearFilters: () => void;
  }) => {
    const columns = useColumns({ onActions });
    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { fromMobileScreen } = useScreen();

    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchAPIKeys({
      ...filters,
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

    const apiKeys = collection?.apiKeys || [];
    const totalCount = Number(collection?.totalCount) || 0;

    const emptyState = (
      <CollectionEmpty
        data={apiKeys}
        isFilter={!isNotEmpty(filters?.disabled ?? filters?.environmentIds)}
        searchQuery={filters.searchQuery}
        onClear={onClearFilters}
        empty={<EmptyCollection onAdd={onAdd} />}
      />
    );

    return isError ? (
      <PageLayout.ErrorState onRetry={refetch} />
    ) : (
      <div className="flex flex-col">
        {fromMobileScreen ? (
          <DataTable
            isLoading={isLoading}
            data={apiKeys}
            columns={columns}
            onSortingChange={onSortingChangeHandler}
            emptyCollection={emptyState}
          />
        ) : (
          <CardCollection
            isLoading={isLoading}
            emptyCollection={emptyState}
            data={apiKeys}
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
      </div>
    );
  }
);

export default CollectionLoader;
