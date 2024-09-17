import axiosClient from '@api/axios-client';
import { OrderBy, OrderDirection, OrganizationsCollection } from '@types';

export interface OrganizationsFetcherParams {
  pageSize: number;
  cursor: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  searchKeyword?: string;
  disabled: boolean;
  archived: boolean;
}

export const organizationsFetcher = async (
  params?: OrganizationsFetcherParams
): Promise<OrganizationsCollection> => {
  return axiosClient
    .post<OrganizationsCollection>('/v1/environment/list_organizations', params)
    .then(response => response.data);
};
