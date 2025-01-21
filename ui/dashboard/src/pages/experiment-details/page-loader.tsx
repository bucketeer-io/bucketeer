import { useNavigate } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const isLoading = false;
  const isErrorState = false;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : (
        <>
          <PageDetailsHeader
            onBack={() =>
              navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`)
            }
          >
            <HeaderDetails />
          </PageDetailsHeader>
          <PageContent />
        </>
      )}
    </>
  );
};

export default PageLoader;
