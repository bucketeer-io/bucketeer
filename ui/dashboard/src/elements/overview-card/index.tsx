import { FunctionComponent, ReactNode } from 'react';
import { cva } from 'class-variance-authority';
import { Color, IconSize } from '@types';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Icon from 'components/icon';

export type OverviewIconColor =
  | 'green'
  | 'brand'
  | 'yellow'
  | 'gray'
  | 'pink'
  | 'red'
  | 'orange'
  | 'blue';

type CardIconProps = {
  icon: FunctionComponent;
  color: OverviewIconColor;
  iconSize?: IconSize;
  iconClassName?: string;
};

type CardDescriptionProps = {
  title: ReactNode;
  count?: string;
  description?: ReactNode;
  highlightText?: string;
  highlightType?: 'increase' | 'decrease';
};

type Props = CardIconProps &
  CardDescriptionProps & {
    className?: string;
    showArrow?: boolean;
    onClick?: () => void;
  };

const cardIconVariants = cva('flex-center size-16 p-5 rounded-lg', {
  variants: {
    color: {
      green: 'bg-accent-green-50',
      brand: 'bg-primary-50',
      yellow: 'bg-accent-yellow-50',
      gray: 'bg-gray-200',
      pink: 'bg-accent-pink-50',
      red: 'bg-accent-red-50',
      orange: 'bg-accent-orange-50',
      blue: 'bg-accent-blue-50'
    }
  }
});

const getIconColor = (color: OverviewIconColor): Color => {
  switch (color) {
    case 'green':
      return 'accent-green-500';
    case 'yellow':
      return 'accent-yellow-500';
    case 'gray':
      return 'gray-200';
    case 'pink':
      return 'accent-pink-500';
    case 'red':
      return 'accent-red-500';
    case 'orange':
      return 'accent-orange-500';
    case 'blue':
      return 'accent-blue-500';
    case 'brand':
    default:
      return 'primary-500';
  }
};

const CardIcon = ({ icon, color, iconSize, iconClassName }: CardIconProps) => {
  return (
    <div className={cn(cardIconVariants({ color }), iconClassName)}>
      <Icon icon={icon} size={iconSize} color={getIconColor(color)} />
    </div>
  );
};

const CardDescription = ({
  title,
  count,
  description,
  highlightText,
  highlightType = 'increase'
}: CardDescriptionProps) => {
  return (
    <div className="flex flex-col flex-1 gap-y-1 overflow-hidden">
      <div className="w-full typo-para-medium text-gray-600 truncate capitalize">
        {title}
      </div>
      {count && (
        <p className="typo-head-bold-huge leading-6 text-gray-900">{count}</p>
      )}
      {(description || highlightText) && (
        <div className="flex items-center gap-x-2">
          {highlightText && (
            <p
              className={cn('typo-head-bold-huge font-extrabold leading-6', {
                'text-accent-green-500': highlightType === 'increase',
                'text-accent-red-500': highlightType === 'decrease'
              })}
            >
              {highlightText}
            </p>
          )}
          {description && (
            <div className="w-full typo-para-small text-gray-600">
              {description}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

const OverviewCard = ({
  icon,
  color,
  showArrow,
  className,
  iconSize = 'fit',
  iconClassName = '',
  onClick,
  ...props
}: Props) => {
  return (
    <div
      className={cn(
        'flex flex-1 items-center p-4 gap-x-4 w-full min-w-[268px] bg-white shadow-card rounded-2xl overflow-hidden cursor-pointer hover:shadow-gray-300',
        className
      )}
      onClick={onClick}
    >
      <CardIcon
        icon={icon}
        color={color}
        iconSize={iconSize}
        iconClassName={iconClassName}
      />
      <CardDescription {...props} />
      {showArrow && (
        <Icon icon={IconChevronRight} size={'md'} color="gray-500" />
      )}
    </div>
  );
};

export default OverviewCard;
