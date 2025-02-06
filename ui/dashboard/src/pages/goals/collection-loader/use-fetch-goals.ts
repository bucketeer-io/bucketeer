import { useQueryGoals } from '@queries/goals';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchGoals = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  disabled,
  environmentId,
  isInUseStatus,
  archived
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  environmentId?: string;
  isInUseStatus?: boolean;
  archived?: boolean;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryGoals({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      disabled,
      environmentId,
      isInUseStatus,
      archived
    }
  });
};
