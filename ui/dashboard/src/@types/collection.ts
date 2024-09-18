import { UseQueryOptions } from '@tanstack/react-query';

export type QueryOptionsRespond<T> = Omit<UseQueryOptions<T>, 'queryKey'>;

export type OrderBy =
  | 'DEFAULT'
  | 'ID'
  | 'CREATED_AT'
  | 'UPDATED_AT'
  | 'NAME'
  | 'URL_CODE';

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
  searchKeyword: string;
  disabled: boolean;
}
