import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from './list';

const webhookSortOptions = [
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
];

export type WebhookSortOption = typeof webhookSortOptions[number];

export function isWebhookSortOption(so: unknown): so is WebhookSortOption {
  return typeof so === 'string' && webhookSortOptions.includes(so);
}

export interface WebhookSearchOptions {
  q?: string;
  sort?: string;
  page?: string;
}
