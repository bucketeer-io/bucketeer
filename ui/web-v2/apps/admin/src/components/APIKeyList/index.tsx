import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { APIKEY_LIST_PAGE_SIZE } from '../../constants/apiKey';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/apiKeys';
import { useCurrentEnvironment, useIsOwner } from '../../modules/me';
import { APIKey } from '../../proto/account/api_key_pb';
import { APIKeySearchOptions } from '../../types/apiKey';
import { classNames } from '../../utils/css';
import { APIKeySearch } from '../APIKeySearch';
import { CopyChip } from '../CopyChip';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { Switch } from '../Switch';

export interface APIKeyListProps {
  searchOptions: APIKeySearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (
    apiKeyId: string,
    apiKeyName: string,
    enabled: boolean
  ) => void;
  onChangeSearchOptions: (options: APIKeySearchOptions) => void;
  onAdd: () => void;
  onUpdate: (a: APIKey.AsObject) => void;
}

export const APIKeyList: FC<APIKeyListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onSwitchEnabled,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
  }) => {
    const { formatMessage: f, formatDate, formatTime } = useIntl();
    const editable = useIsOwner();
    const currentEnvironment = useCurrentEnvironment();
    const apiKeys = useSelector<AppState, APIKey.AsObject[]>(
      (state) => selectAll(state.apiKeys),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.apiKeys.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.apiKeys.totalCount,
      shallowEqual
    );

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <APIKeySearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : apiKeys.length == 0 ? (
          searchOptions.q || searchOptions.enabled ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.apiKey.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(messages.apiKey.list.noResult.searchKeyword),
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
                    title: f(messages.apiKey.list.header.title),
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.apiKey.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/tutorial-top?id=create-your-api-key"
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
                {apiKeys.map((apiKey) => {
                  return (
                    <tr key={apiKey.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(apiKey)}
                          >
                            {apiKey.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(apiKey.createdAt * 1000)}
                            />
                          </div>
                        </div>
                      </td>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <CopyChip text={apiKey.id}>
                          <span
                            className={classNames(
                              'p-1.5 rounded-lg text-xs bg-gray-200',
                              'cursor-pointer'
                            )}
                          >
                            {apiKey.id}
                          </span>
                        </CopyChip>
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <Switch
                          enabled={!apiKey.disabled}
                          onChange={() =>
                            onSwitchEnabled(
                              apiKey.id,
                              apiKey.name,
                              apiKey.disabled
                            )
                          }
                          size={'small'}
                          readOnly={!editable}
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / APIKEY_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
