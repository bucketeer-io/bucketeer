import { attributeKeysFetcher, AttributeKeysParams } from '@api/features';
import { QueryClient, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AttributeKeysResponse, QueryOptionsRespond } from '@types';

type QueryOptions = QueryOptionsRespond<AttributeKeysResponse> & {
  params?: AttributeKeysParams;
};

export const USER_ATTRIBUTE_KEYS = 'user-attribute-keys';

export const useQueryAttributeKeys = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const query = useQuery({
    queryKey: [USER_ATTRIBUTE_KEYS, params],
    queryFn: async () => {
      return attributeKeysFetcher(params);
    },
    ...queryOptions
  });
  return query;
};

export const usePrefetchAttributeKeys = (options?: QueryOptions) => {
  const { params, ...queryOptions } = options || {};
  const queryClient = useQueryClient();
  queryClient.prefetchQuery({
    queryKey: [USER_ATTRIBUTE_KEYS, params],
    queryFn: async () => {
      return attributeKeysFetcher(params);
    },
    ...queryOptions
  });
};

export const prefetchAttributeKeys = (
  queryClient: QueryClient,
  options?: QueryOptions
) => {
  const { params, ...queryOptions } = options || {};
  queryClient.prefetchQuery({
    queryKey: [USER_ATTRIBUTE_KEYS, params],
    queryFn: async () => {
      return attributeKeysFetcher(params);
    },
    ...queryOptions
  });
};

export const invalidateAttributeKeys = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({
    queryKey: [USER_ATTRIBUTE_KEYS]
  });
};
