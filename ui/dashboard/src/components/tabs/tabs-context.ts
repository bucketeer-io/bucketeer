import { createContext, useContext } from 'react';
import type { TabContainedValue } from './types';

export type TabContextValue = {
  space?: 'full';
  overflow?: 'hidden';
  padded?: boolean;
  contained?: TabContainedValue;
};

const tabContextDefaultValue = {
  space: undefined,
  overflow: undefined,
  padded: true,
  contained: false
};

export const TabContext = createContext<TabContextValue>(
  tabContextDefaultValue
);

export const useTabContext = () => {
  const tabContext = useContext(TabContext);
  return tabContext;
};
