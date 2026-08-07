import { getEditorEnvironments, useAuth } from 'auth';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import {
  checkEnvironmentEmptyId,
  onFormatEnvironments,
  resolveEnvironmentEmptyId
} from 'utils/function';
import { PushFilters } from 'pages/pushes/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

const booleanOptions = () => [
  { value: 1, label: t('yes') },
  { value: 0, label: t('no') }
];

export const pushFilterConfig: FilterModalConfig<PushFilters> = {
  mode: 'multi',
  fields: [
    {
      type: FilterTypes.ENVIRONMENT_IDs,
      labelKey: 'environment',
      valueKind: 'searchable',
      emptyValue: [],
      useData: () => {
        const { consoleAccount } = useAuth();
        const { editorEnvironments } = getEditorEnvironments(consoleAccount!);
        const { formattedEnvironments } =
          onFormatEnvironments(editorEnvironments);
        const options = formattedEnvironments.map(item => ({
          label: item.name,
          // Normalize the empty-id environment to a stable placeholder so
          // `fromFilter` can hydrate it without any cross-render state.
          value: resolveEnvironmentEmptyId(checkEnvironmentEmptyId(item.id))
        }));
        return {
          options,
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
      toFilter: filterValue => ({ disabled: !filterValue }),
      fromFilter: filters =>
        filters.disabled === undefined ? undefined : filters.disabled ? 0 : 1
    }
  ]
};
