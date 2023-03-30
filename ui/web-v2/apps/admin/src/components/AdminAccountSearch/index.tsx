import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/accounts';
import { useIsOwner } from '../../modules/me';
import { Account } from '../../proto/account/account_pb';
import { AdminAccountSearchOptions } from '../../types/adminAccount';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { classNames } from '../../utils/css';
import { FilterChip } from '../FilterChip';
import { FilterPopover, Option } from '../FilterPopover';
import { FilterRemoveAllButtonProps } from '../FilterRemoveAllButton';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.account.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.account.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.account.sort.emailAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.account.sort.emailZa),
  },
];

export enum FilterTypes {
  ENABLED = 'enabled',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: intl.formatMessage(messages.account.filter.enabled),
  },
];

export const enabledOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.enabled),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.disabled),
  },
];

export interface AdminAccountSearchProps {
  options: AdminAccountSearchOptions;
  onChange: (options: AdminAccountSearchOptions) => void;
  onAdd: () => void;
}

export const AdminAccountSearch: FC<AdminAccountSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsOwner();
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.accounts.loading,
      shallowEqual
    );
    const accounts = useSelector<AppState, Account.AsObject[]>(
      (state) => selectAll(state.accounts),
      shallowEqual
    );
    const [filterValues, setFilterValues] = useState<Option[]>([]);

    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.ENABLED:
            setFilterValues(enabledOptions);
            return;
        }
      },
      [setFilterValues, accounts]
    );

    const handleUpdateOption = (
      optionPart: Partial<AdminAccountSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const handleFilterAdd = (key: string, value?: string): void => {
      switch (key) {
        case FilterTypes.ENABLED:
          handleUpdateOption({
            enabled: value,
          });
          return;
      }
    };
    return (
      <div
        className={classNames(
          'w-full',
          'px-5 py-5 sticky top-0',
          'z-10 border-b border-gray-300'
        )}
      >
        <div className={classNames('w-full min-w-max', 'flex flex-row')}>
          <div className="flex-none w-72">
            <SearchInput
              placeholder={f(messages.account.search.placeholder)}
              value={options.q}
              onChange={(query: string) =>
                handleUpdateOption({
                  q: query,
                })
              }
            />
          </div>
          <div className="flex-none mx-2">
            <FilterPopover
              keys={filterOptions}
              values={filterValues}
              onChangeKey={handleFilterKeyChange}
              onAdd={handleFilterAdd}
            />
          </div>
          <div className="flex-grow" />
          <div className="flex-none -mr-2">
            <SortSelect
              sortKey={options.sort}
              sortItems={sortItems}
              onChange={(sort: string) =>
                handleUpdateOption({
                  sort: sort,
                })
              }
            />
          </div>
          {editable && (
            <div className="flex-none ml-8">
              <button
                type="button"
                className="btn-submit"
                disabled={false}
                onClick={onAdd}
              >
                <PlusIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
                {f(messages.button.add)}
              </button>
            </div>
          )}
        </div>
        {options.enabled && (
          <div className="flex space-x-2 pt-2">
            {options.enabled && (
              <FilterChip
                label={`${f(messages.account.filter.enabled)}: ${
                  enabledOptions.find(
                    (option) => option.value === options.enabled
                  ).label
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    enabled: null,
                  })
                }
              />
            )}
            {options.enabled && (
              <FilterRemoveAllButtonProps
                onClick={() =>
                  handleUpdateOption({
                    enabled: null,
                  })
                }
              />
            )}
          </div>
        )}
      </div>
    );
  }
);
