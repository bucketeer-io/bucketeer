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
  description: string;
};

type Props = CardIconProps &
  CardDescriptionProps & {
    showArrow?: boolean;
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
  description
}: CardDescriptionProps) => {
  return (
    <div className="flex flex-col flex-1 gap-y-2 overflow-hidden">
      <p className="w-full typo-para-small text-gray-700 truncate uppercase">
        {title}
      </p>
      <p className="text-5xl leading-8 font-medium text-gray-900">{count}</p>
      <p className="w-full typo-para-small text-gray-700 truncate">
        {description}
      </p>
    </div>
  );
};

const OverviewCard = ({
  icon,
  color,
  title,
  count,
  description,
  showArrow
}: Props) => {
  return (
    <div className="flex flex-1 items-center p-4 gap-x-4 w-full min-w-[300px] bg-white border border-gray-200 rounded-2xl overflow-hidden">
      <CardIcon icon={icon} color={color} />
      <CardDescription title={title} count={count} description={description} />
      {showArrow && (
        <Icon icon={IconChevronRight} size={'md'} color="gray-500" />
      )}
    </div>
  );
};

export default OverviewCard;
