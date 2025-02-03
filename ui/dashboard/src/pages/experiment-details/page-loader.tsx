import { useNavigate, useParams } from 'react-router-dom';
import { useQueryExperimentDetails } from '@queries/experiment-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import NotFoundPage from 'pages/not-found';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const params = useParams();

  const {
    data: experimentCollection,
    isLoading,
    isError,
    refetch
  } = useQueryExperimentDetails({
    params: {
      id: params?.experimentId || '',
      environmentId: currentEnvironment.id
    }
  });

  const experiment = experimentCollection?.experiment;
  const isErrorState = isError || !experiment;
  if (params?.tab && !['results', 'settings'].includes(params.tab))
    return <NotFoundPage />;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader
            onBack={() =>
              navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`)
            }
          >
            <HeaderDetails experiment={experiment} />
          </PageDetailsHeader>
          <PageContent experiment={experiment} />
        </>
      )}
    </>
  );
};

export default PageLoader;
