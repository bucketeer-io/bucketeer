import { useMemo } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { ID_NEW, PAGE_PATH_PROJECTS } from 'constants/routing';
import { Project } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchProjects } from './collection-loader/use-fetch-projects';
import PageContent from './page-content';
import AddProjectModal from './project-modal/add-project-modal';
import EditProjectModal from './project-modal/edit-project-modal';

const PageLoader = () => {
  const params = useParams();
  const location = useLocation();
  const navigate = useNavigate();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchProjects({ pageSize: 1 });

  const isAdd = useMemo(() => params['*'] && params['*'] === ID_NEW, [params]);
  const isEdit = useMemo(() => params['*'] && !isAdd, [params, isAdd]);

  const onOpenAddModal = () => navigate(`${location.pathname}/${ID_NEW}`);
  const onOpenEditModal = (project: Project) =>
    navigate(`${location.pathname}/${project.id}`);

  const onCloseActionModal = () =>
    navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}`);

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
        <PageContent onAdd={onOpenAddModal} onEdit={onOpenEditModal} />
      )}
      {isAdd && <AddProjectModal isOpen={isAdd} onClose={onCloseActionModal} />}
      {isEdit && (
        <EditProjectModal isOpen={isEdit} onClose={onCloseActionModal} />
      )}
    </>
  );
};

export default PageLoader;
