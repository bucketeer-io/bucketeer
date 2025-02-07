import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import PageContent from './page-content';

const PageLoader = () => {
  const isLoading = false,
    isError = false,
    isEmpty = false;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={() => {}} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={() => {}} />
      )}
    </>
  );
};

export default PageLoader;
