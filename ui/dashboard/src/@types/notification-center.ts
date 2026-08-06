export enum NotificationCenterStatus {
  DRAFT = 'DRAFT',
  PUBLISHED = 'PUBLISHED'
}

export interface NotificationCenterTag {
  name: string;
  color: string;
}

// One (notification, language) pair. Markdown source lives in `content`.
export interface NotificationCenterLocalization {
  language: string;
  tags: NotificationCenterTag[];
  title: string;
  content: string;
}

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
}

export interface NotificationCenterDraftCollection {
  notifications: NotificationCenterFeedItem[];
  cursor: string;
  totalCount: string;
}

export interface NotificationCenterPublishPayload {
  status: NotificationCenterStatus;
  localizations: NotificationCenterLocalization[];
}
