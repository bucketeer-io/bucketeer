import type { FunctionComponent } from 'react';
import { NavLink } from 'react-router-dom';
import { cn } from 'utils/style';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Icon from 'components/icon';

export type MenuItem = {
  icon: FunctionComponent;
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
  const iconEl = <Icon color="primary-50" size="sm" icon={icon} />;
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
      <DropdownMenu>
        <DropdownMenuTrigger
          loading={loading}
          trigger={
            <div className="flex items-center justify-between w-full">
              <div className="flex items-center w-full gap-x-2 !text-primary-50">
                {iconEl}
                {label}
              </div>
              {!loading && actionIcon}
            </div>
          }
          showArrow={false}
          className="w-full !border-none !shadow-none bg-transparent hover:bg-primary-400 hover:opacity-100 opacity-80 sidebar-menu"
        />
        <DropdownMenuContent side="right" align="start">
          {options?.map((item, index) => (
            <DropdownMenuItem
              key={index}
              label={item.label}
              value={item.value}
              className="[&>div>button]:!cursor-pointer"
              onSelectOption={value => onSelectOption?.(value as string)}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
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
