import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { ExperimentStatus } from '@types';
import { ExperimentFilters } from 'pages/experiments/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string, ns = 'common') => i18n.t(`${ns}:${key}`);

// Matches `experimentStatusOptions` in use-options.
const statusOptions = () => [
  { value: 'WAITING', label: t('experiment.waiting', 'table') },
  { value: 'RUNNING', label: t('experiment.running', 'table') },
  { value: 'STOPPED', label: t('experiment.stopped', 'table') },
  { value: 'FORCE_STOPPED', label: t('experiment.force-stopped', 'table') }
];

export const experimentFilterConfig: FilterModalConfig<ExperimentFilters> = {
  mode: 'multi',
  // Page treats `isFilter`/`filterBySummary` as the "filter is active" markers;
  // only hydrate when one is present, and re-stamp `isFilter` on submit.
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
      useData: ({ enabled }) => {
        const { consoleAccount } = useAuth();
        const currentEnvironment = getCurrentEnvironment(consoleAccount!);
        const { data, isLoading } = useQueryAccounts({
          params: {
            cursor: String(0),
            environmentId: currentEnvironment?.id,
            organizationId: currentEnvironment?.organizationId
          },
          enabled
        });
        const options = (data?.accounts || []).map(item => ({
          label: item.email,
          value: item.email
        }));
        // Maintainer label echoes the stored email verbatim.
        return {
          options,
          isLoading,
          getLabel: value => (value as string) || ''
        };
      },
      toFilter: filterValue => ({ maintainer: filterValue as string }),
      fromFilter: filters => filters.maintainer
    }
  ]
};
