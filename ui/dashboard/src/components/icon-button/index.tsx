import { forwardRef } from 'react';
import type { ButtonHTMLAttributes, FunctionComponent, Ref } from 'react';
import { Link } from 'react-router-dom';
import clsx from 'clsx';
import styles from './styles.module.css';

type IconButtonSize = 'xs' | 'sm' | 'md';

export type IconButtonProps = Omit<
  ButtonHTMLAttributes<HTMLButtonElement>,
  'className'
> & {
  variant?: 'primary' | 'secondary' | 'gray' | 'icon' | 'icon-2';
  size?: IconButtonSize;
  href?: string;
  icon: FunctionComponent;
};

type IconButtonRef = Ref<HTMLButtonElement>;

const IconButton = forwardRef(
  (
    { variant = 'primary', href, ...props }: IconButtonProps,
    ref: IconButtonRef
  ) => {
    const { size = 'md', icon: SvgIcon, ...otherProps } = props;

    const cls = clsx(
      `${styles.button} ${styles[variant]} ${styles[`size-${size}`]}`
    );
    const content = (
      <i className={styles.svg}>
        <SvgIcon />
      </i>
    );

    if (href) {
      return (
        <Link className={cls} to={href}>
          {content}
        </Link>
      );
    }

    return (
      <button {...otherProps} ref={ref} className={cls}>
        {content}
      </button>
    );
  }
);

export default IconButton;
