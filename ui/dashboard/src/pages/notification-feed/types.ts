export type NotificationTab = 'unread' | 'read' | 'publish';

export type SortOption = 'newest' | 'oldest';

export enum NotificationStatus {
  DRAFT = 0,
  PUBLISHED = 1
}

export interface NotificationTag {
  name: string;
  color: string;
}

export interface NotificationRecord {
  id: string;
  status: NotificationStatus;
  createdBy: string; // editor email
  lastEditedBy: string;
  publishedBy?: string; // null while draft
  publishedAt: number; // epoch ms; 0 while draft
  createdAt: number;
  updatedAt: number;
}

// Row in `notification_localization`. One per (notificationId, language).
export interface NotificationLocalization {
  notificationId: string;
  language: string; // BCP 47 code: 'en', 'ja'
  tags: NotificationTag[];
  title: string;
  content: string; // Markdown source
}

// Row in `notification_read`. Read state is per viewer (email), global across
// organizations, so it lives outside the notification record.
export interface NotificationReadRow {
  notificationId: string;
  email: string;
  readAt: number; // epoch ms
}

// ---------------------------------------------------------------------------
// View models (a record projected into one language + one viewer's read state)
// ---------------------------------------------------------------------------

export interface FeedNotification {
  id: string;
  title: string;
  content: string; // Markdown source, from the resolved localization
  tags: NotificationTag[];
  read: boolean; // derived from notification_read for the current viewer
  status: NotificationStatus;
  publishedAt: number; // epoch ms; 0 while draft
  createdAt: number;
  updatedAt: number;
  createdBy: string;
  lastEditedBy: string;
  // Every language version, so it can be shown/edited in the detail panel.
  localizations: NotificationLocalizationInput[];
}

export interface NotificationDraft {
  id: string;
  title: string;
  content: string;
  tags: NotificationTag[];
  status: NotificationStatus;
  createdAt: number;
  updatedAt: number; // epoch ms
  createdBy: string;
  lastEditedBy: string;
  // Every language version, so the publish form can edit them all when editing.
  localizations: NotificationLocalizationInput[];
}

// The shape the detail SlideModal renders. Both FeedNotification and
// NotificationDraft satisfy it, so either can be shown in the panel.
export interface NotificationDetail {
  id: string;
  title: string;
  content: string;
  tags: NotificationTag[];
  status: NotificationStatus;
  createdAt: number;
  updatedAt: number;
  createdBy: string;
  lastEditedBy: string;
  // Every language version, for editing all of them from the detail panel.
  localizations: NotificationLocalizationInput[];
}

// ---------------------------------------------------------------------------
// Publish/draft request payloads (what the frontend sends to the backend)
// ---------------------------------------------------------------------------

// One language's worth of a notification, written to `notification_localization`.
export interface NotificationLocalizationInput {
  language: string; // 'en' | 'ja'
  title: string;
  content: string; // Markdown source
  tags: NotificationTag[];
}

// Request body for publishing or saving a draft. The author identity
// (created_by / published_by) is intentionally absent — the backend fills it
// from the authenticated user, so it cannot be spoofed by the client.
export interface PublishNotificationInput {
  status: NotificationStatus;
  localizations: NotificationLocalizationInput[];
}

export interface NotificationFilters {
  tab: NotificationTab;
  searchQuery: string;
  sort: SortOption;
  days?: number;
  from?: number;
  to?: number;
}

// Query args for the paginated feed. `read` selects the Unread/Read tab; the
// rest mirror NotificationFilters plus pagination.
export interface FeedQuery {
  read: boolean;
  page: number; // 1-based
  pageSize: number;
  searchQuery?: string;
  sort?: SortOption;
  from?: number;
  to?: number;
}

// Paginated feed response. `total` drives "Showing 1–10 of N" and the page
// controls; the two counts drive the Unread/Read tab labels.
export interface FeedPage {
  items: FeedNotification[];
  total: number; // matches the current `read` filter
  unreadCount: number;
  readCount: number;
}
