import { useParams } from 'react-router-dom';
import { useQueryOrganizationDetails } from '@queries/organization-details';
import { useFormatDateTime } from 'utils/date-time';
import InvalidMessage from 'elements/invalid-message';
import PageDetailHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const formatDateTime = useFormatDateTime();
  const { organizationId } = useParams();

  const { data, isLoading, refetch, isError } = useQueryOrganizationDetails({
    params: { id: organizationId! }
  });

  const organization = data?.organization;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : !organization ? (
        <InvalidMessage>{`Invalid data`}</InvalidMessage>
      ) : (
        <>
          <PageDetailHeader
            title={organization.name}
            description={`Created ${formatDateTime(organization.createdAt)}`}
          />
          <PageContent />
        </>
      )}
    </>
  );
};

export default PageLoader;
