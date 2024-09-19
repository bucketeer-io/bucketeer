import { FunctionComponent } from 'react';
import { cn } from 'utils/style';
import { Badge } from 'components/badge';

export type TabItemType = 'icon';

export type TabItemProps = {
  text: string;
  icon?: FunctionComponent;
  badge?: number;
  selected?: boolean;
  value: string;
  onClick?: (value: string) => void;
};

const TabItem = ({
  text,
  icon: SvgIcon,
  badge,
  selected,
  value,
  onClick
}: TabItemProps) => {
  return (
    <li
      className={cn(
        'typo-para-medium text-gray-500 h-10 min-w-[62px] py-1 px-4 flex items-center gap-1 hover:cursor-pointer border-b-2 border-transparent',
        selected && 'text-primary-500 border-primary-500'
      )}
      onClick={() => onClick && onClick(value)}
    >
      {SvgIcon && <SvgIcon />}
      {text}
      {badge && (
        <Badge variant={selected ? 'primary' : 'secondary'}>{badge}</Badge>
      )}
    </li>
  );
};

export default TabItem;
