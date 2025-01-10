import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const GoalsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('goals')}>
      <PageHeader title={t('goals')} description={t('goals-subtitle')} />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default GoalsPage;
