import { useMemo } from 'react';
import type { FunctionComponent } from 'react';
import {
  IconKeyboardArrowLeftFilled,
  IconKeyboardArrowRightFilled,
  IconKeyboardDoubleArrowLeftFilled,
  IconKeyboardDoubleArrowRightFilled
} from 'react-icons-material-design';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import Icon from 'components/icon';

const cellVariant = cva(
  [
    'min-w-8 h-8 p-1 rounded-lg flex items-center justify-center text-gray-500 dark:text-dark-gray-200'
  ],
  {
    variants: {
      variant: {
        number: ['bg-white dark:bg-dark-black-800'],
        next: ['border dark:border-dark-black-700'],
        previous: ['border dark:border-dark-black-700'],
        first: ['border dark:border-dark-black-700'],
        last: ['border dark:border-dark-black-700']
      }
    },
    defaultVariants: {
      variant: 'number'
    }
  }
);

export type PaginationCellType =
  | 'number'
  | 'next'
  | 'first'
  | 'previous'
  | 'last';

export type PaginationCellProps = {
  checked?: boolean;
  value?: number;
  variant?: PaginationCellType;
  disabled?: boolean;
  onClick?: (value?: number) => void;
};

const PaginationIcon = ({ icon }: { icon: FunctionComponent }) => (
  <div className="flex-center text-gray-500 dark:text-dark-gray-200">
    <Icon icon={icon} />
  </div>
);

const PaginationCell = ({
  checked,
  value,
  variant = 'number',
  disabled = false,
  onClick
}: PaginationCellProps) => {
  const variantRender = useMemo(() => {
    switch (variant) {
      case 'number':
        return value;
      case 'next':
        return <PaginationIcon icon={IconKeyboardArrowRightFilled} />;
      case 'previous':
        return <PaginationIcon icon={IconKeyboardArrowLeftFilled} />;
      case 'first':
        return <PaginationIcon icon={IconKeyboardDoubleArrowLeftFilled} />;
      case 'last':
        return <PaginationIcon icon={IconKeyboardDoubleArrowRightFilled} />;
    }
  }, [variant]);

  return (
    <button
      className={cn(
        cellVariant({ variant }),
        checked && 'bg-primary-500 text-white typo-para-medium',
        disabled && 'cursor-not-allowed opacity-80'
      )}
      onClick={() => !disabled && onClick && onClick(value)}
    >
      {variantRender}
    </button>
  );
};

export default PaginationCell;
