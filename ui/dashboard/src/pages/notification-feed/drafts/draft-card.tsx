import { useFormatDateTime } from 'utils/date-time';
import NotificationCard from '../elements/notification-card';
import TagChip from '../elements/tag-chip';
import { NotificationDraft } from '../types';

interface DraftCardProps {
  draft: NotificationDraft;
  active?: boolean;
  onClick?: () => void;
}

const DraftCard = ({ draft, active, onClick }: DraftCardProps) => {
  const formatDateTime = useFormatDateTime();
  return (
    <NotificationCard
      active={active}
      onClick={onClick}
      header={
        <span className="typo-para-medium font-medium text-gray-900">
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
  );
};

export default DraftCard;
