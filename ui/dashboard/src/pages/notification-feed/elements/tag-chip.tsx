import { cn } from 'utils/style';
import { NotificationTag } from '../types';

// Renders a pill for a notification tag: a colored dot (the tag color) followed
// by the tag name on a soft tinted background, matching the Tag select chips.
// Tags without a color fall back to a neutral gray dot and background.
const TagChip = ({
  tag,
  className,
  hideDot = false
}: {
  tag: NotificationTag;
  className?: string;
  hideDot?: boolean;
}) => {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded px-2 py-0.5 typo-para-tiny font-medium',
        !tag.color && 'bg-gray-100 text-gray-700',
        className
      )}
      style={
        tag.color
          ? { color: tag.color, backgroundColor: `${tag.color}1A` }
          : undefined
      }
    >
      {!hideDot && (
        <span
          className={cn('size-1.5 shrink-0 rounded-full', !tag.color && 'bg-gray-400')} // prettier-ignore
          style={tag.color ? { backgroundColor: tag.color } : undefined}
        />
      )}
      {tag.name}
    </span>
  );
};

export default TagChip;
