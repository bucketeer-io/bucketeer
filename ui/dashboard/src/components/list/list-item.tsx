import { FunctionComponent } from 'react';
import { cn } from 'utils/style';
import Icon from 'components/icon';

export type ListItemProps = {
  label: string;
  icon?: FunctionComponent;
  className?: string;
  selected?: boolean;
  onClick?: () => void;
};

const ListItem = ({
  label,
  icon,
  selected,
  className,
  onClick
}: ListItemProps) => {
  return (
    <li
      className={cn(
        'flex items-center justify-between cursor-default',
        'rounded-lg bg-white px-3 py-2 text-gray-700',
        selected && 'bg-gray-100',
        className
      )}
      onClick={onClick}
    >
      <p className="typo-para-medium">{label}</p>
      {icon && <Icon icon={icon} size="sm" />}
    </li>
  );
};

export default ListItem;
