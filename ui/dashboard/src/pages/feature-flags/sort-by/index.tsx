import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import useOptions from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { OrderBy, OrderDirection } from '@types';
import Dropdown from 'components/dropdown';
import { FlagFilters } from '../types';

interface SortedState {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const SortBy = ({
  filters,
  setFilters
}: {
  filters: FlagFilters;
  setFilters: (filters: FlagFilters) => void;
}) => {
  useTranslation(['common', 'table', 'form']);
  const { flagSortByOptions, flagSortDirectionOptions } = useOptions();
  const [sortedState, setSortedState] = useState<SortedState>({
    orderBy: filters.orderBy,
    orderDirection: filters.orderDirection
  });

  const currentOption = useMemo(
    () => flagSortByOptions.find(item => item.value === sortedState.orderBy),
    [sortedState, flagSortByOptions]
  );

  const handleSorting = useCallback(
    (value: Partial<SortedState>) => {
      setSortedState({
        ...sortedState,
        ...value
      });
      setFilters({
        ...filters,
        ...value
      });
    },
    [filters, sortedState]
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
      options={flagSortByOptions}
      value={sortedState.orderBy}
      onChange={value => handleSorting({ orderBy: value as OrderBy })}
      wrapTriggerStyle="w-fit"
      className="w-fit"
      contentClassName="!max-h-fit !divide-y"
      addititonOptions={flagSortDirectionOptions}
      additionalValue={sortedState.orderDirection}
      onChangeAdditional={value =>
        handleSorting({ orderDirection: value as OrderDirection })
      }
    />
  );
};

export default SortBy;
