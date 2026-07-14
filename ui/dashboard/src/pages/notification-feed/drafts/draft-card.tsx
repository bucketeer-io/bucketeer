import { format } from 'timeago.js';
import { cn } from 'utils/style';
import Button from 'components/button';
import TagChip from '../tag-chip';
import { NotificationDraft } from '../types';

interface DraftCardProps {
  draft: NotificationDraft;
  active?: boolean;
  onClick?: () => void;
}

const DraftCard = ({ draft, active, onClick }: DraftCardProps) => {
  return (
    <Button
      type="button"
      variant="text"
      onClick={onClick}
      className={cn(
        // Override the Button base (centered, nowrap, fixed height) so it reads
        // as a full-width, left-aligned card.
        'flex h-auto w-full flex-col items-start justify-start gap-2 whitespace-normal rounded-lg border p-4 text-left transition-colors',
        active
          ? 'border-primary-500 shadow-border-primary-500'
          : 'border-gray-200 hover:border-gray-300'
      )}
    >
      <span className="typo-para-medium font-medium text-gray-900">
        {draft.title}
      </span>
      <div className="flex flex-wrap items-center gap-2">
        {draft.tags.map(tag => (
          <TagChip key={tag.name} tag={tag} />
        ))}
        <span className="typo-para-tiny text-gray-500">
          {format(draft.updatedAt)}
        </span>
      </div>
      <span className="typo-para-small text-gray-500">
        {draft.createdBy.split('@')[0]}
      </span>
    </Button>
  );
};

export default DraftCard;
