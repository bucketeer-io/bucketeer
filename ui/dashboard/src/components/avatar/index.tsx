import type { FunctionComponent, ImgHTMLAttributes } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { BG_COLOR_MAPPING, COLOR_MAPPING } from 'constants/styles';
import type { Color } from '@types';
import { cn } from 'utils/style';

// Avatar icon
const AvatarIconVariants = cva(['flex items-center justify-center'], {
  variants: {
    size: {
      lg: 'size-16 text-3xl',
      md: 'size-10 text-xl',
      sm: 'size-8 text-lg',
      xs: 'size-6 text-md'
    },
    rounded: {
      md: 'rounded-md',
      full: 'rounded-full'
    }
  },
  defaultVariants: {
    size: 'md',
    rounded: 'md'
  }
});

export interface AvatarIconProps
  extends VariantProps<typeof AvatarIconVariants> {
  className?: string;
  icon: FunctionComponent;
  color?: Color;
  bgColor?: Color;
}

const AvatarIcon = ({
  className,
  size,
  color = 'gray-500',
  bgColor = 'gray-100',
  icon: Icon,
  ...props
}: AvatarIconProps) => {
  return (
    <div
      className={cn(
        AvatarIconVariants({ size }),
        COLOR_MAPPING[color],
        BG_COLOR_MAPPING[bgColor],
        className
      )}
      {...props}
    >
      <Icon />
    </div>
  );
};

// Avatar image
const AvatarImageVariants = cva(['rounded-full object-cover'], {
  variants: {
    size: {
      xl: 'w-[120px] h-[120px] min-w-[120px]',
      lg: 'size-[60px] min-w-[60px]',
      md: 'size-8 min-w-8',
      sm: 'size-6 min-w-6'
    }
  },
  defaultVariants: {
    size: 'md'
  }
});

export interface AvatarImageProps
  extends ImgHTMLAttributes<HTMLImageElement>,
    VariantProps<typeof AvatarImageVariants> {
  image: string;
}

const AvatarImage = ({
  className,
  size,
  image,
  ...props
}: AvatarImageProps) => {
  return (
    <img
      src={image}
      className={cn(AvatarImageVariants({ size }), className)}
      {...props}
    />
  );
};

export { AvatarImage, AvatarIcon };
