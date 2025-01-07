import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, OrganizationCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface OrganizationsFetcherParams extends CollectionParams {
  archived?: boolean;
}

export const organizationsFetcher = async (
  _params?: OrganizationsFetcherParams
): Promise<OrganizationCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<OrganizationCollection>(
      `/v1/environment/list_organizations?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
