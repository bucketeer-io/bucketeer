import { type PropsWithChildren, type ReactNode, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { QueryErrorResetBoundary } from '@tanstack/react-query';
import Spinner from 'components/spinner';
import { ErrorState } from '../empty-state/error';
import { PageLayoutProvider } from './context';

export interface PageLayoutProps {
  title: string;
  showNavBar?: boolean;
  noNavLinks?: boolean;
  totalSteps?: number;
  initialStep?: number;
  children: ReactNode;
}

const PageLayoutRoot = ({
  title,
  totalSteps,
  initialStep,
  children
}: PageLayoutProps) => {
  const [step, setStep] = useState<number | undefined>(
    initialStep || undefined
  );

  return (
    <PageLayoutProvider
      value={{ title, totalSteps, step, onChangeStep: setStep }}
    >
      {/* <Flex flexDir='column' h='screen' bg='bg.primary'> */}
      {/* {showNavBar && (
					<Box pos='relative' flex='initial' zIndex={10}>
						<PageNavBar noNavLinks={noNavLinks} />
					</Box>
				)} */}

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
      {/* </Flex> */}
    </PageLayoutProvider>
  );
};

const PageLayoutLoadingState = () => {
  return (
    <div className="h-full flex-center">
      <Spinner />
    </div>
  );
};

const PageLayoutErrorState = ({ onRetry }: { onRetry?: () => void }) => {
  return (
    <div className="h-full flex-center">
      <ErrorState onRetry={onRetry} />
    </div>
  );
};

const PageLayoutEmptyState = ({ children }: PropsWithChildren) => {
  return <div className="h-full flex-center">{children}</div>;
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
