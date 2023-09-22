import { PlusIcon } from '@heroicons/react/solid';
import MUArchiveIcon from '@material-ui/icons/Archive';
import MUFileCopyIcon from '@material-ui/icons/FileCopy';
import MUUnarchiveIcon from '@material-ui/icons/Unarchive';
import dayjs from 'dayjs';
import React, { FC, useState, memo, useCallback, useEffect } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

import { FEATURE_LIST_PAGE_SIZE } from '../../constants/feature';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll as selectAllAccounts } from '../../modules/accounts';
import { selectAll as selectAllFeatures } from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { selectAll as selectAllTags } from '../../modules/tags';
import { Account } from '../../proto/account/account_pb';
import { Feature, Tag } from '../../proto/feature/feature_pb';
import { FeatureSearchOptions } from '../../types/feature';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { classNames } from '../../utils/css';
import { ActionMenu, MenuActions, MenuItem } from '../ActionMenu';
import { FeatureIdChip } from '../FeatureIdChip';
import { FilterChip } from '../FilterChip';
import { Option, FilterPopover } from '../FilterPopover';
import { FilterRemoveAllButtonProps } from '../FilterRemoveAllButton';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';
import { Switch } from '../Switch';
import { TagChips } from '../TagsChips';

export enum FlagStatus {
  NEW, // This flag is new and has not been requested yet.
  RECEIVING_REQUESTS, // It is receiving one more requests in the last 7 days.
  INACTIVE, // It is not receiving requests for 7 days.
}

interface FlagStatusIconProps {
  flagStatus: FlagStatus;
}

export const FlagStatucIcon: FC<FlagStatusIconProps> = ({ flagStatus }) => {
  let msg = '';
  switch (flagStatus) {
    case FlagStatus.NEW:
      msg = intl.formatMessage(messages.feature.flagStatus.new);
      break;
    case FlagStatus.RECEIVING_REQUESTS:
      msg = intl.formatMessage(messages.feature.flagStatus.receivingRequests);
      break;
    case FlagStatus.INACTIVE:
      msg = intl.formatMessage(messages.feature.flagStatus.inactive);
      break;
  }
  return (
    <p className="whitespace-nowrap">
      <span>{msg}</span>
    </p>
  );
};

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

export function getFlagStatus(
  feature: Feature.AsObject,
  relativeDate: Date
): FlagStatus {
  if (!feature.lastUsedInfo) {
    return FlagStatus.NEW;
  }
  const lastUsedAt = dayjs(new Date(feature.lastUsedInfo.lastUsedAt * 1000));
  if (lastUsedAt.diff(dayjs(relativeDate), 'day', true) > -7) {
    return FlagStatus.RECEIVING_REQUESTS;
  }
  return FlagStatus.INACTIVE;
}

function getClientVersion(feature: Feature.AsObject): string {
  if (!feature.lastUsedInfo) {
    return '';
  }
  if (
    !feature.lastUsedInfo.clientOldestVersion ||
    !feature.lastUsedInfo.clientLatestVersion
  ) {
    if (feature.lastUsedInfo.clientOldestVersion) {
      return feature.lastUsedInfo.clientOldestVersion;
    }
    if (feature.lastUsedInfo.clientLatestVersion) {
      return feature.lastUsedInfo.clientLatestVersion;
    }
    return '';
  }
  if (
    feature.lastUsedInfo.clientOldestVersion ===
    feature.lastUsedInfo.clientLatestVersion
  ) {
    return feature.lastUsedInfo.clientOldestVersion;
  }
  return `${feature.lastUsedInfo.clientOldestVersion} ~ ${feature.lastUsedInfo.clientLatestVersion}`;
}

export enum FilterTypes {
  MAINTAINER = 'maintainer',
  HAS_EXPERIMENT = 'has_experiment',
  ENABLED = 'enabled',
  ARCHIVED = 'archived',
  TAGS = 'tags',
  HAS_PREREQUISITES = 'prerequisites',
}

const filterOptions: Option[] = [
  {
    value: FilterTypes.TAGS,
    label: intl.formatMessage(messages.tags),
  },
  {
    value: FilterTypes.HAS_PREREQUISITES,
    label: intl.formatMessage(messages.feature.filter.hasPrerequisites),
  },
  {
    value: FilterTypes.MAINTAINER,
    label: intl.formatMessage(messages.feature.filter.maintainer),
  },
  {
    value: FilterTypes.HAS_EXPERIMENT,
    label: intl.formatMessage(messages.feature.filter.hasExperiment),
  },
  {
    value: FilterTypes.ENABLED,
    label: intl.formatMessage(messages.feature.filter.enabled),
  },
  {
    value: FilterTypes.ARCHIVED,
    label: intl.formatMessage(messages.feature.filter.archived),
  },
];

const enabledOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.enabled),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.disabled),
  },
];

const hasExperimentOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no),
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

const hasPrerequisitesOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes),
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no),
  },
];

export interface FeatureListProps {
  searchOptions: FeatureSearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (featureId: string, enabled: boolean) => void;
  onAdd: () => void;
  onChangeSearchOptions: (options: FeatureSearchOptions) => void;
  onClearSearchOptions: () => void;
  onArchive: (feature: Feature.AsObject) => void;
  onClone: (feature: Feature.AsObject) => void;
}

export const FeatureList: FC<FeatureListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onSwitchEnabled,
    onAdd,
    onChangeSearchOptions,
    onClearSearchOptions,
    onArchive,
    onClone,
  }) => {
    const { formatMessage: f, formatDate } = useIntl();
    const currentEnvironment = useCurrentEnvironment();
    const editable = useIsEditable();
    const isFeatureLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading,
      shallowEqual
    );
    const isAccountLoading = useSelector<AppState, boolean>(
      (state) => state.accounts.loading,
      shallowEqual
    );
    const isLoading = isFeatureLoading || isAccountLoading;
    const features = useSelector<AppState, Feature.AsObject[]>(
      (state) => selectAllFeatures(state.features),
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.features.totalCount,
      shallowEqual
    );
    const accounts = useSelector<AppState, Account.AsObject[]>(
      (state) => selectAllAccounts(state.accounts),
      shallowEqual
    );
    const relativeDate = new Date();
    const createMenuItems = (archived: boolean): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      if (archived) {
        items.push({
          action: MenuActions.ARCHIVE,
          name: intl.formatMessage(messages.feature.action.unarchive),
          iconElement: <MUUnarchiveIcon />,
        });
        return items;
      }
      items.push({
        action: MenuActions.ARCHIVE,
        name: intl.formatMessage(messages.feature.action.archive),
        iconElement: <MUArchiveIcon />,
      });
      items.push({
        action: MenuActions.CLONE,
        name: intl.formatMessage(messages.feature.action.clone),
        iconElement: <MUFileCopyIcon />,
      });
      return items;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <FeatureSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onClear={onClearSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : features.length == 0 ? (
          searchOptions.q ||
          searchOptions.enabled ||
          searchOptions.archived ||
          searchOptions.hasExperiment ||
          searchOptions.maintainerId ||
          searchOptions.tagIds?.length > 0 ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.feature.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(
                          messages.feature.list.noResult.searchKeyword
                        ),
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
                    title: f(messages.feature.list.header.title),
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.feature.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/creating-a-feature-flag?id=create-a-feature-flag"
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
            <table className="table-auto leading-normal">
              <tbody className="text-sm">
                {features.map((feature) => (
                  <tr key={feature.id}>
                    <td className="px-5 py-3 border-b">
                      <div className="flex mb-2">
                        <Link
                          to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_TARGETING}`}
                          className="link text-left"
                        >
                          {feature.name}
                        </Link>
                        <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                          <span className="mr-1">{f(messages.created)}</span>
                          <RelativeDateText
                            date={new Date(feature.createdAt * 1000)}
                          />
                        </div>
                      </div>
                      <div className="mb-3">
                        <FeatureIdChip featureId={feature.id} />
                      </div>
                      <div className="mb-2 text-xs">
                        <p>{getClientVersion(feature)}</p>
                      </div>
                      <TagChips tags={feature.tagsList} />
                    </td>
                    <td
                      className={classNames(
                        'w-[10%] pl-5 pr-2 py-3 border-b border-gray-300 bg-white'
                      )}
                    >
                      <div className="text-gray-700">
                        <FlagStatucIcon
                          flagStatus={getFlagStatus(feature, relativeDate)}
                        />
                      </div>
                    </td>
                    <td
                      className={classNames(
                        'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300 bg-white'
                      )}
                    >
                      <Switch
                        enabled={feature.enabled}
                        onChange={() =>
                          onSwitchEnabled(feature.id, !feature.enabled)
                        }
                        size={'small'}
                        readOnly={feature.archived || !editable}
                      />
                    </td>
                    {editable && (
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
                                onArchive(feature);
                                return;
                              case MenuActions.CLONE:
                                onClone(feature);
                                return;
                            }
                          }}
                          menuItems={createMenuItems(feature.archived)}
                        />
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / FEATURE_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);

interface FeatureSearchProps {
  options: FeatureSearchOptions;
  onChange: (options: FeatureSearchOptions) => void;
  onClear: () => void;
  onAdd: () => void;
}

const FeatureSearch: FC<FeatureSearchProps> = memo(
  ({ options, onChange, onClear, onAdd }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const isAccountLoading = useSelector<AppState, boolean>(
      (state) => state.accounts.loading,
      shallowEqual
    );
    const accounts = useSelector<AppState, Account.AsObject[]>(
      (state) => selectAllAccounts(state.accounts),
      shallowEqual
    );
    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );
    const [filterValues, setFilterValues] = useState<Option[]>([]);

    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.ENABLED:
            setFilterValues(enabledOptions);
            return;
          case FilterTypes.ARCHIVED:
            setFilterValues(archivedOptions);
            return;
          case FilterTypes.HAS_EXPERIMENT:
            setFilterValues(hasExperimentOptions);
            return;
          case FilterTypes.MAINTAINER:
            setFilterValues(
              accounts.map((account) => {
                return {
                  value: account.id,
                  label: account.id,
                };
              })
            );
            return;
          case FilterTypes.HAS_PREREQUISITES:
            setFilterValues(hasPrerequisitesOptions);
            return;
          case FilterTypes.TAGS:
            setFilterValues(
              tagsList.map((tag) => ({
                value: tag.id,
                label: tag.id,
              }))
            );
            return;
        }
      },
      [setFilterValues, accounts, tagsList]
    );

    const handleUpdateOption = (
      optionPart: Partial<FeatureSearchOptions>
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
        case FilterTypes.ARCHIVED:
          handleUpdateOption({
            archived: value,
          });
          return;
        case FilterTypes.HAS_EXPERIMENT:
          handleUpdateOption({
            hasExperiment: value,
          });
          return;
        case FilterTypes.HAS_PREREQUISITES:
          handleUpdateOption({
            hasPrerequisites: value,
          });
          return;
        case FilterTypes.MAINTAINER:
          handleUpdateOption({
            maintainerId: value,
          });
          return;
      }
    };
    const handleMultiFilterAdd = (key: string, value: string[]): void => {
      if (key === FilterTypes.TAGS) {
        handleUpdateOption({
          tagIds: value,
        });
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
              placeholder={f(messages.feature.search.placeholder)}
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
              onAddMulti={handleMultiFilterAdd}
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
        {(options.enabled ||
          options.archived ||
          options.hasExperiment ||
          options.hasPrerequisites ||
          options.maintainerId ||
          options.tagIds?.length > 0) && (
          <div className="flex space-x-2 mt-2">
            {options.enabled && (
              <FilterChip
                label={`${f(messages.feature.filter.enabled)}: ${
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
            {options.archived && (
              <FilterChip
                label={`${f(messages.feature.filter.archived)}: ${
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
            {options.hasExperiment && (
              <FilterChip
                label={`${f(messages.feature.filter.hasExperiment)}: ${
                  hasExperimentOptions.find(
                    (option) => option.value === options.hasExperiment
                  ).label
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    hasExperiment: null,
                  })
                }
              />
            )}
            {options.hasPrerequisites && (
              <FilterChip
                label={`${f(messages.feature.filter.hasPrerequisites)}: ${
                  hasPrerequisitesOptions.find(
                    (option) => option.value === options.hasPrerequisites
                  ).label
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    hasPrerequisites: null,
                  })
                }
              />
            )}
            {options.maintainerId && (
              <FilterChip
                label={`${f(messages.feature.filter.maintainer)}: ${
                  options.maintainerId
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    maintainerId: null,
                  })
                }
              />
            )}
            {typeof options.tagIds === 'string' && (
              <FilterChip
                label={`${f(messages.feature.filter.tags)}: ${options.tagIds}`}
                onRemove={() =>
                  handleUpdateOption({
                    tagIds: null,
                  })
                }
              />
            )}
            {Array.isArray(options.tagIds) &&
              options.tagIds.map((tagId) => (
                <FilterChip
                  key={tagId}
                  label={`${f(messages.feature.filter.tags)}: ${tagId}`}
                  onRemove={() =>
                    handleUpdateOption({
                      tagIds: options.tagIds.filter((tId) => tId !== tagId),
                    })
                  }
                />
              ))}
            {(options.enabled ||
              options.archived ||
              options.hasExperiment ||
              options.hasPrerequisites ||
              options.maintainerId ||
              options.tagIds) && (
              <FilterRemoveAllButtonProps onClick={onClear} />
            )}
          </div>
        )}
      </div>
    );
  }
);
