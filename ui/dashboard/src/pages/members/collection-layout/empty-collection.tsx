import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = ({ onAdd }: { onAdd?: () => void }) => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <EmptyState.Root variant="no-data" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.member-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.member-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        {onAdd && (
          <EmptyState.ActionButton
            type={'new'}
            variant="primary"
            onClick={onAdd}
          >
            <IconAddOutlined />
            {t(`invite-member`)}
          </EmptyState.ActionButton>
        )}
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
