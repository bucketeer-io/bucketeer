import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export type FlagDataType = 'number' | 'string' | 'json' | 'boolean';
export type FlagActionType =
  | 'ARCHIVE'
  | 'UNARCHIVE'
  | 'CLONE'
  | 'ACTIVE'
  | 'INACTIVE';

export enum StatusFilterType {
  NEW = 'NEW',
  ACTIVE = 'ACTIVE',
  NO_ACTIVITY = 'NO_ACTIVITY',
  UNKNOWN = 'UNKNOWN'
}

export interface FlagFilters {
  page: number;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
  archived?: boolean;
  searchQuery: string;
  tab: CollectionStatusType;
  environmentId: string;
  pageSize?: number;
  maintainer?: string;
  hasExperiment?: boolean;
  enabled?: boolean;
  hasPrerequisites?: boolean;
  tags?: string[];
  status?: StatusFilterType;
  hasFeatureFlagAsRule?: boolean;
}

export enum FeatureActivityStatus {
  ACTIVE = 'active',
  NEW = 'new',
  INACTIVE = 'in-active'
}

export enum FlagOperationType {
  ROLLOUT = 'ROLLOUT',
  SCHEDULE = 'SCHEDULE',
  EVENT_RATE = 'EVENT_RATE'
}
