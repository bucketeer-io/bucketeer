import { FunctionComponent } from 'react';
import {
  IconCloseFilled,
  IconKeyboardArrowDownFilled
} from 'react-icons-material-design';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import Icon from 'components/icon';

const tagVariants = cva(
  [
    'typo-para-small flex justify-center items-center w-fit py-[6px] px-2 rounded-[4px] h-[26px]'
  ],
  {
    variants: {
      variant: {
        brand: ['text-primary-500 bg-primary-50'],
        gray: ['text-gray-600 bg-gray-200'],
        green: ['text-accent-green-500 bg-accent-green-50'],
        orange: ['text-accent-orange-500 bg-accent-orange-50'],
        blue: ['text-accent-blue-500 bg-accent-blue-50'],
        pink: ['text-accent-pink-500 bg-accent-pink-50'],
        yellow: ['text-accent-yellow-500 bg-accent-yellow-50']
      }
    },
    defaultVariants: {
      variant: 'gray'
    }
  }
);

export type TagVariant =
  | 'brand'
  | 'gray'
  | 'green'
  | 'orange'
  | 'blue'
  | 'pink'
  | 'yellow';

export type TagType =
  | 'text'
  | 'leftIcon'
  | 'dismissible'
  | 'expandable'
  | 'iconOnly';

export type TagProps = {
  text?: string;
  variant?: TagVariant;
  type?: TagType;
  icon?: FunctionComponent;
  className?: string;
};

const Tag = ({
  text,
  variant,
  type = 'text',
  icon: SvgIcon,
  className
}: TagProps) => {
  return (
    <div className={cn(tagVariants({ variant }), className)}>
      {['iconOnly', 'leftIcon'].includes(type) && SvgIcon && (
        <Icon icon={SvgIcon} size="xs" />
      )}
      {type !== 'iconOnly' &&
        text &&
        `${text[0].toUpperCase()}${text.slice(1, text.length)}`}
      {type === 'dismissible' && <Icon icon={IconCloseFilled} size="xs" />}
      {type === 'expandable' && (
        <Icon icon={IconKeyboardArrowDownFilled} size="xs" />
      )}
    </div>
  );
};

export default Tag;
