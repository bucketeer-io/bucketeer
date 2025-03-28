import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import { Organization } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { OrganizationActionsType, OrganizationFilters } from '../types';
import { useFetchOrganizations } from './use-fetch-organizations';

const CollectionLoader = ({
  onAdd,
  filters,
  setFilters,
  onActions
}: {
  onAdd: () => void;
  filters: OrganizationFilters;
  setFilters: (values: Partial<OrganizationFilters>) => void;
  onActions: (item: Organization, type: OrganizationActionsType) => void;
}) => {
  const columns = useColumns({ onActions });
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ ...filters });

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

  const organizations = collection?.Organizations || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={organizations}
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
        data={organizations}
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
