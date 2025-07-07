import { useQueryAuditLogs } from '@queries/audit-logs';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchAuditLogs = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  from,
  to,
  entityType,
  environmentId,
  enabledFetching = true
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  environmentId?: string;
  from?: string;
  to?: string;
  entityType?: number;
  enabledFetching?: boolean;
}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryAuditLogs({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      from,
      to,
      entityType,
      environmentId
    },
    enabled: enabledFetching,
    gcTime: 0
  });
};
