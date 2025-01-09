import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <EmptyState.Root variant="no-data" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.push-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.push-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        <EmptyState.ActionButton variant="primary" onClick={onAdd}>
          <IconAddOutlined />
          {t(`new-push`)}
        </EmptyState.ActionButton>
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
