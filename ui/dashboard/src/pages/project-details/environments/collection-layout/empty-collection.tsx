import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = ({ onAdd }: { onAdd?: () => void }) => {
  const isProjectManagement = onAdd !== undefined;
  const { t } = useTranslation(['common', 'table']);

  return (
    <EmptyState.Root variant="no-data" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.env-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {isProjectManagement
            ? t(`table:empty.env-desc`)
            : t(`table:empty.env-org-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
      {isProjectManagement && (
        <EmptyState.Actions>
          <EmptyState.ActionButton
            isNeedAdminAccess
            type={'new'}
            variant="primary"
            onClick={onAdd}
          >
            <IconAddOutlined />
            {t(`new-env`)}
          </EmptyState.ActionButton>
        </EmptyState.Actions>
      )}
    </EmptyState.Root>
  );
};
