import { PlusIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback, useState, useEffect } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  listProjects,
  selectAll as selectAllProjects,
} from '../../modules/projects';
import { Project } from '../../proto/environment/project_pb';
import { AppDispatch } from '../../store';
import { EnvironmentSearchOptions } from '../../types/environment';
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
    message: intl.formatMessage(messages.adminEnvironment.sort.newest),
  },
  {
    key: SORT_OPTIONS_CREATED_AT_ASC,
    message: intl.formatMessage(messages.adminEnvironment.sort.oldest),
  },
  {
    key: SORT_OPTIONS_NAME_ASC,
    message: intl.formatMessage(messages.adminEnvironment.sort.nameAz),
  },
  {
    key: SORT_OPTIONS_NAME_DESC,
    message: intl.formatMessage(messages.adminEnvironment.sort.nameZa),
  },
];

export enum FilterTypes {
  PROJECT = 'project',
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.PROJECT,
    label: intl.formatMessage(messages.adminEnvironment.filter.project),
  },
];

export interface EnvironmentSearchProps {
  options: EnvironmentSearchOptions;
  onChange: (options: EnvironmentSearchOptions) => void;
  onAdd: () => void;
}

export const EnvironmentSearch: FC<EnvironmentSearchProps> = memo(
  ({ options, onChange, onAdd }) => {
    const { formatMessage: f } = useIntl();
    const projects = useSelector<AppState, Project.AsObject[]>(
      (state) => selectAllProjects(state.projects),
      shallowEqual
    );
    const [filterValues, setFilterValues] = useState<Option[]>([]);

    const handleFilterKeyChange = useCallback(
      (key: string): void => {
        switch (key) {
          case FilterTypes.PROJECT:
            setFilterValues(
              projects.map((project) => {
                return {
                  value: project.id,
                  label: project.id,
                };
              })
            );
            return;
        }
      },
      [setFilterValues, projects]
    );

    const handleUpdateOption = (
      optionPart: Partial<EnvironmentSearchOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const handleFilterAdd = (key: string, value?: string): void => {
      switch (key) {
        case FilterTypes.PROJECT:
          handleUpdateOption({
            projectId: value,
          });
          return;
      }
    };
    const dispatch = useDispatch<AppDispatch>();
    useEffect(() => {
      dispatch(
        listProjects({
          pageSize: 0,
          cursor: '',
        })
      );
    }, [dispatch]);
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
              placeholder={f(messages.adminEnvironment.search.placeholder)}
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
        </div>
        {options.projectId && (
          <div className="flex space-x-2 pt-2">
            {options.projectId && (
              <FilterChip
                label={`${f(messages.adminEnvironment.filter.project)}: ${
                  options.projectId
                }`}
                onRemove={() =>
                  handleUpdateOption({
                    projectId: null,
                  })
                }
              />
            )}
            {options.projectId && (
              <FilterRemoveAllButtonProps
                onClick={() =>
                  handleUpdateOption({
                    projectId: null,
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
