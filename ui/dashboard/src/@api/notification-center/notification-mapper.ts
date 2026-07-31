import { getLanguage } from 'i18n';
import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization
} from '@types';

// The `Notification` proto message as it comes over the wire: `localizations`
// (all languages) is populated, but the backend never resolves the singular
// `localization` field (see proto/notification/notification.proto), so admin
// endpoints must resolve a display language client-side.
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
