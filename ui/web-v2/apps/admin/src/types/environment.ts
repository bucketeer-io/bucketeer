import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_ID_ASC,
  SORT_OPTIONS_ID_DESC,
} from './list';

const projectSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_ID_ASC,
  SORT_OPTIONS_ID_DESC,
];

export type EnvironmentSortOption = typeof projectSortOptions[number];

export function isEnvironmentSortOption(
  so: unknown
): so is EnvironmentSortOption {
  return typeof so === 'string' && projectSortOptions.includes(so);
}

export interface EnvironmentSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  projectId?: string;
}
