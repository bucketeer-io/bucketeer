import type { HTMLAttributes } from 'react';
import clsx from 'clsx';

type DividerProps = {
  width?: number;
  vertical?: boolean;
  dashed?: boolean;
  darker?: boolean;
} & HTMLAttributes<HTMLDivElement>;

const Divider = ({
  className,
  vertical,
  dashed,
  darker,
  width
}: DividerProps) => {
  return (
    <div
      role="separator"
      className={clsx(
        vertical ? 'h-full w-px border-l' : 'border-t',
        darker
          ? 'border-gray-600 dark:border-dark-black-700'
          : 'border-gray-200 dark:border-dark-black-700',
        className,
        dashed && 'border-dashed'
      )}
      style={{ width }}
    />
  );
};

export default Divider;
