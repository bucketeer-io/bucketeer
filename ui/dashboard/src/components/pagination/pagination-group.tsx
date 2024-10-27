import { PropsWithChildren } from 'react';

const PaginationGroup = ({ children }: PropsWithChildren) => {
  return <div className="flex gap-2">{children}</div>;
};

export default PaginationGroup;
