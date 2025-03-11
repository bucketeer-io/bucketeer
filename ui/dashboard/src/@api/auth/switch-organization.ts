import axiosClient from '@api/axios-client';
import { SwitchOrganizationPayload, SwitchOrganizationResponse } from '@types';

export const switchOrganization = async (
  formValues: SwitchOrganizationPayload
): Promise<SwitchOrganizationResponse> => {
  return axiosClient
    .post(`/v1/auth/switch_organization`, formValues)
    .then(response => response.data);
};
