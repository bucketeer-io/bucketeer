import { forwardRef, Ref, useMemo } from 'react';
import * as CheckboxPrimitive from '@radix-ui/react-checkbox';
import { v4 as uuid } from 'uuid';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import Icon from 'components/icon';

type CheckboxProps = CheckboxPrimitive.CheckboxProps & {
  title?: string;
  description?: string;
  isExpand?: boolean;
  isReverse?: boolean;
};

const Checkbox = forwardRef(
  (
    {
      checked,
      title,
      description,
      isExpand,
      isReverse,
      onCheckedChange,
      ...props
    }: CheckboxProps,
    ref: Ref<HTMLButtonElement>
  ) => {
    const inputId = useMemo(() => uuid(), []);

    return (
      <div
        className={cn('flex w-fit items-center gap-x-2', {
          'w-full justify-between': isExpand,
          'flex-row-reverse': isReverse
        })}
      >
        <div className="flex-center size-5">
          <CheckboxPrimitive.Root
            className={cn(
              'flex-center size-5 rounded border border-gray-500 transition-colors duration-200',
              {
                'border-primary-500 bg-primary-500': checked
              }
            )}
            checked={checked}
            id={inputId}
            ref={ref}
            onCheckedChange={onCheckedChange}
            {...props}
          >
            <CheckboxPrimitive.Indicator
              className={cn('flex-center size-full opacity-0', {
                'opacity-100': checked
              })}
              forceMount={true}
            >
              <Icon icon={IconChecked} size={'fit'} className="text-white" />
            </CheckboxPrimitive.Indicator>
          </CheckboxPrimitive.Root>
        </div>
        {title && (
          <label className="flex flex-col gap-y-2 text-left" htmlFor={inputId}>
            <span className="typo-para-medium text-gray-700">{title}</span>
            {description && (
              <span className="typo-para-small text-additional-gray-500">
                {description}
              </span>
            )}
          </label>
        )}
      </div>
    );
  }
);

export default Checkbox;
