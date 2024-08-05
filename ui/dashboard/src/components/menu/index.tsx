import MenuGroup from './menu-group';
import MenuItem, { MenuItemProps } from './menu-item';
import MenuTitle from './menu-title';

export type MenuProps = {
  className?: string;
  title?: string;
  options?: MenuItemProps[];
};

const Menu = ({ className, title, options = [] }: MenuProps) => {
  return (
    <div className={className}>
      {title && <MenuTitle text={title} />}
      <MenuGroup>
        {options.map((item, index) => (
          <MenuItem key={index} {...item}>
            {item?.children}
          </MenuItem>
        ))}
      </MenuGroup>
    </div>
  );
};

Menu.Group = MenuGroup;
Menu.Item = MenuItem;
Menu.Title = MenuTitle;

export default Menu;
