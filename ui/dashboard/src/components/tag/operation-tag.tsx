import { useMemo } from 'react';
import {
  IconCalendarTodayOutlined,
  IconLoopOutlined,
  IconTrendingUpOutlined
} from 'react-icons-material-design';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import Icon from 'components/icon';

const operationTypes = cva(
  ['h-[26px] w-[26px] rounded-[4px] grid place-items-center'],
  {
    variants: {
      type: {
        up: ['bg-accent-blue-50 text-accent-blue-500'],
        calendar: ['bg-primary-50 text-primary-500'],
        loop: ['bg-accent-pink-50 text-accent-pink-500']
      }
    },
    defaultVariants: {
      type: 'up'
    }
  }
);

export type OperationType = 'up' | 'calendar' | 'loop';

export type OperationProps = {
  type: OperationType;
};

const OperationTag = ({ type }: OperationProps) => {
  const renderType = useMemo(() => {
    switch (type) {
      case 'up':
        return IconTrendingUpOutlined;
      case 'calendar':
        return IconCalendarTodayOutlined;
      case 'loop':
        return IconLoopOutlined;
    }
  }, [type]);

  return (
    <div className={cn(operationTypes({ type }))}>
      <Icon icon={renderType} size="xs" />
    </div>
  );
};

export default OperationTag;
