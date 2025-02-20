import axiosClient from '@api/axios-client';
import { OrganizationRole } from '@types';
import { EnvironmentRoleItem } from './account-creator';

export interface AccountUpdaterParams {
  email: string;
  organizationId: string;
  name?: string;
  avatarImageUrl?: string;
  organizationRole?: {
    role: OrganizationRole;
  };
  environmentRoles?: EnvironmentRoleItem[];
  firstName?: string;
  lastName?: string;
  language?: string;
  lastSeen?: string;
  avatar?: {
    avatarImage: string;
    avatarFileType: string;
  };
  disabled?: boolean;
  tags?: string[];
}

export const accountUpdater = async (params?: AccountUpdaterParams) => {
  return axiosClient
    .post('/v1/account/update_account', params)
    .then(response => response.data);
};
