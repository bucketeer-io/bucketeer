import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Checkbox from 'components/checkbox';
import { markdownToText } from '../elements/markdown-content';
import NotificationCard from '../elements/notification-card';
import TagList from '../elements/tag-list';
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
      <div className="flex-1 w-full">
        <NotificationCard
          onClick={onClick}
          bordered={false}
          header={
            <div className="flex w-full flex-col gap-1">
              <div className="flex w-full min-w-0 items-center gap-2">
                {!notification.read && (
                  <span className="size-2 shrink-0 rounded-full bg-primary-500" />
                )}
                <span
                  className={cn(
                    'min-w-0 flex-1 truncate typo-para-medium text-gray-900',
                    !notification.read && 'font-semibold'
                  )}
                >
                  {notification.title}
                </span>
              </div>
              {notification.content && (
                <p className="line-clamp-2 break-words typo-para-small text-gray-500">
                  {markdownToText(notification.content)}
                </p>
              )}
            </div>
          }
          footer={
            <div className="flex w-full items-center justify-between">
              <span className="typo-para-small text-gray-500">
                {notification.createdBy.split('@')[0]}
              </span>
              <span className="typo-para-tiny text-gray-500">
                {formatDateTime(notification.publishedAt)}
              </span>
            </div>
          }
        >
          <TagList tags={notification.tags} />
        </NotificationCard>
      </div>
    </div>
  );
};

export default NotificationRow;
