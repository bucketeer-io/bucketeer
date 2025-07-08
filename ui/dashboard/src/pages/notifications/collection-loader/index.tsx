import { memo } from 'react';
import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, getEditorEnvironments, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { Notification } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { NotificationActionsType, NotificationFilters } from '../types';
import { useFetchNotifications } from './use-fetch-notifications';

const CollectionLoader = memo(
  ({
    filters,
    setFilters,
    onAdd,
    onActions,
    onClearFilters
  }: {
    filters: NotificationFilters;
    setFilters: (values: Partial<NotificationFilters>) => void;
    onAdd: () => void;
    onActions: (item: Notification, type: NotificationActionsType) => void;
    onClearFilters: () => void;
  }) => {
    const columns = useColumns({ onActions });
    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { editorEnvironmentIDs } = getEditorEnvironments(consoleAccount!);

    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchNotifications({
      ...filters,
      environmentIds: filters?.environmentIds || editorEnvironmentIDs,
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

    const apiKeys = collection?.subscriptions || [];
    const totalCount = Number(collection?.totalCount) || 0;

    const emptyState = (
      <CollectionEmpty
        data={apiKeys}
        isFilter={isNotEmpty(filters?.disabled ?? filters.environmentIds)}
        searchQuery={filters.searchQuery}
        onClear={onClearFilters}
        empty={<EmptyCollection onAdd={onAdd} />}
      />
    );

    return isError ? (
      <PageLayout.ErrorState onRetry={refetch} />
    ) : (
      <TableListContent>
        <DataTable
          isLoading={isLoading}
          data={apiKeys}
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
  }
);

export default CollectionLoader;
