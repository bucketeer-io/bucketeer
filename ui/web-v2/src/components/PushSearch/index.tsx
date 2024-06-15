import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { PushSearchOptions } from '../../types/push';
import { classNames } from '../../utils/css';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.push.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.push.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.push.sort.nameAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.push.sort.nameZa),
  },
];

export interface PushSearchProps {
  options: PushSearchOptions;
  onChange: (options: PushSearchOptions) => void;
  onAdd: () => void;
}

export const PushSearch: FC<PushSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const handleUpdateOption = (
      optionPart: Partial<PushSearchOptions>
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
              placeholder={f(messages.push.search.placeholder)}
              value={options.q}
              onChange={(query: string) =>
                handleUpdateOption({
                  q: query,
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
      </div>
    );
  }
);
