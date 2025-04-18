import { auditLogsFetcher, AuditLogsFetcherParams } from '@api/audit-logs';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AuditLogCollection, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AuditLogCollection> & {
  params?: AuditLogsFetcherParams;
};

export const AUDIT_LOGS_QUERY_KEY = 'audit-logs';

export const useQueryAuditLogs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAuditLogs = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAuditLogs = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUDIT_LOGS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAuditLogs = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUDIT_LOGS_QUERY_KEY]
  });
};
