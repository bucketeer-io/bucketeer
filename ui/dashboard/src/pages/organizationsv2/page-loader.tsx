import { useToggleOpen } from 'hooks/use-toggle-open';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchOrganizations } from './collection-loader/use-fetch-organizations';
import PageContent from './page-content';

const PageLoader = () => {
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ pageSize: 1 });

  const [, onOpenAddModal] = useToggleOpen(false);

  if (isLoading) {
    return <PageLayout.LoadingState />;
  }

  if (isError) {
    return <PageLayout.ErrorState onRetry={refetch} />;
  }

  const isEmpty = collection?.Organizations.length === 0;

  if (isEmpty) {
    return (
      <PageLayout.EmptyState>
        <EmptyCollection onAdd={onOpenAddModal} />
      </PageLayout.EmptyState>
    );
  }
  return (
    <>
      <PageContent onAdd={onOpenAddModal} />
    </>
  );
};

export default PageLoader;
