import { forwardRef, ReactNode, Ref } from 'react';
import { cn } from 'utils/style';

const Card = forwardRef(
  (
    {
      children,
      className
    }: {
      children: ReactNode;
      className?: string;
    },
    ref: Ref<HTMLDivElement>
  ) => {
    return (
      <div
        ref={ref}
        className={cn(
          'flex flex-col w-full p-5 gap-y-5 bg-white rounded-lg shadow-card',
          className
        )}
      >
        {children}
      </div>
    );
  }
);

export default Card;
