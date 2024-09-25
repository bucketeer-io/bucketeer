import { PropsWithChildren } from 'react';
import { AddonSlot } from '@types';
import { cn } from 'utils/style';

type PopoverItemWrapperProps = PropsWithChildren & {
  type: 'trigger' | 'item';
  addonSlot?: AddonSlot;
  onClick?: () => void;
};
const PopoverItemWrapper = ({
  type,
  children,
  addonSlot,
  onClick
}: PopoverItemWrapperProps) => {
  if (type === 'trigger') return <>{children}</>;
  return (
    <div
      className={cn(
        'flex cursor-pointer items-center gap-x-2 p-2 text-gray-700 hover:bg-primary-50 [&>*]:hover:text-primary-500',
        {
          'flex-row-reverse': addonSlot === 'right'
        }
      )}
      onClick={onClick && onClick}
    >
      {children}
    </div>
  );
};

export default PopoverItemWrapper;
