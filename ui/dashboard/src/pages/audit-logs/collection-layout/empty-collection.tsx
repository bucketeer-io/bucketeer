import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = () => {
  const { t } = useTranslation(['table']);

  return (
    <EmptyState.Root variant="no-data" size="lg" className='pt-60'>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{t(`table:empty.audit-logs-title`)}</EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.audit-logs-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
    </EmptyState.Root>
  );
};
