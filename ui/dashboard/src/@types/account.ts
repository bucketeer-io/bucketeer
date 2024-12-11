import { EnvironmentRoleType } from './auth';
import { OrganizationRole } from './organization';

export interface EnvironmentRole {
  environmentId: string;
  role: EnvironmentRoleType;
}

export interface SearchFilter {
  id: string;
  name: string;
  query: string;
  filterTargetType: unknown;
  environmentId: string;
  defaultFilter: boolean;
}

export interface Account {
  email: string;
  name: string;
  firstName: string;
  lastName: string;
  language: string;
  lastSeen: string;
  avatarImageUrl: string;
  avatarImage: string;
  avatarFileType: string;
  organizationId: string;
  organizationRole: OrganizationRole;
  environmentCount: boolean;
  environmentRoles: EnvironmentRole[];
  disabled: boolean;
  createdAt: string;
  updatedAt: string;
  searchFilters: SearchFilter[];
}

export interface AccountCollection {
  accounts: Array<Account>;
  cursor: string;
  totalCount: string;
}
