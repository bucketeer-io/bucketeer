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
    <PageLayout.Root title={t('settings')}>
      <PageHeader title={t(`settings`)} description={t(`setting-subtitle`)} />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default SettingsPage;
