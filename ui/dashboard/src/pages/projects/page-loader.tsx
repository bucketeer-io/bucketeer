import { useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Project } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchProjects } from './collection-loader/use-fetch-projects';
import PageContent from './page-content';
import AddProjectModal from './project-modal/add-project-modal';
import EditProjectModal from './project-modal/edit-project-modal';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchProjects({
    pageSize: 1,
    organizationId: currentEnvironment.organizationId
  });

  const [selectedProject, setSelectedProject] = useState<Project>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const isEmpty = collection?.projects.length === 0;

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
        <PageContent
          onAdd={onOpenAddModal}
          onEdit={value => {
            setSelectedProject(value);
            onOpenEditModal();
          }}
        />
      )}
      {isOpenAddModal && (
        <AddProjectModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenEditModal && (
        <EditProjectModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          project={selectedProject!}
        />
      )}
    </>
  );
};

export default PageLoader;
