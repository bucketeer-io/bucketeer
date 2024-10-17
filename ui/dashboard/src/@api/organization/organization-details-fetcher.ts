import axiosClient from '@api/axios-client';
import { Organization } from '@types';

export interface OrganizationDetailsFetcherParams {
  id: string;
}
export interface OrganizationDetailsResponse {
  organization: Organization;
}

export const organizationDetailsFetcher = async (
  params?: OrganizationDetailsFetcherParams
): Promise<OrganizationDetailsResponse> => {
  return axiosClient
    .post<OrganizationDetailsResponse>(
      '/v1/environment/get_organization',
      params
    )
    .then(response => response.data);
};
