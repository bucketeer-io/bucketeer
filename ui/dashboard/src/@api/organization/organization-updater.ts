import axiosClient from '@api/axios-client';
import { OrganizationResponse } from './organization-creator';

export interface OrganizationUpdatePayload {
  id: string;
  name: string;
  description?: string;
  ownerEmail: string;
}

export const organizationUpdater = async (
  params?: OrganizationUpdatePayload
): Promise<OrganizationResponse> => {
  return axiosClient
    .post<OrganizationResponse>('/v1/environment/update_organization', params)
    .then(response => response.data);
};
