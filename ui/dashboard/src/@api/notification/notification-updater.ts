import axiosClient from '@api/axios-client';
import { OrganizationRole } from '@types';
import { EnvironmentRoleItem } from './notification-creator';

type WriteType =
  | 'WriteType_UNSPECIFIED'
  | 'WriteType_OVERRIDE'
  | 'WriteType_PATCH';

export interface AccountUpdaterParams {
  email: string;
  organizationId: string;
  changeNameCommand?: {
    name: string;
  };
  changeAvatarUrlCommand?: {
    avatarImageUrl: string;
  };
  changeOrganizationRoleCommand: {
    role: OrganizationRole;
  };
  changeEnvironmentRolesCommand: {
    roles: EnvironmentRoleItem[];
    writeType?: WriteType;
  };
  changeFirstNameCommand: {
    firstName: string;
  };
  changeLastNameCommand: {
    lastName: string;
  };
  changeLanguageCommand: {
    language: string;
  };
  changeLastSeenCommand?: {
    lastSeen: string;
  };
  changeAvatarCommand?: {
    avatarImage: string;
    avatarFileType: string;
  };
}

export const accountUpdater = async (params?: AccountUpdaterParams) => {
  return axiosClient
    .post('/v1/account/update_account', params)
    .then(response => response.data);
};
