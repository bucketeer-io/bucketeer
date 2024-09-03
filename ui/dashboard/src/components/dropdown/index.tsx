import type { FunctionComponent, ReactNode } from 'react';
import { IconExpandMoreRound } from 'react-icons-material-design';
import clsx from 'clsx';
import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import type { DropdownMenuContentProps } from '@radix-ui/react-dropdown-menu';
import { cn } from 'utils/style';
import Icon from 'components/icon';
import styles from './styles.module.css';

export type DropdownOption<DropdownValue> = {
  value: DropdownValue | undefined;
  icon?: FunctionComponent;
  label: string;
  description?: string;
};

export type DropdownProps<DropdownValue> = {
  align?: DropdownMenuContentProps['align'];
  expand?: 'full';
  addonSlot?: 'left' | 'right';
  placeholder?: string;
  icon?: FunctionComponent;
  title?: string;
  options: DropdownOption<DropdownValue>[];
  action?: ReactNode;
  disabled?: boolean;
  value?: DropdownValue | undefined;
  onChange?: (value: DropdownValue, event?: Event) => void;
  defaultValue?: DropdownValue | undefined;
  readOnly?: boolean;
  modal?: boolean;
  className?: string;
};

const Dropdown = <DropdownValue extends number | string>({
  align = 'start',
  expand,
  placeholder = 'Select...',
  icon,
  title,
  value,
  onChange,
  addonSlot,
  options,
  action,
  disabled,
  modal = false,
  className
}: DropdownProps<DropdownValue>) => {
  const selected = options.find(o => o.value === value);

  return (
    <DropdownMenu.Root modal={modal}>
      <DropdownMenu.Trigger
        className={clsx(
          styles.trigger,
          addonSlot === 'left' && styles['pad-left'],
          addonSlot === 'right' && styles['pad-right'],
          expand === 'full' && styles.full
        )}
        disabled={disabled}
      >
        {icon && (
          <span className={styles.icon}>
            <Icon icon={icon} color="gray-300" />
          </span>
        )}
        {selected ? (
          <span className={styles.selected}>{selected.label}</span>
        ) : (
          <span className={styles.placeholder}>{placeholder}</span>
        )}
        <span className={styles.arrow}>
          <Icon icon={IconExpandMoreRound} />
        </span>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          className={cn(styles.content, className)}
          align={align}
        >
          {title && (
            <DropdownMenu.Label className={styles.label}>
              {title}
            </DropdownMenu.Label>
          )}
          <DropdownMenu.Group className={styles.group}>
            {options.map((option, index) => (
              <DropdownMenu.Item
                key={index}
                className={clsx(
                  styles.item,
                  value === option.value && styles['item-active'],
                  option.icon && styles['item-icon']
                )}
                onSelect={
                  onChange ? event => onChange(option.value!, event) : undefined
                }
              >
                {option.icon && <Icon icon={option.icon} />}
                {option.label}
                {option.description && (
                  <div className={styles['item-description']}>
                    {option.description}
                  </div>
                )}
              </DropdownMenu.Item>
            ))}
            {action && <div className={styles.action}>{action}</div>}
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
};

export default Dropdown;
