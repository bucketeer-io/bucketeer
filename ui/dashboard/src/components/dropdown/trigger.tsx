import {
  ComponentPropsWithoutRef,
  ElementRef,
  forwardRef,
  ReactNode,
  useCallback,
  useRef
} from 'react';
import { IconExpandMoreRound } from 'react-icons-material-design';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import { IconClose } from '@icons';
import Icon from 'components/icon';
import Spinner from 'components/spinner';

const triggerVariants = cva(
  [
    'flex items-center px-3 py-[11px] gap-x-3 w-fit border rounded-lg bg-white max-[600px]:py-2 max-[600px]:text-sm',
    'disabled:cursor-not-allowed disabled:border-gray-400 disabled:bg-gray-100 disabled:!shadow-none'
  ],
  {
    variants: {
      variant: {
        primary:
          'border-primary-500 hover:shadow-border-primary-500 [&>*]:text-primary-500',
        secondary:
          'border-gray-400 hover:shadow-border-gray-400 [&_div]:text-gray-700 [&_span]:text-gray-600 [&>i]:text-gray-500'
      }
    },
    defaultVariants: {
      variant: 'secondary'
    }
  }
);

export type DropdownMenuTriggerProps = ComponentPropsWithoutRef<
  typeof DropdownMenuPrimitive.Trigger
> & {
  label?: ReactNode;
  description?: string;
  isExpand?: boolean;
  placeholder?: ReactNode;
  variant?: 'primary' | 'secondary';
  showArrow?: boolean;
  showClear?: boolean;
  trigger?: ReactNode;
  ariaLabel?: string;
  loading?: boolean;
  onClear?: () => void;
};

export const DropdownMenuTrigger = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Trigger>,
  DropdownMenuTriggerProps
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
      showClear = false,
      trigger,
      ariaLabel,
      loading,
      onClear,
      ...props
    },
    ref
  ) => {
    const clearRef = useRef<HTMLDivElement>(null);

    const handlePointerDown = useCallback(
      (e: React.MouseEvent) => {
        const target = e.target as HTMLElement;
        if (
          (clearRef.current && clearRef.current.contains(target)) ||
          (ariaLabel && target.ariaLabel === ariaLabel)
        ) {
          e.preventDefault();
        }
      },
      [ariaLabel]
    );

    const handleClearClick = useCallback(
      (e: React.MouseEvent) => {
        e.stopPropagation();
        e.preventDefault();
        onClear?.();
      },
      [onClear]
    );

    return (
      <DropdownMenuPrimitive.Trigger
        type="button"
        ref={ref}
        className={cn(
          triggerVariants({ variant }),
          { 'justify-between w-full': isExpand },
          className,
          'group'
        )}
        onPointerDown={handlePointerDown}
        {...props}
      >
        <>
          <div className="flex items-center w-full justify-between typo-para-medium overflow-hidden">
            {trigger ? (
              trigger
            ) : label ? (
              <div className="max-w-full truncate">
                {label} {description && <span>{description}</span>}
              </div>
            ) : (
              <p className="!text-gray-500">{placeholder}</p>
            )}
          </div>

          {showClear && label && !loading && (
            <div
              ref={clearRef}
              className="size-6 min-w-6 pointer-events-auto"
              onClick={handleClearClick}
            >
              <Icon
                icon={IconClose}
                size="md"
                className="text-gray-500 hover:text-gray-900"
              />
            </div>
          )}

          {showArrow && !loading && (
            <div className="size-6 min-w-6 transition-all duration-200 group-data-[state=closed]:rotate-0 group-data-[state=open]:rotate-180">
              <Icon icon={IconExpandMoreRound} size="md" color="gray-500" />
            </div>
          )}

          {loading && (
            <div className="flex-center size-fit">
              <Spinner className="size-4 border-2" />
            </div>
          )}
        </>
      </DropdownMenuPrimitive.Trigger>
    );
  }
);
