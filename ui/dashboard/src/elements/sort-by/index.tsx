import { useCallback, useMemo, useState } from 'react';
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
  const [sortedState, setSortedState] = useState<SortedState>({
    orderBy: filters.orderBy,
    orderDirection: filters.orderDirection
  });

  const currentOption = useMemo(
    () => sortByOptions.find(item => item.value === sortedState.orderBy),
    [sortedState, sortByOptions]
  );

  const handleSorting = useCallback(
    (value: Partial<SortedState>) => {
      setSortedState(prev => ({ ...prev, ...value }));
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
      value={sortedState.orderBy}
      onChange={value => handleSorting({ orderBy: value as OrderBy })}
      wrapTriggerStyle="w-fit"
      className="w-fit"
      contentClassName="!max-h-fit !divide-y"
      additionalOptions={sortDirectionOptions}
      additionalValue={sortedState.orderDirection}
      onChangeAdditional={value =>
        handleSorting({ orderDirection: value as OrderDirection })
      }
    />
  );
};

export default SortBy;
