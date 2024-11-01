import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export const useFetchOrganizations = ({
  page = 1,
  pageSize,
  orderBy,
  status,
  disabled,
  searchQuery,
  orderDirection
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  disabled?: boolean;
  status?: CollectionStatusType;
  orderDirection?: OrderDirection;
} = {}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryOrganizations({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      disabled,
      searchKeyword: searchQuery,
      archived: status ? status === 'ARCHIVED' : undefined
    }
  });
};
