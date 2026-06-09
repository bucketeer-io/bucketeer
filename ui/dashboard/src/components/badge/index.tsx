import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from 'utils/style';

const badgeVariants = cva(
  [
    'inline-flex justify-center items-center',
    'typo-para-small rounded-full w-5 h-5'
  ],
  {
    variants: {
      variant: {
        primary:
          'text-primary-500 bg-primary-100 dark:text-dark-purple-500 dark:bg-dark-purple-100',
        secondary:
          'text-gray-500 bg-gray-200 dark:text-dark-gray-200 dark:bg-dark-black-700'
      }
    },
    defaultVariants: {
      variant: 'primary'
    }
  }
);

export interface BadgeProps extends VariantProps<typeof badgeVariants> {
  className?: string;
  children: string | number;
}

const Badge = ({ className, variant, ...props }: BadgeProps) => {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
};

export { Badge };
