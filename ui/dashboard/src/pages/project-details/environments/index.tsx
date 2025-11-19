import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { environmentArchive, environmentUnarchive } from '@api/environment';
import { invalidateEnvironments } from '@queries/environments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth, useAuthAccess } from 'auth';
import {
  PAGE_PATH_ENVIRONMENTS,
  PAGE_PATH_NEW,
  PAGE_PATH_PROJECTS
} from 'constants/routing';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import EnvironmentCreateUpdateModal from './environment-modal/environment-create-update-modal';
import PageContent from './page-content';
import { EnvironmentActionsType } from './types';

const ProjectEnvironments = ({
  organizationId
}: {
  organizationId: string;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
  const params = useParams();
  const { consoleAccount, onMeFetcher } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();

  const commonPath = useMemo(
    () =>
      `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${params.projectId}${PAGE_PATH_ENVIRONMENTS}`,
    [currentEnvironment]
  );
  const { isAdd, isEdit, onOpenEditModal, onCloseActionModal, onOpenAddModal } =
    useActionWithURL({
      closeModalPath: `${commonPath}?organizationId=${organizationId}`,
      addPath: `${commonPath}${PAGE_PATH_NEW}?organizationId=${organizationId}`,
      idKey: 'environmentId'
    });

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (id: string) => {
      const archiveMutation = isArchiving
        ? environmentArchive
        : environmentUnarchive;

      return archiveMutation({ id });
    },
    onSuccess: async () => {
      await onMeFetcher({ organizationId: currentEnvironment.organizationId });
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
      if (type === 'EDIT') {
        onOpenEditModal(
          `${commonPath}/${environment.id}?organizationId=${organizationId}`
        );
      }
      if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
        setIsArchiving(type === 'ARCHIVE');
        return onOpenConfirmModal();
      }
    },
    []
  );

  const handleOnCloseModal = useCallback(() => {
    onCloseActionModal();
    onCloseConfirmModal();
    setSelectedEnvironment(undefined);
  }, []);

  return (
    <>
      <PageContent
        organizationId={organizationId}
        onAdd={onOpenAddModal}
        onActionHandler={onHandleActions}
      />
      {(!!isAdd || !!isEdit) && (
        <EnvironmentCreateUpdateModal
          organizationId={organizationId}
          isOpen={!!isAdd || !!isEdit}
          environment={selectedEnvironment}
          onClose={onCloseActionModal}
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
