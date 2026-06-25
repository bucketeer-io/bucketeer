import { useQueryEnvironments } from '@queries/environments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { APIKeysFilters } from 'pages/api-keys/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

// Matches `booleanOptions` in use-options (numeric 1 = yes / 0 = no).
const booleanOptions = () => [
  { value: 1, label: t('yes') },
  { value: 0, label: t('no') }
];

// The empty-id environment is represented in the dropdown by an index-suffixed
// placeholder id (see `onFormatEnvironments`). `useData` resolves it each render;
// `fromFilter` reads it back so a stored "" hydrates to the matching option.
let emptyEnvironmentId = '';

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
        const { emptyEnvironmentId: emptyId, formattedEnvironments } =
          onFormatEnvironments(data?.environments || []);
        emptyEnvironmentId = emptyId;
        const options = formattedEnvironments.map(item => ({
          label: item.name,
          value: item.id
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
          ? filters.environmentIds.map(item => item || emptyEnvironmentId)
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
