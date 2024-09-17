import { PropsWithChildren } from 'react';
import clsx from 'clsx';
import { AddonSlot } from '@types';
import styles from '../styles.module.css';

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
      className={clsx(styles.item, addonSlot === 'right' && styles['reverse'])}
      onClick={onClick && onClick}
    >
      {children}
    </div>
  );
};

export default PopoverItemWrapper;
