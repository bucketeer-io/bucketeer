import { getLanguage } from 'i18n';
import {
  NotificationCenterFeedItem,
  NotificationCenterLocalization
} from '@types';

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
