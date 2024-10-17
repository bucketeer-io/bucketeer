import { EnvironmentRoleType } from './auth';
import { OrganizationRole } from './organization';

export interface Account {
  email: string;
  name: string;
  avatarImageUrl: string;
  organizationId: string;
  organizationRole: OrganizationRole;
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
}

export interface AccountCollection {
  accounts: Array<Account>;
  cursor: string;
  totalCount: string;
}
