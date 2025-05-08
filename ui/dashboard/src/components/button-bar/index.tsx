import { ReactNode } from 'react';
import { cn } from 'utils/style';

export const ButtonBar = ({
  primaryButton,
  secondaryButton,
  className
}: {
  primaryButton: ReactNode;
  secondaryButton?: ReactNode;
  className?: string;
}) => {
  return (
    <div className={cn('p-5 border-t border-gray-200 w-full', className)}>
      <div className="flex items-center gap-4 justify-end">
        {primaryButton}
        {secondaryButton}
      </div>
    </div>
  );
};
