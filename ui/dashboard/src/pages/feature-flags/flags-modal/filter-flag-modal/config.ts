import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useAccountsLoader } from 'hooks/use-accounts-loading-more';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { FlagFilters, StatusFilterType } from 'pages/feature-flags/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

// Matches `booleanOptions` in use-options (numeric 1 = yes / 0 = no).
const booleanOptions = () => [
  { value: 1, label: t('yes') },
  { value: 0, label: t('no') }
];

// Matches `flagStatusOptions` in use-options.
const statusOptions = () => [
  { value: StatusFilterType.NEVER_USED, label: t('never-used') },
  { value: StatusFilterType.RECEIVING_TRAFFIC, label: t('receiving-traffic') },
  { value: StatusFilterType.NO_RECENT_TRAFFIC, label: t('no-recent-traffic') }
];

/**
 * Boolean flag fields all share the same wiring: a yes/no dropdown whose stored
 * value is 1/0, mapped to/from a single boolean key on the page filters.
 */
const booleanField = (
  type: FilterTypes,
  labelKey: string,
  key: keyof FlagFilters
) => ({
  type,
  labelKey,
  valueKind: 'boolean' as const,
  emptyValue: '',
  useData: () => ({ options: booleanOptions() }),
  toFilter: (filterValue: unknown) => ({ [key]: !!filterValue }),
  fromFilter: (filters: Partial<FlagFilters>) => {
    const value = filters[key];
    return value === undefined ? undefined : value ? 1 : 0;
  }
});

export const flagFilterConfig: FilterModalConfig<FlagFilters> = {
  mode: 'multi',
  fields: [
    booleanField(
      FilterTypes.HAS_PREREQUISITES,
      'has-prerequisites',
      'hasPrerequisites'
    ),
    booleanField(
      FilterTypes.HAS_RULE,
      'has-flag-as-rule',
      'hasFeatureFlagAsRule'
    ),
    booleanField(
      FilterTypes.HAS_ACTIVE_AUTO_OPS,
      'has-active-auto-ops',
      'hasActiveAutoOps'
    ),
    booleanField(
      FilterTypes.HAS_FINISHED_AUTO_OPS,
      'has-finished-auto-ops',
      'hasFinishedAutoOps'
    ),
    booleanField(FilterTypes.HAS_EXPERIMENT, 'has-experiment', 'hasExperiment'),
    booleanField(FilterTypes.ENABLED, 'enabled', 'enabled'),
    {
      type: FilterTypes.TAGS,
      labelKey: 'tags',
      valueKind: 'searchable',
      emptyValue: [],
      useData: ({ enabled }) => {
        const { consoleAccount } = useAuth();
        const currentEnvironment = getCurrentEnvironment(consoleAccount!);
        const { data, isLoading } = useQueryTags({
          params: {
            cursor: String(0),
            environmentId: currentEnvironment.id,
            entityType: 'FEATURE_FLAG'
          },
          enabled
        });
        const tags = data?.tags || [];
        const options = tags.map(item => ({
          label: item.name,
          value: item.name
        }));
        return {
          options,
          isLoading,
          getLabel: value =>
            (Array.isArray(value) &&
              tags.length &&
              value
                .map(item => tags.find(tag => tag.name === item)?.name)
                .filter(Boolean)
                .join(', ')) ||
            ''
        };
      },
      toFilter: filterValue => ({ tags: filterValue as string[] }),
      fromFilter: filters => filters.tags
    },
    {
      type: FilterTypes.STATUS,
      labelKey: 'status',
      valueKind: 'enum',
      emptyValue: '',
      useData: () => {
        const options = statusOptions();
        return {
          options,
          getLabel: value =>
            options.find(item => item.value === value)?.label || ''
        };
      },
      toFilter: filterValue => ({ status: filterValue as StatusFilterType }),
      fromFilter: filters => filters.status
    },
    {
      type: FilterTypes.MAINTAINER,
      labelKey: 'maintainer',
      valueKind: 'searchable-paginated',
      emptyValue: '',
      useData: ({ enabled, value }) => {
        const { consoleAccount } = useAuth();
        const currentEnvironment = getCurrentEnvironment(consoleAccount!);
        const {
          emailOptions,
          isInitialLoading,
          isLoadingMore,
          isSearching,
          hasMore,
          loadMore,
          onSearchChange,
          getAccountLabel
        } = useAccountsLoader({
          organizationId: currentEnvironment.organizationId,
          environmentId: currentEnvironment.id,
          enabled,
          preloadEmails: typeof value === 'string' && value ? [value] : []
        });
        return {
          options: emailOptions,
          isLoading: isInitialLoading,
          hasMore,
          isLoadingMore,
          isSearching,
          loadMore,
          onSearchChange,
          getLabel: filterValue => getAccountLabel(filterValue as string)
        };
      },
      toFilter: filterValue => ({ maintainer: filterValue as string }),
      fromFilter: filters => filters.maintainer
    }
  ]
};
