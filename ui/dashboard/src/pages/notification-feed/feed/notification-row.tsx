import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Checkbox from 'components/checkbox';
import { markdownToText } from '../elements/markdown-content';
import NotificationCard from '../elements/notification-card';
import TagChip from '../elements/tag-chip';
import { FeedNotification } from '../types';

interface NotificationRowProps {
  notification: FeedNotification;
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
    <div className="flex min-h-[130px] items-stretch gap-3 rounded-lg border border-gray-200 p-4 transition-colors hover:border-gray-300">
      {selectable && (
        <div className="self-start pt-0.5">
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
            <div className="flex w-full flex-col gap-1">
              <div className="flex w-full items-center gap-2">
                {!notification.read && (
                  <span className="size-2 shrink-0 rounded-full bg-primary-500" />
                )}
                <span
                  className={cn(
                    'truncate typo-para-medium text-gray-900',
                    !notification.read && 'font-semibold'
                  )}
                >
                  {notification.title}
                </span>
                {notification.tags.length > 0 && (
                  <div className="flex shrink-0 items-center gap-1.5">
                    {notification.tags.map(tag => (
                      <TagChip key={tag.name} tag={tag} />
                    ))}
                  </div>
                )}
                <span className="ml-auto shrink-0 typo-para-tiny text-gray-500">
                  {formatDateTime(notification.publishedAt)}
                </span>
              </div>
              {notification.content && (
                <p className="line-clamp-2 typo-para-small text-gray-500">
                  {markdownToText(notification.content)}
                </p>
              )}
            </div>
          }
          footer={
            <span className="typo-para-small text-gray-500">
              {notification.createdBy.split('@')[0]}
            </span>
          }
        />
      </div>
    </div>
  );
};

export default NotificationRow;
