import { useQueryProjects } from '@queries/projects';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchProjects = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  organizationId,
  disabled
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  organizationId?: string;
  disabled?: boolean;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryProjects({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      disabled,
      orderDirection,
      searchKeyword: searchQuery,
      organizationId
    }
  });
};
