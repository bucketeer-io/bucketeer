import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const goalSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];
export type GoalSortOption = typeof goalSortOptions[number];

export function isGoalSortOption(so: unknown): so is GoalSortOption {
  return typeof so === 'string' && goalSortOptions.includes(so);
}

export interface GoalSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  status?: string;
  archived?: string;
}
