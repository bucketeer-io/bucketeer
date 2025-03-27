import { type PropsWithChildren, type ReactNode } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { QueryErrorResetBoundary } from '@tanstack/react-query';
import { cn } from 'utils/style';
import Spinner from 'components/spinner';
import { ErrorState } from '../empty-state/error';
import { PageLayoutProvider } from './context';

export interface PageLayoutProps {
  title: string;
  children: ReactNode;
}

export interface PageLayoutErrorState {
  onRetry?: () => void;
  className?: string;
}

export interface PageLayoutContentProps {
  children: ReactNode;
  className?: string;
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

const PageLayoutLoadingState = ({ className }: { className?: string }) => {
  return (
    <div className={cn('w-full flex-center py-20', className)}>
      <Spinner />
    </div>
  );
};

const PageLayoutErrorState = ({ onRetry, className }: PageLayoutErrorState) => {
  return (
    <div className={cn('h-full flex-grow flex-center', className)}>
      <ErrorState onRetry={onRetry} />
    </div>
  );
};

const PageLayoutEmptyState = ({ children }: PropsWithChildren) => {
  return <div className="h-full flex-grow flex-center">{children}</div>;
};

const PageLayoutHeader = ({ children }: PropsWithChildren) => {
  return <div className="p-6 border-b border-gray-200">{children}</div>;
};

const PageLayoutContent = ({ children, className }: PageLayoutContentProps) => {
  return (
    <div
      className={cn(
        'p-6 flex flex-1 flex-col h-full w-fit min-w-full',
        className
      )}
    >
      {children}
    </div>
  );
};

const PageLayout = {
  Root: PageLayoutRoot,
  Header: PageLayoutHeader,
  Content: PageLayoutContent,
  LoadingState: PageLayoutLoadingState,
  ErrorState: PageLayoutErrorState,
  EmptyState: PageLayoutEmptyState
};

export default PageLayout;
