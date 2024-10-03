import { createContext } from 'utils/create-context';
import type { PageLayoutProps } from '.';

type PageLayoutContextValue = Pick<PageLayoutProps, 'title' | 'totalSteps'> & {
  step?: number;
  onChangeStep?: (step: number) => void;
};
export const [PageLayoutProvider, usePageLayout] =
  createContext<PageLayoutContextValue>({
    name: 'PageLayoutProvider',
    errorMessage: `usePageLayout returned is 'undefined'. Seems you forgot to wrap the components in "<PageLayout.Root />" `
  });
