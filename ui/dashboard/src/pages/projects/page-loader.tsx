import { useLocation } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { ID_NEW, PAGE_PATH_PROJECTS } from 'constants/routing';
import useActionWithURL from 'hooks/use-action-with-url';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchProjects } from './collection-loader/use-fetch-projects';
import PageContent from './page-content';
import AddProjectModal from './project-modal/add-project-modal';
import EditProjectModal from './project-modal/edit-project-modal';

const PageLoader = () => {
  const location = useLocation();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { isAdd, isEdit, onOpenAddModal, onOpenEditModal, onCloseActionModal } =
    useActionWithURL({
      idKey: '*',
      addPath: `${location.pathname}/${ID_NEW}`,
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}`
    });
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchProjects({ pageSize: 1 });

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
          onEdit={project =>
            onOpenEditModal(`${location.pathname}/${project.id}`)
          }
        />
      )}
      {isAdd && <AddProjectModal isOpen={isAdd} onClose={onCloseActionModal} />}
      {isEdit && (
        <EditProjectModal isOpen={isEdit} onClose={onCloseActionModal} />
      )}
    </>
  );
};

export default PageLoader;
