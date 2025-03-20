import { useNavigate, useParams } from 'react-router-dom';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

const PageLoader = () => {
  const params = useParams();

  const navigate = useNavigate();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    isError,
    refetch
  } = useQueryFeature({
    params: {
      environmentId: currentEnvironment?.id,
      id: params?.flagId || ''
    },
    enabled: !!params?.flagId && !!currentEnvironment?.id
  });

  const feature = collection?.feature;

  const isErrorState = isError || !feature;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader onBack={() => navigate(`${PAGE_PATH_FEATURES}`)}>
            <HeaderDetails feature={feature} />
          </PageDetailsHeader>
          <PageContent feature={feature} />
        </>
      )}
    </>
  );
};

export default PageLoader;
