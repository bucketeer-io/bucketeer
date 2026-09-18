import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const NotificationsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('slack')}>
      <PageHeader
        title={t('slack')}
        description={t('slack-subtitle')}
        link={DOCUMENTATION_LINKS.SLACK}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default NotificationsPage;
