import { useQueryFeatures } from '@queries/features';
import { LIST_PAGE_SIZE } from 'constants/app';
import { OrderBy, OrderDirection } from '@types';

export const useFetchFlags = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  environmentId,
  maintainer,
  archived,
  hasExperiment,
  enabled,
  hasPrerequisites,
  tags
}: {
  environmentId: string;
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  maintainer?: string;
  archived?: boolean;
  hasExperiment?: boolean;
  enabled?: boolean;
  hasPrerequisites?: boolean;
  tags?: string[];
}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;

  return useQueryFeatures({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      environmentId,
      maintainer,
      enabled,
      archived,
      hasExperiment,
      hasPrerequisites,
      tags
    }
  });
};
