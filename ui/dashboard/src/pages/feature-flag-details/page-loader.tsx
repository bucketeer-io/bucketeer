import { useNavigate } from 'react-router-dom';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { mockFlags } from 'pages/feature-flags/collection-loader';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
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
          <PageDetailsHeader onBack={() => navigate(`${PAGE_PATH_FEATURES}`)}>
            <HeaderDetails featureFlag={mockFlags[0]} />
          </PageDetailsHeader>
          <PageContent featureFlag={mockFlags[0]} />
        </>
      )}
    </>
  );
};

export default PageLoader;
