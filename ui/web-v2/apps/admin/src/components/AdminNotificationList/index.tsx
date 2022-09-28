import MUDeleteIcon from '@material-ui/icons/Delete';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { NOTIFICATION_LIST_PAGE_SIZE } from '../../constants/notification';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable } from '../../modules/me';
import { selectAll } from '../../modules/notifications';
import { Subscription } from '../../proto/notification/subscription_pb';
import { NotificationSearchOptions } from '../../types/notification';
import { classNames } from '../../utils/css';
import { ListSkeleton } from '../ListSkeleton';
import { NotificationSearch } from '../NotificationSearch';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { Switch } from '../Switch';

export interface AdminNotificationListProps {
  searchOptions: NotificationSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: NotificationSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (notification: Subscription.AsObject) => void;
  onSwitch: (notification: Subscription.AsObject) => void;
  onDelete: (notification: Subscription.AsObject) => void;
}

export const AdminNotificationList: FC<AdminNotificationListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
    onSwitch,
    onDelete,
  }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const notificationList = useSelector<AppState, Subscription.AsObject[]>(
      (state) => selectAll(state.adminNotification),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.adminNotification.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.adminNotification.totalCount,
      shallowEqual
    );

    return (
      <div className="w-full">
        <div className="flex items-stretch mb-8 text-sm">
          <p className="text-gray-700">
            {f(messages.notification.list.header.description)}
          </p>
          <a
            className="link"
            target="_blank"
            href="https://bucketeer.io/docs#/notification"
            rel="noreferrer"
          >
            {f(messages.readMore)}
          </a>
        </div>
        <div className="bg-white border border-gray-300 rounded-md">
          <div>
            <NotificationSearch
              options={searchOptions}
              onChange={onChangeSearchOptions}
              onAdd={onAdd}
            />
          </div>
          {isLoading ? (
            <ListSkeleton />
          ) : notificationList.length == 0 ? (
            searchOptions.q || searchOptions.enabled ? (
              <div className="my-10 flex justify-center">
                <div className="text-gray-700">
                  <h1 className="text-lg">
                    {f(messages.noResult.title, {
                      title: f(messages.notification.list.header.title),
                    })}
                  </h1>
                  <div className="flex justify-center mt-4">
                    <ul className="list-disc">
                      <li>
                        {f(messages.noResult.searchByKeyword, {
                          keyword: f(
                            messages.notification.list.noResult.searchKeyword
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
                      title: f(messages.notification.list.header.title),
                    })}
                  </h1>
                  <p className="mt-5">
                    {f(messages.notification.list.noData.description)}
                  </p>
                </div>
              </div>
            )
          ) : (
            <div>
              <table className="table-auto leading-normal">
                <tbody className="text-sm">
                  {notificationList.map((notification) => (
                    <tr key={notification.id} className={classNames('p-2')}>
                      <td className="px-5 py-2 border-b">
                        <div className="flex pb-1 text-primary">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(notification)}
                          >
                            {notification.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(notification.createdAt * 1000)}
                            />
                          </div>
                        </div>
                      </td>
                      <td
                        className={classNames(
                          'w-[10%] px-5 py-3 border-b border-gray-300',
                          'whitespace-nowrap text-right'
                        )}
                      >
                        <Switch
                          enabled={!notification.disabled}
                          onChange={() => onSwitch(notification)}
                          size={'small'}
                          readOnly={!editable}
                        />
                      </td>
                      {editable && (
                        <td
                          className={classNames(
                            'w-[1%] px-5 py-3 border-b border-gray-300',
                            'whitespace-nowrap text-gray-500'
                          )}
                        >
                          <button
                            className="text-gray-500"
                            onClick={() => onDelete(notification)}
                          >
                            <MUDeleteIcon />
                          </button>
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
              <Pagination
                maxPage={Math.ceil(totalCount / NOTIFICATION_LIST_PAGE_SIZE)}
                currentPage={
                  searchOptions.page ? Number(searchOptions.page) : 1
                }
                onChange={onChangePage}
              />
            </div>
          )}
        </div>
      </div>
    );
  }
);
