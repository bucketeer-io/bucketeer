import { HTMLAttributes } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';

const styles = cva(['size-4 rotate-45 rounded-sm border border-white'], {
  variants: {
    variant: {
      blue: ['bg-accent-blue-500'],
      pink: ['bg-accent-pink-500'],
      green: ['bg-accent-green-500']
    }
  },
  defaultVariants: {
    variant: 'blue'
  }
});

export type VariationProps = HTMLAttributes<HTMLDivElement> & {
  text?: string;
  variant?: 'blue' | 'pink' | 'green';
};

const Variation = ({
  text,
  className,
  variant = 'blue',
  style
}: VariationProps) => {
  return (
    <div className={cn('flex items-center gap-[6px]', className)} style={style}>
      <div className={cn(styles({ variant }))}></div>
      {text && (
        <p className="text-gray-700 typo-para-small truncate max-w-20">
          {text}
        </p>
      )}
    </div>
  );
};

export default Variation;
