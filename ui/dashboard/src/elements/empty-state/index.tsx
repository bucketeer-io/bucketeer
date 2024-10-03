import type { PropsWithChildren, ReactNode } from 'react';
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
      <div>{children}</div>
    </EmptyStateProvider>
  );
};

const EmptyStateIllustration = () => {
  const { variant } = useEmptyState();

  switch (variant) {
    case 'error':
      return <img alt="Error" src="/assets/empty-state-error.svg" />;

    case 'no-data':
      return <img alt="No Data" src="/assets/empty-state-no-data.svg" />;

    case 'no-search':
      return <img alt="No Search" src="/assets/empty-state-no-search.svg" />;
  }
};

const EmptyStateBody = ({ children }: PropsWithChildren) => {
  return <div className="max-w-80 text-center mx-auto mt-10">{children}</div>;
};

const EmptyStateTitle = ({ children }: { children: string }) => {
  return <div>{children}</div>;
};

const EmptyStateDescription = ({ children }: { children: string }) => {
  return <div>{children}</div>;
};

const EmptyStateActions = ({ children }: PropsWithChildren) => {
  return <div className="flex justify-center mt-10 gap-3">{children}</div>;
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
