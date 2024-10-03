// import { CollectionWrapper } from '~/elements/collection/collection-wrapper';
import { OrganizationsFetcherParams } from '@api/organization';
import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import PageLayout from 'elements/page-layout';
// import { getInfiniteCollectionData } from '~/utils/collection';
// import { EmptyCollection } from '../collection-layout/empty-collection';
import { ListCollection } from '../collection-layout/list-collection';

// import type { CompaniesFilters } from '../types';

const CollectionLoader = () =>
  // 	{
  //   filters,
  //   setFilters,
  //   onAdd
  // }: {
  //   filters: CompaniesFilters;
  //   setFilters: (values: Partial<CompaniesFilters>) => void;
  //   onAdd: () => void;
  // }
  {
    const defaultParams: OrganizationsFetcherParams = {
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      archived: false
    };

    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useQueryOrganizations({
      params: defaultParams
    });

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
          <ListCollection organizations={organizations} />
          // </CollectionWrapper>
        )}
      </>
    );
  };

export default CollectionLoader;
