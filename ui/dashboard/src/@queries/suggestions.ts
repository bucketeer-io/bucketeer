import {
  suggestionsFetcher,
  SuggestionsFetcherParams,
  SuggestionsFetcherResponse
} from '@api/ai-chat';
import { useQuery } from '@tanstack/react-query';
import type { QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<SuggestionsFetcherResponse> & {
  params?: SuggestionsFetcherParams;
};

export const SUGGESTIONS_QUERY_KEY = 'ai-chat-suggestions';

export const useQuerySuggestions = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  return useQuery({
    queryKey: [SUGGESTIONS_QUERY_KEY, params],
    queryFn: async () => {
      if (!params) throw new Error('params is required');
      return suggestionsFetcher(params);
    },
    // Rule-based suggestions rarely change; cache for 5 minutes
    staleTime: 5 * 60 * 1000,
    enabled: !!params?.environmentId,
    ...queryOptions
  });
};
