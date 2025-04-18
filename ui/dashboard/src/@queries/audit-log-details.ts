import {
  auditLogDetailsFetcher,
  AuditLogDetailsFetcherParams,
  AuditLogDetailsResponse
} from '@api/audit-logs';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AuditLogDetailsResponse> & {
  params?: AuditLogDetailsFetcherParams;
};

export const AUDIT_LOG_DETAILS_QUERY_KEY = 'audit-log-details';

export const useQueryAuditLogDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [AUDIT_LOG_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAuditLogDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [AUDIT_LOG_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAuditLogDetails = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [AUDIT_LOG_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return auditLogDetailsFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAuditLogDetails = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [AUDIT_LOG_DETAILS_QUERY_KEY]
  });
};
