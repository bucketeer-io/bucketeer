import type { FunctionComponent } from 'react';
import { NavLink } from 'react-router-dom';
import { cn } from 'utils/style';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';

export type MenuItem = {
  icon?: FunctionComponent;
  label: string;
  actIcon?: FunctionComponent;
  href?: string;
  options?: DropdownOption[];
  loading?: boolean;
  onClick?: () => void;
  onSelectOption?: (value: string) => void;
};

const MenuItemComponent = ({
  icon,
  label,
  href,
  actIcon,
  options,
  loading,
  onClick,
  onSelectOption
}: MenuItem) => {
  const iconEl = icon ? (
    <Icon color="primary-50" size="sm" icon={icon} />
  ) : null;
  const actionIcon = actIcon ? (
    <Icon color="primary-50" size="sm" icon={actIcon} />
  ) : null;

  const textClsx = cn(
    'flex items-center gap-x-2 w-full text-primary-50',
    'px-3 py-3 rounded-lg typo-para-medium my-0.5 capitalize',
    'hover:bg-primary-400 hover:opacity-100 opacity-80 sidebar-menu'
  );

  const actionEl =
    options && options?.length > 0 ? (
      <Dropdown
        loading={loading}
        trigger={
          <div className={textClsx}>
            <div className="flex items-center gap-x-2 w-full">
              {iconEl}
              <div className="w-fit truncate">{label}</div>
            </div>
            {!loading && actionIcon}
          </div>
        }
        showArrow={false}
        options={options.map(item => ({
          ...item,
          iconElement: item?.icon ? (
            <div className="flex-center size-fit mt-0.5">
              <Icon size="sm" icon={item?.icon} />
            </div>
          ) : null
        }))}
        onChange={value => onSelectOption?.(value as string)}
        className="w-full !p-0 !border-none !shadow-none [&>div>div>div>div]:text-primary-50 bg-transparent hover:bg-primary-400 hover:opacity-100  sidebar-menu"
        menuContentSide="right"
      />
    ) : href ? (
      <NavLink onClick={onClick} className={textClsx} to={href}>
        {iconEl}
        {label}
      </NavLink>
    ) : (
      <button
        className={cn(textClsx, { 'justify-between': actionIcon })}
        onClick={onClick}
      >
        <div className="flex items-center gap-x-2 truncate">
          {iconEl}
          <div className="w-fit truncate">{label}</div>
        </div>
        {actionIcon}
      </button>
    );

  return actionEl;
};

export default MenuItemComponent;
