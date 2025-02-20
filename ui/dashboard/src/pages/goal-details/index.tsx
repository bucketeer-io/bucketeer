import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const GoalDetailsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('goals')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default GoalDetailsPage;
