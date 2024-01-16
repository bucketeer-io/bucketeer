import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
} from './list';

const auditLogSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
];

export type AuditLogSortOption = typeof auditLogSortOptions[number];

export function isAuditLogSortOption(so: unknown): so is AuditLogSortOption {
  return typeof so === 'string' && auditLogSortOptions.includes(so);
}

export interface AuditLogSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  from?: number;
  to?: number;
  entityType?: number;
}
