import clsx from 'clsx';
import * as CheckboxPrimitive from '@radix-ui/react-checkbox';
import { v4 as uuid } from 'uuid';
import { IconChecked } from '@icons';
import Icon from 'components/icon';
import styles from './styles.module.css';

type CheckboxProps = CheckboxPrimitive.CheckboxProps & {
  title?: string;
  description?: string;
  expand?: 'full';
  reverse?: boolean;
};

const Checkbox = ({
  checked,
  title,
  description,
  expand,
  reverse,
  onCheckedChange,
  ...props
}: CheckboxProps) => {
  return (
    <div
      className={clsx(
        styles.wrapper,
        expand === 'full' && styles.full,
        reverse && styles.reverse
      )}
    >
      <div className={styles['checkbox-wrapper']}>
        <CheckboxPrimitive.Root
          {...props}
          className={clsx(styles.checkbox, checked && styles.checked)}
          checked={checked}
          id={uuid()}
          onCheckedChange={onCheckedChange}
        >
          <CheckboxPrimitive.Indicator
            className={styles.indicator}
            forceMount={true}
          >
            <Icon icon={IconChecked} />
          </CheckboxPrimitive.Indicator>
        </CheckboxPrimitive.Root>
      </div>
      {title && (
        <label className={styles.label} htmlFor={uuid()}>
          <span className={styles.title}>{title}</span>
          {description && (
            <span className={styles.description}>{description}</span>
          )}
        </label>
      )}
    </div>
  );
};

export default Checkbox;
