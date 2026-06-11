import { useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useScreen } from 'hooks/use-screen';
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
  const { t } = useTranslation(['common', 'table', 'form']);
  const { isMobile } = useScreen();

  const triggerWidthCls = isMobile
    ? 'w-[200px] min-w-[200px] max-w-[200px]'
    : 'w-fit';

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
        ) : undefined
      }
      placeholder={t('common:sort-by-placeholder')}
      options={sortByOptions}
      value={filters.orderBy}
      onChange={value => handleSorting({ orderBy: value as OrderBy })}
      wrapTriggerStyle={triggerWidthCls}
      className={triggerWidthCls}
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
