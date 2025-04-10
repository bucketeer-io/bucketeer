import { ReactNode } from 'react';
import { cn } from 'utils/style';

const TableListContainer = ({
  children,
  className
}: {
  children: ReactNode;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        'flex flex-col flex-1 w-full p-6 pt-0 mt-5 overflow-y-hidden overflow-x-auto hidden-scroll',
        className
      )}
    >
      {children}
    </div>
  );
};

export default TableListContainer;
