import { useCallback, useEffect } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryAutoOps } from '@queries/auto-ops';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Feature, FeatureCountByStatus } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from '../collection-layout/empty-collection';
import GridViewCollection from '../collection-layout/grid-view-collection';
import { FlagActionType, FlagFilters } from '../types';
import { useFetchFlags } from './use-fetch-flags';

const CollectionLoader = ({
  filters,
  onAdd,
  setFilters,
  setSummary,
  onHandleActions,
  onClearFilters
}: {
  filters: FlagFilters;
  onAdd: () => void;
  setFilters: (filters: Partial<FlagFilters>) => void;
  setSummary: (summary: FeatureCountByStatus) => void;
  onHandleActions: (item: Feature, type: FlagActionType) => void;
  onClearFilters: () => void;
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

  const { data: accountCollection } = useQueryAccounts({
    params: {
      organizationId: currenEnvironment?.organizationId,
      cursor: String(0)
    }
  });

  const { data: autoOpsCollection } = useQueryAutoOps({
    params: {
      environmentId: currenEnvironment?.id,
      cursor: String(0)
    }
  });

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      environmentId: currenEnvironment?.id,
      cursor: String(0)
    }
  });

  const autoOpsRules = autoOpsCollection?.autoOpsRules || [];
  const rollouts = rolloutCollection?.progressiveRollouts || [];
  const accounts = accountCollection?.accounts || [];
  const features = collection?.features || [];
  const totalCount = Number(collection?.totalCount) || 0;

  const handleTagFilters = useCallback(
    (tag: string) => {
      const tags = filters?.tags as string[];
      const isNotEmptyTag = isNotEmpty(tags);
      if (isNotEmptyTag) {
        const isExistedTag = tags.includes(tag);
        const _tags = isExistedTag
          ? tags.filter(item => item !== tag)
          : [...tags, tag];
        return setFilters({
          ...filters,
          tags: _tags.length ? _tags : undefined
        });
      }

      setFilters({
        ...filters,
        tags: [tag]
      });
    },
    [filters]
  );

  const emptyState = (
    <CollectionEmpty
      data={features}
      isFilter={isNotEmpty(
        filters?.enabled ??
          filters?.hasExperiment ??
          filters?.hasPrerequisites ??
          filters?.maintainer ??
          filters?.tags
      )}
      searchQuery={filters?.searchQuery}
      onClear={onClearFilters}
      empty={<EmptyCollection onAdd={onAdd} />}
    />
  );

  useEffect(() => {
    if (collection) {
      setSummary(collection.featureCountByStatus);
    }
  }, [collection]);

  return isLoading ? (
    <PageLayout.LoadingState />
  ) : isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <div className="flex flex-col gap-y-6">
      <GridViewCollection
        filterTags={filters?.tags}
        autoOpsRules={autoOpsRules}
        rollouts={rollouts}
        accounts={accounts}
        data={features}
        onActions={onHandleActions}
        emptyState={emptyState}
        handleTagFilters={handleTagFilters}
      />

      {!isLoading && (
        <Pagination
          page={filters.page}
          totalCount={totalCount}
          onChange={page => setFilters({ page })}
        />
      )}
    </div>
  );
};

export default CollectionLoader;
