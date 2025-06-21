import React, {
  forwardRef,
  ReactNode,
  Ref,
  RefObject,
  useRef,
  type FunctionComponent
} from 'react';
import * as PopoverPrimitive from '@radix-ui/react-popover';
import type { PopoverContentProps } from '@radix-ui/react-popover';
import { AddonSlot, Color } from '@types';
import { cn } from 'utils/style';
import { IconClose } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import PopoverItem from './popover-item';

export type PopoverOption<PopoverValue> = {
  value: PopoverValue;
  icon?: FunctionComponent;
  label: ReactNode;
  description?: string;
  disabled?: boolean;
  tooltip?: string;
  color?: Color;
  [key: string]:
    | string
    | number
    | boolean
    | ReactNode
    | FunctionComponent
    | PopoverValue
    | undefined;
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
      'max-h-[260px] min-w-[167px] overflow-auto rounded-lg bg-gray-50 p-1 shadow-dropdown',
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
  triggerCls?: string;
  icon?: FunctionComponent;
  options?: PopoverOption<PopoverValue>[];
  disabled?: boolean;
  value?: PopoverValue | undefined;
  modal?: boolean;
  className?: string;
  closeWhenSelected?: boolean;
  children?: ReactNode;
  closeBtnCls?: string;
  sideOffset?: number;
  closeRef?: RefObject<HTMLButtonElement>;
  onClick?: (value: PopoverValue) => void;
  onOpenChange?: (open: boolean) => void;
  onPointerDownOutside?: () => void;
};

const Popover = forwardRef(
  (
    {
      align = 'start',
      expand,
      trigger,
      triggerLabel = '',
      triggerCls,
      icon,
      addonSlot,
      options,
      disabled,
      modal = false,
      className,
      closeWhenSelected = true,
      children,
      closeBtnCls,
      sideOffset = 0,
      closeRef,
      onClick,
      onOpenChange,
      onPointerDownOutside
    }: PopoverProps<PopoverValue>,
    ref: Ref<HTMLDivElement>
  ) => {
    const popoverCloseRef = useRef<HTMLButtonElement>(null);

    const handleSelectItem = (value: PopoverValue) => {
      onClick!(value);
      if (closeWhenSelected) (closeRef ?? popoverCloseRef)?.current?.click();
    };

    return (
      <PopoverRoot modal={modal} onOpenChange={onOpenChange}>
        <PopoverTrigger
          className={cn(
            'typo-para-small flex items-center justify-center gap-x-2 text-gray-700 hover:text-gray-600 hover:drop-shadow disabled:cursor-not-allowed',
            {
              'flex-row-reverse': addonSlot === 'right',
              'w-full justify-between': expand === 'full'
            },
            triggerCls
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
            sideOffset={sideOffset}
            onPointerDownOutside={onPointerDownOutside}
          >
            <PopoverClose
              ref={closeRef ?? popoverCloseRef}
              className={cn('hidden', closeBtnCls)}
            >
              <Icon icon={IconClose} size={'sm'} className="flex-center" />
            </PopoverClose>
            {children
              ? children
              : options?.map((item, index) => (
                  <Tooltip
                    align="end"
                    key={index}
                    trigger={
                      <div>
                        <PopoverItem
                          type="item"
                          addonSlot={addonSlot}
                          icon={item.icon}
                          label={item.label}
                          color={item?.color}
                          disabled={item?.disabled}
                          onClick={() =>
                            onClick && handleSelectItem(item.value)
                          }
                        />
                      </div>
                    }
                    content={item.tooltip}
                    className="max-w-[300px]"
                  />
                ))}
          </PopoverContent>
        </PopoverPrimitive.Portal>
      </PopoverRoot>
    );
  }
);

export { PopoverRoot, PopoverTrigger, PopoverClose, PopoverContent, Popover };
