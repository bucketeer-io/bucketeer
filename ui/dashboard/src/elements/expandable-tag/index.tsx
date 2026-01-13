import { useCallback, useMemo, useState } from 'react';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Icon from 'components/icon';
import TruncationWithTooltip from 'elements/truncation-with-tooltip';

export const Tag = ({
  tagId,
  value,
  className,
  tooltipCls,
  maxSize = 250,
  onTagClick
}: {
  tagId?: string;
  value: string;
  className?: string;
  tooltipCls?: string;
  maxSize?: number;
  onTagClick?: () => void;
}) => (
  <TruncationWithTooltip
    elementId={tagId || ''}
    content={value}
    maxSize={maxSize}
    className={cn('!w-fit !z-10', tooltipCls)}
  >
    <div
      id={tagId}
      className={cn(
        'px-2 py-1.5 bg-primary-100/70 text-primary-500 typo-para-small !w-fit leading-[14px] rounded whitespace-nowrap border border-transparent',
        className
      )}
      onClick={onTagClick}
    >
      {value}
    </div>
  </TruncationWithTooltip>
);

const ExpandableTag = ({
  rowId,
  tags,
  filterTags,
  className,
  wrapperClassName,
  maxSize,
  tooltipCls,
  onTagClick
}: {
  rowId: string;
  tags: string[];
  filterTags?: string[];
  className?: string;
  wrapperClassName?: string;
  maxSize?: number;
  tooltipCls?: string;
  onTagClick?: (tag: string) => void;
}) => {
  const [expandedTags, setExpandedTags] = useState<string[]>([]);

  const isExpanded = useMemo(
    () => expandedTags.includes(rowId),
    [expandedTags, rowId]
  );

  const handleExpandTag = useCallback(() => {
    setExpandedTags(
      isExpanded
        ? expandedTags.filter(item => item !== rowId)
        : [...expandedTags, rowId]
    );
  }, [expandedTags, rowId, isExpanded]);

  return (
    <div
      className={cn(
        'flex flex-wrap items-center w-full gap-x-2',
        {
          'items-start': isExpanded
        },
        wrapperClassName
      )}
    >
      <div className="flex w-full items-center flex-wrap gap-2">
        {(isExpanded ? tags : tags.slice(0, 3))?.map((tag, index) => (
          <Tag
            tagId={`${tag}-${index}`}
            key={index}
            value={tag}
            className={cn(className, {
              'border border-primary-500': filterTags?.includes(tag)
            })}
            maxSize={maxSize}
            tooltipCls={tooltipCls}
            onTagClick={() => onTagClick && onTagClick(tag)}
          />
        ))}
        {tags.length > 3 && !isExpanded && (
          <Tag value={`+${tags.length - 3}`} tooltipCls={tooltipCls} />
        )}
      </div>
      {tags.length > 3 && (
        <div
          className={cn('flex-center cursor-pointer hover:bg-gray-200 rounded')}
          onClick={handleExpandTag}
        >
          <Icon
            icon={IconChevronDown}
            size={'sm'}
            className={cn('flex-center transition-all duration-200 rotate-0', {
              'rotate-180': isExpanded
            })}
          />
        </div>
      )}
    </div>
  );
};
export default ExpandableTag;
