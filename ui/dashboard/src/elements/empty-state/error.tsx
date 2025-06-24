import { translation } from 'constants/message';
import EmptyState, { type EmptyStateProps } from 'elements/empty-state';

interface ErrorStateProps {
  size?: EmptyStateProps['size'];
  title?: string;
  description?: string;
  onRetry?: () => void;
}

export const ErrorState = ({
  size = 'lg',
  title = translation('message:something-went-wrong'),
  description = translation('message:try-again-later'),
  onRetry
}: ErrorStateProps) => {
  return (
    <EmptyState.Root variant="error" size={size}>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{title}</EmptyState.Title>
        <EmptyState.Description>{description}</EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        {onRetry && (
          <EmptyState.ActionButton variant="primary" onClick={onRetry}>
            {`Retry`}
          </EmptyState.ActionButton>
        )}
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
