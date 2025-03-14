import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export type FlagStatusType = 'new' | 'no_activity' | 'active';
export type FlagDataType = 'number' | 'string' | 'json' | 'boolean';
export type FlagActionType = 'ARCHIVE' | 'UNARCHIVE' | 'CLONE' | 'ACTIVE';

export interface FlagFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  archived?: boolean;
  searchQuery: string;
  status: CollectionStatusType;
  environmentId: string;
  pageSize?: number;
  maintainer?: string;
  hasExperiment?: boolean;
  enabled?: boolean;
  hasPrerequisites?: boolean;
  tags?: string[];
}
