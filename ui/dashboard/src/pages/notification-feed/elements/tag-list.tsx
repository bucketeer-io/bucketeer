import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { cn } from 'utils/style';
import { Tooltip } from 'components/tooltip';
import { NotificationTag } from '../types';
import TagChip from './tag-chip';

const MAX_VISIBLE_TAGS = 3;

const TagList = ({
  tags,
  className
}: {
  tags: NotificationTag[];
  className?: string;
}) => {
  const visibleTags = tags.slice(0, MAX_VISIBLE_TAGS);
  const hiddenTags = tags.slice(MAX_VISIBLE_TAGS);

  if (!visibleTags.length) return null;

  const hiddenTagsLabel = hiddenTags.map(tag => tag.name).join(', ');

  return (
    <div className={cn('flex shrink-0 items-center gap-1.5', className)}>
      {visibleTags.map(tag => (
        <TagChip key={tag.name} tag={tag} />
      ))}
      {hiddenTags.length > 0 && (
        <Tooltip
          content={hiddenTagsLabel}
          trigger={
            <span
              aria-label={hiddenTagsLabel}
              className="inline-flex items-center gap-0.5 rounded bg-gray-100 dark:bg-dark-black-700 px-2 py-0.5 typo-para-tiny font-medium text-gray-700 dark:text-dark-gray-400"
            >
              <IconMoreHorizOutlined style={{ fontSize: 14 }} />+
              {hiddenTags.length}
            </span>
          }
        />
      )}
    </div>
  );
};

export default TagList;
