import { ErrorComponentProps } from '@tanstack/react-router';
import { handleGetErrorMessage } from 'utils/function';
import EmptyState, { type EmptyStateProps } from 'elements/empty-state';

interface ErrorStateProps {
  size?: EmptyStateProps['size'];
  title?: string;
  description?: string;
  error?: ErrorComponentProps;
  onRetry?: () => void;
}

export const ErrorState = ({
  size = 'lg',
  title = `Oops! Something went wrong`,
  description = `We're on it. Please try again later.`,
  error,
  onRetry
}: ErrorStateProps) => {
  const _error = error ? handleGetErrorMessage(error) : null;

  return (
    <EmptyState.Root variant="error" size={size}>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{title}</EmptyState.Title>
        <EmptyState.Description>
          {_error?.message || description}
        </EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        {(onRetry || _error?.reset) && (
          <EmptyState.ActionButton
            variant="primary"
            onClick={_error?.reset || onRetry}
          >
            {`Retry`}
          </EmptyState.ActionButton>
        )}
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
