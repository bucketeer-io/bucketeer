import { useState } from 'react';
import { Trans } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { environmentArchive, environmentUnarchive } from '@api/environment';
import { invalidateEnvironments } from '@queries/environments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
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
  const { projectId } = useParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchEnvironments({
    pageSize: 1,
    projectId,
    organizationId: currentEnvironment.organizationId
  });

  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (id: string) => {
      const archiveMutation = isArchiving
        ? environmentArchive
        : environmentUnarchive;

      return archiveMutation({ id });
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
      onOpenEditModal();
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
