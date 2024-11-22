export type OrganizationRole =
  | 'Organization_UNASSIGNED'
  | 'Organization_MEMBER'
  | 'Organization_ADMIN'
  | 'Organization_OWNER';

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
  environmentCount: number;
  projectCount: number;
  userCount: number;
  ownerEmail: string;
}

export interface OrganizationCollection {
  Organizations: Array<Organization>;
  cursor: string;
  totalCount: string;
}
