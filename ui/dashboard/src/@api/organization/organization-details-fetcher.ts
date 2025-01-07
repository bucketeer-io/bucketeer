import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Organization } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface OrganizationDetailsFetcherParams {
  id: string;
}
export interface OrganizationDetailsResponse {
  organization: Organization;
}

export const organizationDetailsFetcher = async (
  _params?: OrganizationDetailsFetcherParams
): Promise<OrganizationDetailsResponse> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<OrganizationDetailsResponse>(
      `/v1/environment/get_organization?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
