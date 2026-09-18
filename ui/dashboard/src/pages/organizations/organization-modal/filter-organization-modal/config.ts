import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { OrganizationFilters } from 'pages/organizations/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

const enabledOptions = () => [
  { value: 'yes', label: t('yes') },
  { value: 'no', label: t('no') }
];

export const organizationFilterConfig: FilterModalConfig<OrganizationFilters> =
  {
    mode: 'single',
    fields: [
      {
        type: FilterTypes.ENABLED,
        labelKey: 'enabled',
        valueKind: 'enum',
        emptyValue: '',
        useData: () => ({ options: enabledOptions() }),
        toFilter: filterValue => ({ disabled: filterValue === 'no' }),
        fromFilter: filters =>
          filters.disabled === undefined
            ? undefined
            : filters.disabled
              ? 'no'
              : 'yes'
      }
    ]
  };
