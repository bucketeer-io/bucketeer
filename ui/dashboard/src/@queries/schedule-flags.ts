import {
  featuresScheduleFetcher,
  FeaturesScheduleFetcherParams
} from '@api/features';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { QueryOptionsRespond, ScheduleFlagCollection } from '@types';

type QueryOptions = QueryOptionsRespond<ScheduleFlagCollection> & {
  params?: FeaturesScheduleFetcherParams;
};

export const SCHEDULE_FLAGS_QUERY_KEY = 'schedule-flags-query-key';

export const useQueryScheduleFlags = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [SCHEDULE_FLAGS_QUERY_KEY, params],
    queryFn: async () => {
      return featuresScheduleFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchScheduleFlags = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [SCHEDULE_FLAGS_QUERY_KEY, params],
    queryFn: async () => {
      return featuresScheduleFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchScheduleFlags = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [SCHEDULE_FLAGS_QUERY_KEY, params],
    queryFn: async () => {
      return featuresScheduleFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateScheduleFlags = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [SCHEDULE_FLAGS_QUERY_KEY]
  });
};
