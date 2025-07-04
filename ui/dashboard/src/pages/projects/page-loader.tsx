import { useCallback, useState } from 'react';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Project } from '@types';
import PageContent from './page-content';
import ProjectCreateUpdateModal from './project-modal/project-create-update-modal/index.tsx';

const PageLoader = () => {
  const [selectedProject, setSelectedProject] = useState<Project>();

  const [
    isOpenCreateUpdateModal,
    onOpenCreateUpdateModal,
    onCloseCreateUpdateModal
  ] = useToggleOpen(false);

  const handleOnCloseCreateUpdateModal = useCallback(() => {
    onCloseCreateUpdateModal();
    setSelectedProject(undefined);
  }, []);

  const handleOnEditProject = useCallback((value: Project) => {
    setSelectedProject(value);
    onOpenCreateUpdateModal();
  }, []);

  return (
    <>
      <PageContent
        onAdd={onOpenCreateUpdateModal}
        onEdit={handleOnEditProject}
      />

      {isOpenCreateUpdateModal && (
        <ProjectCreateUpdateModal
          isOpen={isOpenCreateUpdateModal}
          project={selectedProject}
          onClose={handleOnCloseCreateUpdateModal}
        />
      )}
    </>
  );
};

export default PageLoader;
