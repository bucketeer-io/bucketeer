import { PlusIcon } from '@heroicons/react/solid';
import MUArchiveIcon from '@material-ui/icons/Archive';
import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { GOAL_LIST_PAGE_SIZE } from '../../constants/goal';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/goals';
import { useIsEditable } from '../../modules/me';
import { Goal } from '../../proto/experiment/goal_pb';
import { GoalSearchOptions } from '../../types/goal';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { classNames } from '../../utils/css';
import { ActionMenu, MenuActions, MenuItem } from '../ActionMenu';
import { FilterChip } from '../FilterChip';
import { FilterPopover, Option } from '../FilterPopover';
import { FilterRemoveAllButtonProps } from '../FilterRemoveAllButton';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

export interface GoalListProps {
  searchOptions: GoalSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: GoalSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (g: Goal.AsObject) => void;
  onArchive: (g: Goal.AsObject) => void;
}

export const GoalList: FC<GoalListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
    onArchive,
  }) => {
    const { formatMessage: f, formatDate, formatTime } = useIntl();
    const editable = useIsEditable();
    const goals = useSelector<AppState, Goal.AsObject[]>(
      (state) => selectAll(state.goals),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.goals.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.goals.totalCount,
      shallowEqual
    );

    const createMenuItems = (): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      items.push({
        action: MenuActions.ARCHIVE,
        name: intl.formatMessage(messages.feature.action.archive),
        iconElement: <MUArchiveIcon />,
      });
      return items;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <GoalSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <div className="p-9">
            <ListSkeleton />
          </div>
        ) : goals.length == 0 ? (
          searchOptions.q || searchOptions.status ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.goal.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(messages.goal.list.noResult.searchKeyword),
                      })}
                    </li>
                    <li>{f(messages.noResult.changeFilterSelection)}</li>
                    <li>{f(messages.noResult.checkTypos)}</li>
                  </ul>
                </div>
              </div>
            </div>
          ) : (
            <div className="my-10 flex justify-center">
              <div className="w-[600px] text-gray-700 text-center">
                <h1 className="text-lg">
                  {f(messages.noData.title, {
                    title: f(messages.goal.list.header.title),
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.goal.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/running-abn-tests?id=create-goals"
                  target="_blank"
                  rel="noreferrer"
                  className="link"
                >
                  {f(messages.readMore)}
                </a>
              </div>
            </div>
          )
        ) : (
          <div>
            <table className="min-w-full table-auto leading-normal">
              <tbody className="text-sm">
                {goals.map((goal) => {
                  return (
                    <tr key={goal.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(goal)}
                          >
                            {goal.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(goal.createdAt * 1000)}
                            />
                          </div>
                        </div>
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'text-gray-700',
                          'whitespace-nowrap'
                        )}
                      >
                        {
                          statusOptions.find(
                            (option) =>
                              option.value == goal.isInUseStatus.toString()
                          ).label
                        }
                      </td>
                      {editable && !goal.archived && (
                        <td
                          className={classNames(
                            'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                            'whitespace-nowrap'
                          )}
                        >
                          <ActionMenu
                            onClickAction={(action) => {
                              switch (action) {
                                case MenuActions.ARCHIVE:
                                  onArchive(goal);
                                  return;
                              }
                            }}
                            menuItems={createMenuItems()}
                          />
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / GOAL_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);

const sortItems: SortItem[] = [
  {
    key: SORT_OPTIONS_CREATED_AT_DESC,
    message: intl.formatMessage(messages.feature.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.feature.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.feature.sort.nameAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.feature.sort.nameZa),
  },
];

export enum FilterTypes {
  STATUS = 'status',
  ARCHIVED = 'archived',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.STATUS,
    label: intl.formatMessage(messages.goal.filter.status),
  },
  {
    value: FilterTypes.ARCHIVED,
    label: intl.formatMessage(messages.goal.filter.archived),
  },
];

export const statusOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.goal.status.inUse),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.goal.status.notInUse),
  },
];

const archivedOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no),
  },
];

export interface GoalSearchProps {
  options: GoalSearchOptions;
  onChange: (options: GoalSearchOptions) => void;
  onAdd: () => void;
}

export const GoalSearch: FC<GoalSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const [filterValues, setFilterValues] = useState<Option[]>([]);

    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.STATUS:
            setFilterValues(statusOptions);
            return;
          case FilterTypes.ARCHIVED:
            setFilterValues(archivedOptions);
            return;
        }
      },
      [setFilterValues]
    );

    const handleUpdateOption = (
      optionPart: Partial<GoalSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const handleFilterAdd = (key: string, value?: string): void => {
      switch (key) {
        case FilterTypes.STATUS:
          handleUpdateOption({
            status: value,
          });
          return;
        case FilterTypes.ARCHIVED:
          handleUpdateOption({
            archived: value,
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
              placeholder={f(messages.experiment.search.placeholder)}
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
        {(options.status || options.archived) && (
          <div className="flex space-x-2 pt-2">
            {options.status && (
              <FilterChip
                label={`${f(messages.goal.status.status)}: ${
                  statusOptions.find(
                    (option) => option.value === options.status
                  ).label
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    status: null,
                  })
                }
              />
            )}
            {options.archived && (
              <FilterChip
                label={`${f(messages.goal.filter.archived)}: ${
                  archivedOptions.find(
                    (option) => option.value === options.archived
                  ).label
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    archived: null,
                  })
                }
              />
            )}
            {(options.status || options.archived) && (
              <FilterRemoveAllButtonProps
                onClick={() =>
                  handleUpdateOption({
                    status: null,
                    archived: null,
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
