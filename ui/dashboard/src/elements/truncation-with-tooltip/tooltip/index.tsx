import { ReactNode } from 'react';
import { cn } from 'utils/style';
import { IconAngleDown } from '@icons';
import Icon from 'components/icon';

const absoluteCenterCls = `absolute left-1/2 -translate-x-1/2`;

const TooltipWrapper = ({
  children,
  className
}: {
  children: ReactNode;
  className?: string;
}) => (
  <div
    className={cn(
      'flex-center w-[250px] lg:w-[500px] bottom-[calc(100%+16px)]',
      'opacity-0 transition-all delay-500 group-hover:opacity-100',
      absoluteCenterCls,
      className
    )}
  >
    {children}
  </div>
);

const TooltipContent = ({
  content,
  className
}: {
  content: ReactNode;
  className?: string;
}) => (
  <div
    className={cn(
      'relative flex-center',
      'w-fit h-0 p-0 max-w-full',
      'typo-para-small text-center break-words',
      'pointer-events-none rounded-md',
      'bg-transparent text-transparent',
      'group-hover:h-fit group-hover:p-2 group-hover:bg-gray-700 group-hover:text-white',
      className
    )}
  >
    {content}
  </div>
);

const TooltipArrow = () => (
  <div
    className={cn(
      'hidden size-fit top-[calc(100%-1px)]',
      'group-hover:flex-center',
      absoluteCenterCls
    )}
  >
    <Icon icon={IconAngleDown} className="!text-gray-700" size={'md'} />
  </div>
);

export const Tooltip = ({
  tooltipWrapperCls,
  tooltipContentCls,
  content
}: {
  tooltipWrapperCls?: string;
  tooltipContentCls?: string;
  content: ReactNode;
}) => (
  <TooltipWrapper className={tooltipWrapperCls}>
    <TooltipContent content={content} className={tooltipContentCls} />
    <TooltipArrow />
  </TooltipWrapper>
);
