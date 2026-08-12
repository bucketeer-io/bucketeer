import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Icon from 'components/icon';
import { markdownToText } from '../elements/markdown-content';
import NotificationCard from '../elements/notification-card';
import TagList from '../elements/tag-list';
import { NotificationDraft } from '../types';

interface DraftCardProps {
  draft: NotificationDraft;
  active?: boolean;
  onClick?: () => void;
  onDelete?: () => void;
}

const DraftCard = ({ draft, active, onClick, onDelete }: DraftCardProps) => {
  const { t } = useTranslation(['common']);
  const formatDateTime = useFormatDateTime();
  return (
    <div className="relative">
      <NotificationCard
        active={active}
        onClick={onClick}
        header={
          <div className="flex w-full flex-col gap-1">
            <span
              className={cn(
                'truncate typo-para-medium font-medium text-gray-900',
                onDelete && 'pr-8'
              )}
            >
              {draft.title}
            </span>
            {draft.content && (
              <p className="line-clamp-2 typo-para-small text-gray-500">
                {markdownToText(draft.content)}
              </p>
            )}
          </div>
        }
        footer={
          <div className="flex w-full items-center justify-between">
            <span className="typo-para-small text-gray-500">
              {draft.createdBy.split('@')[0]}
            </span>
            <span className="typo-para-tiny text-gray-500">
              {formatDateTime(draft.updatedAt)}
            </span>
          </div>
        }
      >
        <TagList tags={draft.tags} />
      </NotificationCard>
      {onDelete && (
        <button
          type="button"
          aria-label={t('delete')}
          className="absolute right-4 top-4 shrink-0 text-gray-500 hover:text-accent-red-500"
          onClick={e => {
            e.stopPropagation();
            onDelete();
          }}
        >
          <Icon icon={IconTrash} size="sm" />
        </button>
      )}
    </div>
  );
};

export default DraftCard;
