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
  return <div className="p-6 border-b border-gray-200">{children}</div>;
};

const PageLayoutContent = ({ children }: PropsWithChildren) => {
  return <div className="p-6 flex flex-1 flex-col h-full">{children}</div>;
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
