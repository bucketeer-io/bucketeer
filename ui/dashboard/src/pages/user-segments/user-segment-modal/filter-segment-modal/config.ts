import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { UserSegmentsFilters } from 'pages/user-segments/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

const statusOptions = () => [
  { value: FilterTypes.IN_USE, label: t('in-use') },
  { value: FilterTypes.NOT_IN_USE, label: t('not-in-use') }
];

export const userSegmentFilterConfig: FilterModalConfig<UserSegmentsFilters> = {
  mode: 'single',
  fields: [
    {
      type: FilterTypes.STATUS,
      labelKey: 'status',
      valueKind: 'enum',
      emptyValue: '',
      useData: () => ({ options: statusOptions() }),
      toFilter: filterValue => ({
        isInUseStatus: filterValue === FilterTypes.IN_USE
      }),
      fromFilter: filters =>
        filters.isInUseStatus === undefined
          ? undefined
          : filters.isInUseStatus
            ? FilterTypes.IN_USE
            : FilterTypes.NOT_IN_USE
    }
  ]
};
