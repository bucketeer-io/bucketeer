import { useEffect } from 'react';
import { useToggleOpen } from 'hooks';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import AddExperimentModal from './experiments-modal/add-experiment-modal';
import PageContent from './page-content';

const PageLoader = ({
  setTotalCount
}: {
  setTotalCount: (value: string | number) => void;
}) => {
  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const isLoading = false,
    isError = false,
    isEmpty = true;

  useEffect(() => {
    setTotalCount(25);
  }, []);

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={() => {}} />
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
