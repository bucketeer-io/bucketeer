import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useScreen } from 'hooks';
import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const AuditLogsPage = () => {
  const { t } = useTranslation(['common', 'form', 'table']);
  const { fromTabletScreen } = useScreen();
  return (
    <PageLayout.Root title={t('audit-logs')}>
      <PageHeader
        title={t('audit-logs')}
        description={t('audit-logs-subtitle')}
        link={!fromTabletScreen ? DOCUMENTATION_LINKS.AUDIT_LOGS : undefined}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default AuditLogsPage;
