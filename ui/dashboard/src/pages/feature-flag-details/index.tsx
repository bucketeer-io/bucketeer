import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';

const FeatureFlagDetailsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('feature-flags')}>
      Feature Flag Details
    </PageLayout.Root>
  );
};

export default FeatureFlagDetailsPage;
