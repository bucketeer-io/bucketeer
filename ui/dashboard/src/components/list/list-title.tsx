import { cn } from 'utils/style';

export type ListTitleProps = {
  children: string;
  className?: string;
};

const ListTitle = ({ children, className }: ListTitleProps) => {
  return (
    <h3 className={cn('typo-head-bold-medium text-gray-700', className)}>
      {children}
    </h3>
  );
};

export default ListTitle;
