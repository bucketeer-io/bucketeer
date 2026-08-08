import { useCallback, useEffect, useState } from 'react';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CheckCheck } from 'lucide-react';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import Button from 'components/button';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Dropdown from 'components/dropdown';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import {
  useFetchDrafts,
  useFetchNotification,
  useFetchTabCounts,
  useMarkAllAsRead
} from './collection-loader/use-fetch-notifications';
import DraftsPanel from './drafts/drafts-panel';
import NotificationDetailModal from './elements/notification-detail';
import NotificationList from './feed/notification-list';
import PublishForm from './publisher/publish-form';
import {
  NotificationDetail,
  NotificationFilters,
  NotificationTab,
  SortOption
} from './types';

const PageContent = ({
  disabled,
  isSystemAdmin
}: {
  disabled?: boolean;
  isSystemAdmin?: boolean;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<NotificationFilters> = searchOptions;

  const defaultFilters = {
    tab: 'unread',
    searchQuery: '',
    sort: 'newest',
    ...searchFilters
  } as NotificationFilters;

  const [filters, setFilters] =
    usePartialState<NotificationFilters>(defaultFilters);

  // Sourced independently of whichever NotificationList is mounted, so both
  // tab badges stay live across tab switches and mutations.
  const { unreadCount, readCount } = useFetchTabCounts();
  const markAllAsRead = useMarkAllAsRead();

  const selectedId = filters.notificationId;
  const { data: detail } = useFetchNotification(selectedId);

  const onSelectDetail = (notification: NotificationDetail) =>
    onChangeFilters({ notificationId: notification.id });

  const [editingId, setEditingId] = useState<string>();
  const { data: draftsData } = useFetchDrafts(
    !!isSystemAdmin && (filters.tab === 'publish' || !!editingId)
  );
  const editingDraft = draftsData?.notifications.find(d => d.id === editingId);

  const onChangeFilters = useCallback(
    (values: Partial<NotificationFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  // Resets search + date range only, keeping the current tab and sort order.
  const onClearFilters = useCallback(
    () => onChangeFilters({ searchQuery: '', from: undefined, to: undefined }),
    [onChangeFilters]
  );

  const onEditDraft = (draft: NotificationDetail) => {
    setEditingId(draft.id);
    onChangeFilters({ notificationId: undefined, tab: 'publish' });
  };

  // Leaving edit mode (Clear or a successful submit): drop the draft so the
  // form reverts to "new notification" mode.
  const onClearEdit = () => setEditingId(undefined);

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  useEffect(() => {
    if (!isSystemAdmin && filters.tab === 'publish') {
      onChangeFilters({ tab: 'unread' });
    }
  }, [isSystemAdmin, filters.tab]);

  const onCloseDetail = () => onChangeFilters({ notificationId: undefined });

  const sortOptions = [
    { label: t('sort-by-newest'), value: 'newest' },
    { label: t('sort-by-oldest'), value: 'oldest' }
  ];

  return (
    <PageLayout.Content>
      <Filter
        className="mb-6"
        isShowDocumentation={false}
        placeholder={t('form:search-notifications')}
        name="notifications-search"
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
        action={
          <>
            <Dropdown
              className="w-[200px]"
              wrapTriggerStyle="w-fit"
              isTruncate={false}
              value={filters.sort}
              options={sortOptions}
              onChange={value => onChangeFilters({ sort: value as SortOption })}
            />
            <ReactDateRangePicker
              from={filters.from}
              to={filters.to}
              isAllTime={!filters.from && !filters.to}
              onChange={(startDate, endDate) => {
                onChangeFilters({
                  from: startDate ? startDate.toString() : undefined,
                  to: endDate ? endDate.toString() : undefined
                });
              }}
            />
          </>
        }
      />

      <Tabs
        value={filters.tab}
        onValueChange={v => onChangeFilters({ tab: v as NotificationTab })}
        className="flex flex-1 flex-col"
      >
        <div
          className={cn('p-6 grid grid-cols-1 gap-8', {
            'lg:grid-cols-[1fr_360px]': filters.tab === 'publish'
          })}
        >
          <div className="flex flex-col relative">
            <div className="flex items-center justify-between">
              <TabsList className="justify-start">
                <TabsTrigger value="unread">
                  {t('unread')} ({unreadCount})
                </TabsTrigger>
                <TabsTrigger value="read">
                  {t('read')} ({readCount})
                </TabsTrigger>
                {isSystemAdmin && (
                  <TabsTrigger value="publish">
                    {t('publish-notification')}
                  </TabsTrigger>
                )}
              </TabsList>

              {filters.tab === 'unread' && (
                <Button
                  variant="text"
                  size="sm"
                  className="absolute right-0 top-0"
                  onClick={() => markAllAsRead.mutate()}
                  disabled={unreadCount === 0}
                  loading={markAllAsRead.isPending}
                >
                  <CheckCheck size={16} />
                  {t('mark-all-as-read')}
                </Button>
              )}
            </div>

            <TabsContent value="unread">
              <NotificationList
                read={false}
                filters={filters}
                onSelect={onSelectDetail}
                onClearFilters={onClearFilters}
              />
            </TabsContent>
            <TabsContent value="read">
              <NotificationList
                read
                filters={filters}
                onSelect={onSelectDetail}
                onClearFilters={onClearFilters}
              />
            </TabsContent>
            {isSystemAdmin && (
              <TabsContent value="publish">
                <PublishForm
                  disabled={disabled}
                  initialDraft={editingDraft}
                  onClear={onClearEdit}
                />
              </TabsContent>
            )}
          </div>

          {isSystemAdmin && filters.tab === 'publish' && (
            <aside className="lg:border-l lg:border-gray-200 lg:pl-8">
              <DraftsPanel
                filters={filters}
                onSelect={onSelectDetail}
                onClearFilters={onClearFilters}
              />
            </aside>
          )}
        </div>
      </Tabs>

      <NotificationDetailModal
        notification={detail}
        isOpen={!!selectedId}
        onClose={onCloseDetail}
        onEditDraft={onEditDraft}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
