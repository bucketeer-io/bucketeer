import { useCallback, useMemo } from 'react';
import { useAuth } from 'auth/auth-context.tsx';
import { getCurrentEnvironment } from 'auth/utils.ts';
import { PAGE_PATH_PROJECTS } from 'constants/routing.ts';
import useActionWithURL from 'hooks/use-action-with-url.tsx';
import { Project } from '@types';
import PageContent from './page-content';
import ProjectCreateUpdateModal from './project-modal/project-create-update-modal/index.tsx';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const commonPath = useMemo(
    () => `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}`,
    [currentEnvironment]
  );
  const { isAdd, isEdit, onOpenEditModal, onCloseActionModal, onOpenAddModal } =
    useActionWithURL({
      closeModalPath: commonPath
    });

  const handleOnCloseCreateUpdateModal = useCallback(() => {
    onCloseActionModal();
  }, []);

  const handleOnEditProject = useCallback(
    (value: Project) => {
      onOpenEditModal(
        `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${value.id}`
      );
    },
    [currentEnvironment]
  );

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onEdit={handleOnEditProject} />

      {(!!isAdd || !!isEdit) && (
        <ProjectCreateUpdateModal
          isOpen={!!isAdd || !!isEdit}
          onClose={handleOnCloseCreateUpdateModal}
        />
      )}
    </>
  );
};

export default PageLoader;
