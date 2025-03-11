import { useState } from 'react';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Project } from '@types';
import PageContent from './page-content';
import AddProjectModal from './project-modal/add-project-modal';
import EditProjectModal from './project-modal/edit-project-modal';

const PageLoader = () => {
  const [selectedProject, setSelectedProject] = useState<Project>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  return (
    <>
      <PageContent
        onAdd={onOpenAddModal}
        onEdit={value => {
          setSelectedProject(value);
          onOpenEditModal();
        }}
      />
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
