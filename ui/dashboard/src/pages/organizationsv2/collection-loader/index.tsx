// import { CollectionWrapper } from '~/elements/collection/collection-wrapper';
import { useState } from 'react';
// import { useNavigate, useLocation } from 'react-router-dom';
import { OrganizationsFetcherParams } from '@api/organization';
import { useQueryOrganizations } from '@queries/organizations';
import { SortingState } from '@tanstack/react-table';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';
import { orderDirectionType, sortingListFields } from 'utils/collection';
import PageLayout from 'elements/page-layout';
// import { getInfiniteCollectionData } from '~/utils/collection';
// import { EmptyCollection } from '../collection-layout/empty-collection';
import { ListCollection } from '../collection-layout/list-collection';

// import type { OrganizationFilters } from '../types';

// interface CollectionLoaderProps {
//   filters?: OrganizationFilters;
//   setFilters?: (values: Partial<OrganizationFilters>) => void;
//   onAdd?: () => void;
// }

interface OrganizationParams {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword: string;
  archived: boolean;
}

const CollectionLoader = () => {
  // const navigate = useNavigate();
  // const { pathname } = useLocation();

  const [params, seParams] = useState<OrganizationParams>({
    orderBy: sortingListFields.default,
    orderDirection: orderDirectionType.asc,
    searchKeyword: '',
    archived: false
  });

  const defaultParams: OrganizationsFetcherParams = {
    pageSize: LIST_PAGE_SIZE,
    cursor: String(0),
    disabled: false,
    ...params
  };

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useQueryOrganizations({
    params: defaultParams
  });

  // const onUpdateURL = useCallback(
  //   (options: Record<string, string | number | boolean | undefined>) => {
  //     navigate(`${pathname}?${stringifySearchParams(options)}`, {
  //       replace: true
  //     });
  //   },
  //   [navigate]
  // );

  const onSortingChangeHandler = (sorting: SortingState) => {
    if (sorting.length > 0) {
      seParams({
        ...params,
        orderBy: sortingListFields[sorting[0].id],
        orderDirection: sorting[0].desc ? 'DESC' : 'ASC'
      });
    } else {
      seParams({
        ...params,
        orderBy: sortingListFields.default,
        orderDirection: orderDirectionType.asc
      });
    }
  };

  const organizations = collection?.Organizations || [];

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
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
          onSortingChange={onSortingChangeHandler}
        />
        // </CollectionWrapper>
      )}
    </>
  );
};

export default CollectionLoader;
