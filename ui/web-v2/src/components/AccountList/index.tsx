import { getRoleV1 } from '../../pages/account';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { ACCOUNT_LIST_PAGE_SIZE } from '../../constants/account';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/accounts';
import { useCurrentEnvironment, useIsOwner } from '../../modules/me';
import { AccountV2 } from '../../proto/account/account_pb';
import { AccountSearchOptions } from '../../types/account';
import { classNames } from '../../utils/css';
import { AccountSearch } from '../AccountSearch';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { Switch } from '../Switch';
import { TagChips } from '../TagsChips';

export interface AccountListProps {
  searchOptions: AccountSearchOptions;
  onChangePage: (page: number) => void;
  onSwitchEnabled: (accountId: string, enabled: boolean) => void;
  onChangeSearchOptions: (options: AccountSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (a: AccountV2.AsObject) => void;
}

export const AccountList: FC<AccountListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onSwitchEnabled,
    onChangeSearchOptions,
    onAdd,
    onUpdate
  }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsOwner();
    const accounts = useSelector<AppState, AccountV2.AsObject[]>(
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
    const currentEnvironment = useCurrentEnvironment();

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
                    title: f(messages.account.list.header.title)
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(messages.account.list.noResult.searchKeyword)
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
                    title: f(messages.account.list.header.title)
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.account.list.noData.description)}
                </p>
                <a
                  href="https://docs.bucketeer.io/getting-started/create-an-account"
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
                  const envRole = account.environmentRolesList.find(
                    (r) => r.environmentId === currentEnvironment.id
                  );
                  const accountDisabled =
                    !envRole ||
                    envRole.role ===
                      AccountV2.Role.Environment.ENVIRONMENT_UNASSIGNED;
                  const orgRole = account.organizationRole;
                  return (
                    <tr key={account.email} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex pb-1.5">
                          <button
                            className="link whitespace-nowrap"
                            onClick={() => onUpdate(account)}
                          >
                            {account.email}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(account.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        <TagChips tags={account.tagsList} />
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
                              getRoleV1(
                                orgRole,
                                envRole?.role ??
                                  AccountV2.Role.Environment
                                    .ENVIRONMENT_UNASSIGNED
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
                          enabled={!accountDisabled}
                          onChange={() =>
                            onSwitchEnabled(account.email, accountDisabled)
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
