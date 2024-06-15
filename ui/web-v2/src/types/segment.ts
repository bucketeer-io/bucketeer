import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const segmentSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type SegmentSortOption = typeof segmentSortOptions[number];

export function isSegmentSortOption(so: unknown): so is SegmentSortOption {
  return typeof so === 'string' && segmentSortOptions.includes(so);
}

export interface SegmentSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  inUse?: string;
}
