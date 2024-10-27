import { SortingState } from '@tanstack/react-table';
import { LIST_PAGE_SIZE } from 'constants/app';
import { sortingListFields } from 'constants/collection';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { ProjectsFilters } from '../types';
import { useFetchProjects } from './use-fetch-projects';

const CollectionLoader = ({
  onAdd,
  filters,
  setFilters
}: {
  onAdd?: () => void;
  filters: ProjectsFilters;
  setFilters: (values: Partial<ProjectsFilters>) => void;
}) => {
  const columns = useColumns();
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchProjects({ ...filters });

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

  const projects = collection?.projects || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={projects}
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
        data={projects}
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
    </>
  );
};

export default CollectionLoader;
