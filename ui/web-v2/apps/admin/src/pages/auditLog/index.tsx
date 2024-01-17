import { Option } from '@/components/Select';
import { intl } from '@/lang';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory, useRouteMatch } from 'react-router-dom';

import { AuditLogList } from '../../components/AuditLogList';
import { Header } from '../../components/Header';
import { AUDITLOG_LIST_PAGE_SIZE } from '../../constants/auditLog';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  listAuditLogs,
  OrderBy,
  OrderDirection,
} from '../../modules/auditLogs';
import { useCurrentEnvironment } from '../../modules/me';
import { ListAuditLogsRequest } from '../../proto/auditlog/service_pb';
import { Event } from '../../proto/event/domain/event_pb';
import { AppDispatch } from '../../store';
import { AuditLogSortOption, isAuditLogSortOption } from '../../types/auditLog';
import { SORT_OPTIONS_CREATED_AT_ASC } from '../../types/list';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../utils/search-params';

const entityTypeOptions: Option[] = [
  {
    value: Event.EntityType.FEATURE.toString(),
    label: intl.formatMessage(messages.sourceType.featureFlag),
  },
  {
    value: Event.EntityType.GOAL.toString(),
    label: intl.formatMessage(messages.sourceType.goal),
  },
  {
    value: Event.EntityType.EXPERIMENT.toString(),
    label: intl.formatMessage(messages.sourceType.experiment),
  },
  {
    value: Event.EntityType.SEGMENT.toString(),
    label: intl.formatMessage(messages.sourceType.segment),
  },
  {
    value: Event.EntityType.ACCOUNT.toString(),
    label: intl.formatMessage(messages.sourceType.account),
  },
  {
    value: Event.EntityType.APIKEY.toString(),
    label: intl.formatMessage(messages.sourceType.apiKey),
  },
  {
    value: Event.EntityType.AUTOOPS_RULE.toString(),
    label: intl.formatMessage(messages.sourceType.autoOperation),
  },
  {
    value: Event.EntityType.PROGRESSIVE_ROLLOUT.toString(),
    label: intl.formatMessage(messages.sourceType.progressiveRollout),
  },
  {
    value: Event.EntityType.PUSH.toString(),
    label: intl.formatMessage(messages.sourceType.push),
  },
  {
    value: Event.EntityType.SUBSCRIPTION.toString(),
    label: intl.formatMessage(messages.sourceType.subscription),
  },
];

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

export const AuditLogIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const currentEnvironment = useCurrentEnvironment();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const isLoading = useSelector<AppState, boolean>(
    (state) => state.auditLog.loading,
    shallowEqual
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
    },
    [updateURL]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
    },
    [updateURL, searchOptions]
  );

  useEffect(() => {
    const sort = createSort(
      isAuditLogSortOption(searchOptions.sort)
        ? searchOptions.sort
        : 'createdAt'
    );
    const page = searchOptions.page ? Number(searchOptions.page) : 1;
    const cursor = (page - 1) * AUDITLOG_LIST_PAGE_SIZE;
    const from = searchOptions.from ? Number(searchOptions.from) : null;
    const to = searchOptions.to ? Number(searchOptions.to) : null;
    const entityType = searchOptions.entityType
      ? Number(searchOptions.entityType)
      : null;

    dispatch(
      listAuditLogs({
        environmentNamespace: currentEnvironment.id,
        pageSize: AUDITLOG_LIST_PAGE_SIZE,
        cursor: String(cursor),
        searchKeyword: searchOptions.q as string,
        orderBy: sort.orderBy,
        orderDirection: sort.orderDirection,
        from: from,
        to: to,
        entityType,
      })
    );
  }, [dispatch, searchOptions, currentEnvironment]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.auditLog.list.header.title)}
          description={f(messages.auditLog.list.header.description)}
        />
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
