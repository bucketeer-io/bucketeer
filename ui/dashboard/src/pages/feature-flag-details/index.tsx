import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const FeatureFlagDetailsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('feature-flags')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default FeatureFlagDetailsPage;
