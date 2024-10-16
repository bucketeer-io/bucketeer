import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const OrganizationDetails = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('organizations')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default OrganizationDetails;
