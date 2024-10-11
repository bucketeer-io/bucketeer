import type { ReactNode } from 'react';
import clsx from 'clsx';
import * as RadixTabs from '@radix-ui/react-tabs';
import styles from './styles.module.css';

export type TabsTriggerProps = {
  value: string;
  disabled?: boolean;
  onClick?: () => void;
  icon?: ReactNode;
} & (
  | {
      children: ReactNode;
      title?: never;
      subtitle?: never;
    }
  | {
      title?: string;
      subtitle?: string | number;
      children?: never;
    }
);

const TabsTrigger = ({
  value,
  disabled,
  title,
  subtitle,
  icon,
  children,
  onClick
}: TabsTriggerProps) => {
  return (
    <RadixTabs.Trigger
      className={clsx(
        styles.trigger,
        disabled && styles.disabled,
        icon && 'flex items-center'
      )}
      value={value}
      onClick={onClick}
    >
      {icon && <span className="mr-1.5 flex items-center">{icon}</span>}
      {title || children}
      {subtitle !== undefined && (
        <div className={styles.subtitle}>{subtitle}</div>
      )}
    </RadixTabs.Trigger>
  );
};

export default TabsTrigger;
