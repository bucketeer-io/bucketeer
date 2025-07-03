import { CollectionStatusType, OrderBy, OrderDirection } from '@types';

export type FlagDataType = 'number' | 'string' | 'json' | 'boolean';
export type FlagActionType =
  | 'ARCHIVE'
  | 'UNARCHIVE'
  | 'CLONE'
  | 'ACTIVE'
  | 'INACTIVE';

export enum StatusFilterType {
  NEVER_USED = 'NEW',
  RECEIVING_TRAFFIC = 'ACTIVE', // The backend uses ACTIVE for receiving traffic
  NO_RECENT_TRAFFIC = 'NO_ACTIVITY', // The backend uses NO_ACTIVITY for no recent traffic
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
  NEVER_USED = 'never-used',
  RECEIVING_TRAFFIC = 'receiving-traffic',
  NO_RECENT_TRAFFIC = 'no-recent-traffic'
}

export enum FlagOperationType {
  ROLLOUT = 'ROLLOUT',
  SCHEDULE = 'SCHEDULE',
  EVENT_RATE = 'EVENT_RATE'
}
