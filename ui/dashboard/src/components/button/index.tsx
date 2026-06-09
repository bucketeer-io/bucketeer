import { forwardRef } from 'react';
import type { ButtonHTMLAttributes } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from 'utils/style';
import Spinner from 'components/spinner';

const buttonVariants = cva(
  'inline-flex animate-fade gap-2 disabled:cursor-not-allowed items-center justify-center duration-300 ease-out whitespace-nowrap',
  {
    variants: {
      variant: {
        primary: [
          'bg-primary-500 text-gray-50',
          'rounded-lg px-6 py-2',
          'hover:bg-primary-700',
          'disabled:bg-primary-200 disabled:text-primary-50 dark:disabled:bg-dark-black-700 dark:disabled:text-dark-gray-200'
        ],
        secondary: [
          'text-primary-500 dark:text-dark-purple-500 shadow-border-primary-500 dark:shadow-border-primary-glow',
          'rounded-lg px-6 py-2',
          'hover:text-primary-700 dark:hover:text-dark-purple-600 hover:shadow-border-primary-700 dark:hover:shadow-border-primary-glow',
          'disabled:text-gray-500 dark:disabled:text-dark-gray-200 disabled:shadow-border-gray-500 dark:disabled:shadow-border-dark-black-700'
        ],
        'secondary-2': [
          'text-gray-700 dark:text-dark-gray-300 shadow-border-gray-300 dark:shadow-border-gray-500',
          'rounded-lg px-6 py-2',
          'hover:text-gray-900 dark:hover:text-dark-gray-400 hover:shadow-border-gray-500',
          'disabled:text-gray-500 dark:disabled:text-dark-gray-200 disabled:shadow-border-gray-300 disabled:bg-gray-100 dark:disabled:bg-dark-black-700'
        ],
        negative: [
          'bg-accent-red-500 text-gray-50 shadow-border-accent-red-500',
          'rounded-lg px-6 py-2',
          'hover:bg-accent-red-600 hover:shadow-border-accent-red-600',
          'disabled:bg-accent-red-50 disabled:shadow-border-gray-400 disabled:text-gray-500 dark:disabled:bg-dark-black-700 dark:disabled:shadow-border-dark-black-700 dark:disabled:text-dark-gray-200'
        ],
        text: [
          'text-primary-500 dark:text-dark-purple-500 px-2',
          'hover:text-primary-600 dark:hover:text-dark-purple-600',
          'disabled:text-gray-500 dark:disabled:text-dark-gray-200'
        ],
        grey: [
          'text-gray-600 dark:text-dark-gray-200 rounded-lg',
          'hover:text-gray-700 dark:hover:text-dark-gray-400',
          'disabled:text-gray-400 dark:disabled:text-dark-gray-200'
        ]
      },
      size: {
        xs: 'typo-para-tiny h-8',
        sm: 'typo-para-small h-10',
        md: 'typo-para-medium h-12',
        lg: 'typo-para-big h-14',
        icon: 'px-1 w-12 h-12',
        'icon-sm': 'px-1 w-9 h-9'
      }
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md'
    }
  }
);

export interface ButtonProps
  extends
    ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  loading?: boolean;
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    { className, variant, size, loading, children, disabled, ...props },
    ref
  ) => {
    return (
      <button
        className={cn(
          buttonVariants({ variant, size }),
          { relative: loading },
          className
        )}
        disabled={loading || disabled}
        ref={ref}
        {...props}
      >
        {children}
        {loading && (
          <div className="absolute inset-0 flex h-full w-full items-center justify-center">
            <Spinner />
          </div>
        )}
      </button>
    );
  }
);

export default Button;
