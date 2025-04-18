import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const AuditLogsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('audit-logs')}>
      <PageHeader
        title={t('audit-logs')}
        description={t('audit-logs-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default AuditLogsPage;
