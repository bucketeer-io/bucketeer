import axiosClient from '@api/axios-client';
import { OrderBy, OrderDirection, ProjectCollection } from '@types';

export interface ProjectFetcherParams {
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword?: string;
  disabled: boolean;
  organizationIds: string[];
}

export const projectsFetcher = async (
  params?: ProjectFetcherParams
): Promise<ProjectCollection> => {
  return axiosClient
    .post<ProjectCollection>('/v1/environment/list_projects', params)
    .then(response => response.data);
};
