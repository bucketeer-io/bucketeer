import MUCloudDownloadIcon from '@material-ui/icons/CloudDownload';
import MUDeleteIcon from '@material-ui/icons/Delete';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { SEGMENT_LIST_PAGE_SIZE } from '../../constants/segment';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable } from '../../modules/me';
import { selectAll } from '../../modules/segments';
import { Segment } from '../../proto/feature/segment_pb';
import { SegmentSearchOptions } from '../../types/segment';
import { classNames } from '../../utils/css';
import { ActionMenu, MenuActions, MenuItem } from '../ActionMenu';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { SegmentSearch } from '../SegmentSearch';

export interface SegmentListProps {
  searchOptions: SegmentSearchOptions;
  onChangePage: (page: number) => void;
  onDelete: (s: Segment.AsObject) => void;
  onDownload: (s: Segment.AsObject) => void;
  onChangeSearchOptions: (options: SegmentSearchOptions) => void;
  onAdd: () => void;
  onUpdate: (s: Segment.AsObject) => void;
}

export const SegmentList: FC<SegmentListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onDelete,
    onDownload,
    onChangeSearchOptions,
    onAdd,
    onUpdate
  }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const segments = useSelector<AppState, Segment.AsObject[]>(
      (state) => selectAll(state.segments),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.segments.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.segments.totalCount,
      shallowEqual
    );
    const userCount = (
      count: number,
      status: Segment.StatusMap[keyof Segment.StatusMap]
    ): React.ReactNode => {
      switch (status) {
        case Segment.Status.UPLOADING:
          return (
            <div className="flex items-center mt-2">
              <div className="spinner"></div>
              <span className={classNames('text-xs text-gray-700 ml-2')}>
                {f(messages.segment.status.uploading)}
              </span>
            </div>
          );
        case Segment.Status.FAILED:
          return (
            <div className="flex items-center mt-1 text-xs text-red-600">
              {f(messages.segment.status.uploadFailed)}
            </div>
          );
        default:
          return (
            <div className="flex items-center mt-1 text-xs text-gray-700">
              {f(messages.segment.userCount)}: {count}
            </div>
          );
      }
    };
    const createMenuItems = (includedUserCount: number): Array<MenuItem> => {
      const items: Array<MenuItem> = [];
      items.push({
        action: MenuActions.DOWNLOAD,
        name: intl.formatMessage(messages.segment.action.download),
        iconElement: <MUCloudDownloadIcon />,
        disabled: includedUserCount === 0
      });
      items.push({
        action: MenuActions.DELETE,
        name: intl.formatMessage(messages.segment.action.delete),
        iconElement: <MUDeleteIcon />
      });
      return items;
    };

    return (
      <div className="w-full bg-white border border-gray-300 rounded-md">
        <div>
          <SegmentSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            onAdd={onAdd}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : segments.length == 0 ? (
          searchOptions.q || searchOptions.inUse ? (
            <div className="my-10 flex justify-center">
              <div className="text-gray-700">
                <h1 className="text-lg">
                  {f(messages.noResult.title, {
                    title: f(messages.segment.list.header.title)
                  })}
                </h1>
                <div className="flex justify-center mt-4">
                  <ul className="list-disc">
                    <li>
                      {f(messages.noResult.searchByKeyword, {
                        keyword: f(messages.segment.list.noResult.searchKeyword)
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
                    title: f(messages.segment.list.header.title)
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.segment.list.noData.description)}
                </p>
                <a
                  href="https://bucketeer.io/docs#/creating-a-feature-flag?id=create-a-feature-flag"
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
                {segments.map((segment) => {
                  return (
                    <tr key={segment.id} className={classNames('p-2')}>
                      <td className="pl-5 pr-2 py-3 border-b">
                        <div className="flex">
                          <button
                            className="link text-left"
                            onClick={() => onUpdate(segment)}
                          >
                            {segment.name}
                          </button>
                          <div className="flex items-center ml-2 text-xs text-gray-700 whitespace-nowrap">
                            <span className="mr-1">{f(messages.created)}</span>
                            <RelativeDateText
                              date={new Date(segment.createdAt * 1000)}
                            />
                          </div>
                        </div>
                        {userCount(segment.includedUserCount, segment.status)}
                      </td>
                      <td className="w-[10%] pl-5 px-5 border-b">
                        <div className="flex justify-end">
                          <div className="flex justify-center items-center rounded-md w-16 h-6 bg-gray-200">
                            <span className="text-xs text-gray-700">
                              {segment.isInUseStatus
                                ? f(messages.segment.filterOptions.inUse)
                                : f(messages.segment.filterOptions.notInUse)}
                            </span>
                          </div>
                        </div>
                      </td>
                      {editable && (
                        <td
                          className={classNames(
                            'w-[1%] pl-2 pr-5 py-3 border-b border-gray-300',
                            'whitespace-nowrap'
                          )}
                        >
                          <ActionMenu
                            onClickAction={(action) => {
                              switch (action) {
                                case MenuActions.DOWNLOAD:
                                  onDownload(segment);
                                  return;
                                case MenuActions.DELETE:
                                  onDelete(segment);
                                  return;
                              }
                            }}
                            menuItems={createMenuItems(
                              segment.includedUserCount
                            )}
                          />
                        </td>
                      )}
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / SEGMENT_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);
