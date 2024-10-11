import type { ReactNode } from 'react';

export interface TabsPanesProps {
  className?: string;
  children: ReactNode | ReactNode[];
}

const TabsPanes = ({ className, children }: TabsPanesProps) => {
  return <div className={className}>{children}</div>;
};

export default TabsPanes;
