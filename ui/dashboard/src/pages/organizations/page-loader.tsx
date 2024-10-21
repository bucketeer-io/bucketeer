import { useToggleOpen } from 'hooks/use-toggle-open';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchOrganizations } from './collection-loader/use-fetch-organizations';
import AddOrganizationModal from './organization-modal/add-organization-modal';
import PageContent from './page-content';

const PageLoader = () => {
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ pageSize: 1 });

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const isEmpty = collection?.Organizations.length === 0;

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
        <>
          <PageContent onAdd={onOpenAddModal} />

          {isOpenAddModal && (
            <AddOrganizationModal
              isOpen={isOpenAddModal}
              onClose={onCloseAddModal}
            />
          )}
        </>
      )}
    </>
  );
};

export default PageLoader;
