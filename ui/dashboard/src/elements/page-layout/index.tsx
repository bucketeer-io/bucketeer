import { type PropsWithChildren, type ReactNode } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { QueryErrorResetBoundary } from '@tanstack/react-query';
import Spinner from 'components/spinner';
import { ErrorState } from '../empty-state/error';
import { PageLayoutProvider } from './context';

export interface PageLayoutProps {
  title: string;
  children: ReactNode;
}

const PageLayoutRoot = ({ title, children }: PageLayoutProps) => {
  return (
    <PageLayoutProvider value={{ title }}>
      <div className="flex flex-col min-h-screen">
        <QueryErrorResetBoundary>
          {({ reset }) => (
            <ErrorBoundary
              fallbackRender={({ resetErrorBoundary }) => (
                <PageLayout.ErrorState onRetry={resetErrorBoundary} />
              )}
              onReset={reset}
            >
              {children}
            </ErrorBoundary>
          )}
        </QueryErrorResetBoundary>
      </div>
    </PageLayoutProvider>
  );
};

const PageLayoutLoadingState = () => {
  return (
    <div className="w-full flex-center py-20">
      <Spinner />
    </div>
  );
};

const PageLayoutErrorState = ({ onRetry }: { onRetry?: () => void }) => {
  return (
    <div className="h-full flex-grow flex-center">
      <ErrorState onRetry={onRetry} />
    </div>
  );
};

const PageLayoutEmptyState = ({ children }: PropsWithChildren) => {
  return <div className="h-full flex-grow flex-center">{children}</div>;
};

const PageLayoutHeader = ({ children }: PropsWithChildren) => {
  return <div>{children}</div>;
};

const PageLayoutContent = ({ children }: PropsWithChildren) => {
  return <div className="h-full flex flex-col">{children}</div>;
};

const PageLayoutBody = ({
  children,
  withContainer = true
}: PropsWithChildren & { withContainer?: boolean }) => {
  if (withContainer) {
    return <div className="flex-1 container">{children}</div>;
  }

  return <div className="flex-1 pb-4">{children}</div>;
};

const PageLayoutFooter = ({ children }: PropsWithChildren) => {
  return <div className="flex-initial">{children}</div>;
};

const PageLayout = {
  Root: PageLayoutRoot,
  Header: PageLayoutHeader,
  Content: PageLayoutContent,
  Body: PageLayoutBody,
  Footer: PageLayoutFooter,
  LoadingState: PageLayoutLoadingState,
  ErrorState: PageLayoutErrorState,
  EmptyState: PageLayoutEmptyState
};

export default PageLayout;
