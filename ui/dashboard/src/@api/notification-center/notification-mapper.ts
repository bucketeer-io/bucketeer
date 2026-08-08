import { getLanguage } from 'i18n';
import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization
} from '@types';

// The `Notification` proto message as it comes over the wire: `localizations`
// (all languages) is always populated, but only ListNotifications,
// GetNotification, and the drafts list resolve the singular `localization`
// field server-side (see pkg/notification/api/api.go). Create/update/publish
// responses leave it unset, so this must still resolve a display language
// client-side as a fallback for those.
export interface NotificationWire extends Omit<
  NotificationCenterFeedItem,
  'title' | 'content' | 'tags'
> {
  localization?: NotificationCenterLocalization;
}

const resolveLocalization = (
  localizations: NotificationCenterLocalization[]
): NotificationCenterLocalization | undefined => {
  const lang = getLanguage();
  return (
    localizations.find(l => l.language === lang) ??
    localizations.find(l => l.language === 'en') ??
    localizations[0]
  );
};

// Flattens a raw `Notification` wire object into the shape the notification
// center UI renders, resolving the display language the backend doesn't.
export const toFeedItem = (
  notification: NotificationWire
): NotificationCenterFeedItem => {
  const loc =
    notification.localization ??
    resolveLocalization(notification.localizations);
  return {
    ...notification,
    title: loc?.title ?? '',
    content: loc?.content ?? '',
    tags: loc?.tags ?? []
  };
};
