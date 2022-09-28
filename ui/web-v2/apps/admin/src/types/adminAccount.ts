import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_ID_ASC,
  SORT_OPTIONS_ID_DESC,
} from './list';

const accountSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_ID_ASC,
  SORT_OPTIONS_ID_DESC,
];

export type AdminAccountSortOption = typeof accountSortOptions[number];

export function isAdminAccountSortOption(
  so: unknown
): so is AdminAccountSortOption {
  return typeof so === 'string' && accountSortOptions.includes(so);
}

export interface AdminAccountSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  enabled?: string;
}
