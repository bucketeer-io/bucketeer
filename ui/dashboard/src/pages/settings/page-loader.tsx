import { useQueryOrganizationDetails } from '@queries/organization-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data, isLoading, refetch, isError } = useQueryOrganizationDetails({
    params: { id: currenEnvironment.organizationId }
  });

  const organization = data?.organization;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError || !organization ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <PageContent organization={organization} />
      )}
    </>
  );
};

export default PageLoader;
