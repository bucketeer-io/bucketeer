import { memo, ReactNode } from 'react';
import { cn } from 'utils/style';

const TableListContent = memo(
  ({ children, className }: { children: ReactNode; className?: string }) => {
    return (
      <div className={cn('flex flex-col w-full', className)}>{children}</div>
    );
  }
);

export default TableListContent;
