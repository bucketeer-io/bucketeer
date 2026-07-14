import { useCallback, useEffect, useMemo, useState } from 'react';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CalendarDays, CheckCheck } from 'lucide-react';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import SearchInput from 'components/search-input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import PageLayout from 'elements/page-layout';
import {
  useFetchDrafts,
  useMarkAllAsRead
} from './collection-loader/use-fetch-notifications';
import DraftsPanel from './drafts/drafts-panel';
import NotificationList from './feed/notification-list';
import NotificationDetailModal from './notification-detail';
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
  environmentId
}: {
  disabled?: boolean;
  environmentId: string;
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

  // Tab counts are reported up by whichever NotificationList is mounted; every
  // feed response carries both totals, so a single active list keeps them fresh.
  const [counts, setCounts] = useState({ unreadCount: 0, readCount: 0 });
  const markAllAsRead = useMarkAllAsRead(environmentId);

  // The notification/draft shown in the detail SlideModal.
  const [detail, setDetail] = useState<NotificationDetail>();

  // Which draft is open for editing in the publish form, tracked by id only.
  // The draft object itself is derived from the live drafts query below, so
  // the form always edits current data instead of a stale snapshot taken when
  // "Edit Draft" was clicked.
  const [editingId, setEditingId] = useState<string>();
  const { data: drafts = [] } = useFetchDrafts(environmentId);
  const editingDraft = drafts.find(d => d.id === editingId);

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

  const dateFilters = useMemo(() => {
    if (!filters.days) return {};
    const to = Date.now();
    return { from: to - filters.days * DAY_MS, to };
  }, [filters.days]);

  const { unreadCount, readCount } = counts;

  const sortOptions = [
    { label: t('sort-by-newest'), value: 'newest' },
    { label: t('sort-by-oldest'), value: 'oldest' }
  ];

  const dateOptions = [
    { label: t('last-7-days'), value: 7 },
    { label: t('last-30-days'), value: 30 },
    { label: t('last-90-days'), value: 90 }
  ];

  return (
    <PageLayout.Content>
      <div className="mb-6 px-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="w-full md:max-w-[440px]">
          <SearchInput
            placeholder={t('form:search-notifications')}
            value={filters.searchQuery}
            onChange={searchQuery => onChangeFilters({ searchQuery })}
          />
        </div>
        <div className="flex items-center gap-3">
          <Dropdown
            className="w-[200px]"
            isTruncate={false}
            value={filters.sort}
            options={sortOptions}
            onChange={value => onChangeFilters({ sort: value as SortOption })}
          />
          <Dropdown
            className="w-[180px]"
            value={filters.days}
            placeholder={
              <span className="flex items-center gap-2">
                <CalendarDays size={16} />
                {t('last-30-days')}
              </span>
            }
            options={dateOptions}
            onChange={days => onChangeFilters({ days: Number(days) })}
          />
        </div>
      </div>

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
          <div className="flex flex-col">
            <div className="flex items-center justify-between">
              <TabsList className="justify-start">
                <TabsTrigger value="unread">
                  {t('unread')} ({unreadCount})
                </TabsTrigger>
                <TabsTrigger value="read">
                  {t('read')} ({readCount})
                </TabsTrigger>
                <TabsTrigger value="publish">
                  {t('publish-notification')}
                </TabsTrigger>
              </TabsList>

              {filters.tab === 'unread' && (
                <Button
                  variant="text"
                  size="sm"
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
                onCounts={setCounts}
                onSelect={setDetail}
              />
            </TabsContent>
            <TabsContent value="read">
              <NotificationList
                filters={{ ...filters, ...dateFilters }}
                environmentId={environmentId}
                onCounts={setCounts}
                onSelect={setDetail}
              />
            </TabsContent>
            <TabsContent value="publish">
              <PublishForm
                disabled={disabled}
                environmentId={environmentId}
                initialDraft={editingDraft}
                onClear={onClearEdit}
              />
            </TabsContent>
          </div>

          {filters.tab === 'publish' && (
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
