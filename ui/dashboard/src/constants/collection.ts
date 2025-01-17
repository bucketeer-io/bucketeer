import { OrderBy } from '@types';

interface SortingListFields {
  [x: string]: OrderBy;
}

export const sortingListFields: SortingListFields = {
  default: 'DEFAULT',
  id: 'ID',
  name: 'NAME',
  email: 'EMAIL',
  state: 'STATE',
  role: 'ROLE',
  urlCode: 'URL_CODE',
  createdAt: 'CREATED_AT',
  updatedAt: 'UPDATED_AT',
  userCount: 'USER_COUNT',
  environmentCount: 'ENVIRONMENT_COUNT',
  projectCount: 'PROJECT_COUNT',
  organizationRole: 'ORGANIZATION_ROLE',
  featureFlagCount: 'FEATURE_COUNT',
  creatorEmail: 'CREATOR_EMAIL',
  lastSeen: 'LAST_SEEN',
  environment: 'ENVIRONMENT',
  tags: 'TAGS'
};
