import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization,
  NotificationCenterPublishPayload,
  NotificationCenterStatus,
  NotificationCenterTag
} from '@types';

export type NotificationTab = 'unread' | 'read' | 'publish';

export type SortOption = 'newest' | 'oldest';

export {
  NotificationCenterStatus as NotificationStatus,
  type NotificationCenterTag as NotificationTag,
  type NotificationCenterLocalization as NotificationLocalizationInput,
  type NotificationCenterFeedItem as FeedNotification,
  type NotificationCenterFeedItem as NotificationDraft,
  type NotificationCenterPublishPayload as PublishNotificationInput
};

export type NotificationDetail = NotificationCenterFeedItem;

export interface NotificationFilters {
  tab: NotificationTab;
  searchQuery: string;
  sort: SortOption;
  from?: string;
  to?: string;
  notificationId?: string;
}
