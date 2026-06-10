import { useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { OrderBy, OrderDirection } from '@types';
import Dropdown, { DropdownOption } from 'components/dropdown';

interface SortedState {
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
}

interface SortByProps<T extends SortedState> {
  filters: T;
  setFilters: (filters: Partial<T>) => void;
  sortByOptions: DropdownOption[];
  sortDirectionOptions: DropdownOption[];
}

const SortBy = <T extends SortedState>({
  filters,
  setFilters,
  sortByOptions,
  sortDirectionOptions
}: SortByProps<T>) => {
  useTranslation(['common', 'table', 'form']);

  const currentOption = useMemo(
    () => sortByOptions.find(item => item.value === filters.orderBy),
    [filters.orderBy, sortByOptions]
  );

  const handleSorting = useCallback(
    (value: Partial<SortedState>) => {
      setFilters(value as Partial<T>);
    },
    [setFilters]
  );

  return (
    <Dropdown
      labelCustom={
        currentOption ? (
          <Trans
            i18nKey={'common:sort-by'}
            values={{
              sortBy: currentOption.label
            }}
          />
        ) : (
          ''
        )
      }
      options={sortByOptions}
      value={filters.orderBy}
      onChange={value => handleSorting({ orderBy: value as OrderBy })}
      wrapTriggerStyle="w-fit"
      className="w-fit"
      contentClassName="!max-h-fit !divide-y"
      additionalOptions={sortDirectionOptions}
      additionalValue={filters.orderDirection}
      onChangeAdditional={value =>
        handleSorting({ orderDirection: value as OrderDirection })
      }
    />
  );
};

export default SortBy;
