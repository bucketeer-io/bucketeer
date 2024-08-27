import { cva, VariantProps } from 'class-variance-authority';
import { cn } from 'utils/style';

const spinnerVariants = cva(
  ['relative inline-block h-12 w-12 rounded-full border-light-100'],
  {
    variants: {
      size: {
        sm: 'h-6 w-6 border-4',
        md: 'h-12 w-12 border-8'
      },
      indicator: {
        sm: '-top-1 -left-1 h-6 w-6 border-4',
        md: '-top-2 -left-2 h-12 w-12 border-8'
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
    <div className={cn(spinnerVariants({ size }), className)} {...props}>
      <div
        className={cn(
          spinnerVariants({ indicator: size }),
          'absolute rounded-full animate-[spin_0.6s_linear_infinite] border-r-transparent',
          'border-b-transparent border-l-transparent border-t-primary-500'
        )}
      />
    </div>
  );
};

export default Spinner;
