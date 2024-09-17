import { useRef, type FunctionComponent } from 'react';
import clsx from 'clsx';
import * as PopoverPrimitive from '@radix-ui/react-popover';
import type { PopoverContentProps } from '@radix-ui/react-popover';
import { AddonSlot } from '@types';
import { cn } from 'utils/style';
import PopoverItem from './_elements/popover-item';
import styles from './styles.module.css';

export type PopoverOption<PopoverValue> = {
  value: PopoverValue;
  icon?: FunctionComponent;
  label: string;
  description?: string;
};

export type PopoverValue = number | string;

export type PopoverProps<PopoverValue> = {
  align?: PopoverContentProps['align'];
  expand?: 'full';
  addonSlot?: AddonSlot;
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

const Popover = <PopoverValue,>({
  align = 'start',
  expand,
  triggerLabel = '',
  icon,
  addonSlot,
  options,
  disabled,
  modal = false,
  className,
  closeWhenSelected = true,
  onClick
}: PopoverProps<PopoverValue>) => {
  const popoverCloseRef = useRef<HTMLButtonElement>(null);

  const handleSelectItem = (value: PopoverValue) => {
    onClick!(value);
    if (closeWhenSelected) popoverCloseRef?.current?.click();
  };

  return (
    <PopoverPrimitive.Root modal={modal}>
      <PopoverPrimitive.Trigger
        className={clsx(
          styles.trigger,
          addonSlot === 'right' && styles['reverse'],
          expand === 'full' && styles.expand
        )}
        disabled={disabled}
      >
        <PopoverItem
          type="trigger"
          addonSlot={addonSlot}
          icon={icon}
          label={triggerLabel}
        />
      </PopoverPrimitive.Trigger>
      <PopoverPrimitive.Portal>
        <PopoverPrimitive.Content
          hideWhenDetached={true}
          className={cn(styles.content, className)}
          align={align}
        >
          <PopoverPrimitive.Close ref={popoverCloseRef} className="hidden" />
          {options.map((item, index) => (
            <PopoverItem
              type="item"
              key={index}
              addonSlot={addonSlot}
              icon={item.icon}
              label={item.label}
              onClick={() => onClick && handleSelectItem(item.value)}
            />
          ))}
        </PopoverPrimitive.Content>
      </PopoverPrimitive.Portal>
    </PopoverPrimitive.Root>
  );
};

export default Popover;
