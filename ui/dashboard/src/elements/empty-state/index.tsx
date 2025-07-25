import type { PropsWithChildren, ReactNode } from 'react';
import emptyStateCode from 'assets/empty-state/code.svg';
import emptyStateError from 'assets/empty-state/error.svg';
import emptyStateNoData from 'assets/empty-state/no-data.svg';
import emptyStateNoSearch from 'assets/empty-state/no-search.svg';
import { hasEditable, useAuth } from 'auth';
import { createContext } from 'utils/create-context';
import { cn } from 'utils/style';
import Button, { type ButtonProps } from 'components/button';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

export interface EmptyStateProps {
  variant: 'error' | 'no-data' | 'no-search' | 'invalid';
  size: 'sm' | 'md' | 'lg';
  children: ReactNode;
  className?: string;
}

export interface EmptyStateActionButtonProps
  extends Omit<ButtonProps, 'size' | 'type'> {
  type?: 'retry' | 'new';
  isNeedAdminAccess?: boolean;
}

type EmptyStateContextValue = Omit<EmptyStateProps, 'children'>;

const [EmptyStateProvider, useEmptyState] =
  createContext<EmptyStateContextValue>({
    name: 'EmptyStateProvider',
    errorMessage: `useEmptyState returned is 'undefined'. Seems you forgot to wrap the components in "<EmptyState.Root />" `
  });

const EmptyStateRoot = ({
  variant,
  size,
  children,
  className
}: EmptyStateProps) => {
  return (
    <EmptyStateProvider value={{ variant, size }}>
      <div className={cn('flex flex-col items-center gap-4', className)}>
        {children}
      </div>
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

    case 'invalid':
      return <img alt="Invalid" className="w-fit" src={emptyStateCode} />;
  }
};

const EmptyStateBody = ({ children }: PropsWithChildren) => {
  return (
    <div className="max-w-[380px] flex flex-col gap-2 text-center mx-auto">
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

const EmptyStateActionButton = ({
  type = 'retry',
  isNeedAdminAccess = false,
  ...props
}: EmptyStateActionButtonProps) => {
  const { size } = useEmptyState();
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);
  const isOrganizationAdmin =
    consoleAccount?.organizationRole === 'Organization_ADMIN' ||
    consoleAccount?.organizationRole === 'Organization_OWNER';
  const isRetry = type === 'retry';

  return (
    <DisabledButtonTooltip
      align="center"
      type={isNeedAdminAccess && !isOrganizationAdmin ? 'admin' : 'editor'}
      hidden={
        (editable && (isNeedAdminAccess ? isOrganizationAdmin : true)) ||
        isRetry
      }
      trigger={
        <Button
          variant="primary"
          size={size === 'lg' ? 'md' : 'sm'}
          disabled={
            (!editable || (isNeedAdminAccess ? !isOrganizationAdmin : false)) &&
            !isRetry
          }
          {...props}
        />
      }
    />
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
