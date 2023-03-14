import MUDeleteIcon from '@material-ui/icons/Delete';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { WEBHOOK_LIST_PAGE_SIZE } from '../../constants/webhook';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable } from '../../modules/me';
import { selectAll } from '../../modules/webhooks';
import { Webhook } from '../../proto/autoops/webhook_pb';
import { WebhookSearchOptions } from '../../types/webhook';
import { classNames } from '../../utils/css';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { WebhookSearch } from '../WebhookSearch';

export interface WebhookListProps {
  searchOptions: WebhookSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: WebhookSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (webhook: Webhook.AsObject) => void;
  onDelete: (webhook: Webhook.AsObject) => void;
}

export const WebhookList: FC<WebhookListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onChangeSearchOptions,
    onAdd,
    onUpdate,
    onDelete,
  }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const webhookList = useSelector<AppState, Webhook.AsObject[]>(
      (state) => selectAll(state.webhook),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.webhook.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.webhook.totalCount,
      shallowEqual
    );

    return (
      <div className="w-full">
        <div className="flex items-stretch mb-8 text-sm">
          <p className="text-gray-700">
            {f(messages.webhook.list.header.description)}
          </p>
        </div>
        <div className="min-w-max bg-white border border-gray-300 rounded-md">
          <div>
            <WebhookSearch
              options={searchOptions}
              onChange={onChangeSearchOptions}
              onAdd={onAdd}
            />
          </div>
          {isLoading ? (
            <ListSkeleton />
          ) : webhookList.length == 0 ? (
            searchOptions.q ? (
              <div className="my-10 flex justify-center">
                <div className="text-gray-700">
                  <h1 className="text-lg">
                    {f(messages.noResult.title, {
                      title: f(messages.webhook.list.header.title),
                    })}
                  </h1>
                  <div className="flex justify-center mt-4">
                    <ul className="list-disc">
                      <li>
                        {f(messages.noResult.searchByKeyword, {
                          keyword: f(
                            messages.webhook.list.noResult.searchKeyword
                          ),
                        })}
                      </li>
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
                      title: f(messages.webhook.list.header.title),
                    })}
                  </h1>
                  <p className="mt-5">
                    {f(messages.webhook.list.noData.description)}
                  </p>
                </div>
              </div>
            )
          ) : (
            <div>
              <table className="table-auto leading-normal">
                <tbody className="text-sm">
                  {webhookList.map((webhook) => (
                    <tr key={webhook.id} className={classNames('p-2')}>
                      <td className="px-5 py-2 border-b">
                        <div className="flex pb-1 text-primary">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(webhook)}
                          >
                            {webhook.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(webhook.createdAt * 1000)}
                            />
                          </div>
                        </div>
                      </td>
                      {editable && (
                        <td
                          className={classNames(
                            'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                            'whitespace-nowrap text-gray-500'
                          )}
                        >
                          <button
                            className="text-gray-500"
                            onClick={() => onDelete(webhook)}
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
                maxPage={Math.ceil(totalCount / WEBHOOK_LIST_PAGE_SIZE)}
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
