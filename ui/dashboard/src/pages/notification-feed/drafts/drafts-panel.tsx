import { useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import Pagination from 'components/pagination';
import Spinner from 'components/spinner';
import { NoResultsCollection } from 'elements/collection/collection-empty';
import ConfirmModal from 'elements/confirm-modal';
import EmptyState from 'elements/empty-state';
import {
  useDeleteNotification,
  useFetchDrafts
} from '../collection-loader/use-fetch-notifications';
import { DRAFTS_PAGE_SIZE } from '../constants';
import { markdownToText } from '../elements/markdown-content';
import { NotificationDraft, NotificationFilters } from '../types';
import DraftCard from './draft-card';

interface DraftsPanelProps {
  environmentId: string;
  filters?: Partial<NotificationFilters>;
  onSelect?: (draft: NotificationDraft) => void;
  onClearFilters?: () => void;
}

const DraftsPanel = ({
  environmentId,
  filters,
  onSelect,
  onClearFilters
}: DraftsPanelProps) => {
  const { t } = useTranslation(['common', 'message']);
  const { notify, errorNotify } = useToast();
  const { data, isLoading } = useFetchDrafts(environmentId);
  const drafts = data?.notifications ?? [];
  const [page, setPage] = useState(1);
  const [activeId, setActiveId] = useState<string>();
  const [deletingDraft, setDeletingDraft] = useState<NotificationDraft>();
  const deleteMutation = useDeleteNotification(environmentId);

  const onCloseDeleteModal = () => {
    setDeletingDraft(undefined);
    deleteMutation.reset();
  };

  const onConfirmDelete = () => {
    if (!deletingDraft) return;
    deleteMutation.mutate(deletingDraft.id, {
      onSuccess: () => {
        onCloseDeleteModal();
        notify({
          message: t('message:collection-action-success', {
            collection: t('common:notification'),
            action: t('common:deleted')
          })
        });
      },
      onError: error => errorNotify(error)
    });
  };

  const filtered = useMemo(() => {
    const query = (filters?.searchQuery ?? '').trim().toLowerCase();
    return drafts
      .filter(d =>
        filters?.from ? Number(d.updatedAt) * 1000 >= filters.from : true
      )
      .filter(d =>
        filters?.to ? Number(d.updatedAt) * 1000 <= filters.to : true
      )
      .filter(d =>
        query
          ? d.title.toLowerCase().includes(query) ||
            markdownToText(d.content).toLowerCase().includes(query)
          : true
      )
      .sort((a, b) =>
        filters?.sort === 'oldest'
          ? Number(a.updatedAt) - Number(b.updatedAt)
          : Number(b.updatedAt) - Number(a.updatedAt)
      );
  }, [drafts, filters]);

  // Reset to the first page whenever the filtered result set changes so the
  // current page never points past the available drafts.
  useEffect(() => setPage(1), [filters]);

  const paged = useMemo(() => {
    const start = (page - 1) * DRAFTS_PAGE_SIZE;
    return filtered.slice(start, start + DRAFTS_PAGE_SIZE);
  }, [filtered, page]);

  return (
    <div className="flex flex-col gap-4">
      <h2 className="typo-head-bold-small text-gray-900">
        {t('drafts')} ({filtered.length})
      </h2>

      {isLoading ? (
        <div className="flex justify-center py-10">
          <Spinner />
        </div>
      ) : filtered.length === 0 ? (
        <div className="py-10">
          {drafts.length > 0 ? (
            <NoResultsCollection onClear={onClearFilters} />
          ) : (
            <EmptyState.Root variant="no-data" size="sm">
              <EmptyState.Illustration />
              <EmptyState.Body>
                <EmptyState.Title>{t('no-drafts')}</EmptyState.Title>
                <EmptyState.Description>
                  {t('no-drafts-desc')}
                </EmptyState.Description>
              </EmptyState.Body>
            </EmptyState.Root>
          )}
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {paged.map(draft => (
            <DraftCard
              key={draft.id}
              draft={draft}
              active={draft.id === activeId}
              onClick={() => {
                setActiveId(draft.id);
                onSelect?.(draft);
              }}
              onDelete={() => setDeletingDraft(draft)}
            />
          ))}
        </div>
      )}

      <Pagination
        page={page}
        pageSize={DRAFTS_PAGE_SIZE}
        totalCount={filtered.length}
        onChange={setPage}
      />

      {deletingDraft && (
        <ConfirmModal
          isOpen={!!deletingDraft}
          onClose={onCloseDeleteModal}
          onSubmit={onConfirmDelete}
          title={t('delete-draft-title')}
          description={
            <Trans
              t={t}
              i18nKey="delete-draft-desc"
              values={{ name: deletingDraft.title }}
              components={{ bold: <strong /> }}
            />
          }
          loading={deleteMutation.isPending}
        />
      )}
    </div>
  );
};

export default DraftsPanel;
