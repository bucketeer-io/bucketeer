import { OrderBy, OrderDirection } from '@types';

export interface AuditLogsFilters {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  from?: string;
  to?: string;
  entityType?: number;
  environmentId?: string;
  range?: string;
}

export enum ExpandOrCollapse {
  EXPAND = 'EXPAND',
  COLLAPSE = 'COLLAPSE'
}

export enum AuditLogTab {
  CHANGES = 'CHANGES',
  SNAPSHOT = 'SNAPSHOT'
}
