import { useQueryAdminAuditLogs } from '@queries/admin-audit-logs';
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
  isSystemAdmin
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
  isSystemAdmin?: boolean;
}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return (isSystemAdmin ? useQueryAdminAuditLogs : useQueryAuditLogs)({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      from,
      to,
      entityType,
      ...(isSystemAdmin ? {} : { environmentId })
    }
  });
};
