import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const accountSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type AccountSortOption = typeof accountSortOptions[number];

export function isAccountSortOption(so: unknown): so is AccountSortOption {
  return typeof so === 'string' && accountSortOptions.includes(so);
}

export interface AccountSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  role?: string;
  enabled?: string;
}
