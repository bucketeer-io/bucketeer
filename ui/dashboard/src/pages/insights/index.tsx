import { useTranslation } from 'react-i18next';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const InsightsPage = () => {
  const { t } = useTranslation();
  return (
    <PageLayout.Root title={t('navigation.insights')}>
      <PageHeader
        title={t('navigation.insights')}
        description={t('insights.subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default InsightsPage;
