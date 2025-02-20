import { EnvironmentRoleType } from './auth';
import { OrganizationRole } from './organization';

export interface Account {
  email: string;
  name: string;
  firstName: string;
  lastName: string;
  language: string;
  lastSeen: string;
  avatarImageUrl: string;
  organizationId: string;
  organizationRole: OrganizationRole;
  environmentCount: boolean;
  environmentRoles: [
    {
      environmentId: string;
      role: EnvironmentRoleType;
    }
  ];
  disabled: boolean;
  createdAt: string;
  updatedAt: string;
  searchFilters: [
    {
      id: string;
      name: string;
      query: string;
      filterTargetType: unknown;
      environmentId: string;
      defaultFilter: boolean;
    }
  ];
  tags: string[];
}

export interface AccountCollection {
  accounts: Array<Account>;
  cursor: string;
  totalCount: string;
}
