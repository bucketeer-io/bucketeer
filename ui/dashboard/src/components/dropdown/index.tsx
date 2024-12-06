import {
  ComponentPropsWithoutRef,
  ElementRef,
  forwardRef,
  FunctionComponent,
  ReactNode
} from 'react';
import { IconExpandMoreRound } from 'react-icons-material-design';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import { IconSearch } from '@icons';
import Icon from 'components/icon';
import Input, { InputProps } from 'components/input';

export type DropdownValue = number | string;

export type DropdownOption = {
  label: string;
  value: DropdownValue;
  icon?: FunctionComponent;
  description?: boolean;
  haveCheckbox?: boolean;
};

const DropdownMenu = DropdownMenuPrimitive.Root;

const DropdownMenuGroup = DropdownMenuPrimitive.Group;

const DropdownMenuPortal = DropdownMenuPrimitive.Portal;

const triggerVariants = cva(
  [
    'flex items-center px-3 py-[11px] gap-x-3 w-fit border rounded-lg bg-white',
    'disabled:cursor-not-allowed disabled:border-gray-400 disabled:bg-gray-100 disabled:!shadow-none'
  ],
  {
    variants: {
      variant: {
        primary:
          'border-primary-500 hover:shadow-border-primary-500 [&>*]:text-primary-500',
        secondary:
          'border-gray-400 hover:shadow-border-gray-400 [&_p]:text-gray-700 [&_span]:text-gray-600 [&>i]:text-gray-500'
      }
    },
    defaultVariants: {
      variant: 'secondary'
    }
  }
); 

const DropdownMenuTrigger = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Trigger>,
  ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Trigger> & {
    label?: string;
    description?: string;
    isExpand?: boolean;
    placeholder?: string;
    variant?: 'primary' | 'secondary';
    showArrow?: boolean;
    trigger?: ReactNode;
  }
>(
  (
    {
      className,
      variant,
      label,
      description,
      isExpand,
      placeholder = '',
      showArrow = true,
      trigger,
      ...props
    },
    ref
  ) => (
    <DropdownMenuPrimitive.Trigger
      ref={ref}
      className={cn(
        triggerVariants({
          variant
        }),
        {
          'justify-between w-full': isExpand
        },
        className
      )}
      {...props}
    >
      <div className="flex items-center w-full justify-between typo-para-medium">
        {trigger ? (
          trigger
        ) : label ? (
          <p>
            {label} {description && <span>{description}</span>}
          </p>
        ) : (
          <p className={'!text-gray-500'}>{placeholder}</p>
        )}
      </div>

      {showArrow && <Icon icon={IconExpandMoreRound} size={'md'} />}
    </DropdownMenuPrimitive.Trigger>
  )
);

const DropdownMenuContent = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Content>,
  ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Content> & {
    isExpand?: boolean;
  }
>(({ className, sideOffset = 4, isExpand, ...props }, ref) => (
  <DropdownMenuPrimitive.Portal>
    <DropdownMenuPrimitive.Content
      ref={ref}
      sideOffset={sideOffset}
      className={cn(
        'z-50 min-w-[196px] max-h-[252px] overflow-x-hidden overflow-y-auto rounded-lg border bg-white p-1 shadow-dropdown small-scroll',
        'data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2',
        { 'dropdown-menu-expand': isExpand },
        className
      )}
      {...props}
    />
  </DropdownMenuPrimitive.Portal>
));

const DropdownMenuItem = forwardRef<
  React.ElementRef<typeof DropdownMenuPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Item> & {
    icon?: FunctionComponent;
    isMultiselect?: boolean;
    selected?: boolean;
    label?: string;
    value: DropdownValue;
    description?: string;
    closeWhenSelected?: boolean;
    onSelectOption?: (value: DropdownValue, event: Event) => void;
  }
>(
  (
    {
      className,
      icon,
      label,
      value,
      description,
      closeWhenSelected = true,
      onSelectOption,
      ...props
    },
    ref
  ) => (
    <DropdownMenuPrimitive.Item
      ref={ref}
      className={cn(
        'relative flex items-center w-full cursor-pointer select-none rounded-[5px] p-2 gap-x-2 outline-none transition-colors hover:bg-gray-100 data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
        className
      )}
      onSelect={
        onSelectOption
          ? event => {
              if (!closeWhenSelected) event.preventDefault();
              return onSelectOption(value, event);
            }
          : undefined
      }
      {...props}
    >
      {icon && (
        <div className="flex-center size-5">
          <Icon icon={icon} size={'xs'} color="gray-600" />
        </div>
      )}

      <div className="flex flex-col gap-y-1.5">
        <p className="typo-para-medium leading-5 text-gray-700">{label}</p>
        {description && (
          <p className="typo-para-small leading-[14px] text-gray-500">
            {description}
          </p>
        )}
      </div>
    </DropdownMenuPrimitive.Item>
  )
);

type DropdownSearchProps = InputProps;

const DropdownMenuSearch = ({
  value,
  onChange,
  ...props
}: DropdownSearchProps) => {
  return (
    <div className="sticky top-0 left-0 right-0 flex items-center w-full px-3 py-[11.5px] gap-x-2 border-b border-gray-200 bg-white z-50">
      <div className="flex-center size-5">
        <Icon icon={IconSearch} size={'xs'} color="gray-500" />
      </div>
      <Input
        {...props}
        value={value}
        onChange={onChange}
        className="p-0 border-none"
        autoFocus={true}
      />
    </div>
  );
};

export {
  DropdownMenu,
  DropdownMenuGroup,
  DropdownMenuPortal,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSearch
};
