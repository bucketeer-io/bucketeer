import React, { FC, memo, useEffect, useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { AuditLogList } from '../../components/AuditLogList';
import { AUDITLOG_LIST_PAGE_SIZE } from '../../constants/auditLog';
import {
  listFeatureHistory,
  OrderBy,
  OrderDirection,
} from '../../modules/auditLogs';
import { useCurrentEnvironment } from '../../modules/me';
import { ListAuditLogsRequest } from '../../proto/auditlog/service_pb';
import { AppDispatch } from '../../store';
import { AuditLogSortOption, isAuditLogSortOption } from '../../types/auditLog';
import { SORT_OPTIONS_CREATED_AT_ASC } from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: AuditLogSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAuditLogsRequest.OrderBy.TIMESTAMP,
        orderDirection: ListAuditLogsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListAuditLogsRequest.OrderBy.TIMESTAMP,
        orderDirection: ListAuditLogsRequest.OrderDirection.DESC,
      };
  }
};

interface FeatureHistoryPageProps {
  featureId: string;
}

export const FeatureHistoryPage: FC<FeatureHistoryPageProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const history = useHistory();
    const { url } = useRouteMatch();
    const searchOptions = useSearchParams();
    searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';

    const updateHistoryList = useCallback(
      (options, page: number) => {
        const sort = createSort(
          isAuditLogSortOption(options && options.sort)
            ? options.sort
            : '-createdAt'
        );
        const cursor = (page - 1) * AUDITLOG_LIST_PAGE_SIZE;
        const from = options.from ? Number(options.from) : null;
        const to = options.to ? Number(options.to) : null;
        dispatch(
          listFeatureHistory({
            featureId: featureId,
            environmentNamespace: currentEnvironment.id,
            pageSize: AUDITLOG_LIST_PAGE_SIZE,
            cursor: String(cursor),
            searchKeyword: options.q as string,
            orderBy: sort.orderBy,
            orderDirection: sort.orderDirection,
            from: from,
            to: to,
          })
        );
      },
      [dispatch]
    );

    const updateURL = useCallback(
      (options: Record<string, string | number | boolean | undefined>) => {
        history.replace(
          `${url}?${stringifySearchParams({
            ...options,
          })}`
        );
      },
      [history]
    );

    const handleSearchOptionsChange = useCallback(
      (options) => {
        updateURL({ ...options, page: 1 });
        updateHistoryList(options, 1);
      },
      [updateURL]
    );

    const handlePageChange = useCallback(
      (page: number) => {
        updateURL({ ...searchOptions, page });
        updateHistoryList(searchOptions, page);
      },
      [updateURL, updateHistoryList, searchOptions]
    );

    useEffect(() => {
      updateHistoryList(
        searchOptions,
        searchOptions.page ? Number(searchOptions.page) : 1
      );
    }, [updateHistoryList]);

    return (
      <div className="p-10 bg-gray-100">
        <AuditLogList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onChangeSearchOptions={handleSearchOptionsChange}
        />
      </div>
    );
  }
);
