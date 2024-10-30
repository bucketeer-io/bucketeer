import { FunctionComponent } from 'react';
import { PropsWithChildren } from 'react';
import { AddonSlot } from '@types';
import { cn } from 'utils/style';
import Icon from 'components/icon';

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
        'flex cursor-pointer items-center gap-x-2 p-2 text-gray-700',
        'hover:bg-primary-50 [&>*]:hover:text-primary-500',
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

export type PopoverItemProps = {
  type: 'trigger' | 'item';
  addonSlot?: AddonSlot;
  icon?: FunctionComponent;
  label?: string;
  onClick?: () => void;
};

const PopoverItem = ({
  type,
  addonSlot,
  icon,
  label,
  onClick
}: PopoverItemProps) => {
  return (
    <PopoverItemWrapper type={type} addonSlot={addonSlot} onClick={onClick}>
      {icon && (
        <span
          className={'flex size-5 items-center justify-center text-gray-600'}
        >
          <Icon icon={icon} size={type === 'item' ? 'xxs' : 'sm'} />
        </span>
      )}
      {label && <span className={'typo-para-small select-none'}>{label}</span>}
    </PopoverItemWrapper>
  );
};

export default PopoverItem;
