import { formatDate } from '@/utils/date';
import { BanIcon } from '@heroicons/react/solid';
import MUArchiveIcon from '@material-ui/icons/Archive';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { EXPERIMENT_LIST_PAGE_SIZE } from '../../constants/experiment';
import {
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURES,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/experiments';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { ExperimentSearchOptions } from '../../types/experiment';
import { classNames } from '../../utils/css';
import { ActionMenu, MenuActions, MenuItem } from '../ActionMenu';
import { ExperimentSearch, statusOptions } from '../ExperimentSearch';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';

export interface ExperimentListProps {
  searchOptions: ExperimentSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: ExperimentSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (e: Experiment.AsObject) => void;
  onArchive: (e: Experiment.AsObject) => void;
  onStop: (e: Experiment.AsObject) => void;
}

export const ExperimentList: FC<ExperimentListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
    onArchive,
    onStop,
  }) => {
    const { formatMessage: f } = useIntl();
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();
    const editable = useIsEditable();
    const experiments = useSelector<AppState, Experiment.AsObject[]>(
      (state) => selectAll(state.experiments),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.experiments.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.experiments.totalCount,
      shallowEqual
    );
    const createMenuItems = (
      experimentSatus: Experiment.StatusMap[keyof Experiment.StatusMap]
    ): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      const isExperimentWaitingRunnning =
        experimentSatus === Experiment.Status.WAITING ||
        experimentSatus === Experiment.Status.RUNNING;

      if (isExperimentWaitingRunnning) {
        items.push({
          action: MenuActions.STOP,
          name: intl.formatMessage(messages.experiment.stop.button),
          iconElement: <BanIcon />,
        });
      }
      items.push({
        action: MenuActions.ARCHIVE,
        name: intl.formatMessage(messages.experiment.action.archive),
        iconElement: <MUArchiveIcon />,
        disabled: isExperimentWaitingRunnning,
        tooltipMessage: isExperimentWaitingRunnning
          ? intl.formatMessage(messages.experiment.action.archiveTooltip)
          : null,
        alignRight: true,
      });
      return items;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <ExperimentSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : experiments.length == 0 ? (
          searchOptions.q ||
          searchOptions.status ||
          searchOptions.maintainerId ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.experiment.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(
                          messages.experiment.list.noResult.searchKeyword
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
                    title: f(messages.experiment.list.header.title),
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.experiment.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/running-abn-tests?id=running-abn-tests"
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
                {experiments.map((experiment) => {
                  const startAt = new Date(experiment.startAt * 1000);
                  const endAt =
                    experiment.status === Experiment.Status.FORCE_STOPPED
                      ? new Date(Number(experiment.stoppedAt) * 1000)
                      : new Date(experiment.stopAt * 1000);
                  return (
                    <tr key={experiment.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(experiment)}
                          >
                            {experiment.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(experiment.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        <div className="text-xs text-gray-700">
                          {`${f(messages.experiment.period)}: `}
                          {`${formatDate({ date: startAt })} - ${formatDate({
                            date: endAt,
                          })}`}
                        </div>
                      </td>
                      <td
                        className={classNames(
                          'w-[10%] pl-5 pr-2 py-3 border-b border-gray-300',
                          'text-gray-700',
                          'whitespace-nowrap'
                        )}
                      >
                        {
                          statusOptions.find(
                            (option) =>
                              option.value == experiment.status.toString()
                          ).label
                        }
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <button
                          type="button"
                          className="btn-cancel"
                          disabled={false}
                          onClick={() =>
                            history.push(
                              `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${experiment.featureId}${PAGE_PATH_EXPERIMENTS}?experimentId=${experiment.id}`
                            )
                          }
                        >
                          {f(messages.button.result)}
                        </button>
                      </td>
                      {editable && !experiment.archived && (
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
                                  onArchive(experiment);
                                  break;
                                case MenuActions.STOP:
                                  onStop(experiment);
                                  break;
                                default:
                                  return;
                              }
                            }}
                            menuItems={createMenuItems(experiment.status)}
                          />
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / EXPERIMENT_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
