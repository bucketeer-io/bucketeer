import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <EmptyState.Root variant="no-data" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.experiment-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.experiment-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        <EmptyState.ActionButton type={'new'} variant="primary" onClick={onAdd}>
          <IconAddOutlined />
          {t(`new-experiment`)}
        </EmptyState.ActionButton>
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
