import type { PropsWithChildren, ReactNode } from 'react';
import emptyStateError from 'assets/empty-state/error.svg';
import emptyStateNoData from 'assets/empty-state/no-data.svg';
import emptyStateNoSearch from 'assets/empty-state/no-search.svg';
import { createContext } from 'utils/create-context';
import Button, { type ButtonProps } from 'components/button';

export interface EmptyStateProps {
  variant: 'error' | 'no-data' | 'no-search';
  size: 'sm' | 'md' | 'lg';
  children: ReactNode;
}

type EmptyStateContextValue = Omit<EmptyStateProps, 'children'>;

const [EmptyStateProvider, useEmptyState] =
  createContext<EmptyStateContextValue>({
    name: 'EmptyStateProvider',
    errorMessage: `useEmptyState returned is 'undefined'. Seems you forgot to wrap the components in "<EmptyState.Root />" `
  });

const EmptyStateRoot = ({ variant, size, children }: EmptyStateProps) => {
  return (
    <EmptyStateProvider value={{ variant, size }}>
      <div className="flex flex-col items-center gap-4">{children}</div>
    </EmptyStateProvider>
  );
};

const EmptyStateIllustration = () => {
  const { variant } = useEmptyState();

  switch (variant) {
    case 'error':
      return <img alt="Error" className="w-fit" src={emptyStateError} />;

    case 'no-data':
      return <img alt="No Data" className="w-fit" src={emptyStateNoData} />;

    case 'no-search':
      return <img alt="No Search" className="w-fit" src={emptyStateNoSearch} />;
  }
};

const EmptyStateBody = ({ children }: PropsWithChildren) => {
  return (
    <div className="max-w-[350px] flex flex-col gap-2 text-center mx-auto">
      {children}
    </div>
  );
};

const EmptyStateTitle = ({ children }: { children: string }) => {
  return <div className="text-gray-900 typo-head-bold-medium">{children}</div>;
};

const EmptyStateDescription = ({ children }: { children: string }) => {
  return <div className="text-gray-600 typo-para-small">{children}</div>;
};

const EmptyStateActions = ({ children }: PropsWithChildren) => {
  return <div className="flex justify-center mt-3 gap-3">{children}</div>;
};

const EmptyStateActionButton = (props: Omit<ButtonProps, 'size'>) => {
  const { size } = useEmptyState();
  return (
    <Button variant="primary" size={size === 'lg' ? 'md' : 'sm'} {...props} />
  );
};

const Root = EmptyStateRoot;
const Illustration = EmptyStateIllustration;
const Body = EmptyStateBody;
const Actions = EmptyStateActions;
const Title = EmptyStateTitle;
const Description = EmptyStateDescription;
const ActionButton = EmptyStateActionButton;

const EmptyState = {
  Root,
  Illustration,
  Body,
  Actions,
  Title,
  Description,
  ActionButton
};

export default EmptyState;
