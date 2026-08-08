import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Icon from 'components/icon';
import NotificationCard from '../elements/notification-card';
import TagChip from '../elements/tag-chip';
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
          <span
            className={cn(
              'typo-para-medium font-medium text-gray-900',
              onDelete && 'pr-8'
            )}
          >
            {draft.title}
          </span>
        }
        footer={
          <span className="typo-para-small text-gray-500">
            {draft.createdBy.split('@')[0]}
          </span>
        }
      >
        {draft.tags.map(tag => (
          <TagChip key={tag.name} tag={tag} />
        ))}
        <span className="typo-para-tiny text-gray-500">
          {formatDateTime(draft.updatedAt)}
        </span>
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
