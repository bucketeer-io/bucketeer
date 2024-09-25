import { useMemo } from 'react';
import { IconInfoFilled } from 'react-icons-material-design';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import Icon from 'components/icon';

const subStatus = cva(
  [
    'h-[26px] px-2 grid place-items-center typo-para-small rounded-[4px] flex items-center gap-1 w-fit'
  ],
  {
    variants: {
      variant: {
        new: ['bg-accent-blue-50 text-accent-blue-500'],
        activity: ['bg-accent-green-50 text-accent-green-500'],
        noActivity: ['bg-accent-yellow-50 text-accent-yellow-500'],
        waiting: ['bg-accent-orange-50 text-accent-orange-500'],
        inUse: ['bg-accent-green-50 text-accent-green-500']
      }
    },
    defaultVariants: {
      variant: 'new'
    }
  }
);

export type StatusTagType =
  | 'new'
  | 'activity'
  | 'noActivity'
  | 'waiting'
  | 'inUse';

export type StatusTagProps = {
  variant: StatusTagType;
  label?: string;
};

const StatusTag = ({ variant, label }: StatusTagProps) => {
  const name = useMemo(() => {
    switch (variant) {
      case 'inUse':
        return 'In Use';
      case 'noActivity':
        return 'No Activity';
      default:
        return `${variant[0]?.toUpperCase()}${variant?.slice(1, variant.length)}`;
    }
  }, [variant]);

  return (
    <div className={cn(subStatus({ variant }))}>
      {variant === 'noActivity' && <Icon icon={IconInfoFilled} size="xs" />}
      {label || name?.replace('_', ' ')}
    </div>
  );
};

export default StatusTag;
