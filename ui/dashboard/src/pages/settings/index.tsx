import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

export interface SettingsPageForm {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail: string;
}

const SettingsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('organization-settings')}>
      <PageHeader
        title={t(`organization-settings`)}
        description={t(`setting-subtitle`)}
        isShowApiEndpoint
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default SettingsPage;
