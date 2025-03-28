import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import { Account } from '@types';
import { OrganizationMembersFilters } from 'pages/organization-details/types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { useFetchMembers } from './use-fetch-members';

const CollectionLoader = ({
  filters,
  setFilters,
  onActions
}: {
  filters: OrganizationMembersFilters;
  setFilters: (values: Partial<OrganizationMembersFilters>) => void;
  onActions: (value: Account) => void;
}) => {
  const columns = useColumns({ onActions });
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchMembers({ ...filters });

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
      searchQuery={filters.searchQuery}
      onClear={() => setFilters({ searchQuery: '' })}
      empty={<EmptyCollection />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <>
      <DataTable
        isLoading={isLoading}
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
    </>
  );
};

export default CollectionLoader;
