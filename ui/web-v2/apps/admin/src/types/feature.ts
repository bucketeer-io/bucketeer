import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from './list';

const featureSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
];
export type FeatureSortOption = typeof featureSortOptions[number];

export function isFeatureSortOption(so: unknown): so is FeatureSortOption {
  return typeof so === 'string' && featureSortOptions.includes(so);
}

export interface FeatureSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  maintainerId?: string;
  enabled?: string;
  archived?: string;
  hasExperiment?: string;
  tagIds?: string[];
  hasPrerequisites?: string;
}
