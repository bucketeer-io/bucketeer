import {
  scheduledFlagChangeCreator,
  ScheduledFlagChangeCreatorParams
} from '@api/features/scheduled-flag-change-creator';
import {
  scheduledFlagChangeDelete,
  ScheduledFlagChangeDeleteParams
} from '@api/features/scheduled-flag-change-delete';
import {
  scheduledFlagChangeExecutor,
  ScheduledFlagChangeExecuteParams
} from '@api/features/scheduled-flag-change-executor';
import {
  scheduledFlagChangeGet,
  ScheduledFlagChangeGetParams
} from '@api/features/scheduled-flag-change-get';
import {
  scheduledFlagChangeSummaryFetcher,
  ScheduledFlagChangeSummaryParams
} from '@api/features/scheduled-flag-change-summary';
import {
  scheduledFlagChangeUpdater,
  ScheduledFlagChangeUpdaterParams
} from '@api/features/scheduled-flag-change-updater';
import {
  scheduledFlagChangesFetcher,
  ScheduledFlagChangesFetcherParams
} from '@api/features/scheduled-flag-changes-fetch';
import {
  QueryClient,
  useMutation,
  useQuery,
  useQueryClient
} from '@tanstack/react-query';
import type {
  QueryOptionsRespond,
  ScheduledFlagChangeCollection,
  GetScheduledFlagChangeResponse,
  GetScheduledFlagChangeSummaryResponse,
  CreateScheduledFlagChangeResponse,
  UpdateScheduledFlagChangeResponse,
  ExecuteScheduledFlagChangeResponse
} from '@types';

export const SCHEDULED_FLAG_CHANGES_QUERY_KEY =
  'scheduled-flag-changes-query-key';
export const SCHEDULED_FLAG_CHANGE_QUERY_KEY =
  'scheduled-flag-change-query-key';
export const SCHEDULED_FLAG_CHANGE_SUMMARY_QUERY_KEY =
  'scheduled-flag-change-summary-query-key';

type ListQueryOptions = QueryOptionsRespond<ScheduledFlagChangeCollection> & {
  params?: ScheduledFlagChangesFetcherParams;
};

export const useQueryScheduledFlagChanges = (options?: ListQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [SCHEDULED_FLAG_CHANGES_QUERY_KEY, params],
    queryFn: async () => {
      return scheduledFlagChangesFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchScheduledFlagChanges = (
  queryClient: QueryClient,
  options?: ListQueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [SCHEDULED_FLAG_CHANGES_QUERY_KEY, params],
    queryFn: async () => {
      return scheduledFlagChangesFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateScheduledFlagChanges = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [SCHEDULED_FLAG_CHANGES_QUERY_KEY]
  });
};

type GetQueryOptions = QueryOptionsRespond<GetScheduledFlagChangeResponse> & {
  params?: ScheduledFlagChangeGetParams;
};

export const useGetScheduledFlagChange = (options?: GetQueryOptions) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [SCHEDULED_FLAG_CHANGE_QUERY_KEY, params],
    queryFn: async () => {
      return scheduledFlagChangeGet(params);
    },
    ...queryOptions
  });
};

type SummaryQueryOptions =
  QueryOptionsRespond<GetScheduledFlagChangeSummaryResponse> & {
    params?: ScheduledFlagChangeSummaryParams;
  };

export const useGetScheduledFlagChangeSummary = (
  options?: SummaryQueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [SCHEDULED_FLAG_CHANGE_SUMMARY_QUERY_KEY, params],
    queryFn: async () => {
      return scheduledFlagChangeSummaryFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateScheduledFlagChangeSummary = (
  queryClient: QueryClient
) => {
  queryClient.invalidateQueries({
    queryKey: [SCHEDULED_FLAG_CHANGE_SUMMARY_QUERY_KEY]
  });
};

export const invalidateScheduledFlagChange = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [SCHEDULED_FLAG_CHANGE_QUERY_KEY]
  });
};

const invalidateAllScheduledFlagChangeQueries = (queryClient: QueryClient) => {
  invalidateScheduledFlagChanges(queryClient);
  invalidateScheduledFlagChangeSummary(queryClient);
  invalidateScheduledFlagChange(queryClient);
};

export const useCreateScheduledFlagChange = () => {
  const queryClient = useQueryClient();
  return useMutation<
    CreateScheduledFlagChangeResponse,
    Error,
    ScheduledFlagChangeCreatorParams
  >({
    mutationFn: async (params: ScheduledFlagChangeCreatorParams) => {
      return scheduledFlagChangeCreator(params);
    },
    onSuccess: () => {
      invalidateAllScheduledFlagChangeQueries(queryClient);
    }
  });
};

export const useUpdateScheduledFlagChange = () => {
  const queryClient = useQueryClient();
  return useMutation<
    UpdateScheduledFlagChangeResponse,
    Error,
    ScheduledFlagChangeUpdaterParams
  >({
    mutationFn: async (params: ScheduledFlagChangeUpdaterParams) => {
      return scheduledFlagChangeUpdater(params);
    },
    onSuccess: () => {
      invalidateAllScheduledFlagChangeQueries(queryClient);
    }
  });
};

export const useDeleteScheduledFlagChange = () => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, ScheduledFlagChangeDeleteParams>({
    mutationFn: async (params: ScheduledFlagChangeDeleteParams) => {
      return scheduledFlagChangeDelete(params);
    },
    onSuccess: () => {
      invalidateAllScheduledFlagChangeQueries(queryClient);
    }
  });
};

export const useExecuteScheduledFlagChange = () => {
  const queryClient = useQueryClient();
  return useMutation<
    ExecuteScheduledFlagChangeResponse,
    Error,
    ScheduledFlagChangeExecuteParams
  >({
    mutationFn: async (params: ScheduledFlagChangeExecuteParams) => {
      return scheduledFlagChangeExecutor(params);
    },
    onSuccess: () => {
      invalidateAllScheduledFlagChangeQueries(queryClient);
    }
  });
};
