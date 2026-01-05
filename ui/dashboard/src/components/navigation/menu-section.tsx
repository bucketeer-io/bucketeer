import { cn } from 'utils/style';
import MenuItemComponent, { MenuItem } from './menu-item';

export type MenuProps = {
  isExpanded?: boolean;
  className?: string;
  title: string;
  items: MenuItem[];
  onClickNavLink?: () => void;
};

const SectionMenu = ({
  isExpanded = true,
  className,
  title,
  items = [],
  onClickNavLink
}: MenuProps) => {
  return (
    <div className={cn('flex flex-col', className)}>
      <div
        className={cn(
          'px-3 uppercase typo-head-bold-tiny text-primary-50 mb-3 opacity-70',
          isExpanded ? 'block' : 'hidden md:block'
        )}
      >
        {title}
      </div>

      {items.map((item, index) => (
        <MenuItemComponent
          {...item}
          key={index}
          onClick={onClickNavLink}
          isExpanded={isExpanded}
        />
      ))}
    </div>
  );
};

export default SectionMenu;
