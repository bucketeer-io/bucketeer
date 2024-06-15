import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { AUDITLOG_LIST_PAGE_SIZE } from '../../constants/auditLog';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectAll } from '../../modules/auditLogs';
import { AuditLog } from '../../proto/auditlog/auditlog_pb';
import { AuditLogSearchOptions } from '../../types/auditLog';
import { classNames } from '../../utils/css';
import { AuditLogSearch } from '../AuditLogSearch';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';
import { RelativeDateText } from '../RelativeDateText';
import { Option } from '../Select';

export interface AuditLogListProps {
  searchOptions: AuditLogSearchOptions;
  onChangePage: (page: number) => void;
  onChangeSearchOptions: (options: AuditLogSearchOptions) => void;
  showEntityTypeFilter?: boolean;
  entityTypeOptions?: Option[];
}

export const AuditLogList: FC<AuditLogListProps> = memo(
  ({
    searchOptions,
    onChangePage,
    onChangeSearchOptions,
    showEntityTypeFilter,
    entityTypeOptions,
  }) => {
    const auditLogs = useSelector<AppState, AuditLog.AsObject[]>(
      (state) => selectAll(state.auditLog),
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.auditLog.loading,
      shallowEqual
    );
    const totalCount = useSelector<AppState, number>(
      (state) => state.auditLog.totalCount,
      shallowEqual
    );

    return (
      <div className="min-w-max bg-white border border-gray-300 rounded-md">
        <div>
          <AuditLogSearch
            options={searchOptions}
            onChange={onChangeSearchOptions}
            showEntityTypeFilter={showEntityTypeFilter}
            entityTypeOptions={entityTypeOptions}
          />
        </div>
        {isLoading ? (
          <ListSkeleton />
        ) : auditLogs.length === 0 ? (
          <NoData searchOptions={searchOptions} />
        ) : (
          <div>
            <table className="min-w-full table-auto leading-normal">
              <tbody className="text-sm">
                {auditLogs.map((auditLog) => (
                  <tr key={auditLog.id} className={classNames('p-2')}>
                    <td className="px-5 py-2 border-b">
                      <span className={classNames('text-primary mr-2')}>
                        {auditLog.editor.email}
                      </span>
                      <span className={classNames('text-gray-700')}>
                        {auditLog.localizedMessage.message}
                      </span>
                      <div className="flex my-1 items-center text-xs text-gray-700">
                        <RelativeDateText
                          date={new Date(auditLog.timestamp * 1000)}
                        />
                      </div>
                      {auditLog.options && (
                        <div className="border-l-4 my-1 border-gray-300 px-3 text-gray-700 text-xs">
                          {nl2br(auditLog.options.comment)}
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <Pagination
              maxPage={Math.ceil(totalCount / AUDITLOG_LIST_PAGE_SIZE)}
              currentPage={searchOptions.page ? Number(searchOptions.page) : 1}
              onChange={onChangePage}
            />
          </div>
        )}
      </div>
    );
  }
);

interface NoDataProps {
  searchOptions: AuditLogSearchOptions;
}

const NoData: FC<NoDataProps> = ({ searchOptions }) => {
  const { formatMessage: f } = useIntl();

  if (searchOptions.from) {
    return (
      <div className="my-10 text-center">
        <h1 className="text-lg">{f(messages.noResult.dateRange.title)}</h1>
        <p className="mt-2">{f(messages.noResult.dateRange.description)}</p>
      </div>
    );
  } else if (searchOptions.q || searchOptions.entityType) {
    return (
      <div className="my-10 flex justify-center">
        <div className="text-gray-700">
          <h1 className="text-lg">
            {f(messages.noResult.title, {
              title: f(messages.auditLog.list.header.title),
            })}
          </h1>
          <div className="flex justify-center mt-4">
            <ul className="list-disc">
              <li>
                {f(messages.noResult.searchByKeyword, {
                  keyword: f(messages.auditLog.list.noResult.searchKeyword),
                })}
              </li>
              <li>{f(messages.noResult.checkTypos)}</li>
            </ul>
          </div>
        </div>
      </div>
    );
  } else {
    return (
      <div className="my-10 flex justify-center">
        <div className="w-[600px] text-gray-700 text-center">
          <h1 className="text-lg">
            {f(messages.noData.title, {
              title: f(messages.auditLog.list.header.title),
            })}
          </h1>
          <p className="mt-5">{f(messages.auditLog.list.noData.description)}</p>
          <a
            href="https://bucketeer.io/docs#/audit-logs?id=environment-audit-logs"
            target="_blank"
            rel="noreferrer"
            className="link"
          >
            {f(messages.readMore)}
          </a>
        </div>
      </div>
    );
  }
};

const nl2br = (text: string): Array<React.ReactNode> => {
  const regex = /(\n)/g;
  return text.split(regex).map((line, i) => {
    if (line.match(regex)) {
      return <br key={i} />;
    } else {
      return line;
    }
  });
};
