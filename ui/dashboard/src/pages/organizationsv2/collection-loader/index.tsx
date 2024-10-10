// import { CollectionWrapper } from '~/elements/collection/collection-wrapper';
import { useState } from 'react';
// import { useNavigate, useLocation } from 'react-router-dom';
import { SortingState } from '@tanstack/react-table';
import { sortingListFields } from 'constants/collection';
import { OrderBy, OrderDirection } from '@types';
import PageLayout from 'elements/page-layout';
// import { getInfiniteCollectionData } from '~/utils/collection';
// import { EmptyCollection } from '../collection-layout/empty-collection';
import { ListCollection } from '../collection-layout/list-collection';
import { useFetchOrganizations } from './use-fetch-organizations';

// import type { OrganizationFilters } from '../types';

interface OrganizationParams {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
  archived: boolean;
}

const CollectionLoader = () => {
  // const navigate = useNavigate();
  // const { pathname } = useLocation();

  const [params, setParams] = useState<OrganizationParams>({
    orderBy: 'DEFAULT',
    orderDirection: 'ASC',
    searchKeyword: '',
    archived: false
  });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ ...params });

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

    setParams({
      ...params,
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
