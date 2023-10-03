import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { ENVIRONMENT_LIST_PAGE_SIZE } from '../../constants/environment';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/environments';
import { EnvironmentV2 } from '../../proto/environment/environment_pb';
import { EnvironmentSearchOptions } from '../../types/environment';
import { classNames } from '../../utils/css';
import { EnvironmentSearch } from '../EnvironmentSearch';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';

export interface EnvironmentListProps {
  searchOptions: EnvironmentSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: EnvironmentSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (e: EnvironmentV2.AsObject) => void;
}

export const EnvironmentList: FC<EnvironmentListProps> = memo(
  ({ searchOptions, onChangePage, onChangeSearchOptions, onAdd, onUpdate }) => {
    const { formatMessage: f } = useIntl();
    const environments = useSelector<AppState, EnvironmentV2.AsObject[]>(
      (state) => selectAll(state.environments),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.environments.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.environments.totalCount,
      shallowEqual
    );
    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <EnvironmentSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : environments.length == 0 ? (
          searchOptions.q ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.adminEnvironment.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(
                          messages.adminEnvironment.list.noResult.searchKeyword
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
                    title: f(messages.adminEnvironment.list.header.title),
                  })}
                </h1>
                <a
                  href="https://bucketeer.io/docs/#/environments"
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
            <table className="table-auto leading-normal w-full">
              <tbody className="text-sm">
                {environments.map((environment) => {
                  return (
                    <tr key={environment.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(environment)}
                          >
                            {environment.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(environment.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        <div className="text-xs text-gray-700">
                          {environment.projectId}
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / ENVIRONMENT_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
