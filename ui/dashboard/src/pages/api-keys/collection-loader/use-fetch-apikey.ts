import { useQueryAPIKeys } from '@queries/api-keys';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchAPIKeys = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  disabled,
  environmentNamespace
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  environmentNamespace?: string;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryAPIKeys({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      disabled,
      environmentNamespace
    }
  });
};
