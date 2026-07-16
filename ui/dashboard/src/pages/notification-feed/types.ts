import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization,
  NotificationCenterPublishPayload,
  NotificationCenterStatus,
  NotificationCenterTag
} from '@types';

export type NotificationTab = 'unread' | 'read' | 'publish';

export type SortOption = 'newest' | 'oldest';

// Re-exported so the rest of this page imports everything from one place.
export {
  NotificationCenterStatus as NotificationStatus,
  type NotificationCenterTag as NotificationTag,
  type NotificationCenterLocalization as NotificationLocalizationInput,
  type NotificationCenterFeedItem as FeedNotification,
  type NotificationCenterFeedItem as NotificationDraft,
  type NotificationCenterPublishPayload as PublishNotificationInput
};

// The shape the detail SlideModal renders. Both a feed item and a draft
// satisfy it.
export type NotificationDetail = NotificationCenterFeedItem;

export interface NotificationFilters {
  tab: NotificationTab;
  searchQuery: string;
  sort: SortOption;
  days?: number;
  from?: number;
  to?: number;
}
