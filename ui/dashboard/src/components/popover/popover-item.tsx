import { FunctionComponent, ReactNode } from 'react';
import { PropsWithChildren } from 'react';
import { AddonSlot, Color } from '@types';
import { cn } from 'utils/style';
import Icon from 'components/icon';

type PopoverItemWrapperProps = PropsWithChildren & {
  type: 'trigger' | 'item';
  addonSlot?: AddonSlot;
  disabled?: boolean;
  onClick?: () => void;
};
const PopoverItemWrapper = ({
  type,
  children,
  addonSlot,
  disabled,
  onClick
}: PopoverItemWrapperProps) => {
  if (type === 'trigger') return <>{children}</>;
  return (
    <div
      className={cn(
        'flex cursor-pointer items-center gap-x-2 p-2 text-gray-700',
        'hover:bg-primary-50 [&>*]:hover:text-primary-500',
        {
          'flex-row-reverse': addonSlot === 'right',
          '!bg-transparent !text-gray-400 [&>*]:hover:!text-gray-400 cursor-not-allowed':
            disabled
        }
      )}
      onClick={() => {
        if (!disabled && onClick) onClick();
      }}
    >
      {children}
    </div>
  );
};

export type PopoverItemProps = {
  type: 'trigger' | 'item';
  addonSlot?: AddonSlot;
  icon?: FunctionComponent;
  label?: ReactNode;
  disabled?: boolean;
  color?: Color;
  onClick?: () => void;
};

const PopoverItem = ({
  type,
  addonSlot,
  icon,
  label,
  disabled,
  color,
  onClick
}: PopoverItemProps) => {
  return (
    <PopoverItemWrapper
      disabled={disabled}
      type={type}
      addonSlot={addonSlot}
      onClick={onClick}
    >
      {icon && (
        <div
          className={cn(
            'flex size-5 items-center justify-center',
            disabled ? 'text-gray-400' : 'text-gray-600'
          )}
        >
          <Icon
            icon={icon}
            size={type === 'item' ? 'xxs' : 'sm'}
            color={color}
          />
        </div>
      )}
      {label && <div className={'typo-para-small select-none'}>{label}</div>}
    </PopoverItemWrapper>
  );
};

export default PopoverItem;
