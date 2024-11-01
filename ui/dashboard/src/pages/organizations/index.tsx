import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const OrganizationsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('organizations')}>
      <PageHeader
        title={t('organizations')}
        description={t('organization-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default OrganizationsPage;
