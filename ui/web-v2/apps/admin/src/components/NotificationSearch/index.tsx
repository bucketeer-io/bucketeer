import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';

import { FilterChip } from '../../components/FilterChip';
import { FilterPopover } from '../../components/FilterPopover';
import { FilterRemoveAllButtonProps } from '../../components/FilterRemoveAllButton';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { NotificationSearchOptions } from '../../types/notification';
import { classNames } from '../../utils/css';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

export interface Option {
  value: string;
  label: string;
}

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.notification.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.notification.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.notification.sort.nameAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.notification.sort.nameZa),
  },
];

export enum FilterTypes {
  ENABLED = 'enabled',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: intl.formatMessage(messages.notification.filter.enabled),
  },
];

export const enabledOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.notification.filterOptions.enabled),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.notification.filterOptions.disabled),
  },
];

export interface NotificationSearchProps {
  options: NotificationSearchOptions;
  onChange: (options: NotificationSearchOptions) => void;
  onAdd: () => void;
}

export const NotificationSearch: FC<NotificationSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const [filterValues, setFilterValues] = useState<Option[]>([]);
    const handleUpdateOption = (
      optionPart: Partial<NotificationSearchOptions>
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
    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.ENABLED:
            setFilterValues(enabledOptions);
            return;
        }
      },
      [setFilterValues]
    );

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
              placeholder={f(messages.apiKey.search.placeholder)}
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
              <button type="button" className="btn-submit" onClick={onAdd}>
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
                label={`${f(messages.notification.filter.enabled)}: ${
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
