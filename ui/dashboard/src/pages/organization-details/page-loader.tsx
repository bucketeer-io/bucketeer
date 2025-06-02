import { useNavigate, useParams } from 'react-router-dom';
import { useQueryOrganizationDetails } from '@queries/organization-details';
import { PAGE_PATH_ORGANIZATIONS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();
  const { organizationId } = useParams();

  const { data, isLoading, refetch, isError } = useQueryOrganizationDetails({
    params: { id: organizationId! }
  });

  const organization = data?.organization;
  const isErrorState = isError || !organization;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader
            title={organization.name}
            description={t('created-at-time', {
              time: formatDateTime(organization.createdAt)
            })}
            onBack={() => navigate(`${PAGE_PATH_ORGANIZATIONS}`)}
          />
          <PageContent organization={organization} />
        </>
      )}
    </>
  );
};

export default PageLoader;
