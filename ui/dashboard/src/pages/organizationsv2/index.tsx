import PageHeader from 'containers/page-header';
import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const OrganizationPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title="Organizations">
      <PageHeader
        title={t('organizations')}
        description={t('organization-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default OrganizationPage;
