import { Children, useCallback } from 'react';
import type { ReactNode, ReactElement } from 'react';
import clsx from 'clsx';
import * as RadixTabs from '@radix-ui/react-tabs';
import TabsList from './list';
import TabsPane from './pane';
import TabsPanes from './panes';
import { TabContext } from './tabs-context';
import TabsTrigger from './trigger';
import type { TabContainedValue } from './types';

export type TabsProps<TabValue> = {
  className?: string;
  space?: 'full';
  overflow?: 'hidden';
  padded?: boolean;
  contained?: TabContainedValue;
  children: ReactNode[];
} & (
  | {
      defaultValue: TabValue;
      value?: never;
      onChange?: never;
    }
  | {
      defaultValue?: never;
      value: TabValue | undefined;
      onChange: (tab: TabValue) => void;
    }
);

const Tabs = <TabValue extends string>({
  className,
  space,
  overflow,
  padded = true,
  contained,
  defaultValue,
  value,
  onChange: _onChange,
  children
}: TabsProps<TabValue>) => {
  const childs = Children.toArray(children) as ReactElement[];
  const tabList = childs.find(item => item.type === TabsList);
  const tabPanes = childs.filter(item => item.type === TabsPanes);
  const tabPaneItems = childs.filter(item => item.type === TabsPane);
  const onChange = useCallback((v: string) => {
    if (!_onChange) return;
    _onChange(v as TabValue);
  }, []);

  return (
    <TabContext.Provider value={{ space, overflow, padded, contained }}>
      <RadixTabs.Root
        className={clsx(
          space === 'full' ? 'flex h-full flex-col' : '',
          className
        )}
        defaultValue={defaultValue}
        value={value}
        onValueChange={onChange}
        orientation="vertical"
      >
        {tabList}
        {tabPanes.length ? tabPanes : tabPaneItems}
      </RadixTabs.Root>
    </TabContext.Provider>
  );
};

Tabs.List = TabsList;
Tabs.Panes = TabsPanes;
Tabs.Pane = TabsPane;
Tabs.Trigger = TabsTrigger;

export default Tabs;
