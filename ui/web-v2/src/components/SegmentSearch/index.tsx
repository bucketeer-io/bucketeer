import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable } from '../../modules/me';
import { selectAll } from '../../modules/segments';
import { Segment } from '../../proto/feature/segment_pb';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { SegmentSearchOptions } from '../../types/segment';
import { classNames } from '../../utils/css';
import { FilterChip } from '../FilterChip';
import { FilterPopover, Option } from '../FilterPopover';
import { FilterRemoveAllButtonProps } from '../FilterRemoveAllButton';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.segment.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.segment.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.segment.sort.nameAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.segment.sort.nameZa),
  },
];

export enum FilterTypes {
  STATUS = 'status',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.STATUS,
    label: intl.formatMessage(messages.segment.filter.status),
  },
];

export const inUseOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.segment.filterOptions.inUse),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.segment.filterOptions.notInUse),
  },
];

export interface SegmentSearchProps {
  options: SegmentSearchOptions;
  onChange: (options: SegmentSearchOptions) => void;
  onAdd: () => void;
}

export const SegmentSearch: FC<SegmentSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.segments.loading,
      shallowEqual
    );
    const segments = useSelector<AppState, Segment.AsObject[]>(
      (state) => selectAll(state.segments),
      shallowEqual
    );
    const [filterValues, setFilterValues] = useState<Option[]>([]);
    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.STATUS:
            setFilterValues(inUseOptions);
            return;
        }
      },
      [setFilterValues, segments]
    );

    const handleUpdateOption = (
      optionPart: Partial<SegmentSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const handleFilterAdd = (key: string, value?: string): void => {
      switch (key) {
        case FilterTypes.STATUS:
          handleUpdateOption({
            inUse: value,
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
              placeholder={f(messages.segment.search.placeholder)}
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
        {options.inUse && (
          <div className="flex space-x-2 pt-2">
            <FilterChip
              label={`${f(messages.segment.filter.status)}: ${
                inUseOptions.find((option) => option.value === options.inUse)
                  .label
              }`}
              onRemove={() =>
                handleUpdateOption({
                  inUse: null,
                })
              }
            />
            <FilterRemoveAllButtonProps
              onClick={() =>
                handleUpdateOption({
                  inUse: null,
                })
              }
            />
          </div>
        )}
      </div>
    );
  }
);
