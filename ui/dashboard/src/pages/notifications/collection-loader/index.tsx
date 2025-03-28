import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import { Notification } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { NotificationActionsType, NotificationFilters } from '../types';
import { useFetchNotifications } from './use-fetch-notifications';

const CollectionLoader = ({
  filters,
  setFilters,
  onAdd,
  onActions
}: {
  filters: NotificationFilters;
  setFilters: (values: Partial<NotificationFilters>) => void;
  onAdd: () => void;
  onActions: (item: Notification, type: NotificationActionsType) => void;
}) => {
  const columns = useColumns({ onActions });
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchNotifications({
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

  const apiKeys = collection?.subscriptions || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={apiKeys}
      searchQuery={filters.searchQuery}
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
        data={apiKeys}
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
    </TableListContent>
  );
};

export default CollectionLoader;
