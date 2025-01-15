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
  | 'ORGANIZATION_ROLE'
  | 'EMAIL'
  | 'STATE'
  | 'LAST_SEEN'
  | 'CREATOR_EMAIL'
  | 'ENVIRONMENT'
  | 'TAGS'
  | 'CONNECTIONS'
  | 'USERS';

export type OrderDirection = 'ASC' | 'DESC';

export type CollectionStatusType = 'ACTIVE' | 'ARCHIVED';

export interface Collection<T> {
  data: T[];
  cursor: string;
  totalCount: string;
}

export interface CollectionParams {
  pageSize: number;
  cursor: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  searchKeyword?: string;
  disabled?: boolean;
}
