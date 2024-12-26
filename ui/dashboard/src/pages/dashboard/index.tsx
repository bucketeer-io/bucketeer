import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const DashboardPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('dashboard')}>
      <PageHeader
        title={t('dashboard')}
        description={t('dashboard-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default DashboardPage;
