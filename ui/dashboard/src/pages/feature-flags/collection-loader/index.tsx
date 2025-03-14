import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { Feature } from '@types';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from '../collection-layout/empty-collection';
import GridViewCollection from '../collection-layout/grid-view-collection';
import { FlagActionType, FlagFilters } from '../types';
import { useFetchFlags } from './use-fetch-flags';

const CollectionLoader = ({
  onAdd,
  filters,
  setFilters,
  onHandleActions
}: {
  onAdd: () => void;
  filters: FlagFilters;
  setFilters: (filters: Partial<FlagFilters>) => void;
  onHandleActions: (item: Feature, type: FlagActionType) => void;
}) => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchFlags({
    ...filters,
    environmentId: currenEnvironment?.id
  });

  const features = collection?.features || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const emptyState = (
    <CollectionEmpty
      data={features}
      searchQuery={filters?.searchQuery}
      onClear={() => {}}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );
  return isLoading ? (
    <PageLayout.LoadingState />
  ) : isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <>
      <GridViewCollection
        data={features}
        onActions={onHandleActions}
        emptyState={emptyState}
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
