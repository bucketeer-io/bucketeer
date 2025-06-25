import axiosClient from '@api/axios-client';
import { Account, EnvironmentRoleType, OrganizationRole } from '@types';

export interface EnvironmentRoleItem {
  environmentId: string;
  role: EnvironmentRoleType;
}

export interface AccountCreatorPayload {
  organizationId: string;
  email: string;
  organizationRole: OrganizationRole;
  environmentRoles?: EnvironmentRoleItem[];
  tags: string[];
}

export interface AccountCreatorResponse {
  account: Array<Account>;
}

export const accountCreator = async (
  params?: AccountCreatorPayload
): Promise<AccountCreatorResponse> => {
  return axiosClient
    .post<AccountCreatorResponse>('/v1/account/create_account', params)
    .then(response => response.data);
};
