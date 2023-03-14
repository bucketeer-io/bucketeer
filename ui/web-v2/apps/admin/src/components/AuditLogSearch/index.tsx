import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { AuditLogSearchOptions } from '../../types/auditLog';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
} from '../../types/list';
import { classNames } from '../../utils/css';
import { DateRangePopover } from '../DateRangePopover';
import { Option } from '../FilterPopover';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.auditLog.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.auditLog.sort.oldest),
  },
];

export enum FilterTypes {
  DATES = 'dates',
  TYPE = 'type',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.DATES,
    label: intl.formatMessage(messages.auditLog.filter.dates),
  },
  {
    value: FilterTypes.TYPE,
    label: intl.formatMessage(messages.auditLog.filter.type),
  },
];

export interface AuditLogSearchProps {
  options: AuditLogSearchOptions;
  onChange: (options: AuditLogSearchOptions) => void;
}

export const AuditLogSearch: FC<AuditLogSearchProps> = memo(
  ({ options, onChange }) => {
    const { formatMessage: f } = useIntl();
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.auditLog.loading,
      shallowEqual
    );
    const handleUpdateOption = (
      optionPart: Partial<AuditLogSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
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
              onChange={(query: string) =>
                handleUpdateOption({
                  q: query,
                })
              }
            />
          </div>
          <div className="flex-none mx-2 relative">
            <DateRangePopover
              options={options}
              onChange={(from: number, to: number) =>
                handleUpdateOption({
                  from,
                  to,
                })
              }
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
        </div>
      </div>
    );
  }
);
