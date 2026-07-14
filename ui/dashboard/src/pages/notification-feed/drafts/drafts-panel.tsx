import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'i18n';
import Pagination from 'components/pagination';
import Spinner from 'components/spinner';
import { useFetchDrafts } from '../collection-loader/use-fetch-notifications';
import { DRAFTS_PAGE_SIZE } from '../constants';
import { markdownToText } from '../markdown-content';
import { NotificationDraft, NotificationFilters } from '../types';
import DraftCard from './draft-card';

interface DraftsPanelProps {
  environmentId: string;
  filters?: Partial<NotificationFilters>;
  onSelect?: (draft: NotificationDraft) => void;
}

const DraftsPanel = ({
  environmentId,
  filters,
  onSelect
}: DraftsPanelProps) => {
  const { t } = useTranslation(['common']);
  const { data: drafts = [], isLoading } = useFetchDrafts(environmentId);
  const [page, setPage] = useState(1);
  const [activeId, setActiveId] = useState<string>();

  const filtered = useMemo(() => {
    const query = (filters?.searchQuery ?? '').trim().toLowerCase();
    return drafts
      .filter(d => (filters?.from ? d.updatedAt >= filters.from : true))
      .filter(d => (filters?.to ? d.updatedAt <= filters.to : true))
      .filter(d =>
        query
          ? d.title.toLowerCase().includes(query) ||
            markdownToText(d.content).toLowerCase().includes(query)
          : true
      )
      .sort((a, b) =>
        filters?.sort === 'oldest'
          ? a.updatedAt - b.updatedAt
          : b.updatedAt - a.updatedAt
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
        <p className="py-10 text-center typo-para-medium text-gray-500">
          {t('no-drafts')}
        </p>
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
    </div>
  );
};

export default DraftsPanel;
