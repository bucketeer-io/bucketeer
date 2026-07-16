import { useState } from 'react';
import { useNavigate } from 'react-router';
import * as Popover from '@radix-ui/react-popover';
import * as ROUTING from 'constants/routing';
import { useTranslation } from 'i18n';
import { ExternalLink } from 'lucide-react';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconNotifications } from '@icons';
import {
  useFetchUnreadCount,
  useFetchFeed,
  useMarkAsRead,
  useMarkAllAsRead
} from 'pages/notification-feed/collection-loader/use-fetch-notifications';
import {
  firstMarkdownLink,
  markdownToText
} from 'pages/notification-feed/elements/markdown-content';
import TagChip from 'pages/notification-feed/elements/tag-chip';
import {
  FeedNotification,
  NotificationFilters
} from 'pages/notification-feed/types';
import { Badge } from 'components/badge';
import Icon from 'components/icon';

const PREVIEW_PAGE_SIZE = 5;

const previewFilters: NotificationFilters = {
  tab: 'unread',
  searchQuery: '',
  sort: 'newest'
};

const NotificationBell = ({
  environmentId,
  envUrlCode
}: {
  // The environment's internal id, used for API calls.
  environmentId: string;
  // The environment's URL code, used only to build the route link.
  envUrlCode: string;
}) => {
  const { t } = useTranslation(['common']);
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();

  const [isOpen, setIsOpen] = useState(false);

  const { data: unread } = useFetchUnreadCount(environmentId);
  const { data } = useFetchFeed(environmentId, false, 1, previewFilters);
  const markAsRead = useMarkAsRead(environmentId);
  const markAllAsRead = useMarkAllAsRead(environmentId);

  const unreadCount = Number(unread?.count ?? 0);
  const items = (data?.notifications ?? []).slice(0, PREVIEW_PAGE_SIZE);

  const goToFeed = () => {
    setIsOpen(false);
    navigate(`/${envUrlCode}${ROUTING.PAGE_PATH_NOTIFICATION_FEED}`);
  };

  const onSelectNotification = (notification: FeedNotification) => {
    if (!notification.read) markAsRead.mutate(notification.id);
    setIsOpen(false);
    navigate(`/${envUrlCode}${ROUTING.PAGE_PATH_NOTIFICATION_FEED}`, {
      state: { notification }
    });
  };

  return (
    <Popover.Root open={isOpen} onOpenChange={setIsOpen}>
      <Popover.Trigger asChild>
        <button type="button" className="relative flex bottom-1">
          <Icon icon={IconNotifications} color="primary-50" />
          {unreadCount > 0 && (
            <Badge
              variant="primary"
              className="absolute -right-1.5 -top-1.5 h-4 w-auto min-w-4 whitespace-nowrap bg-accent-red-500 px-1 text-white typo-para-tiny"
            >
              {unreadCount > 99 ? '99+' : unreadCount}
            </Badge>
          )}
        </button>
      </Popover.Trigger>

      <Popover.Content
        align="start"
        side="top"
        sideOffset={8}
        className="w-[380px] rounded-lg border-none bg-white p-0 shadow-menu"
      >
        <div className="flex items-center justify-between px-4 py-3.5">
          <span className="typo-head-bold-medium text-gray-900">
            {t('notifications')}
          </span>
          {unreadCount > 0 && (
            <button
              type="button"
              className="typo-para-small font-medium text-primary-500"
              onClick={() => markAllAsRead.mutate()}
            >
              {t('mark-all-as-read')}
            </button>
          )}
        </div>

        <div className="max-h-[400px] overflow-y-auto">
          {items.length === 0 ? (
            <p className="px-4 py-8 text-center typo-para-medium text-gray-500">
              {t('no-notifications')}
            </p>
          ) : (
            items.map(notification => {
              const link = firstMarkdownLink(notification.content);
              return (
                <button
                  key={notification.id}
                  type="button"
                  onClick={() => onSelectNotification(notification)}
                  className="flex w-full flex-col items-start gap-1.5 border-t border-gray-100 px-4 py-3.5 text-left first:border-t-0 hover:bg-gray-50"
                >
                  <div className="flex w-full items-start gap-2">
                    {!notification.read && (
                      <span className="mt-1.5 size-1.5 shrink-0 rounded-full bg-primary-500" />
                    )}
                    <span
                      className={cn(
                        'typo-para-medium text-gray-900',
                        !notification.read ? 'font-semibold' : 'ml-3.5'
                      )}
                    >
                      {notification.title}
                    </span>
                    {notification.tags[0] && (
                      <TagChip tag={notification.tags[0]} hideDot />
                    )}
                    <span className="ml-auto shrink-0 typo-para-tiny text-gray-500">
                      {formatDateTime(notification.publishedAt)}
                    </span>
                  </div>
                  <span
                    className={cn(
                      'line-clamp-2 typo-para-small text-gray-500',
                      !notification.read && 'pl-3.5'
                    )}
                  >
                    {markdownToText(notification.content)}
                  </span>
                  {link && (
                    <a
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      onClick={e => e.stopPropagation()}
                      className={cn(
                        'flex items-center gap-1 typo-para-small font-medium text-primary-500 hover:underline',
                        !notification.read && 'pl-3.5'
                      )}
                    >
                      {link.label}
                      <ExternalLink size={12} />
                    </a>
                  )}
                </button>
              );
            })
          )}
        </div>

        <button
          type="button"
          onClick={goToFeed}
          className="block w-full border-t border-gray-100 px-4 py-3 text-center typo-para-medium font-medium text-primary-500 hover:bg-gray-50"
        >
          {t('view-all-notifications')}
        </button>
      </Popover.Content>
    </Popover.Root>
  );
};

export default NotificationBell;
