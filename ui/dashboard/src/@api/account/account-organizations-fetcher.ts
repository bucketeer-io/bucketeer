import axiosClient from '@api/axios-client';
import { Organization } from '@types';

export interface AccountOrganizationsResponse {
  organizations: Array<Organization>;
}

export const accountOrganizationFetcher =
  async (): Promise<AccountOrganizationsResponse> => {
    return axiosClient
      .get<AccountOrganizationsResponse>('/v1/account/my_organizations')
      .then(response => response.data);
  };
