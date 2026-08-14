import { useEffect, useState } from 'react';
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
  useFetchDraftsPage
} from '../collection-loader/use-fetch-notifications';
import { DRAFTS_PAGE_SIZE } from '../constants';
import { NotificationDraft, NotificationFilters } from '../types';
import DraftCard from './draft-card';

interface DraftsPanelProps {
  filters?: Partial<NotificationFilters>;
  onSelect?: (draft: NotificationDraft) => void;
  onClearFilters?: () => void;
}

const DraftsPanel = ({
  filters,
  onSelect,
  onClearFilters
}: DraftsPanelProps) => {
  const { t } = useTranslation(['common', 'message']);
  const { notify, errorNotify } = useToast();
  const [page, setPage] = useState(1);
  const searchQuery = (filters?.searchQuery ?? '').trim();
  const { data, isLoading } = useFetchDraftsPage(
    page,
    searchQuery,
    DRAFTS_PAGE_SIZE,
    filters?.sort ?? 'newest'
  );
  const drafts = data?.notifications ?? [];
  const totalCount = Number(data?.totalCount ?? 0);
  const [activeId, setActiveId] = useState<string>();
  const [deletingDraft, setDeletingDraft] = useState<NotificationDraft>();
  const deleteMutation = useDeleteNotification();

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

  useEffect(() => setPage(1), [searchQuery, filters?.sort]);

  return (
    <div className="flex flex-col gap-4">
      <h2 className="typo-head-bold-small text-gray-900">
        {t('drafts')} ({totalCount})
      </h2>

      {isLoading ? (
        <div className="flex justify-center py-10">
          <Spinner />
        </div>
      ) : drafts.length === 0 ? (
        <div className="py-10">
          {searchQuery ? (
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
          {drafts.map(draft => (
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
        totalCount={totalCount}
        onChange={setPage}
        className="flex-col gap-y-2"
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
          submitText={t('delete')}
          loading={deleteMutation.isPending}
        />
      )}
    </div>
  );
};

export default DraftsPanel;
