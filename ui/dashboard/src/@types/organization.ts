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
  createdAt: number;
  updatedAt: number;
  systemAdmin: boolean;
}
