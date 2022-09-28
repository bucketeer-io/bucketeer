import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { ACCOUNT_LIST_PAGE_SIZE } from '../../constants/account';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/accounts';
import { useIsOwner } from '../../modules/me';
import { Account } from '../../proto/account/account_pb';
import { AccountSearchOptions } from '../../types/account';
import { classNames } from '../../utils/css';
import { roleOptions } from '../AccountAddForm';
import { AccountSearch } from '../AccountSearch';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { Switch } from '../Switch';

export interface AccountListProps {
  searchOptions: AccountSearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (accountId: string, enabled: boolean) => void;
  onChangeSearchOptions: (options: AccountSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (a: Account.AsObject) => void;
}

export const AccountList: FC<AccountListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onSwitchEnabled,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
  }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsOwner();
    const accounts = useSelector<AppState, Account.AsObject[]>(
      (state) => selectAll(state.accounts),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.accounts.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.accounts.totalCount,
      shallowEqual
    );

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <AccountSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : accounts.length == 0 ? (
          searchOptions.q || searchOptions.role || searchOptions.enabled ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.account.list.header.title),
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(
                          messages.account.list.noResult.searchKeyword
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
                    title: f(messages.account.list.header.title),
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.account.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/managing-teams?id=environment-account"
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
                {accounts.map((account) => {
                  return (
                    <tr key={account.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1">
                          <button
                            className="link whitespace-nowrap"
                            onClick={() => onUpdate(account)}
                          >
                            {account.id}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(account.createdAt * 1000)}
                            />
                          </div>
                        </div>
                      </td>
                      <td
                        className={classNames(
                          'w-[10%] pl-5 pr-2 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <div className="flex justify-center items-center rounded-md w-16 h-6 bg-gray-200">
                          <span className="text-xs text-gray-700">
                            {
                              roleOptions.find(
                                (o) => o.value === account.role.toString()
                              ).label
                            }
                          </span>
                        </div>
                      </td>
                      <td
                        className={classNames(
                          'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap'
                        )}
                      >
                        <Switch
                          enabled={!account.disabled}
                          onChange={() =>
                            onSwitchEnabled(account.id, account.disabled)
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
              maxPage={Math.ceil(totalCount / ACCOUNT_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
