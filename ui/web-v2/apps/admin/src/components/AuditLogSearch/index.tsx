import { FC, memo, useState } from 'react';
import { useIntl } from 'react-intl';
import { components } from 'react-select';

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
import { SearchInput } from '../SearchInput';
import { Option, Select } from '../Select';
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

export interface AuditLogSearchProps {
  options: AuditLogSearchOptions;
  onChange: (options: AuditLogSearchOptions) => void;
  showEntityTypeFilter?: boolean;
  entityTypeOptions?: Option[];
}

export const AuditLogSearch: FC<AuditLogSearchProps> = memo(
  ({ options, onChange, showEntityTypeFilter, entityTypeOptions }) => {
    const [selectedEntityType, setSelectedEntityType] = useState<Option>(
      options.entityType
        ? entityTypeOptions.find(
            (e) => Number(e.value) === Number(options.entityType)
          )
        : null
    );
    const { formatMessage: f } = useIntl();

    const handleUpdateOption = (
      optionPart: Partial<AuditLogSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const handleEntityType = (option: Option) => {
      setSelectedEntityType(option);
      onChange({
        ...options,
        entityType: option ? Number(option.value) : null,
      });
    };

    const ControlComponent = ({ children, ...props }) => (
      <components.Control {...props}>
        <span className="ml-2">{f(messages.action)}:</span> {children}
      </components.Control>
    );

    return (
      <div
        className={classNames(
          'w-full',
          'px-5 py-5 sticky top-0',
          'z-10 border-b border-gray-300'
        )}
      >
        <div className="flex justify-between">
          <div className="flex space-x-2">
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
            {showEntityTypeFilter && (
              <Select
                placeholder={f(messages.all)}
                clearable
                options={entityTypeOptions}
                className={classNames('flex-none w-[262px]')}
                value={selectedEntityType}
                onChange={handleEntityType}
                customControl={ControlComponent}
              />
            )}
            <div className="flex-none relative">
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
          </div>
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
