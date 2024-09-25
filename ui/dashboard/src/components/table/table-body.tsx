import { HTMLAttributes } from 'react';
import { TableCommonType } from './root';

const TableBody = ({
  children,
  className,
  ...props
}: TableCommonType & HTMLAttributes<HTMLTableSectionElement>) => {
  return (
    <tbody {...props} className={className}>
      {children}
    </tbody>
  );
};

export default TableBody;
