import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Icon from 'components/icon';

export type ListItemProps = {
  label: string;
  expanded?: boolean;
  className?: string;
  value: string;
  selected?: boolean;
  onSelect?: (value: string) => void;
};

const ListItem = ({
  label,
  expanded,
  value,
  selected,
  className,
  onSelect
}: ListItemProps) => {
  return (
    <li
      className={cn(
        'flex items-center justify-between cursor-default',
        'rounded-lg bg-white px-3 py-2 text-gray-700',
        selected && 'bg-gray-100',
        className
      )}
      onClick={() => onSelect?.(value)}
    >
      <p className="typo-para-medium">{label}</p>
      {expanded && <Icon icon={IconChevronRight} size="sm" />}
    </li>
  );
};

export default ListItem;
