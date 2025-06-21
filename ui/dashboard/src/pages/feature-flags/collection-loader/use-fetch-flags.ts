import { useQueryFeatures } from '@queries/features';
import { LIST_PAGE_SIZE } from 'constants/app';
import { CollectionStatusType, OrderBy, OrderDirection } from '@types';
import { StatusFilterType } from '../types';

export const useFetchFlags = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  environmentId,
  maintainer,
  hasExperiment,
  enabled,
  hasPrerequisites,
  tags,
  tab,
  status,
  hasFeatureFlagAsRule
}: {
  environmentId: string;
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  maintainer?: string;
  hasExperiment?: boolean;
  enabled?: boolean;
  hasPrerequisites?: boolean;
  tags?: string[];
  tab: CollectionStatusType;
  status?: StatusFilterType;
  hasFeatureFlagAsRule?: boolean;
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
      archived: tab === 'ARCHIVED',
      hasExperiment,
      hasPrerequisites,
      tags,
      status,
      hasFeatureFlagAsRule
    }
  });
};
