export enum NotificationCenterStatus {
  DRAFT = 0,
  PUBLISHED = 1
}

export interface NotificationCenterTag {
  name: string;
  color: string;
}

// One (notification, language) pair. Markdown source lives in `content`.
export interface NotificationCenterLocalization {
  language: string; // BCP 47 code, e.g. 'en', 'ja'
  tags: NotificationCenterTag[];
  title: string;
  content: string;
}

// A feed/draft row projected into one language: the resolved localization's
// fields are flattened, matching what the list/drafts endpoints return for
// display. `localizations` carries every language version so the detail
// panel and publish form can edit them all.
export interface NotificationCenterFeedItem {
  id: string;
  title: string;
  content: string; // Markdown source, from the resolved localization
  tags: NotificationCenterTag[];
  read: boolean;
  status: NotificationCenterStatus;
  publishedAt: string; // epoch seconds; "0" while draft
  createdAt: string;
  updatedAt: string;
  createdBy: string;
  lastEditedBy: string;
  localizations: NotificationCenterLocalization[];
}

export interface NotificationCenterFeedCollection {
  notifications: NotificationCenterFeedItem[];
  cursor: string;
  totalCount: string;
  unreadCount: string;
  readCount: string;
}

export interface NotificationCenterDraftCollection {
  notifications: NotificationCenterFeedItem[];
  cursor: string;
  totalCount: string;
}

export interface NotificationCenterUnreadCount {
  count: string;
}

// ---------------------------------------------------------------------------
// Write payloads
// ---------------------------------------------------------------------------

// The author identity (created_by / published_by) is intentionally absent —
// the backend fills it from the authenticated user.
export interface NotificationCenterPublishPayload {
  status: NotificationCenterStatus;
  localizations: NotificationCenterLocalization[];
}

export interface NotificationCenterUpdatePayload {
  id: string;
  status: NotificationCenterStatus;
  localizations: NotificationCenterLocalization[];
}
