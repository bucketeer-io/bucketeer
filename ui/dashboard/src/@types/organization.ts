export interface OrganizationRole {
  ORGANIZATION_UNASSIGNED: 0;
  ORGANIZATION_MEMBER: 1;
  ORGANIZATION_ADMIN: 2;
  ORGANIZATION_OWNER: 3;
}

export interface Organization {
  id: string;
  name: string;
  urlCode: string;
  description: string;
  disabled: boolean;
  archived: boolean;
  trial: boolean;
  createdAt: string;
  updatedAt: string;
  systemAdmin: boolean;
}

export interface OrganizationCollection {
  Organizations: Array<Organization>;
  cursor: string;
  totalCount: string;
}
