import { cn } from 'utils/style';
import MenuItemComponent, { MenuItem } from './menu-item';

export type MenuProps = {
  className?: string;
  title: string;
  items: MenuItem[];
  onClickNavLink?: () => void;
};

const SectionMenu = ({
  className,
  title,
  items = [],
  onClickNavLink
}: MenuProps) => {
  return (
    <div className={cn('flex flex-col', className)}>
      <div className="px-3 uppercase typo-head-bold-tiny text-primary-50 mb-3 opacity-70">
        {title}
      </div>

      {items.map((item, index) => (
        <MenuItemComponent {...item} key={index} onClick={onClickNavLink} />
      ))}
    </div>
  );
};

export default SectionMenu;
