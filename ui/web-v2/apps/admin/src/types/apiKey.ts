import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const apiKeySortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type APIKeySortOption = typeof apiKeySortOptions[number];

export function isAPIKeySortOption(so: unknown): so is APIKeySortOption {
  return typeof so === 'string' && apiKeySortOptions.includes(so);
}

export interface APIKeySearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  enabled?: string;
}
