import { ReactNode } from 'react';
import { cn } from 'utils/style';

const TableListContent = ({
  children,
  className
}: {
  children: ReactNode;
  className?: string;
}) => {
  return (
    <div className={cn('flex flex-col w-full min-w-[900px]', className)}>
      {children}
    </div>
  );
};

export default TableListContent;
