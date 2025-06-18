import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import EmptyState, { EmptyStateProps } from 'elements/empty-state';

const EmptyCollection = ({
  variant,
  onAdd
}: {
  variant: EmptyStateProps['variant'];
  onAdd?: () => void;
}) => {
  const { t } = useTranslation(['table', 'common']);

  const isEmpty = variant === 'no-data';

  return (
    <EmptyState.Root variant={variant} size="lg" className="mt-10">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>
          {t(`code-refs.${isEmpty ? 'empty' : 'enable'}`)}
        </EmptyState.Title>
        <EmptyState.Description>
          {t(`code-refs.${isEmpty ? 'empty' : 'enable'}-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
      {!isEmpty && onAdd && (
        <EmptyState.Actions>
          <EmptyState.ActionButton
            isNeedAdminAccess
            type={'new'}
            variant="primary"
            onClick={onAdd}
          >
            <IconAddOutlined />
            {t(`common:create-api-key`)}
          </EmptyState.ActionButton>
        </EmptyState.Actions>
      )}
    </EmptyState.Root>
  );
};

export default EmptyCollection;
