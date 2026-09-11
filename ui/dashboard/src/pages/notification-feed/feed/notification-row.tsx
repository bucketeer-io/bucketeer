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
    <div className="flex min-h-[130px] items-stretch gap-3 rounded-lg border border-gray-200 dark:border-dark-black-700 p-4 transition-colors hover:border-gray-300 dark:hover:border-dark-black-600">
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
                    'truncate typo-para-medium text-gray-900 dark:text-dark-gray-400',
                    !notification.read && 'font-semibold'
                  )}
                >
                  {notification.title}
                </span>
                <TagList tags={notification.tags} />
                <span className="ml-auto shrink-0 typo-para-tiny text-gray-500 dark:text-dark-gray-200">
                  {formatDateTime(notification.publishedAt)}
                </span>
              </div>
              {notification.content && (
                <p className="line-clamp-2 typo-para-small text-gray-500 dark:text-dark-gray-200">
                  {markdownToText(notification.content)}
                </p>
              )}
            </div>
          }
          footer={
            <span className="typo-para-small text-gray-500 dark:text-dark-gray-200">
              {notification.createdBy.split('@')[0]}
            </span>
          }
        />
      </div>
    </div>
  );
};

export default NotificationRow;
