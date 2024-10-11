import { CollectionStatusType } from '@types';

export interface OrganizationFilters {
  searchQuery: string;
  status: CollectionStatusType;
}
