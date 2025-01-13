import React, {
  forwardRef,
  ReactNode,
  Ref,
  useRef,
  type FunctionComponent
} from 'react';
import * as PopoverPrimitive from '@radix-ui/react-popover';
import type { PopoverContentProps } from '@radix-ui/react-popover';
import { AddonSlot } from '@types';
import { cn } from 'utils/style';
import PopoverItem from './popover-item';

export type PopoverOption<PopoverValue> = {
  value: PopoverValue;
  icon?: FunctionComponent;
  label: string;
  description?: string;
  disabled?: boolean;
};

export type PopoverValue = number | string;

const PopoverRoot = PopoverPrimitive.Root;
const PopoverTrigger = PopoverPrimitive.Trigger;
const PopoverContent = React.forwardRef<
  React.ElementRef<typeof PopoverPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof PopoverPrimitive.Content>
>(({ className, align = 'center', sideOffset = 4, ...props }, ref) => (
  <PopoverPrimitive.Content
    ref={ref}
    align={align}
    sideOffset={sideOffset}
    className={cn(
      'max-h-[260px] min-w-[167px] overflow-auto rounded-lg bg-gray-50 p-1 shadow-menu',
      className
    )}
    {...props}
  />
));
PopoverContent.displayName = PopoverPrimitive.Content.displayName;

const PopoverClose = PopoverPrimitive.Close;

export type PopoverProps<PopoverValue> = {
  align?: PopoverContentProps['align'];
  expand?: 'full';
  addonSlot?: AddonSlot;
  trigger?: ReactNode;
  triggerLabel?: string;
  icon?: FunctionComponent;
  options: PopoverOption<PopoverValue>[];
  disabled?: boolean;
  value?: PopoverValue | undefined;
  modal?: boolean;
  className?: string;
  closeWhenSelected?: boolean;
  onClick?: (value: PopoverValue) => void;
};

const Popover = forwardRef(
  (
    {
      align = 'start',
      expand,
      trigger,
      triggerLabel = '',
      icon,
      addonSlot,
      options,
      disabled,
      modal = false,
      className,
      closeWhenSelected = true,
      onClick
    }: PopoverProps<PopoverValue>,
    ref: Ref<HTMLDivElement>
  ) => {
    const popoverCloseRef = useRef<HTMLButtonElement>(null);

    const handleSelectItem = (value: PopoverValue) => {
      onClick!(value);
      if (closeWhenSelected) popoverCloseRef?.current?.click();
    };

    return (
      <PopoverRoot modal={modal}>
        <PopoverTrigger
          className={cn(
            'typo-para-small flex items-center justify-center gap-x-2 text-gray-700 hover:text-gray-600 hover:drop-shadow',
            {
              'flex-row-reverse': addonSlot === 'right',
              'w-full justify-between': expand === 'full'
            }
          )}
          disabled={disabled}
        >
          {trigger ? (
            trigger
          ) : (
            <PopoverItem
              type="trigger"
              addonSlot={addonSlot}
              icon={icon}
              label={triggerLabel}
            />
          )}
        </PopoverTrigger>
        <PopoverPrimitive.Portal>
          <PopoverContent
            ref={ref}
            hideWhenDetached={true}
            className={className}
            align={align}
          >
            <PopoverClose ref={popoverCloseRef} className="hidden" />
            {options.map((item, index) => (
              <PopoverItem
                type="item"
                key={index}
                addonSlot={addonSlot}
                icon={item.icon}
                label={item.label}
                disabled={item?.disabled}
                onClick={() => onClick && handleSelectItem(item.value)}
              />
            ))}
          </PopoverContent>
        </PopoverPrimitive.Portal>
      </PopoverRoot>
    );
  }
);

export { PopoverRoot, PopoverTrigger, PopoverClose, PopoverContent, Popover };
