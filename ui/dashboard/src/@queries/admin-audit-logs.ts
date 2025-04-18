import {
  adminAuditLogsFetcher,
  AdminAuditLogsFetcherParams
} from '@api/audit-logs/admin-audit-logs-fetcher';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AuditLogCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AuditLogCollection> & {
  params?: AdminAuditLogsFetcherParams;
};

export const ADMIN_AUDIT_LOGS_QUERY_KEY = 'admin-audit-logs';

export const useQueryAdminAuditLogs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ADMIN_AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return adminAuditLogsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAdminAuditLogs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [ADMIN_AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return adminAuditLogsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAdminAuditLogs = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [ADMIN_AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return adminAuditLogsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAdminAuditLogs = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ADMIN_AUDIT_LOGS_QUERY_KEY]
  });
};
