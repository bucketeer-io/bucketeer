import { ReactNode, useState } from 'react';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';

const DiscardChangeItems = ({
  children,
  title
}: {
  children: ReactNode;
  title?: string;
}) => {
  const [isExpanded, setIsExpanded] = useState(true);
  return (
    <div className={cn('flex flex-col w-full bg-white  min-w-fit')}>
      <div
        className=" flex gap-1 pb-2 cursor-pointer"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <p className="typo-para-medium leading-4 text-gray-700">{title}</p>
        <Icon
          icon={IconChevronDown}
          className={cn('rotate-0 transition-all duration-100', {
            'rotate-180': isExpanded
          })}
          size={'sm'}
        />
      </div>
      <Divider className="border-gray-300" />
      <div
        className={cn(' pt-2 opacity-1 h-fit transition-all duration-100', {
          'opacity-0 h-0 overflow-hidden': !isExpanded
        })}
      >
        {children}
      </div>
    </div>
  );
};

export default DiscardChangeItems;
