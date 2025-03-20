import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const FeatureFlagsPage = () => {
  const { t } = useTranslation(['common']);
  return (
    <PageLayout.Root title={t('feature-flags')}>
      <PageHeader
        title={t('feature-flags')}
        description={t('feature-flags-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default FeatureFlagsPage;
