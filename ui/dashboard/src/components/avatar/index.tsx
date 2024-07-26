import type { FunctionComponent } from 'react';
import clsx from 'clsx';
import type { AvatarColor, Color } from '@types';
import { BG_COLOR_MAPPING, COLOR_MAPPING } from 'app/constants';
import styles from './styles.module.css';

export type AvatarSize = 'xxs' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';

export interface AvatarBaseProps {
  size?: AvatarSize;
  rounded?: 'medium' | 'full';
}

export interface AvatarIconProps extends AvatarBaseProps {
  icon: FunctionComponent;
  color?: Color;
  bgColor?: Color;
}

export const AvatarIcon = ({
  size = 'md',
  rounded = 'medium',
  icon: Icon,
  color = 'gray-500',
  bgColor = 'gray-100'
}: AvatarIconProps) => {
  return (
    <div
      className={clsx(
        'flex items-center justify-center',
        styles[`size-${size}`],
        COLOR_MAPPING[color],
        BG_COLOR_MAPPING[bgColor],
        rounded === 'full' && 'rounded-full'
      )}
    >
      <Icon />
    </div>
  );
};

export interface AvatarImageProps extends AvatarBaseProps {
  image: string;
}

export const AvatarImage = ({
  size = 'md',
  rounded = 'medium',
  image
}: AvatarImageProps) => {
  return (
    <img
      className={clsx(
        styles[`size-${size}`],
        'object-cover',
        rounded === 'full' && 'rounded-full'
      )}
      src={image}
    />
  );
};

export interface AvatarPlaceholderProps extends AvatarBaseProps {
  color?: AvatarColor;
}
