import { format } from 'timeago.js';
import { cn } from 'utils/style';
import Button from 'components/button';
import Checkbox from 'components/checkbox';
import TagChip from '../tag-chip';
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
      <Button
        type="button"
        variant="text"
        onClick={onClick}
        className="flex h-auto w-full flex-col items-start justify-start gap-2 whitespace-normal px-0 text-left"
      >
        <div className="flex items-center gap-2">
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
            {format(notification.createdAt)}
          </span>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {notification.tags.map(tag => (
            <TagChip key={tag.name} tag={tag} />
          ))}
        </div>
        <span className="typo-para-small text-gray-500">
          {notification.createdBy.split('@')[0]}
        </span>
      </Button>
    </div>
  );
};

export default NotificationRow;
