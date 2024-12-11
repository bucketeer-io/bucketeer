import { SearchFilter } from './account';
import { Environment } from './environment';
import { Organization, OrganizationRole } from './organization';
import { Project } from './project';

export interface ServerError {
  code: number;
  details: [];
  message: string;
}

export type EnvironmentRoleType =
  | 'Environment_UNASSIGNED'
  | 'Environment_EDITOR'
  | 'Environment_VIEWER';

export interface AuthTypeMap {
  AUTH_TYPE_UNSPECIFIED: 0;
  AUTH_TYPE_USER_PASSWORD: 1;
  AUTH_TYPE_GOOGLE: 2;
  AUTH_TYPE_GITHUB: 3;
}

export interface SignInForm {
  email: string;
  password: string;
}

export interface AuthToken {
  accessToken: string;
  tokenType: string;
  refreshToken: string;
  expiry: number;
}

export interface AuthResponse {
  token: AuthToken;
}

export interface UserInfoForm {
  first_name: string;
  last_name: string;
  language: string;
}

export interface EnvironmentRole {
  environment: Environment;
  project: Project;
  role: EnvironmentRoleType;
}

export interface ConsoleAccount {
  email: string;
  name: string;
  avatarUrl: string;
  isSystemAdmin: boolean;
  organization: Organization;
  environmentRoles: EnvironmentRole[];
  organizationRole: OrganizationRole;
  firstName: string;
  lastName: string;
  language: string;
  lastSeen: string;
  avatarImageUrl: string;
  avatarImage: string;
  avatarFileType: string;
  searchFilters: SearchFilter[];
}

export interface ConsoleAccountResponse {
  account: ConsoleAccount;
}

export interface AuthUrlResponse {
  url: string;
}
