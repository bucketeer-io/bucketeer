import {
  ComponentPropsWithoutRef,
  ElementRef,
  FunctionComponent,
  ReactNode,
  forwardRef,
  useMemo
} from 'react';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';
import NameWithTooltip from 'elements/name-with-tooltip';
import { DropdownValue } from './types';

export type DropdownMenuItemProps = ComponentPropsWithoutRef<
  typeof DropdownMenuPrimitive.Item
> & {
  icon?: FunctionComponent;
  iconElement?: ReactNode;
  isMultiselect?: boolean;
  isSelected?: boolean;
  isSelectedItem?: boolean;
  label?: ReactNode;
  value: DropdownValue;
  description?: string;
  closeWhenSelected?: boolean;
  additionalElement?: ReactNode;
  disabled?: boolean;
  /** When true, renders a tooltip when the label text overflows. */
  withTooltip?: boolean;
  onSelectOption?: (value: DropdownValue, event: Event) => void;
};

export const DropdownMenuItem = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Item>,
  DropdownMenuItemProps
>(
  (
    {
      className,
      icon,
      iconElement,
      label,
      value,
      description,
      isMultiselect,
      isSelected,
      isSelectedItem,
      closeWhenSelected = true,
      additionalElement,
      disabled,
      withTooltip = false,
      onSelectOption,
      ...props
    },
    ref
  ) => {
    // Use only value for the ID — label may be a ReactNode (not serialisable)
    const itemId = useMemo(() => `dropdown-item-${value}`, [value]);

    return (
      <DropdownMenuPrimitive.Item
        ref={ref}
        disabled={disabled}
        className={cn(
          'relative flex items-center w-full cursor-pointer select-none rounded-[5px] p-2 gap-x-2 outline-none transition-colors hover:bg-gray-100 data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
          { '!bg-gray-100': isSelectedItem },
          className
        )}
        onSelect={
          onSelectOption
            ? event => {
                if (!closeWhenSelected || isMultiselect) event.preventDefault();
                onSelectOption(value, event);
              }
            : undefined
        }
        {...props}
      >
        {isMultiselect && <Checkbox checked={isSelected} />}

        {(iconElement || icon) &&
          (iconElement ? (
            iconElement
          ) : (
            <div className="flex-center size-5">
              <Icon
                icon={icon as FunctionComponent}
                size="xs"
                color="gray-600"
              />
            </div>
          ))}

        <div className="flex flex-col gap-y-1.5 w-full overflow-hidden">
          {withTooltip ? (
            <NameWithTooltip
              id={itemId}
              content={<NameWithTooltip.Content content={label} id={itemId} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={itemId}
                  name={label as string}
                  haveAction={false}
                  maxLines={1}
                  className="cursor-pointer"
                />
              }
              maxLines={1}
            />
          ) : (
            <div className="typo-para-medium leading-5 text-gray-700 truncate">
              {label}
            </div>
          )}
          {description && (
            <p className="typo-para-small text-gray-500">{description}</p>
          )}
        </div>

        {additionalElement}
        {isSelectedItem && <IconChecked className="text-primary-500 w-6" />}
      </DropdownMenuPrimitive.Item>
    );
  }
);
