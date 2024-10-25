import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = () => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <EmptyState.Root variant="no-data" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.user-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.user-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
    </EmptyState.Root>
  );
};
