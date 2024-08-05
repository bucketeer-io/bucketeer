import { cn } from 'utils/style';
import MenuItemComponent, { MenuItem } from './menu-item';

export type MenuProps = {
  className?: string;
  title: string;
  items: MenuItem[];
};

const SectionMenu = ({ className, title, items = [] }: MenuProps) => {
  return (
    <div className={cn('flex flex-col', className)}>
      <div className="px-3 uppercase typo-head-bold-tiny text-primary-50 mb-3">
        {title}
      </div>

      {items.map((item, index) => (
        <MenuItemComponent {...item} key={index} />
      ))}
    </div>
  );
};

export default SectionMenu;
