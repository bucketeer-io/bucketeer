import { useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import { environmentArchive, environmentUnarchive } from '@api/environment';
import { invalidateEnvironments } from '@queries/environments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  ID_NEW,
  PAGE_PATH_ENVIRONMENTS,
  PAGE_PATH_PROJECTS
} from 'constants/routing';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchEnvironments } from './collection-loader/use-fetch-environments';
import AddEnvironmentModal from './environment-modal/add-environment-modal';
import EditEnvironmentModal from './environment-modal/edit-environment-modal';
import PageContent from './page-content';
import { EnvironmentActionsType } from './types';

const ProjectEnvironments = () => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();

  const params = useParams();
  const location = useLocation();
  const navigate = useNavigate();

  const isAdd = useMemo(() => params['*'] === 'new', [params]);
  const isEdit = useMemo(() => params['*'] && !isAdd, [params, isAdd]);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchEnvironments({ pageSize: 1 });

  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();

  const onOpenAddModal = () => navigate(`${location.pathname}/${ID_NEW}`);

  const onCloseActionModal = () =>
    navigate(
      `/${params?.envUrlCode}${PAGE_PATH_PROJECTS}/${params?.projectId}${PAGE_PATH_ENVIRONMENTS}`
    );

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (id: string) => {
      const archiveMutation = isArchiving
        ? environmentArchive
        : environmentUnarchive;

      return archiveMutation({ id, command: {} });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateEnvironments(queryClient);
      mutation.reset();
    }
  });

  const onHandleArchive = () => {
    if (selectedEnvironment?.id) {
      mutation.mutate(selectedEnvironment?.id);
    }
  };

  const onHandleActions = (
    environment: Environment,
    type: EnvironmentActionsType
  ) => {
    if (type === 'ARCHIVE') {
      setIsArchiving(true);
      onOpenConfirmModal();
    } else if (type === 'UNARCHIVE') {
      setIsArchiving(false);
      onOpenConfirmModal();
    } else {
      navigate(`${location.pathname}/${environment.id}`);
    }
    setSelectedEnvironment(environment);
  };

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
            onActionHandler={onHandleActions}
          />

          {isAdd && (
            <AddEnvironmentModal isOpen={isAdd} onClose={onCloseActionModal} />
          )}
          {isEdit && (
            <EditEnvironmentModal
              isOpen={isEdit}
              onClose={onCloseActionModal}
            />
          )}
          {openConfirmModal && (
            <ConfirmModal
              isOpen={openConfirmModal}
              onClose={onCloseConfirmModal}
              onSubmit={onHandleArchive}
              title={
                isArchiving
                  ? t(`table:popover.archive-env`)
                  : t(`table:popover.unarchive-env`)
              }
              description={
                <Trans
                  i18nKey={
                    isArchiving
                      ? 'table:environment.confirm-archive-desc'
                      : 'table:environment.confirm-unarchive-desc'
                  }
                  values={{ name: selectedEnvironment?.name }}
                  components={{ bold: <strong /> }}
                />
              }
              loading={mutation.isPending}
            />
          )}
        </>
      )}
    </>
  );
};

export default ProjectEnvironments;
