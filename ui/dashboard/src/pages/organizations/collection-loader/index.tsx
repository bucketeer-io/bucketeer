import { memo } from 'react';
import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import { Organization } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { OrganizationActionsType, OrganizationFilters } from '../types';
import { useFetchOrganizations } from './use-fetch-organizations';

const CollectionLoader = memo(
  ({
    onAdd,
    filters,
    setFilters,
    onActions,
    onClearFilters
  }: {
    onAdd: () => void;
    filters: OrganizationFilters;
    setFilters: (values: Partial<OrganizationFilters>) => void;
    onActions: (item: Organization, type: OrganizationActionsType) => void;
    onClearFilters: () => void;
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
        isFilter={isNotEmpty(filters.disabled)}
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
      </TableListContent>
    );
  }
);

export default CollectionLoader;
