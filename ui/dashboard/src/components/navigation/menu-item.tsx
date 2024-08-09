import type { FunctionComponent } from 'react';
import { Link } from 'react-router-dom';
import type { Route } from '@types';
import { cn } from 'utils/style';
import Icon from 'components/icon';

export type MenuItem = {
  icon: FunctionComponent;
  label: string;
  actIcon?: FunctionComponent;
} & (
  | {
      href: Route;
      onClick?: never;
    }
  | {
      href?: never;
      onClick: () => void;
    }
);

const MenuItemComponent = ({
  icon,
  label,
  href,
  actIcon,
  onClick
}: MenuItem) => {
  const iconEl = <Icon color="primary-50" size="sm" icon={icon} />;
  const textClsx = cn(
    'flex items-center gap-x-2 w-full text-primary-50',
    'px-3 py-3 rounded-lg typo-para-medium',
    'hover:bg-primary-400 hover:opacity-80'
  );

  const actionEl = href ? (
    <Link className={textClsx} to={href}>
      {iconEl}
      {label}
    </Link>
  ) : (
    <button
      className={cn(textClsx, { 'justify-between': actIcon })}
      onClick={onClick}
    >
      <div className="flex items-center gap-x-2">
        {iconEl}
        {label}
      </div>
      {actIcon && <Icon color="primary-50" icon={actIcon} />}
    </button>
  );

  return actionEl;
};

export default MenuItemComponent;
