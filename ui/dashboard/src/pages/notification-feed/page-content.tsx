import { useCallback, useEffect, useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CheckCheck } from 'lucide-react';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconThreeLines } from '@icons';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import {
  useFetchDrafts,
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

const DAY_MS = 24 * 60 * 60 * 1000;

const PageContent = ({
  disabled,
  isSystemAdmin,
  environmentId
}: {
  disabled?: boolean;
  isSystemAdmin?: boolean;
  environmentId: string;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<NotificationFilters> = searchOptions;
  const location = useLocation();
  const navigate = useNavigate();

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
  const { unreadCount, readCount } = useFetchTabCounts(environmentId);
  const markAllAsRead = useMarkAllAsRead(environmentId);

  // The notification/draft shown in the detail SlideModal.
  const [detail, setDetail] = useState<NotificationDetail>();

  // Which draft is open for editing in the publish form, tracked by id only.
  // The draft object itself is derived from the live drafts query below, so
  // the form always edits current data instead of a stale snapshot taken when
  // "Edit Draft" was clicked.
  const [editingId, setEditingId] = useState<string>();
  const { data: draftsData } = useFetchDrafts(environmentId, isSystemAdmin);
  const editingDraft = draftsData?.notifications.find(d => d.id === editingId);

  const onEditDraft = (draft: NotificationDetail) => {
    setDetail(undefined);
    setEditingId(draft.id);
    onChangeFilters({ tab: 'publish' });
  };

  // Leaving edit mode (Clear or a successful submit): drop the draft so the
  // form reverts to "new notification" mode.
  const onClearEdit = () => setEditingId(undefined);

  const onChangeFilters = useCallback(
    (values: Partial<NotificationFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  // Publishing is system admin only: bounce anyone else off the "publish"
  // tab (e.g. a stale bookmark or shared link with ?tab=publish) so the tab
  // bar and content never disagree.
  useEffect(() => {
    if (!isSystemAdmin && filters.tab === 'publish') {
      onChangeFilters({ tab: 'unread' });
    }
  }, [isSystemAdmin, filters.tab]);

  // Opened via NotificationBell: the clicked notification is handed off
  // through router state so its detail opens here without an extra fetch.
  // Cleared right after so a refresh or back-navigation doesn't reopen it.
  useEffect(() => {
    const openNotification = (
      location.state as { notification?: NotificationDetail } | null
    )?.notification;
    if (openNotification) {
      setDetail(openNotification);
      navigate(location.pathname + location.search, {
        replace: true,
        state: null
      });
    }
  }, [location.state]);

  const dateFilters = useMemo(() => {
    if (!filters.days) return {};
    const to = Date.now();
    return { from: to - filters.days * DAY_MS, to };
  }, [filters.days]);

  const sortOptions = [
    { label: t('sort-by-newest'), value: 'newest' },
    { label: t('sort-by-oldest'), value: 'oldest' }
  ];

  const dateOptions = [
    { label: t('last-7-days'), value: 7 },
    { label: t('last-30-days'), value: 30 },
    { label: t('last-90-days'), value: 90 }
  ];

  const dateLabel =
    dateOptions.find(item => item.value === filters.days)?.label ||
    t('last-30-days');

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
            <Dropdown
              trigger={
                <div className="flex items-center gap-x-2">
                  <Icon icon={IconThreeLines} size="sm" />
                  <p className="text-gray-600">{dateLabel}</p>
                </div>
              }
              value={filters.days}
              options={dateOptions}
              showArrow={false}
              alignContent="end"
              className="w-full px-4 py-[11px] justify-center"
              wrapTriggerStyle="w-fit"
              onChange={days => onChangeFilters({ days: Number(days) })}
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
                filters={{ ...filters, ...dateFilters }}
                environmentId={environmentId}
                onSelect={setDetail}
              />
            </TabsContent>
            <TabsContent value="read">
              <NotificationList
                filters={{ ...filters, ...dateFilters }}
                environmentId={environmentId}
                onSelect={setDetail}
              />
            </TabsContent>
            {isSystemAdmin && (
              <TabsContent value="publish">
                <PublishForm
                  disabled={disabled}
                  environmentId={environmentId}
                  initialDraft={editingDraft}
                  onClear={onClearEdit}
                />
              </TabsContent>
            )}
          </div>

          {isSystemAdmin && filters.tab === 'publish' && (
            <aside className="lg:border-l lg:border-gray-200 lg:pl-8">
              <DraftsPanel
                environmentId={environmentId}
                filters={{ ...filters, ...dateFilters }}
                onSelect={setDetail}
              />
            </aside>
          )}
        </div>
      </Tabs>

      <NotificationDetailModal
        notification={detail}
        isOpen={!!detail}
        onClose={() => setDetail(undefined)}
        onEditDraft={onEditDraft}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
