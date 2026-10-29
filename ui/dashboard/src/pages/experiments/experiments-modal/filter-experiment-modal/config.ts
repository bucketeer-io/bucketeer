import { getCurrentEnvironment, useAuth } from 'auth';
import { useAccountsLoader } from 'hooks/use-accounts-loading-more';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { ExperimentStatus } from '@types';
import { ExperimentFilters } from 'pages/experiments/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string, ns = 'common') => i18n.t(`${ns}:${key}`);

const statusOptions = () => [
  { value: 'WAITING', label: t('experiment.waiting', 'table') },
  { value: 'RUNNING', label: t('experiment.running', 'table') },
  { value: 'STOPPED', label: t('experiment.stopped', 'table') },
  { value: 'FORCE_STOPPED', label: t('experiment.force-stopped', 'table') }
];

export const experimentFilterConfig: FilterModalConfig<ExperimentFilters> = {
  mode: 'multi',

  shouldHydrate: filters => !!filters.isFilter || !!filters.filterBySummary,
  submitExtra: { isFilter: true },
  fields: [
    {
      type: FilterTypes.STATUSES,
      labelKey: 'status',
      valueKind: 'multiselect',
      emptyValue: [],
      useData: () => {
        const options = statusOptions();
        return {
          options,
          getLabel: value =>
            (Array.isArray(value) &&
              value
                .map(v => options.find(item => item.value === v)?.label)
                .join(', ')
                .toLowerCase()) ||
            ''
        };
      },
      toFilter: filterValue => ({
        statuses: filterValue as ExperimentStatus[]
      }),
      fromFilter: filters => filters.statuses
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
