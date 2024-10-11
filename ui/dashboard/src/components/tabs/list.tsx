import type { ReactNode } from 'react';
import clsx from 'clsx';
import * as RadixTabs from '@radix-ui/react-tabs';
import styles from './styles.module.css';
import { useTabContext } from './tabs-context';
import type { TabContainedValue } from './types';
import { getTabListContainerCls } from './utils';

export interface TabsListProps {
  className?: string;
  contained?: TabContainedValue;
  children: ReactNode | ReactNode[];
}

const TabsList = ({
  className,
  contained: contained1,
  children
}: TabsListProps) => {
  const { contained: contained2 } = useTabContext();
  const contained = contained1 !== undefined ? contained1 : contained2;

  return (
    <RadixTabs.List aria-label="tabs" className={styles.list}>
      <div
        className={clsx(
          'flex flex-nowrap overflow-x-auto',
          className,
          getTabListContainerCls(contained)
        )}
      >
        {children}
      </div>
    </RadixTabs.List>
  );
};

export default TabsList;
