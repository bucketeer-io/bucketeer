import { OrderBy, OrderDirection } from '@types';

interface SortingListFields {
  [x: string]: OrderBy;
}

interface OrderDirectionType {
  [x: string]: OrderDirection;
}

export const orderDirectionType: OrderDirectionType = {
  asc: 'ASC',
  desc: 'DESC'
};

export const sortingListFields: SortingListFields = {
  default: 'DEFAULT',
  id: 'ID',
  name: 'NAME',
  email: 'EMAIL',
  role: 'ROLE',
  urlCode: 'URL_CODE',
  createdAt: 'CREATED_AT',
  updatedAt: 'UPDATED_AT',
  userCount: 'USER_COUNT',
  environmentCount: 'ENVIRONMENT_COUNT',
  projectCount: 'PROJECT_COUNT',
  featureFlagCount: 'FEATURE_COUNT'
};
