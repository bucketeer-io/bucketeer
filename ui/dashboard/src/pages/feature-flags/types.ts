import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export type FlagDataType = 'number' | 'string' | 'json' | 'boolean';
export type FlagActionType =
  | 'ARCHIVE'
  | 'UNARCHIVE'
  | 'CLONE'
  | 'ACTIVE'
  | 'INACTIVE';
export type SummaryType = 'TOTAL' | 'ACTIVE' | 'INACTIVE';

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

export enum FeatureActivityStatus {
  ACTIVE = 'active',
  NEW = 'new',
  INACTIVE = 'in-active'
}
