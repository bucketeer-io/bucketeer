import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickby';
import { CollectionParams, OrganizationCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface OrganizationsFetcherParams extends CollectionParams {
  archived?: boolean;
}

export const organizationsFetcher = async (
  _params?: OrganizationsFetcherParams
): Promise<OrganizationCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .post<OrganizationCollection>('/v1/environment/list_organizations', params)
    .then(response => response.data);
};
