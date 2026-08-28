import type { FunctionComponent } from 'react';
import { NavLink } from 'react-router';
import { cn } from 'utils/style';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

export type MenuItem = {
  icon?: FunctionComponent;
  label: string;
  actIcon?: FunctionComponent;
  href?: string;
  options?: DropdownOption[];
  loading?: boolean;
  tourId?: string;
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
  tourId,
  onClick,
  onSelectOption,
  isCollapsed
}: MenuItem & { isCollapsed?: boolean }) => {
  const iconEl = icon ? (
    <Icon color="primary-50" size="sm" icon={icon} />
  ) : null;
  const actionIcon = actIcon ? (
    <Icon color="primary-50" size="sm" icon={actIcon} />
  ) : null;

  const textClsx = cn(
    'flex items-center gap-x-2 w-full text-primary-50',
    'py-3 rounded-lg typo-para-medium my-0.5 capitalize',
    'hover:bg-primary-400 hover:opacity-100 opacity-80 sidebar-menu',
    isCollapsed ? 'justify-center px-0' : 'px-3'
  );

  const dropdownTrigger = (
    <div className={textClsx}>
      <div
        className={cn('flex items-center gap-x-2 w-full', {
          'w-fit': isCollapsed
        })}
      >
        {iconEl}
        {!isCollapsed && <div className="w-fit truncate">{label}</div>}
      </div>
      {!loading && !isCollapsed && actionIcon}
    </div>
  );

  const actionEl =
    options && options?.length > 0 ? (
      <Dropdown
        loading={loading}
        trigger={
          isCollapsed ? (
            <Tooltip
              trigger={dropdownTrigger}
              content={label}
              side="right"
              triggerCls="w-full"
            />
          ) : (
            dropdownTrigger
          )
        }
        ariaLabel={isCollapsed ? label : undefined}
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
      <NavLink
        onClick={onClick}
        className={textClsx}
        to={href}
        data-tour={tourId}
        aria-label={isCollapsed ? label : undefined}
      >
        {iconEl}
        {!isCollapsed && label}
      </NavLink>
    ) : (
      <button
        className={cn(textClsx, { 'justify-between': actionIcon })}
        onClick={onClick}
        data-tour={tourId}
        aria-label={isCollapsed ? label : undefined}
      >
        <div className="flex items-center gap-x-2 truncate">
          {iconEl}
          {!isCollapsed && <div className="w-fit truncate">{label}</div>}
        </div>
        {!isCollapsed && actionIcon}
      </button>
    );

  if (isCollapsed && (!options || options.length === 0)) {
    return (
      <Tooltip
        trigger={actionEl}
        content={label}
        side="right"
        triggerCls="w-full"
      />
    );
  }

  return actionEl;
};

export default MenuItemComponent;
