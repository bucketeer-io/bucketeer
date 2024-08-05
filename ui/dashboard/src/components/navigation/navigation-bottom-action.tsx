import { FunctionComponent, PropsWithChildren } from 'react';
import { IconButton, IconButtonProps, Popover } from '@radix-ui/themes';
import { cn } from 'utils/style';

export type NavigationBottomActionProps = PropsWithChildren<
  {
    icon: FunctionComponent;
    popoverSide?: 'center' | 'end' | 'start';
  } & IconButtonProps
>;

const NavigationBottomAction = ({
  icon: SvgIcon,
  className,
  children,
  popoverSide = 'start',
  ...props
}: NavigationBottomActionProps) => {
  return (
    <Popover.Root>
      {children && (
        <Popover.Content align={popoverSide} className="border-none p-0">
          <div>{children}</div>
        </Popover.Content>
      )}
      <Popover.Trigger>
        <IconButton
          radius="full"
          className={cn(
            'h-6 w-6 bg-transparent hover:cursor-pointer',
            className
          )}
          {...props}
        >
          <SvgIcon />
        </IconButton>
      </Popover.Trigger>
    </Popover.Root>
  );
};

export default NavigationBottomAction;
