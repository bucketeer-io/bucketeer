import { cva, VariantProps } from 'class-variance-authority';
import { cn } from 'utils/style';

const spinnerVariants = cva(
  ['animate-[spin_0.6s_linear_infinite] rounded-full border-t-primary-500'],
  {
    variants: {
      size: {
        sm: 'h-6 w-6 border-4',
        md: 'h-10 w-10 border-[6px]'
      }
    },
    defaultVariants: {
      size: 'sm'
    }
  }
);

export interface SpinnerProps extends VariantProps<typeof spinnerVariants> {
  className?: string;
  size?: 'sm' | 'md';
}

const Spinner = ({ className, size, ...props }: SpinnerProps) => {
  return (
    <div className={cn(spinnerVariants({ size }), className)} {...props} />
  );
};

export default Spinner;
