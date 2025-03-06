import type { FunctionComponent } from 'react';
import { NavLink } from 'react-router-dom';
import { cn } from 'utils/style';
import Icon from 'components/icon';

export type MenuItem = {
  icon: FunctionComponent;
  label: string;
  actIcon?: FunctionComponent;
  href?: string;
  onClick?: () => void;
};

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
    'px-3 py-3 rounded-lg typo-para-medium my-0.5 capitalize',
    'hover:bg-primary-400 hover:opacity-100 opacity-80 sidebar-menu'
  );

  const actionEl = href ? (
    <NavLink onClick={onClick} className={textClsx} to={href}>
      {iconEl}
      {label}
    </NavLink>
  ) : (
    <button
      className={cn(textClsx, { 'justify-between': actIcon })}
      onClick={onClick}
    >
      <div className="flex items-center gap-x-2 truncate">
        {iconEl}
        <div className="w-fit truncate">{label}</div>
      </div>
      {actIcon && <Icon color="primary-50" size="sm" icon={actIcon} />}
    </button>
  );

  return actionEl;
};

export default MenuItemComponent;
