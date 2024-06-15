import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const projectSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type ProjectSortOption = typeof projectSortOptions[number];

export function isProjectSortOption(so: unknown): so is ProjectSortOption {
  return typeof so === 'string' && projectSortOptions.includes(so);
}

export interface ProjectSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  enabled?: string;
}
