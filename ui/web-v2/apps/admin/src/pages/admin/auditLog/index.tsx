import { Option } from '@/components/Select';
import { intl } from '@/lang';
import { Event } from '@/proto/event/domain/event_pb';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory, useRouteMatch } from 'react-router-dom';

import { AuditLogList } from '../../../components/AuditLogList';
import { Header } from '../../../components/Header';
import { AUDITLOG_LIST_PAGE_SIZE } from '../../../constants/auditLog';
import { messages } from '../../../lang/messages';
import { AppState } from '../../../modules';
import {
  listAdminAuditLogs,
  AdminOrderBy,
  AdminOrderDirection,
} from '../../../modules/auditLogs';
import { ListAuditLogsRequest } from '../../../proto/auditlog/service_pb';
import { AppDispatch } from '../../../store';
import {
  AuditLogSortOption,
  isAuditLogSortOption,
} from '../../../types/auditLog';
import { SORT_OPTIONS_CREATED_AT_ASC } from '../../../types/list';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../../utils/search-params';

interface Sort {
  orderBy: AdminOrderBy;
  orderDirection: AdminOrderDirection;
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

const entityTypeOptions: Option[] = [
  {
    value: Event.EntityType.ENVIRONMENT.toString(),
    label: intl.formatMessage(messages.sourceType.environment),
  },
  {
    value: Event.EntityType.ADMIN_ACCOUNT.toString(),
    label: intl.formatMessage(messages.sourceType.adminAccount),
  },
  {
    value: Event.EntityType.ADMIN_SUBSCRIPTION.toString(),
    label: intl.formatMessage(messages.sourceType.adminSubscription),
  },
  {
    value: Event.EntityType.PROJECT.toString(),
    label: intl.formatMessage(messages.sourceType.project),
  },
];

export const AdminAuditLogIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const isLoading = useSelector<AppState, boolean>(
    (state) => state.auditLog.loading,
    shallowEqual
  );

  const updateAuditLogList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isAuditLogSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * AUDITLOG_LIST_PAGE_SIZE;
      const from = options && options.from ? Number(options.from) : null;
      const to = options && options.to ? Number(options.to) : null;
      const resource =
        options && options.resource ? Number(options.resource) : null;
      dispatch(
        listAdminAuditLogs({
          pageSize: AUDITLOG_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          from: from,
          to: to,
          resource: resource,
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
      updateAuditLogList(options, 1);
    },
    [updateURL, updateAuditLogList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateAuditLogList(searchOptions, page);
    },
    [updateURL, searchOptions, updateAuditLogList]
  );

  useEffect(() => {
    updateAuditLogList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateAuditLogList]);

  return (
    <>
      <div className="flex items-stretch m-10 text-sm">
        <p className="text-gray-700">
          {f(messages.adminAuditLog.list.header.description)}
        </p>
        <a
          className="link"
          target="_blank"
          href="https://bucketeer.io/docs/#/audit-logs"
          rel="noreferrer"
        >
          {f(messages.readMore)}
        </a>
      </div>
      <div className="m-10">
        <AuditLogList
          showEntityTypeFilter
          entityTypeOptions={entityTypeOptions}
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onChangeSearchOptions={handleSearchOptionsChange}
        />
      </div>
    </>
  );
});
