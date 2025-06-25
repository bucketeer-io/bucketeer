import { useTranslation } from 'i18n';
import EmptyState, { type EmptyStateProps } from 'elements/empty-state';

interface ErrorStateProps {
  size?: EmptyStateProps['size'];
  title?: string;
  description?: string;
  onRetry?: () => void;
}

export const ErrorState = ({
  size = 'lg',
  title = '',
  description = '',
  onRetry
}: ErrorStateProps) => {
  const { t } = useTranslation(['message', 'table']);
  const defaultTitle = t('message:something-went-wrong');
  const defaultDescription = t('message:try-again-later');
  return (
    <EmptyState.Root variant="error" size={size}>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{title || defaultTitle}</EmptyState.Title>
        <EmptyState.Description>
          {description || defaultDescription}
        </EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        {onRetry && (
          <EmptyState.ActionButton variant="primary" onClick={onRetry}>
            {t(`table:retry`)}
          </EmptyState.ActionButton>
        )}
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
