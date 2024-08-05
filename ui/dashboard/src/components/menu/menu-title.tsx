import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';

const menuTitleVariants = cva([
  'flex h-12 items-center px-3 text-[10px] font-extrabold text-primary-50'
]);

export type MenuTitleProps = {
  text: string;
  className?: string;
};

const MenuTitle = ({ text, className }: MenuTitleProps) => {
  return <title className={cn(menuTitleVariants(), className)}>{text}</title>;
};

export default MenuTitle;
