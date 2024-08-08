import { FunctionComponent } from 'react';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';

export type ListItemProps = {
  text: string;
  icon?: FunctionComponent;
  selected?: boolean;
  type?: 'text' | 'icon';
  onClick?: () => void;
};

const ListItem = ({
  text,
  icon: SvgIcon,
  selected,
  type,
  onClick
}: ListItemProps) => {
  return (
    <li
      className={cn(
        'flex h-10 min-w-[200px] cursor-pointer items-center justify-between',
        'rounded-lg bg-white px-3 py-2 text-gray-700',
        selected && 'bg-gray-100'
      )}
      onClick={onClick}
    >
      <p className="typo-para-medium">{text}</p>
      {type === 'icon' && (SvgIcon ? <SvgIcon /> : <IconChevronRight />)}
    </li>
  );
};

export default ListItem;
