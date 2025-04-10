import { DomainEventEntityType, OrderBy, OrderDirection } from '@types';

export interface AuditLogsFilters {
  pageSize?: number;
  page?: number;
  searchQuery?: string;
  orderBy?: OrderBy;
  orderDirection?: OrderDirection;
  disabled?: boolean;
  from?: string;
  to?: string;
  entityType?: DomainEventEntityType;
  environmentId?: string
}
