import {
  environmentResultDetailsFetcher,
  EnvironmentResultDetailsFetcherParams
} from '@api/environment-result/environment-result-details';
import { QueryClient, useQuery } from '@tanstack/react-query';
import type { EnvironmentResponse, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<EnvironmentResponse> & {
  params?: EnvironmentResultDetailsFetcherParams;
};

export const ENVIRONMENTS_DETAILS_QUERY_KEY = 'environment-details';

export const useQueryEnvironmentDetails = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [ENVIRONMENTS_DETAILS_QUERY_KEY, params],
    queryFn: async () => {
      return environmentResultDetailsFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const invalidateEnvironmentDetails = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [ENVIRONMENTS_DETAILS_QUERY_KEY]
  });
};
