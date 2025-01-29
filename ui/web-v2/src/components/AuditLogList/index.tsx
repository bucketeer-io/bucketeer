import { FC, memo, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import * as Diff from 'diff';

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
    entityTypeOptions
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

    console.log({ auditLogs });

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
                      {auditLog.entityData && (
                        <AuditLogDetail auditLog={auditLog} />
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

interface AuditLogDetailProps {
  auditLog: AuditLog.AsObject;
}

const AuditLogDetail: FC<AuditLogDetailProps> = ({ auditLog }) => {
  const { formatMessage: f } = useIntl();
  const [showChanges, setShowChanges] = useState(false);
  const [showSnapshot, setShowSnapshot] = useState(false);
  return (
    <div>
      <div className={classNames('text-primary text-xs')}>
        <button onClick={() => setShowChanges(!showChanges)}>
          {showChanges
            ? f(messages.auditLog.detail.hideChanges)
            : `> ${f(messages.auditLog.detail.showChanges)}`}
        </button>
        {showChanges && (
          <div className={classNames('p-3')}>
            <DiffView
              oldStr={auditLog.previousEntityData}
              newStr={auditLog.entityData}
            />
          </div>
        )}
      </div>
      <div>
        <button
          onClick={() => setShowSnapshot(!showSnapshot)}
          className={classNames('text-primary text-xs')}
        >
          {showSnapshot
            ? f(messages.auditLog.detail.hideSnapshot)
            : `> ${f(messages.auditLog.detail.showSnapshot)}`}
        </button>
        {showSnapshot && (
          <div className={classNames('p-3')}>
            <p
              className={classNames('bg-gray-100 whitespace-pre-wrap text-xs')}
            >{`${auditLog.entityData}`}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export interface DiffViewProps {
  oldStr: string;
  newStr: string;
}

export const DiffView: FC<DiffViewProps> = memo(function DiffView({
  oldStr,
  newStr
}) {
  return (
    <div>
      {Diff.createTwoFilesPatch('old version', 'new version', oldStr, newStr)
        .split('\n')
        .filter((line) => {
          if (line.startsWith('\\') || line.startsWith('=')) {
            return false;
          }
          return true;
        })
        .map((line) => {
          if (line.startsWith('@@')) {
            return '...';
          }
          return line;
        })
        .map((line, i) => {
          switch (line[0]) {
            case '+':
              return (
                <div
                  key={i}
                  className={classNames('bg-green-50 text-green-800')}
                >
                  <span key={i}>{line}</span>
                </div>
              );
            case '-':
              return (
                <div
                  key={i}
                  className={classNames('bg-red-50 text-red-800')}
                  data-testid="deleted-line"
                >
                  <span>{line}</span>
                </div>
              );
            default:
              return <div key={i}>{line}</div>;
          }
        })}
    </div>
  );
});

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
              title: f(messages.auditLog.list.header.title)
            })}
          </h1>
          <div className="flex justify-center mt-4">
            <ul className="list-disc">
              <li>
                {f(messages.noResult.searchByKeyword, {
                  keyword: f(messages.auditLog.list.noResult.searchKeyword)
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
              title: f(messages.auditLog.list.header.title)
            })}
          </h1>
          <p className="mt-5">{f(messages.auditLog.list.noData.description)}</p>
          <a
            href="https://docs.bucketeer.io/feature-flags/audit-logs"
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
