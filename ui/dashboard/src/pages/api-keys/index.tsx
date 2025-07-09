import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const APIKeysPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('api-keys')}>
      <PageHeader
        title={t('api-keys')}
        description={t('api-keys-subtitle')}
        isShowApiEndpoint
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default APIKeysPage;
