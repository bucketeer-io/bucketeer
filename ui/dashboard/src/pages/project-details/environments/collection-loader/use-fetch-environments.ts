import { useQueryEnvironments } from '@queries/environments';
import { LIST_PAGE_SIZE } from 'constants/app';
import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export const useFetchEnvironments = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  projectId,
  status,
  organizationId,
  enabledQuery = true
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  projectId?: string;
  status?: CollectionStatusType;
  organizationId?: string;
  enabledQuery?: boolean;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryEnvironments({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      projectId,
      archived: status ? status === 'ARCHIVED' : undefined,
      organizationId
    },
    enabled: enabledQuery
  });
};
