import MUDeleteIcon from '@material-ui/icons/Delete';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { TagChips } from '../../components/TagsChips';
import { PUSH_LIST_PAGE_SIZE } from '../../constants/push';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable } from '../../modules/me';
import { selectAll } from '../../modules/pushes';
import { Push } from '../../proto/push/push_pb';
import { PushSearchOptions } from '../../types/push';
import { classNames } from '../../utils/css';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { PushSearch } from '../PushSearch';
import { RelativeDateText } from '../RelativeDateText';

export interface PushListProps {
  searchOptions: PushSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: PushSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (push: Push.AsObject) => void;
  onDelete: (push: Push.AsObject) => void;
}

export const PushList: FC<PushListProps> = memo(
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
    const pushList = useSelector<AppState, Push.AsObject[]>(
      (state) => selectAll(state.push),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.push.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.push.totalCount,
      shallowEqual
    );

    return (
      <div className="w-full">
        <div className="flex items-stretch mb-8 text-sm">
          <p className="text-gray-700">
            {f(messages.push.list.header.description)}
          </p>
          <a
            className="link"
            target="_blank"
            href="https://bucketeer.io/docs#/push-propagation"
            rel="noreferrer"
          >
            {f(messages.readMore)}
          </a>
        </div>
        <div className="min-w-max bg-white border border-gray-300 rounded-md">
          <div>
            <PushSearch
              options={searchOptions}
              onChange={onChangeSearchOptions}
              onAdd={onAdd}
            />
          </div>
          {isLoading ? (
            <ListSkeleton />
          ) : pushList.length == 0 ? (
            searchOptions.q ? (
              <div className="my-10 flex justify-center">
                <div className="text-gray-700">
                  <h1 className="text-lg">
                    {f(messages.noResult.title, {
                      title: f(messages.push.list.header.title),
                    })}
                  </h1>
                  <div className="flex justify-center mt-4">
                    <ul className="list-disc">
                      <li>
                        {f(messages.noResult.searchByKeyword, {
                          keyword: f(messages.push.list.noResult.searchKeyword),
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
                      title: f(messages.push.list.header.title),
                    })}
                  </h1>
                  <p className="mt-5">
                    {f(messages.push.list.noData.description)}
                  </p>
                </div>
              </div>
            )
          ) : (
            <div>
              <table className="table-auto leading-normal">
                <tbody className="text-sm">
                  {pushList.map((push) => (
                    <tr key={push.id} className={classNames('p-2')}>
                      <td className="px-5 py-2 border-b">
                        <div className="flex pb-1 text-primary">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(push)}
                          >
                            {push.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(push.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        <div>
                          <TagChips tags={push.tagsList} />
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
                            onClick={() => onDelete(push)}
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
                maxPage={Math.ceil(totalCount / PUSH_LIST_PAGE_SIZE)}
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
