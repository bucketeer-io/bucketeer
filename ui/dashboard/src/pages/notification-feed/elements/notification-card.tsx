import { ReactNode } from 'react';
import { cn } from 'utils/style';
import Button from 'components/button';

interface NotificationCardProps {
  active?: boolean;
  onClick?: () => void;
  header: ReactNode;
  children?: ReactNode;
  footer: ReactNode;
  bordered?: boolean;
}

const NotificationCard = ({
  active,
  onClick,
  header,
  children,
  footer,
  bordered = true
}: NotificationCardProps) => (
  <Button
    type="button"
    variant="text"
    onClick={onClick}
    className={cn(
      'flex h-full w-full flex-col items-start justify-start gap-2 whitespace-normal text-left transition-colors',
      bordered
        ? cn(
            'rounded-lg border p-4',
            active
              ? 'border-primary-500 shadow-border-primary-500'
              : 'border-gray-200 dark:border-dark-black-700 hover:border-gray-300 dark:hover:border-dark-black-600'
          )
        : 'p-0'
    )}
  >
    {header}
    {children && (
      <div className="flex flex-wrap items-center gap-2">{children}</div>
    )}
    <div className="mt-auto w-full">{footer}</div>
  </Button>
);

export default NotificationCard;
