import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll as selectAllAccounts } from '../../modules/accounts';
import { useIsEditable } from '../../modules/me';
import { Account } from '../../proto/account/account_pb';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { ExperimentSearchOptions } from '../../types/experiment';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { classNames } from '../../utils/css';
import { FilterChip } from '../FilterChip';
import { Option, FilterPopover } from '../FilterPopover';
import { FilterRemoveAllButtonProps } from '../FilterRemoveAllButton';
import { SearchInput } from '../SearchInput';
import { SortItem, SortSelect } from '../SortSelect';

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
  MAINTAINER = 'maintainer',
  STATUS = 'status',
  ARCHIVED = 'archived',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.MAINTAINER,
    label: intl.formatMessage(messages.experiment.filter.maintainer),
  },
  {
    value: FilterTypes.STATUS,
    label: intl.formatMessage(messages.experiment.filter.status),
  },
  {
    value: FilterTypes.ARCHIVED,
    label: intl.formatMessage(messages.experiment.filter.archived),
  },
];

export const statusOptions: Option[] = [
  {
    value: Experiment.Status.WAITING.toString(),
    label: intl.formatMessage(messages.experiment.status.waiting),
  },
  {
    value: Experiment.Status.RUNNING.toString(),
    label: intl.formatMessage(messages.experiment.status.running),
  },
  {
    value: Experiment.Status.STOPPED.toString(),
    label: intl.formatMessage(messages.experiment.status.stopped),
  },
  {
    value: Experiment.Status.FORCE_STOPPED.toString(),
    label: intl.formatMessage(messages.experiment.status.forceStopped),
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

export interface ExperimentSearchProps {
  options: ExperimentSearchOptions;
  onChange: (options: ExperimentSearchOptions) => void;
  onAdd: () => void;
}

export const ExperimentSearch: FC<ExperimentSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
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
        }
      },
      [setFilterValues, accounts]
    );

    const handleUpdateOption = (
      optionPart: Partial<ExperimentSearchOptions>
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
        case FilterTypes.MAINTAINER:
          handleUpdateOption({
            maintainerId: value,
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
        {(options.status || options.archived || options.maintainerId) && (
          <div className="flex space-x-2 pt-2">
            {options.status && (
              <FilterChip
                label={`${f(messages.experiment.filter.status)}: ${
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
                label={`${f(messages.experiment.filter.archived)}: ${
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
            {options.maintainerId && (
              <FilterChip
                label={`${f(messages.experiment.filter.maintainer)}: ${
                  options.maintainerId
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    maintainerId: null,
                  })
                }
              />
            )}
            {(options.status || options.archived || options.maintainerId) && (
              <FilterRemoveAllButtonProps
                onClick={() =>
                  handleUpdateOption({
                    status: null,
                    maintainerId: null,
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
