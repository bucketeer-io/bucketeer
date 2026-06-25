import { useQueryEnvironments } from '@queries/environments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import {
  checkEnvironmentEmptyId,
  onFormatEnvironments,
  resolveEnvironmentEmptyId
} from 'utils/function';
import { APIKeysFilters } from 'pages/api-keys/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

// Matches `booleanOptions` in use-options (numeric 1 = yes / 0 = no).
const booleanOptions = () => [
  { value: 1, label: t('yes') },
  { value: 0, label: t('no') }
];

export const apiKeyFilterConfig: FilterModalConfig<APIKeysFilters> = {
  mode: 'multi',
  fields: [
    {
      type: FilterTypes.ENVIRONMENT_IDs,
      labelKey: 'environment',
      valueKind: 'searchable',
      emptyValue: [],
      useData: ({ enabled }) => {
        const { consoleAccount } = useAuth();
        const currentEnvironment = getCurrentEnvironment(consoleAccount!);
        const { data, isLoading } = useQueryEnvironments({
          params: {
            cursor: '0',
            organizationId: currentEnvironment.organizationId
          },
          enabled: !!currentEnvironment.organizationId && enabled
        });
        const { formattedEnvironments } = onFormatEnvironments(
          data?.environments || []
        );
        const options = formattedEnvironments.map(item => ({
          label: item.name,
          // Normalize the empty-id environment to a stable placeholder so
          // `fromFilter` can hydrate it without any cross-render state.
          value: resolveEnvironmentEmptyId(checkEnvironmentEmptyId(item.id))
        }));
        return {
          options,
          isLoading,
          getLabel: value =>
            (Array.isArray(value) &&
              value
                .map(item => options.find(env => env.value === item)?.label)
                .filter(Boolean)
                .join(', ')) ||
            ''
        };
      },
      toFilter: filterValue => ({
        environmentIds: Array.isArray(filterValue)
          ? (filterValue as string[]).map(item => checkEnvironmentEmptyId(item))
          : []
      }),
      fromFilter: filters =>
        Array.isArray(filters.environmentIds)
          ? filters.environmentIds.map(resolveEnvironmentEmptyId)
          : undefined
    },
    {
      type: FilterTypes.ENABLED,
      labelKey: 'enabled',
      valueKind: 'boolean',
      emptyValue: '',
      useData: () => ({ options: booleanOptions() }),
      // Stored value is 1 (enabled) / 0 (disabled); `disabled` is the inverse.
      toFilter: filterValue => ({ disabled: !filterValue }),
      fromFilter: filters =>
        filters.disabled === undefined ? undefined : filters.disabled ? 0 : 1
    }
  ]
};
