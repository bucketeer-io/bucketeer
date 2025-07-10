import { useQueryOrganizationDetails } from '@queries/organization-details';
import { getCurrentEnvironment, useAuth } from 'auth';
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
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data, isLoading, refetch, isError } = useQueryOrganizationDetails({
    params: { id: currentEnvironment.organizationId }
  });

  const organization = data?.organization;

  return (
    <PageLayout.Root title={t('organization-settings')}>
      <PageHeader
        title={t(`organization-settings`)}
        createdAt={organization?.createdAt}
        description={t(`setting-subtitle`)}
        isShowApiEndpoint
      />
      <PageLoader
        isLoading={isLoading}
        isError={isError}
        organization={organization}
        onRetry={refetch}
      />
    </PageLayout.Root>
  );
};

export default SettingsPage;
