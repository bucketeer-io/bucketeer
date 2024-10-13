// import { CollectionWrapper } from '~/elements/collection/collection-wrapper';
// import { useNavigate, useLocation } from 'react-router-dom';
import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import PageLayout from 'elements/page-layout';
// import { getInfiniteCollectionData } from '~/utils/collection';
// import { EmptyCollection } from '../collection-layout/empty-collection';
import { ListCollection } from '../collection-layout/list-collection';
import { OrganizationFilters } from '../types';
import { useFetchOrganizations } from './use-fetch-organizations';

// import type { OrganizationFilters } from '../types';

const CollectionLoader = ({
  filters,
  setFilters
}: {
  filters: OrganizationFilters;
  setFilters: (values: Partial<OrganizationFilters>) => void;
}) => {
  // const navigate = useNavigate();
  // const { pathname } = useLocation();

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ ...filters });

  // const onUpdateURL = useCallback(
  //   (options: Record<string, string | number | boolean | undefined>) => {
  //     navigate(`${pathname}?${stringifySearchParams(options)}`, {
  //       replace: true
  //     });
  //   },
  //   [navigate]
  // );

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

  return (
    <>
      {isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        // <CollectionWrapper
        // 	items={companies}
        // 	empty={<EmptyCollection onAdd={onAdd} />}
        // 	filtersTypes={COMPANIES_FILTERS_TYPES}
        // 	searchQuery={filters.searchQuery}
        // 	onClear={() => setFilters({ searchQuery: '' })}
        // 	infiniteLoadMore
        // 	isLoadingMore={isFetchingNextPage}
        // 	canLoadMore={!!hasNextPage}
        // 	onLoadMore={fetchNextPage}
        // >
        <ListCollection
          organizations={organizations}
          isLoading={isLoading}
          onSortingChange={onSortingChangeHandler}
        />
        // </CollectionWrapper>
      )}
    </>
  );
};

export default CollectionLoader;
