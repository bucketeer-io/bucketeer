import { ReactNode } from 'react';
import { cn } from 'utils/style';

const Card = ({
  children,
  className
}: {
  children: ReactNode;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        'flex flex-col w-full p-5 gap-y-5 bg-white rounded-lg shadow-card',
        className
      )}
    >
      {children}
    </div>
  );
};

export default Card;
