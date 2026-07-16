import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Checkbox from 'components/checkbox';
import NotificationCard from '../elements/notification-card';
import TagChip from '../elements/tag-chip';
import { FeedNotification } from '../types';

interface NotificationRowProps {
  notification: FeedNotification;
  // Selection (checkbox + bulk actions) applies only to unread notifications.
  selectable?: boolean;
  selected?: boolean;
  onSelectedChange?: (selected: boolean) => void;
  onClick?: () => void;
}

const NotificationRow = ({
  notification,
  selectable = false,
  selected = false,
  onSelectedChange,
  onClick
}: NotificationRowProps) => {
  const formatDateTime = useFormatDateTime();
  return (
    <div className="flex items-start gap-3 rounded-lg border border-gray-200 p-4 transition-colors hover:border-gray-300">
      {selectable && (
        <div className="pt-0.5">
          <Checkbox
            checked={selected}
            onCheckedChange={value => onSelectedChange?.(Boolean(value))}
          />
        </div>
      )}
      <div className="flex-1">
        <NotificationCard
          onClick={onClick}
          bordered={false}
          header={
            <div className="flex w-full items-center gap-2">
              {!notification.read && (
                <span className="size-2 shrink-0 rounded-full bg-primary-500" />
              )}
              <span
                className={cn(
                  'typo-para-medium text-gray-900',
                  !notification.read && 'font-semibold'
                )}
              >
                {notification.title}
              </span>
              <span className="ml-auto typo-para-tiny text-gray-500">
                {formatDateTime(notification.publishedAt)}
              </span>
            </div>
          }
          footer={
            <span className="typo-para-small text-gray-500">
              {notification.createdBy.split('@')[0]}
            </span>
          }
        >
          {notification.tags.map(tag => (
            <TagChip key={tag.name} tag={tag} />
          ))}
        </NotificationCard>
      </div>
    </div>
  );
};

export default NotificationRow;
