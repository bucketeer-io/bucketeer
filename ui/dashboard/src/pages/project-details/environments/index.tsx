import { useCallback, useState } from 'react';
import { Trans } from 'react-i18next';
import { environmentArchive, environmentUnarchive } from '@api/environment';
import { invalidateEnvironments } from '@queries/environments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuthAccess } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import EnvironmentCreateUpdateModal from './environment-modal/environment-create-update-modal';
import PageContent from './page-content';
import { EnvironmentActionsType } from './types';

const ProjectEnvironments = () => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();
  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();

  const [
    isOpenCreateUpdateModal,
    onOpenCreateUpdateModal,
    onCloseCreateUpdateModal
  ] = useToggleOpen(false);

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

  const onHandleArchive = useCallback(() => {
    if (selectedEnvironment?.id) {
      mutation.mutate(selectedEnvironment?.id);
    }
  }, [selectedEnvironment]);

  const onHandleActions = useCallback(
    (environment: Environment, type: EnvironmentActionsType) => {
      setSelectedEnvironment(environment);
      if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
        setIsArchiving(type === 'ARCHIVE');
        return onOpenConfirmModal();
      }
      return onOpenCreateUpdateModal();
    },
    []
  );

  const handleOnCloseModal = useCallback(() => {
    onCloseCreateUpdateModal();
    onCloseConfirmModal();
    setSelectedEnvironment(undefined);
  }, []);

  return (
    <>
      <PageContent
        onAdd={onOpenCreateUpdateModal}
        onActionHandler={onHandleActions}
      />
      {isOpenCreateUpdateModal && (
        <EnvironmentCreateUpdateModal
          isOpen={isOpenCreateUpdateModal}
          environment={selectedEnvironment}
          onClose={handleOnCloseModal}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={handleOnCloseModal}
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
          disabled={!envEditable || !isOrganizationAdmin}
        />
      )}
    </>
  );
};

export default ProjectEnvironments;
