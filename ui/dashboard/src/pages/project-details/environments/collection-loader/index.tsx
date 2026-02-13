import { memo } from 'react';
import { useParams } from 'react-router-dom';
import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import { useScreen } from 'hooks';
import { Environment } from '@types';
import { useSearchParams } from 'utils/search-params';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { CardCollection } from '../collection-layout/card-collection';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { EnvironmentActionsType, EnvironmentFilters } from '../types';
import { useFetchEnvironments } from './use-fetch-environments';

const CollectionLoader = memo(
  ({
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
    const { searchOptions } = useSearchParams();
    const { fromMobileScreen } = useScreen();
    const organizationId = searchOptions?.organizationId as string;

    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchEnvironments({
      ...filters,
      projectId,
      organizationId,
      enabledQuery: !!projectId && !!organizationId
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
      <TableListContent>
        {fromMobileScreen ? (
          <DataTable
            isLoading={isLoading}
            data={environments}
            columns={columns}
            onSortingChange={onSortingChangeHandler}
            emptyCollection={emptyState}
          />
        ) : (
          <CardCollection
            data={environments}
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
