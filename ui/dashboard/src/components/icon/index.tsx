import type { FunctionComponent } from 'react';
import clsx from 'clsx';
import { COLOR_MAPPING, ICON_SIZE_MAPPING } from 'constants/styles';
import type { Color, IconSize } from '@types';

export interface IconProps {
  color?: Color;
  size?: IconSize | IconSize[];
  icon: FunctionComponent;
}

const getSizeCls = (size: IconSize | IconSize[]) => {
  const arr = Array.isArray(size) ? size : [size];
  const cls = arr.map(item => ICON_SIZE_MAPPING[item]);
  return cls.join(' ');
};

const Icon = ({ color, size = 'md', icon: SvgIcon }: IconProps) => {
  return (
    <i
      className={clsx(
        'inline-flex',
        color && COLOR_MAPPING[color],
        getSizeCls(size)
      )}
    >
      <SvgIcon />
    </i>
  );
};

export default Icon;
