import { useEffect } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchExperiments } from './collection-loader/use-fetch-experiment';
import AddExperimentModal from './experiments-modal/add-experiment-modal';
import PageContent from './page-content';

const PageLoader = ({
  setTotalCount
}: {
  setTotalCount: (value: string | number) => void;
}) => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchExperiments({ environmentId: currenEnvironment.id });

  const isEmpty = collection?.experiments?.length === 0;

  useEffect(() => {
    if (collection?.experiments) setTotalCount(collection?.experiments?.length);
  }, [collection]);

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={onOpenAddModal} onHandleActions={() => {}} />
      )}
      {isOpenAddModal && (
        <AddExperimentModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
    </>
  );
};

export default PageLoader;
