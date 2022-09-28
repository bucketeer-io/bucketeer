import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const experimentSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];
export type ExperimentSortOption = typeof experimentSortOptions[number];

export function isExperimentSortOption(
  so: unknown
): so is ExperimentSortOption {
  return typeof so === 'string' && experimentSortOptions.includes(so);
}

export interface ExperimentSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  maintainerId?: string;
  status?: string;
  archived?: string;
}
