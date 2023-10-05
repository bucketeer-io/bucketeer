import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const environmentSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type EnvironmentSortOption = typeof environmentSortOptions[number];

export function isEnvironmentSortOption(
  so: unknown
): so is EnvironmentSortOption {
  return typeof so === 'string' && environmentSortOptions.includes(so);
}

export interface EnvironmentSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  projectId?: string;
}
