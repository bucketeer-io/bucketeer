import { useEffect, useState } from 'react';
import { useTranslation } from 'i18n';
import { Mail } from 'lucide-react';
import Button from 'components/button';
import Pagination from 'components/pagination';
import Spinner from 'components/spinner';
import {
  useFetchFeed,
  useMarkAsRead,
  useMarkManyAsRead
} from '../collection-loader/use-fetch-notifications';
import { FeedNotification, NotificationFilters } from '../types';
import NotificationRow from './notification-row';

const PAGE_SIZE = 10;

interface NotificationListProps {
  read?: boolean;
  filters: NotificationFilters;
  environmentId: string;
  onCounts?: (counts: { unreadCount: number; readCount: number }) => void;
  onSelect?: (notification: FeedNotification) => void;
}

const NotificationList = ({
  read = true,
  filters,
  environmentId,
  onCounts,
  onSelect
}: NotificationListProps) => {
  const { t } = useTranslation(['common']);
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<Set<string>>(new Set());

  // Reset paging and selection whenever the tab or filters change so we never
  // land on a now-empty page or carry a stale selection across contexts.
  useEffect(() => {
    setPage(1);
    setSelected(new Set());
  }, [read, filters.searchQuery, filters.sort, filters.from, filters.to]);

  const { data, isLoading } = useFetchFeed(environmentId, read, page, filters);
  const markAsRead = useMarkAsRead(environmentId);
  const markManyAsRead = useMarkManyAsRead(environmentId);

  const items = data?.items ?? [];
  const total = data?.total ?? 0;

  // Surface both tab counts to the parent for the Unread/Read labels.
  useEffect(() => {
    if (data) {
      onCounts?.({
        unreadCount: data.unreadCount,
        readCount: data.readCount
      });
    }
  }, [data, onCounts]);

  const toggle = (id: string, isSelected: boolean) =>
    setSelected(prev => {
      const next = new Set(prev);
      if (isSelected) next.add(id);
      else next.delete(id);
      return next;
    });

  const clearSelection = () => setSelected(new Set());

  const markSelectedAsRead = () =>
    markManyAsRead.mutate([...selected], {
      onSuccess: () => clearSelection()
    });

  if (isLoading) {
    return (
      <div className="flex justify-center py-10">
        <Spinner />
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <p className="py-10 text-center typo-para-medium text-gray-500">
        {t('no-notifications')}
      </p>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-3">
        {items.map(notification => (
          <NotificationRow
            key={notification.id}
            notification={notification}
            selectable={!read}
            selected={selected.has(notification.id)}
            onSelectedChange={value => toggle(notification.id, value)}
            onClick={() => {
              if (!notification.read) markAsRead.mutate(notification.id);
              onSelect?.(notification);
            }}
          />
        ))}
      </div>

      {/* Selection action bar — unread tab only, shown once a row is selected. */}
      {!read && selected.size > 0 && (
        <div className="flex items-center justify-between rounded-lg border border-gray-200 bg-gray-50 px-4 py-3">
          <div className="flex items-center gap-3">
            <span className="typo-para-medium text-gray-700">
              {t('selected', { count: selected.size })}
            </span>
            <Button variant="text" size="sm" onClick={clearSelection}>
              {t('clear-selection')}
            </Button>
          </div>
          <Button
            variant="secondary-2"
            size="sm"
            onClick={markSelectedAsRead}
            loading={markManyAsRead.isPending}
          >
            <Mail size={16} />
            {t('mark-selected-as-read')}
          </Button>
        </div>
      )}

      <Pagination
        page={page}
        pageSize={PAGE_SIZE}
        totalCount={total}
        onChange={setPage}
      />
    </div>
  );
};

export default NotificationList;
