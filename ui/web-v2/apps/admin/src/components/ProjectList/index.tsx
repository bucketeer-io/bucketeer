import MUAutorenewIcon from '@material-ui/icons/Autorenew';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { PROJECT_LIST_PAGE_SIZE } from '../../constants/project';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/projects';
import { Project } from '../../proto/environment/project_pb';
import { ProjectSearchOptions } from '../../types/project';
import { classNames } from '../../utils/css';
import { ActionMenu, MenuActions, MenuItem } from '../ActionMenu';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { ProjectSearch } from '../ProjectSearch';
import { RelativeDateText } from '../RelativeDateText';
import { Switch } from '../Switch';

export interface ProjectListProps {
  searchOptions: ProjectSearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (project: Project.AsObject) => void;
  onChangeSearchOptions: (options: ProjectSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (p: Project.AsObject) => void;
  onConvert: (p: Project.AsObject) => void;
}

export const ProjectList: FC<ProjectListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onSwitchEnabled,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
    onConvert,
  }) => {
    const { formatMessage: f, formatDate } = useIntl();
    const projects = useSelector<AppState, Project.AsObject[]>(
      (state) => selectAll(state.projects),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.projects.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.projects.totalCount,
      shallowEqual
    );
    const createMenuItems = (project: Project.AsObject): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      items.push({
        action: MenuActions.CONVERT_PROJECT,
        name: intl.formatMessage(messages.adminProject.action.convertProject),
        iconElement: <MUAutorenewIcon />,
        disabled: !project.trial,
      });
      return items;
    };
    const addDays = (date: Date, days: number): Date => {
      const result = new Date(date);
      result.setDate(result.getDate() + days);
      return result;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <ProjectSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : projects.length == 0 ? (
          searchOptions.q || searchOptions.enabled ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.adminProject.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(
                          messages.adminProject.list.noResult.searchKeyword
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
                    title: f(messages.adminProject.list.header.title),
                  })}
                </h1>
                <a
                  href="https://bucketeer.io/docs/#/projects"
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
                {projects.map((project) => {
                  return (
                    <tr key={project.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(project)}
                          >
                            {project.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(project.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        {project.trial && (
                          <div className="text-xs text-gray-700">
                            {`${f(messages.adminProject.trialPeriod)}: `}
                            {`${formatDate(new Date(project.createdAt * 1000))}
                          - ${formatDate(
                            addDays(new Date(project.createdAt * 1000), 179)
                          )}`}
                          </div>
                        )}
                      </td>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <span className="mr-1 text-sm">
                          {project.creatorEmail}
                        </span>
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <Switch
                          enabled={!project.disabled}
                          onChange={() => onSwitchEnabled(project)}
                          size={'small'}
                        />
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <ActionMenu
                          onClickAction={(action) => {
                            switch (action) {
                              case MenuActions.CONVERT_PROJECT:
                                onConvert(project);
                                return;
                            }
                          }}
                          menuItems={createMenuItems(project)}
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / PROJECT_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
