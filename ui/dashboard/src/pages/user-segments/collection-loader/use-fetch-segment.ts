import { useQueryUserSegments } from '@queries/user-segments';
import { LIST_PAGE_SIZE } from 'constants/app';
import { FeatureSegmentStatus, OrderBy, OrderDirection } from '@types';

export const useFetchSegments = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  status,
  isInUseStatus,
  environmentId
}: {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  status?: FeatureSegmentStatus;
  isInUseStatus?: boolean;
  environmentId: string;
}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryUserSegments({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      status,
      isInUseStatus,
      environmentId
    }
  });
};
