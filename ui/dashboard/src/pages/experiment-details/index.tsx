import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const ExperimentDetailsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('navigation.experiments')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default ExperimentDetailsPage;
