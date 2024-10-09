import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchOrganizations = ({
  orderBy,
  searchQuery,
  orderDirection
}: {
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
} = {}) => {
  return useQueryOrganizations({
    params: {
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      disabled: false,
      orderBy: orderBy || 'DEFAULT',
      orderDirection: orderDirection || 'ASC',
      searchKeyword: searchQuery,
      archived: false
    }
  });
};
