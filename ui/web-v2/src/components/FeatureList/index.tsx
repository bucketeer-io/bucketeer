import { DotsVerticalIcon, PlusIcon } from '@heroicons/react/solid';
import MUArchiveIcon from '@material-ui/icons/Archive';
import MUFileCopyIcon from '@material-ui/icons/FileCopy';
import MUUnarchiveIcon from '@material-ui/icons/Unarchive';
import dayjs from 'dayjs';
import React, {
  FC,
  useState,
  memo,
  useCallback,
  Fragment,
  useEffect
} from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { Link, useLocation, useHistory, Prompt } from 'react-router-dom';

import { FEATURE_LIST_PAGE_SIZE } from '../../constants/feature';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  createSearchFilter,
  selectAll as selectAllAccounts,
  changeSearchFilterName,
  changeSearchFilterQuery,
  deleteSearchFilter
} from '../../modules/accounts';
import { selectAll as selectAllFeatures } from '../../modules/features';
import {
  fetchMe,
  useCurrentEnvironment,
  useIsEditable,
  useMe
} from '../../modules/me';
import { selectAll as selectAllTags } from '../../modules/tags';
import { AccountV2 } from '../../proto/account/account_pb';
import { Feature, Tag } from '../../proto/feature/feature_pb';
import { FeatureSearchOptions } from '../../types/feature';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC
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
import { Dialog, Popover, Transition } from '@headlessui/react';
import {
  InformationCircleIcon,
  PencilIcon,
  TrashIcon,
  XIcon
} from '@heroicons/react/outline';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { shortcutFormSchema } from '../../pages/feature/formSchema';
import { AppDispatch } from '../../store';
import {
  FilterTargetType,
  FilterTargetTypeMap
} from '../../proto/account/search_filter_pb';
import {
  stringifySearchParams,
  useSearchParams
} from '../../utils/search-params';
import SaveSvg from '../../assets/svg/save.svg';
import SaveGraySvg from '../../assets/svg/save-gray.svg';
import SaveLargeSvg from '../../assets/svg/save-large.svg';
import { parse, stringify } from 'query-string';
import { HoverPopover } from '../HoverPopover';
import {
  getSearchFilterDialogShown,
  setSearchFilterDialogShown
} from '../../storage/searchFilterDialog';

export enum FlagStatus {
  NEW, // This flag is new and has not been requested yet.
  RECEIVING_REQUESTS, // It is receiving one more requests in the last 7 days.
  INACTIVE // It is not receiving requests for 7 days.
}

enum ConfirmType {
  YES,
  NO
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
    message: intl.formatMessage(messages.feature.sort.newest)
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.feature.sort.oldest)
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.feature.sort.nameAz)
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.feature.sort.nameZa)
  }
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
  HAS_PREREQUISITES = 'prerequisites'
}

const filterOptions: Option[] = [
  {
    value: FilterTypes.TAGS,
    label: intl.formatMessage(messages.tags)
  },
  {
    value: FilterTypes.HAS_PREREQUISITES,
    label: intl.formatMessage(messages.feature.filter.hasPrerequisites)
  },
  {
    value: FilterTypes.MAINTAINER,
    label: intl.formatMessage(messages.feature.filter.maintainer)
  },
  {
    value: FilterTypes.HAS_EXPERIMENT,
    label: intl.formatMessage(messages.feature.filter.hasExperiment)
  },
  {
    value: FilterTypes.ENABLED,
    label: intl.formatMessage(messages.feature.filter.enabled)
  },
  {
    value: FilterTypes.ARCHIVED,
    label: intl.formatMessage(messages.feature.filter.archived)
  }
];

const enabledOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.enabled)
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.disabled)
  }
];

const hasExperimentOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes)
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no)
  }
];

const archivedOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes)
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no)
  }
];

const hasPrerequisitesOptions: Option[] = [
  {
    value: 'true',
    label: intl.formatMessage(messages.yes)
  },
  {
    value: 'false',
    label: intl.formatMessage(messages.no)
  }
];

export interface FeatureListProps {
  searchOptions: FeatureSearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (featureId: string, enabled: boolean) => void;
  onAdd: () => void;
  onChangeSearchOptions: (options: FeatureSearchOptions) => void;
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
    onArchive,
    onClone
  }) => {
    const { formatMessage: f } = useIntl();
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

    const relativeDate = new Date();
    const createMenuItems = (archived: boolean): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      if (archived) {
        items.push({
          action: MenuActions.ARCHIVE,
          name: intl.formatMessage(messages.feature.action.unarchive),
          iconElement: <MUUnarchiveIcon />
        });
        return items;
      }
      items.push({
        action: MenuActions.ARCHIVE,
        name: intl.formatMessage(messages.feature.action.archive),
        iconElement: <MUArchiveIcon />
      });
      items.push({
        action: MenuActions.CLONE,
        name: intl.formatMessage(messages.feature.action.clone),
        iconElement: <MUFileCopyIcon />
      });
      return items;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <FeatureSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
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
                    title: f(messages.feature.list.header.title)
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(messages.feature.list.noResult.searchKeyword)
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
                    title: f(messages.feature.list.header.title)
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.feature.list.noData.description)}
                </p>
                <a
                  href="https://docs.bucketeer.io/feature-flags/creating-feature-flags/create-feature-flag"
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

export interface SearchFilter {
  id: string;
  name: string;
  query: string;
  filterTargetType: FilterTargetTypeMap[keyof FilterTargetTypeMap];
  environmentId: string;
  defaultFilter: boolean;
  selected: boolean;
  saveChanges: boolean;
}

export interface SelectedSearchFilter {
  id: string;
  name: string;
  query: string;
}

interface FeatureSearchProps {
  options: FeatureSearchOptions;
  onChange: (options: FeatureSearchOptions) => void;
  onAdd: () => void;
}

const defaultOptions = {
  page: '1',
  sort: '-createdAt'
};

const FeatureSearch: FC<FeatureSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const me = useMe();
    const searchOptions = useSearchParams();

    const accounts = useSelector<AppState, AccountV2.AsObject[]>(
      (state) => selectAllAccounts(state.accounts),
      shallowEqual
    );
    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );
    const [filterValues, setFilterValues] = useState<Option[]>([]);
    const [open, setOpen] = useState(false);

    const [searchFiltersList, setSearchFiltersList] = useState<SearchFilter[]>(
      []
    );
    const [selectedSearchFilter, setSelectedSearchFilter] =
      useState<SelectedSearchFilter>();
    const [isAddNewFilterActive, setIsAddNewFilterActive] = useState(false);

    const organizationId = currentEnvironment.organizationId;

    const location = useLocation();
    const [unsavedChanges, setUnsavedChanges] = useState(false);
    const [nextLocation, setNextLocation] = useState(null);
    const [isFilterLoading, setIsFilterLoading] = useState(false);
    const [showSaveChangesDialog, setShowSaveChangesDialog] = useState(false);
    const [unsavedSearchFilterId, setUnsavedSearchFilterId] = useState(null);

    const history = useHistory();
    const isSearchFilterDialogShown = getSearchFilterDialogShown();

    const filteredSearchFiltersList =
      me.consoleAccount.searchFiltersList.filter(
        (s) => s.environmentId === currentEnvironment.id
      );

    useEffect(() => {
      // set default options
      if (Object.keys(options).length === 0) {
        onChange(defaultOptions);
      }
    }, [options]);

    useEffect(() => {
      if (filteredSearchFiltersList.length > 0) {
        let updatedFiltersList = [];
        // Make last search filter selected if new search filter is added
        if (
          filteredSearchFiltersList.length > searchFiltersList.length &&
          Object.keys(options).length > 0 &&
          JSON.stringify(options) !== JSON.stringify(defaultOptions)
        ) {
          updatedFiltersList = filteredSearchFiltersList.map((s, i) => ({
            ...s,
            selected: i + 1 === filteredSearchFiltersList.length
          }));
          setSelectedSearchFilter(
            updatedFiltersList[updatedFiltersList.length - 1]
          );
        } else {
          // Save the changes
          const oldSelectedSearchFilter = searchFiltersList.find(
            (s) => s.selected
          );
          updatedFiltersList = filteredSearchFiltersList.map((s) => ({
            ...s,
            selected: oldSelectedSearchFilter?.id === s.id,
            saveChanges: false
          }));
          setSelectedSearchFilter(updatedFiltersList.find((s) => s.selected));
        }
        setSearchFiltersList(updatedFiltersList);
      } else {
        setSearchFiltersList([]);
      }
    }, [me.consoleAccount.searchFiltersList, onChange]);

    useEffect(() => {
      const stringyOptions = stringifySearchParams(options);
      const query = selectedSearchFilter?.query ?? '';

      if (selectedSearchFilter) {
        if (stringyOptions && query) {
          // If the user has clicked on a Feature flag navbar
          if (stringyOptions === stringify(defaultOptions)) {
            setSelectedSearchFilter(null);
            setSearchFiltersList(
              searchFiltersList.map((s) => ({
                ...s,
                selected: false,
                saveChanges: false
              }))
            );
          } else {
            setSearchFiltersList(
              searchFiltersList.map((s) => ({
                ...s,
                saveChanges:
                  s.id === selectedSearchFilter.id && query !== stringyOptions
                    ? true
                    : false
              }))
            );
          }
        }
      }
    }, [selectedSearchFilter, options]);

    useEffect(() => {
      // If there are any search filters with changes, mark the page as having unsaved changes
      const findSaveChanges = searchFiltersList.find((s) => s.saveChanges);
      setUnsavedChanges(findSaveChanges ? true : false);
    }, [searchFiltersList]);

    useEffect(() => {
      if (!isSearchFilterDialogShown) {
        // If there are unsaved changes when the user tries to leave the page, show the confirmation dialog
        const handleBeforeUnload = (e) => {
          if (unsavedChanges) {
            // The standard way to trigger the confirmation dialog
            e.preventDefault();
            e.returnValue = ''; // For most browsers, this will trigger the default confirmation dialog
          }
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        return () => {
          window.removeEventListener('beforeunload', handleBeforeUnload);
        };
      }
    }, [unsavedChanges, isSearchFilterDialogShown]);

    useEffect(() => {
      if (
        !selectedSearchFilter &&
        JSON.stringify(options) !== JSON.stringify(defaultOptions)
      ) {
        setUnsavedChanges(true);
      } else {
        setUnsavedChanges(false);
      }
    }, [options, selectedSearchFilter]);

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
                  value: account.email,
                  label: account.email
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
                value: tag.name,
                label: tag.name
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
            enabled: value
          });
          return;
        case FilterTypes.ARCHIVED:
          handleUpdateOption({
            archived: value
          });
          return;
        case FilterTypes.HAS_EXPERIMENT:
          handleUpdateOption({
            hasExperiment: value
          });
          return;
        case FilterTypes.HAS_PREREQUISITES:
          handleUpdateOption({
            hasPrerequisites: value
          });
          return;
        case FilterTypes.MAINTAINER:
          handleUpdateOption({
            maintainerId: value
          });
          return;
      }
    };

    const handleMultiFilterAdd = (key: string, value: string[]): void => {
      if (key === FilterTypes.TAGS) {
        handleUpdateOption({
          tagIds: value
        });
      }
    };

    const refetchMe = () => {
      dispatch(
        fetchMe({
          organizationId: currentEnvironment.organizationId
        })
      ).then(() => setIsFilterLoading(false));
    };

    const handleFormSubmit = useCallback(
      async (data) => {
        setIsFilterLoading(true);
        setOpen(false);

        const query = stringifySearchParams({
          ...options,
          page: 1
        });
        const filterTargetType = FilterTargetType.FEATURE_FLAG;
        const defaultFilter = false;
        const environmentId = currentEnvironment.id;

        if (selectedSearchFilter && !isAddNewFilterActive) {
          dispatch(
            changeSearchFilterName({
              id: selectedSearchFilter.id,
              name: data.name,
              email: me.consoleAccount.email,
              environmentId,
              organizationId: me.consoleAccount.organization.id
            })
          ).then(() => {
            if (unsavedChanges) {
              handleSaveChanges(selectedSearchFilter.id);
            } else {
              refetchMe();
            }
          });
        } else {
          dispatch(
            createSearchFilter({
              name: data.name,
              query,
              filterTargetType,
              environmentId,
              defaultFilter,
              organizationId: currentEnvironment.organizationId,
              email: me.consoleAccount.email
            })
          ).then(() => {
            refetchMe();
          });
        }
        setIsAddNewFilterActive(false);
      },
      [searchOptions, selectedSearchFilter, isAddNewFilterActive]
    );

    const handleDeleteSearchFilter = (id: string) => {
      setIsFilterLoading(true);
      if (id === selectedSearchFilter?.id || searchFiltersList.length === 1) {
        onChange(defaultOptions);
        setSelectedSearchFilter(null);
      }

      dispatch(
        deleteSearchFilter({
          id,
          email: me.consoleAccount.email,
          organizationId,
          environmentId: currentEnvironment.id
        })
      ).then(() => refetchMe());
    };

    const handleSearchFilter = (id: string) => {
      if (selectedSearchFilter?.id !== id) {
        const findSearchFilter = searchFiltersList.find((s) => s.id === id);
        setSelectedSearchFilter(findSearchFilter);
        setSearchFiltersList(
          searchFiltersList.map((s) => ({ ...s, selected: s.id === id }))
        );
        onChange(parse(findSearchFilter.query));
        setUnsavedSearchFilterId(null);
      }
    };

    const handleSaveChanges = (id: string) => {
      setIsFilterLoading(true);
      dispatch(
        changeSearchFilterQuery({
          id,
          query: stringifySearchParams(options),
          email: me.consoleAccount.email,
          organizationId,
          environmentId: currentEnvironment.id
        })
      ).then(() => refetchMe());
    };

    const handleClose = () => {
      setOpen(false);
      setIsAddNewFilterActive(false);
    };

    const handleConfirm = (confirmType: ConfirmType) => {
      setSearchFilterDialogShown(true);
      setShowSaveChangesDialog(false);
      if (confirmType === ConfirmType.YES) {
        if (unsavedSearchFilterId) {
          handleSearchFilter(unsavedSearchFilterId);
          return;
        }

        setUnsavedChanges(false);

        const { pathname: nextPathname } = nextLocation;

        const isNew =
          `/${nextPathname.substring(nextPathname.lastIndexOf('/') + 1)}` ==
          PAGE_PATH_NEW;

        onChange(defaultOptions);
        setSelectedSearchFilter(null);
        setSearchFiltersList(
          filteredSearchFiltersList.map((s) => ({
            ...s,
            selected: false,
            saveChanges: false
          }))
        );

        // If the user clicked on the "Add" button, we should show add form
        if (isNew) {
          onAdd();
        } else if (location.pathname !== nextPathname) {
          history.push({
            pathname: nextPathname + nextLocation.search
          });
        }
      }
    };

    const handleNavigation = (nextLocation) => {
      const newPathname = location.pathname.split('/').slice(1, 3).join('/');
      const newNextPathname = nextLocation.pathname
        .split('/')
        .slice(1, 3)
        .join('/');

      const targetingPathname = `/${nextLocation.pathname.split('/')?.[4]}`;

      if (
        newNextPathname && // If the user is logged out and re login
        (newPathname !== newNextPathname || // If the user is trying to go to a different page
          targetingPathname === PAGE_PATH_FEATURE_TARGETING) // If the user is trying to go to a targeting page
      ) {
        setNextLocation(nextLocation); // Save the location the user is trying to go to
        setShowSaveChangesDialog(true); // Show custom confirmation popup
        return false; // Block navigation for now
      }
      return true; // Allow navigation if no unsaved changes
    };

    return (
      <div
        className={classNames(
          'w-full',
          'px-5 py-5 sticky top-0',
          'z-10 border-b border-gray-300'
        )}
      >
        <div className={classNames('w-full', 'flex flex-row')}>
          <div className="flex-none w-72">
            <SearchInput
              placeholder={f(messages.feature.search.placeholder)}
              value={options.q}
              onChange={(query: string) =>
                handleUpdateOption({
                  q: query
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
          <div className="flex gap-2 flex-wrap w-full items-center">
            {searchFiltersList.map((searchFilter) => (
              <Popover className="relative flex" key={searchFilter.id}>
                <div
                  className={classNames(
                    'flex items-center space-x-1.5 rounded p-[6px] text-sm cursor-pointer transition-colors duration-200',
                    searchFilter.selected ? 'bg-purple-100' : 'bg-gray-50'
                  )}
                  onClick={() => handleSearchFilter(searchFilter.id)}
                >
                  <span className="text-primary truncate max-w-[180px]">
                    {searchFilter.name}
                  </span>
                  {searchFilter.saveChanges && (
                    <div className="text-primary opacity-80">
                      <SaveSvg />
                    </div>
                  )}
                  {unsavedChanges &&
                  searchFilter.id !== selectedSearchFilter?.id ? (
                    <DotsVerticalIcon width={14} className="text-primary" />
                  ) : (
                    <Popover.Button className="h-full outline-none">
                      <DotsVerticalIcon width={14} className="text-primary" />
                    </Popover.Button>
                  )}
                </div>
                <Popover.Panel className="absolute bg-white left-0 rounded-lg p-1 whitespace-nowrap drop-shadow z-50 top-10">
                  {({ close }) => (
                    <div>
                      {searchFilter.saveChanges && (
                        <button
                          onClick={() => {
                            close();
                            handleSaveChanges(searchFilter.id);
                          }}
                          className="flex w-full space-x-2 px-2 py-1.5 items-center hover:bg-gray-100"
                        >
                          <SaveGraySvg />
                          <span className="text-sm pl-[2px]">Save Changes</span>
                        </button>
                      )}
                      <button
                        onClick={() => {
                          setOpen(true);
                          handleSearchFilter(searchFilter.id);
                        }}
                        className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                      >
                        <PencilIcon width={18} />
                        <span className="text-sm">
                          {f(messages.saveChanges.editShortcut)}
                        </span>
                      </button>
                      <button
                        onClick={() => {
                          close();
                          handleDeleteSearchFilter(searchFilter.id);
                        }}
                        className="flex space-x-3 w-full px-2 py-1.5 items-center hover:bg-gray-100"
                      >
                        <TrashIcon width={18} className="text-red-500" />
                        <span className="text-red-500 text-sm">
                          {f(messages.saveChanges.deleteShortcut)}
                        </span>
                      </button>
                    </div>
                  )}
                </Popover.Panel>
              </Popover>
            ))}
            {unsavedChanges && (
              <div className="flex space-x-1">
                <button
                  className="bg-gray-50 p-[6px] rounded hover:bg-purple-100 transition-colors disabled:cursor-not-allowed disabled:opacity-50"
                  onClick={() => {
                    setOpen(true);
                    setIsAddNewFilterActive(true);
                  }}
                >
                  <PlusIcon width={20} className="text-primary" />
                </button>
                <HoverPopover
                  render={() => {
                    return (
                      <div
                        className={classNames(
                          'border shadow-sm bg-white text-gray-500 p-1',
                          'text-xs rounded whitespace-normal break-words w-64'
                        )}
                      >
                        {f(messages.saveChanges.addShortcutTooltip)}
                      </div>
                    );
                  }}
                >
                  <div className={classNames('hover:text-gray-500')}>
                    <InformationCircleIcon
                      className="w-5 h-5 text-gray-400"
                      aria-hidden="true"
                    />
                  </div>
                </HoverPopover>
              </div>
            )}
            {isFilterLoading && <div className="spinner" />}
          </div>
          {open && (
            <AddEditShortcutModal
              open={open}
              close={handleClose}
              name={isAddNewFilterActive ? '' : selectedSearchFilter?.name}
              handleFormSubmit={handleFormSubmit}
            />
          )}
          {!isSearchFilterDialogShown && (
            <>
              <Prompt
                when={unsavedChanges && showSaveChangesDialog === false}
                message={(location) => {
                  return handleNavigation(location);
                }}
              />
              <SaveChangesDialog
                open={showSaveChangesDialog}
                close={() => {
                  setShowSaveChangesDialog(false);
                  setUnsavedSearchFilterId(null);
                  setSearchFilterDialogShown(true);
                }}
                onConfirm={handleConfirm}
              />
            </>
          )}
          <div className="flex-none -mr-2">
            <SortSelect
              sortKey={options.sort}
              sortItems={sortItems}
              onChange={(sort: string) =>
                handleUpdateOption({
                  sort: sort
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
                    enabled: null
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
                    archived: null
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
                    hasExperiment: null
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
                    hasPrerequisites: null
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
                    maintainerId: null
                  })
                }
              />
            )}
            {typeof options.tagIds === 'string' && (
              <FilterChip
                label={`${f(messages.feature.filter.tags)}: ${options.tagIds}`}
                onRemove={() =>
                  handleUpdateOption({
                    tagIds: null
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
                      tagIds: options.tagIds.filter((tId) => tId !== tagId)
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
              <FilterRemoveAllButtonProps
                onClick={() =>
                  handleUpdateOption({
                    enabled: null,
                    archived: null,
                    hasExperiment: null,
                    hasPrerequisites: null,
                    maintainerId: null,
                    tagIds: null
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

interface SaveChangesDialogProps {
  open: boolean;
  close: () => void;
  onConfirm: (confirmType: ConfirmType) => void;
}

const SaveChangesDialog = ({
  open,
  close,
  onConfirm
}: SaveChangesDialogProps) => {
  const { formatMessage: f } = useIntl();

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={close}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>
        <form className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-[500px]">
                <div className="flex items-center justify-between px-4 py-5 border-b">
                  <p className="text-xl font-bold">
                    {f(messages.saveChanges.saveChangesBeforeExiting.title)}
                  </p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={close}
                  />
                </div>
                <div className="p-4 space-y-4 py-5 px-11 flex flex-col items-center">
                  <SaveLargeSvg />
                  <p className="text-gray-500 text-center">
                    {f(
                      messages.saveChanges.saveChangesBeforeExiting.description
                    )}
                  </p>
                </div>
                <div className="p-4 flex justify-end border-t space-x-4">
                  <button
                    type="button"
                    className="btn-cancel !min-w-max"
                    disabled={false}
                    onClick={() => onConfirm(ConfirmType.NO)}
                  >
                    {f(messages.no)}
                  </button>
                  <button
                    type="button"
                    className="btn-submit"
                    onClick={() => onConfirm(ConfirmType.YES)}
                  >
                    {f(messages.yes)}
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </form>
      </Dialog>
    </Transition.Root>
  );
};

interface AddEditShortcutModalProps {
  open: boolean;
  close: () => void;
  name: string;
  handleFormSubmit: (data: { name: string }) => void;
}

const AddEditShortcutModal = ({
  open,
  close,
  name,
  handleFormSubmit
}: AddEditShortcutModalProps) => {
  const { formatMessage: f } = useIntl();

  const {
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
    reset: resetShortcut
  } = useForm({
    resolver: yupResolver(shortcutFormSchema),
    defaultValues: {
      name: ''
    },
    mode: 'onChange'
  });

  useEffect(() => {
    if (name) {
      resetShortcut({
        name
      });
    }
  }, [name]);

  const handleSave = (data) => {
    handleFormSubmit({
      name: data.name
    });
  };

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={close}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>
        <form className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-[542px]">
                <div className="flex items-center justify-between px-4 py-5 border-b">
                  <p className="text-xl font-medium">
                    {name
                      ? f(messages.saveChanges.editShortcut)
                      : f(messages.saveChanges.addShortcut)}
                  </p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={close}
                  />
                </div>
                <div className="p-4 space-y-4">
                  <p className="text-gray-500">
                    {f(messages.saveChanges.shortcutDescription)}
                  </p>
                  <div className="space-y-1">
                    <label htmlFor="name" className="flex items-center">
                      <span className="input-label">{f({ id: 'name' })}</span>
                      <span className="text-red-500">*</span>
                    </label>
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="p-4 flex justify-end border-t space-x-4">
                  <button
                    type="button"
                    className="btn-cancel"
                    disabled={false}
                    onClick={close}
                  >
                    {f(messages.button.cancel)}
                  </button>
                  <button
                    type="submit"
                    className="btn-submit"
                    disabled={!isValid || isSubmitting}
                    onClick={handleSubmit(handleSave)}
                  >
                    {f(messages.button.save)}
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </form>
      </Dialog>
    </Transition.Root>
  );
};
