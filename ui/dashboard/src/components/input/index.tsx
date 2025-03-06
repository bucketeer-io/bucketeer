import { forwardRef, InputHTMLAttributes } from 'react';
import type { Ref, ChangeEvent } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import { useInputGroupContext } from 'components/input-group/context';

export interface InputProps
  extends Omit<
    InputHTMLAttributes<HTMLInputElement>,
    'value' | 'size' | 'onChange' | 'onBlur'
  > {
  size?: 'sm' | 'md' | 'lg';
  value?: string | number | undefined;
  variant?: 'primary' | 'secondary';
  onChange?: (value: string, event: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (value: string, event: ChangeEvent<HTMLInputElement>) => void;
}

const inputVariants = cva(
  [
    'typo-para-medium border rounded-lg w-full',
    'disabled:cursor-not-allowed disabled:bg-gray-100'
  ],
  {
    variants: {
      variant: {
        primary: 'border-gray-400 text-gray-700 disabled:border-gray-400',
        secondary:
          '!border-primary-200 text-primary-500 disabled:!border-primary-100 bg-white'
      },
      size: {
        sm: 'px-4 py-2',
        md: 'px-4 py-[11px]',
        lg: 'px-4 py-4'
      },
      addonSlot: {
        left: 'pl-10',
        right: 'pr-10'
      }
    }
  }
);

const Input = forwardRef(
  (
    {
      className,
      size = 'md',
      value: _value,
      onChange,
      onBlur,
      role = 'presentation',
      autoComplete = 'off',
      variant = 'primary',
      ...props
    }: InputProps,
    ref: Ref<HTMLInputElement>
  ) => {
    const { addonSlot } = useInputGroupContext();
    const value = _value === undefined && onChange ? '' : _value;

    return (
      <input
        ref={ref}
        className={cn(
          inputVariants({ className, size, addonSlot, variant }),
          className
        )}
        role={role}
        autoComplete={autoComplete}
        value={value}
        onChange={event => {
          onChange?.(event.target.value, event);
        }}
        onBlur={event => {
          onBlur?.(event.target.value, event);
        }}
        {...props}
      />
    );
  }
);

export default Input;
