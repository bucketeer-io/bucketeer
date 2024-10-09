import { UseQueryOptions } from '@tanstack/react-query';

export type QueryOptionsRespond<T> = Omit<UseQueryOptions<T>, 'queryKey'>;

export type OrderBy =
  | 'DEFAULT'
  | 'ID'
  | 'CREATED_AT'
  | 'UPDATED_AT'
  | 'NAME'
  | 'URL_CODE'
  | 'FEATURE_COUNT'
  | 'ENVIRONMENT_COUNT'
  | 'PROJECT_COUNT'
  | 'USER_COUNT'
  | 'ROLE'
  | 'EMAIL';

export type OrderDirection = 'ASC' | 'DESC';

export interface Collection<T> {
  data: T[];
  cursor: string;
  totalCount: string;
}

export interface CollectionParams {
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  disabled: boolean;
  searchKeyword?: string;
}
