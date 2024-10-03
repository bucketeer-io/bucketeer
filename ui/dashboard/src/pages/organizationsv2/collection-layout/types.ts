import { Organization } from '@types';

export interface ItemProps {
  organization: Organization;
}

// export interface ItemLayoutProps {
// 	layout: CollectionLayoutType;
// 	organization: Organization;
// }

export interface CollectionProps {
  isLoading?: boolean;
  organizations: Organization[];
}
