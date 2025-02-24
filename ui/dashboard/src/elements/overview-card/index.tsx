import { FunctionComponent } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Icon from 'components/icon';

type CardIconProps = {
  icon: FunctionComponent;
  color:
    | 'green'
    | 'brand'
    | 'yellow'
    | 'gray'
    | 'pink'
    | 'red'
    | 'orange'
    | 'blue';
};

type CardDescriptionProps = {
  title: string;
  count: number;
  description?: string;
  highlightText?: string;
  highlightType?: 'increase' | 'decrease';
};

type Props = CardIconProps &
  CardDescriptionProps & {
    showArrow?: boolean;
    onClick?: () => void;
  };

const cardIconVariants = cva('flex-center size-[88px] p-5 rounded-lg', {
  variants: {
    color: {
      green: 'bg-accent-green-50',
      brand: 'bg-primary-50',
      yellow: 'bg-accent-yellow-50',
      gray: 'bg-accent-gray-200',
      pink: 'bg-accent-pink-50',
      red: 'bg-accent-red-50',
      orange: 'bg-accent-orange-50',
      blue: 'bg-accent-blue-50'
    }
  }
});

const CardIcon = ({ icon, color }: CardIconProps) => {
  return (
    <div className={cn(cardIconVariants({ color }))}>
      <Icon icon={icon} size={'fit'} />
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
    <div className="flex flex-col flex-1 gap-y-3 overflow-hidden">
      <p className="w-full typo-para-medium leading-5 text-gray-600 truncate capitalize">
        {title}
      </p>
      <p className="typo-head-bold-huge leading-6 text-gray-900">{count}</p>
      <div className="flex items-center gap-x-2">
        {highlightText && (
          <p
            className={cn('typo-head-bold-huge leading-6', {
              'text-accent-green-500': highlightType === 'increase',
              'text-accent-red-500': highlightType === 'decrease'
            })}
          >
            {highlightText}
          </p>
        )}
        {description && (
          <p className="w-full typo-para-small leading-5 text-gray-600 truncate">
            {description}
          </p>
        )}
      </div>
    </div>
  );
};

const OverviewCard = ({ icon, color, showArrow, onClick, ...props }: Props) => {
  return (
    <div
      className="flex flex-1 items-center p-4 gap-x-4 w-full min-w-[300px] bg-white shadow-card rounded-2xl overflow-hidden cursor-pointer hover:shadow-gray-300"
      onClick={onClick}
    >
      <CardIcon icon={icon} color={color} />
      <CardDescription {...props} />
      {showArrow && (
        <Icon icon={IconChevronRight} size={'md'} color="gray-500" />
      )}
    </div>
  );
};

export default OverviewCard;
