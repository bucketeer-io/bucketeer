import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { Account } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { MemberActionsType, MembersFilters } from '../types';
import { useFetchMembers } from './use-fetch-members';
import { useFetchTags } from './use-fetch-tags';

export * from './use-fetch-tags';

const CollectionLoader = ({
  filters,
  setFilters,
  onAdd,
  onActions,
  onClearFilters
}: {
  filters: MembersFilters;
  setFilters: (values: Partial<MembersFilters>) => void;
  onAdd?: () => void;
  onActions: (item: Account, type: MemberActionsType) => void;
  onClearFilters: () => void;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    organizationId: currentEnvironment.organizationId
  });
  const tagList = tagCollection?.tags || [];
  const columns = useColumns({ onActions, tags: tagList });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchMembers({
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

  const accounts = collection?.accounts || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={accounts}
      isFilter={isNotEmpty(
        filters?.disabled ?? filters?.organizationRole ?? filters?.tags
      )}
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
        data={accounts}
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
