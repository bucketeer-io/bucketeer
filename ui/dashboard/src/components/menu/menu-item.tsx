import { FunctionComponent, PropsWithChildren } from 'react';
import { Popover } from '@radix-ui/themes';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';

const menuItemVariants = cva([
  'flex h-12 min-w-[200px] items-center justify-between rounded-lg bg-transparent px-3 text-primary-50 hover:cursor-pointer hover:bg-primary-400 hover:opacity-80'
]);

export type MenuItemProps = PropsWithChildren<{
  text: string;
  iconLeft?: FunctionComponent;
  iconRight?: FunctionComponent;
  selected?: boolean;
}>;

const MenuItem = ({
  text,
  iconLeft: SvgIconLeft,
  iconRight: SvgIconRight,
  selected,
  children
}: MenuItemProps) => {
  return (
    <Popover.Root>
      <Popover.Trigger>
        <li className={cn(menuItemVariants(), selected && 'bg-primary-400')}>
          <div className="flex items-center gap-2 text-primary-50">
            {SvgIconLeft && <SvgIconLeft />}
            <p className="typo-para-medium">{text}</p>
          </div>
          {SvgIconRight && <SvgIconRight />}
        </li>
      </Popover.Trigger>
      {children && (
        <Popover.Content className="max-w-none border-none p-0">
          <div>{children}</div>
        </Popover.Content>
      )}
    </Popover.Root>
  );
};

export default MenuItem;
