import { useQueryOrganizationDetails } from '@queries/organization-details';
import { getOrgIdStorage } from 'storage/organization';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const organizationId = getOrgIdStorage();

  const { data, isLoading, refetch, isError } = useQueryOrganizationDetails({
    params: { id: organizationId! }
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
