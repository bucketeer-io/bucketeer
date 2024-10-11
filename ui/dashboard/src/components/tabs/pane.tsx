import type { ReactNode } from 'react';
import clsx from 'clsx';
import * as RadixTabs from '@radix-ui/react-tabs';
import styles from './styles.module.css';
import { useTabContext } from './tabs-context';
import type { TabContainedValue } from './types';
import { getTabPanesContainerCls } from './utils';

export interface TabsPaneProps {
  className?: string;
  contained?: TabContainedValue;
  value: string;
  children: ReactNode;
}

const TabsPane = ({
  className,
  contained: contained1,
  value,
  children
}: TabsPaneProps) => {
  const { space, overflow, padded, contained: contained2 } = useTabContext();
  const contained = contained1 !== undefined ? contained1 : contained2;

  return (
    <RadixTabs.Content
      className={clsx(
        space === 'full' && 'min-h-0 flex-1',
        overflow === 'hidden' && 'overflow-hidden',
        padded && 'pt-6',
        getTabPanesContainerCls(contained),
        styles.pane,
        className
      )}
      value={value}
      forceMount
    >
      {children}
    </RadixTabs.Content>
  );
};

export default TabsPane;
