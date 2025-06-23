import { useTranslation } from 'i18n';
import EmptyState, { type EmptyStateProps } from 'elements/empty-state';

interface NoDataStateProps {
  size?: EmptyStateProps['size'];
  title?: string;
  description?: string;
  onAdd?: () => void;
}

export const NoDataState = ({
  size = 'md',
  title = '',
  description = '',
  onAdd
}: NoDataStateProps) => {
  const { t } = useTranslation(['common', 'message']);
  const defaultTitle = t('message:no-data');
  const defaultDescription = t('data-appear');

  return (
    <EmptyState.Root variant="no-data" size={size}>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{title || defaultTitle}</EmptyState.Title>
        {description && (
          <EmptyState.Description>
            {description || defaultDescription}
          </EmptyState.Description>
        )}
        {onAdd && (
          <EmptyState.Actions>
            <EmptyState.ActionButton onClick={onAdd}>
              {t(`add`)}
            </EmptyState.ActionButton>
          </EmptyState.Actions>
        )}
      </EmptyState.Body>
    </EmptyState.Root>
  );
};
