import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const UserSegmentsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('user-segments')}>
      <PageHeader
        title={t('user-segments')}
        description={t('user-segments-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default UserSegmentsPage;
