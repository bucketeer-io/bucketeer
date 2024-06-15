import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const pushSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type PushSortOption = typeof pushSortOptions[number];

export function isPushSortOption(so: unknown): so is PushSortOption {
  return typeof so === 'string' && pushSortOptions.includes(so);
}

export interface PushSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
}
