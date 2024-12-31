import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const PushesPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('fcm')}>
      <PageHeader title={t('fcm')} description={t('fcm-subtitle')} />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default PushesPage;
