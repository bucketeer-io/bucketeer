import { useState } from 'react';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Environment } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchEnvironments } from './collection-loader/use-fetch-environments';
import AddEnvironmentModal from './environment-modal/add-environment-modal';
import EditEnvironmentModal from './environment-modal/edit-environment-modal';
import PageContent from './page-content';

const ProjectEnvironments = () => {
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchEnvironments({ pageSize: 1 });

  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const isEmpty = collection?.environments.length === 0;

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
          <PageContent
            onAdd={onOpenAddModal}
            onEdit={value => {
              setSelectedEnvironment(value);
              onOpenEditModal();
            }}
          />

          {isOpenAddModal && (
            <AddEnvironmentModal
              isOpen={isOpenAddModal}
              onClose={onCloseAddModal}
            />
          )}
          {isOpenEditModal && (
            <EditEnvironmentModal
              isOpen={isOpenEditModal}
              onClose={onCloseEditModal}
              environment={selectedEnvironment!}
            />
          )}
        </>
      )}
    </>
  );
};

export default ProjectEnvironments;
