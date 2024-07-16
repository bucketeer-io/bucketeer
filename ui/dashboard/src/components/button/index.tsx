import { forwardRef } from 'react';
import type { ButtonHTMLAttributes, FunctionComponent, Ref } from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

export type ButtonProps = Omit<
  ButtonHTMLAttributes<HTMLButtonElement>,
  'className'
> & {
  variant?:
    | 'primary'
    | 'secondary'
    | 'secondary-2'
    | 'negative'
    | 'grey'
    | 'text';
  size?: 'xs' | 'sm' | 'md' | 'lg';
  icon?: FunctionComponent;
  iconSlot?: 'left' | 'right';
  expand?: 'full';
  width?: number;
  loading?: boolean;
};

type ButtonRef = Ref<HTMLButtonElement>;

const Button = forwardRef(
  (
    {
      variant = 'primary',
      size = 'md',
      icon: SvgIcon,
      iconSlot = 'right',
      expand,
      width,
      loading,
      disabled,
      onClick,
      children,
      ...otherProps
    }: ButtonProps,
    ref: ButtonRef
  ) => {
    return (
      <button
        {...otherProps}
        ref={ref}
        className={clsx(
          styles.btn,
          styles[`btn-${variant}`],
          styles[`size-${size}`],
          expand === 'full' && styles.fluid,
          loading && styles.loading
        )}
        style={{ width }}
        disabled={loading || disabled}
        onClick={onClick}
      >
        <div className={styles.content}>
          {SvgIcon && iconSlot === 'left' && (
            <i className={clsx(styles.icon, styles['icon-left'])}>
              <SvgIcon />
            </i>
          )}
          {children}
          {SvgIcon && iconSlot === 'right' && (
            <i className={clsx(styles.icon, styles['icon-right'])}>
              <SvgIcon />
            </i>
          )}
        </div>
      </button>
    );
  }
);

export default Button;
