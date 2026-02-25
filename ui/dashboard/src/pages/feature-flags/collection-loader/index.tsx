import { memo, useCallback, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  IconArchiveOutlined,
  IconSaveAsFilled
} from 'react-icons-material-design';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { useScreen } from 'hooks';
import { compact } from 'lodash';
import { Feature, FeatureCountByStatus } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { CardCollection } from '../collection-layout/card-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import GridViewCollection from '../collection-layout/grid-view-collection';
import { FlagActionType, FlagFilters } from '../types';
import { useFetchFlags } from './use-fetch-flags';

const CollectionLoader = memo(
  ({
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
    const { fromMobileScreen } = useScreen();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { t } = useTranslation(['common', 'table']);
    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchFlags({
      ...filters,
      environmentId: currentEnvironment?.id
    });

    const { data: accountCollection } = useQueryAccounts({
      params: {
        organizationId: currentEnvironment?.organizationId,
        cursor: String(0)
      }
    });

    const { data: autoOpsCollection } = useQueryAutoOpsRules({
      params: {
        environmentId: currentEnvironment?.id,
        cursor: String(0)
      }
    });

    const { data: rolloutCollection } = useQueryRollouts({
      params: {
        environmentId: currentEnvironment?.id,
        cursor: String(0)
      }
    });

    const autoOpsRules = autoOpsCollection?.autoOpsRules || [];
    const rollouts = rolloutCollection?.progressiveRollouts || [];
    const accounts = accountCollection?.accounts || [];
    const features = collection?.features || [];
    const totalCount = Number(collection?.totalCount) || 0;

    const { searchOptions } = useSearchParams();
    const editable = hasEditable(consoleAccount!);

    const popoverOptions = useMemo(
      () =>
        compact([
          searchOptions.tab === 'ARCHIVED'
            ? {
                label: `${t('unarchive-flag')}`,
                icon: IconArchiveOutlined,
                value: 'UNARCHIVE'
              }
            : {
                label: `${t('archive-flag')}`,
                icon: IconArchiveOutlined,
                value: 'ARCHIVE'
              },
          {
            label: `${t('clone-flag')}`,
            icon: IconSaveAsFilled,
            value: 'CLONE'
          }
        ]),
      [searchOptions]
    );

    const handleGetMaintainerInfo = useCallback(
      (email: string) => {
        const existedAccount = accounts?.find(
          account => account.email === email
        );
        if (
          !existedAccount ||
          !existedAccount?.firstName ||
          !existedAccount?.lastName
        )
          return email;
        return `${existedAccount.firstName} ${existedAccount.lastName}`;
      },
      [accounts]
    );

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
            filters?.tags ??
            filters?.status ??
            filters?.hasFeatureFlagAsRule
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
      <TableListContent className="gap-y-6">
        {fromMobileScreen ? (
          <GridViewCollection
            filterTags={filters?.tags}
            autoOpsRules={autoOpsRules}
            popoverOptions={popoverOptions}
            currentEnvironment={currentEnvironment}
            handleGetMaintainerInfo={handleGetMaintainerInfo}
            editable={editable}
            rollouts={rollouts}
            data={features}
            onActions={onHandleActions}
            emptyState={emptyState}
            handleTagFilters={handleTagFilters}
          />
        ) : (
          <CardCollection
            isLoading={isLoading}
            emptyCollection={emptyState}
            filterTags={filters?.tags}
            autoOpsRules={autoOpsRules}
            rollouts={rollouts}
            currentEnvironment={currentEnvironment}
            popoverOptions={popoverOptions}
            handleGetMaintainerInfo={handleGetMaintainerInfo}
            editable={editable}
            data={features}
            onActions={onHandleActions}
            handleTagFilters={handleTagFilters}
          />
        )}

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
