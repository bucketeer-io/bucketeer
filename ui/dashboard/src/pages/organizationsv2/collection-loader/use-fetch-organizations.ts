import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export const useFetchOrganizations = ({
  orderBy,
  status,
  searchQuery,
  orderDirection
}: {
  searchQuery?: string;
  orderBy?: OrderBy;
  status?: CollectionStatusType;
  orderDirection?: OrderDirection;
} = {}) => {
  return useQueryOrganizations({
    params: {
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: orderBy || 'DEFAULT',
      orderDirection: orderDirection || 'ASC',
      searchKeyword: searchQuery,
      archived: status === 'ARCHIVED'
    }
  });
};
