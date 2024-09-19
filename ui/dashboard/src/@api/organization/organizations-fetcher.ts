import axiosClient from '@api/axios-client';
import { CollectionParams, OrganizationCollection } from '@types';

export interface OrganizationsFetcherParams extends CollectionParams {
  archived: boolean;
}

export const organizationsFetcher = async (
  params?: OrganizationsFetcherParams
): Promise<OrganizationCollection> => {
  return axiosClient
    .post<OrganizationCollection>('/v1/environment/list_organizations', params)
    .then(response => response.data);
};
