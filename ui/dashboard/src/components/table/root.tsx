import { ReactNode, TableHTMLAttributes } from 'react';
import { cn } from 'utils/style';

export type TableCommonType = {
  children?: ReactNode;
  className?: string;
};

export type TableProps = TableCommonType &
  TableHTMLAttributes<HTMLTableElement>;
const TableRoot = ({ children, className, ...props }: TableProps) => {
  return (
    <table
      className={cn(
        'border-separate border-spacing-y-3 w-full mb-6',
        className
      )}
      {...props}
    >
      {children}
    </table>
  );
};

export default TableRoot;
