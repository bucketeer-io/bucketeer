import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { OrderBy, OrderDirection } from '@types';
import { IconChecked } from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import { FlagFilters } from '../types';

interface SortedState {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const ActiveItem = ({ isActive }: { isActive: boolean }) =>
  isActive ? (
    <Icon
      icon={IconChecked}
      color="primary-500"
      size={'sm'}
      className="flex-center"
    />
  ) : (
    <></>
  );

const SortBy = ({
  filters,
  setFilters
}: {
  filters: FlagFilters;
  setFilters: (filters: FlagFilters) => void;
}) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const [sortedState, setSortedState] = useState<SortedState>({
    orderBy: filters.orderBy,
    orderDirection: filters.orderDirection
  });

  const sortByOptions = useMemo(
    () => [
      {
        label: t('name'),
        value: 'NAME'
      },
      {
        label: t('tags'),
        value: 'TAGS'
      },
      {
        label: t('table:created-at'),
        value: 'CREATED_AT'
      },
      {
        label: t('table:updated-at'),
        value: 'UPDATED_AT'
      },
      {
        label: t('enabled'),
        value: 'ENABLED'
      }
    ],
    []
  );

  const sortDirectionOptions = useMemo(
    () => [
      {
        label: t('sort-asc'),
        value: 'ASC'
      },
      {
        label: t('sort-desc'),
        value: 'DESC'
      }
    ],
    []
  );

  const currentOption = useMemo(
    () => sortByOptions.find(item => item.value === sortedState.orderBy),
    [sortedState]
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
    <DropdownMenu>
      <DropdownMenuTrigger
        label={
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
      />
      <DropdownMenuContent className="!max-h-fit divide-y">
        <div className="pb-1">
          {sortByOptions.map(({ label, value }, index) => (
            <DropdownMenuItem
              label={label}
              value={value}
              key={index}
              additionalElement={
                <ActiveItem isActive={sortedState.orderBy === value} />
              }
              onSelectOption={value =>
                handleSorting({
                  orderBy: value as OrderBy
                })
              }
            />
          ))}
        </div>
        <div className="pt-1">
          {sortDirectionOptions.map(({ label, value }, index) => (
            <DropdownMenuItem
              key={index}
              label={label}
              value={value}
              additionalElement={
                <ActiveItem isActive={sortedState.orderDirection === value} />
              }
              onSelectOption={value =>
                handleSorting({ orderDirection: value as OrderDirection })
              }
            />
          ))}
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default SortBy;
