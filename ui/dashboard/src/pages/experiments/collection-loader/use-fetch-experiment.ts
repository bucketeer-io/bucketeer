import { useQueryExperiments } from '@queries/experiments';
import { LIST_PAGE_SIZE } from 'constants/app';
import { pickBy } from 'lodash';
import { ExperimentStatus, OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { ExperimentTab } from '../types';

export const useFetchExperiments = ({
  page = 1,
  pageSize,
  orderBy,
  searchQuery,
  orderDirection,
  environmentId,
  status,
  ...params
}: {
  environmentId: string;
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  featureId?: string;
  featureVersion?: number;
  from?: string;
  to?: string;
  statuses?: ExperimentStatus[];
  maintainer?: string;
  status?: ExperimentTab;
}) => {
  const cursor = (page - 1) * LIST_PAGE_SIZE;
  const _params = pickBy(
    { ...params },
    (v, k) =>
      !['filterByTab', 'filterBySummary', 'isFilter'].includes(k) &&
      isNotEmpty(v)
  );

  return useQueryExperiments({
    params: {
      pageSize: pageSize || LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy,
      orderDirection,
      searchKeyword: searchQuery,
      environmentId,
      archived: status === 'ARCHIVED',
      ..._params
    }
  });
};
