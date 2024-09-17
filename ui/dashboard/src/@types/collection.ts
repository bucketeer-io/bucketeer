import { UseQueryOptions } from '@tanstack/react-query';

export type QueryOptionsRespond<T> = Omit<UseQueryOptions<T>, 'queryKey'>;

export interface Collection<T> {
  data: T[];
  cursor: string;
  totalCount: string;
}
