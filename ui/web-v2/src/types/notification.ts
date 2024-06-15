import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const notificationSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type NotificationSortOption = typeof notificationSortOptions[number];

export function isNotificationSortOption(
  so: unknown
): so is NotificationSortOption {
  return typeof so === 'string' && notificationSortOptions.includes(so);
}

export interface NotificationSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
  enabled?: string;
}
