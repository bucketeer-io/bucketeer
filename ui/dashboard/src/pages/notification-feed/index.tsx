import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const NotificationFeedPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('notifications')}>
      <PageHeader
        title={t('notifications')}
        description={t('notifications-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default NotificationFeedPage;
