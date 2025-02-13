import { useQueryAccounts } from '@queries/accounts';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchMembers = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  organizationId,
  environmentId,
  disabled,
  organizationRole
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  organizationId?: string;
  environmentId?: string;
  disabled?: boolean;
  organizationRole?: number;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryAccounts({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      organizationId,
      environmentId,
      disabled,
      organizationRole
    }
  });
};
