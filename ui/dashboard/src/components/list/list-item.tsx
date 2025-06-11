import { FunctionComponent } from 'react';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Icon from 'components/icon';

export type ListItemProps = {
  label: string;
  expanded?: boolean;
  className?: string;
  value: string;
  selected?: boolean;
  icon?: FunctionComponent;
  onSelect?: (value: string) => void;
};

const ListItem = ({
  label,
  expanded,
  value,
  selected,
  icon,
  className,
  onSelect
}: ListItemProps) => {
  return (
    <li
      id={value}
      className={cn(
        'flex items-center justify-between cursor-default',
        'rounded-lg px-3 py-2 text-gray-700 hover:bg-gray-100',
        selected && 'bg-gray-100',
        className
      )}
      onClick={() => onSelect?.(value)}
    >
      <p className="typo-para-medium">{label}</p>
      {(expanded || (selected && icon)) && (
        <Icon icon={icon || IconChevronRight} size="sm" />
      )}
    </li>
  );
};

export default ListItem;
