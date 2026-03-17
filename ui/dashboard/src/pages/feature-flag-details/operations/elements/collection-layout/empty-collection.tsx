import { useTranslation } from 'i18n';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = () => {
  const { t } = useTranslation(['table']);

  return (
    <EmptyState.Root variant="no-data" size="lg" className="py-20 sm:pt-52">
      <div className="size-16">
        <EmptyState.Illustration />
      </div>
      <EmptyState.Body>
        <EmptyState.Title>
          {t(`table:empty.flag-operations-title`)}
        </EmptyState.Title>
        <EmptyState.Description>
          {t(`table:empty.flag-operations-desc`)}
        </EmptyState.Description>
      </EmptyState.Body>
    </EmptyState.Root>
  );
};
