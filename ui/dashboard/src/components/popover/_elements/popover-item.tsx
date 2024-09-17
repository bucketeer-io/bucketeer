import { FunctionComponent } from 'react';
import { AddonSlot } from '@types';
import Icon from 'components/icon';
import styles from '../styles.module.css';
import PopoverItemWrapper from './popover-item-wrapper';

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
        <span className={styles.icon}>
          <Icon icon={icon} size={type === 'item' ? 'xxs' : 'sm'} />
        </span>
      )}
      {label && <span className={styles.label}>{label}</span>}
    </PopoverItemWrapper>
  );
};

export default PopoverItem;
