import { useQueryPushes } from '@queries/pushes';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchPushes = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  disabled,
  organizationId,
  environmentIds
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  organizationId?: string;
  environmentIds?: string[];
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryPushes({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      disabled,
      organizationId,
      environmentIds
    }
  });
};
