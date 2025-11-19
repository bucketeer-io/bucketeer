import { useCallback, useState } from 'react';
import { Trans } from 'react-i18next';
import { organizationArchive, organizationUnarchive } from '@api/organization';
import { invalidateOrganizations } from '@queries/organizations';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { PAGE_PATH_ORGANIZATIONS } from 'constants/routing';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import OrganizationCreateUpdateModal from './organization-modal/organization-create-update-modal';
import PageContent from './page-content';
import { OrganizationActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();
  const [selectedOrganization, setSelectedOrganization] =
    useState<Organization>();

  const [isArchiving, setIsArchiving] = useState<boolean>();

  const commonPath = `${PAGE_PATH_ORGANIZATIONS}`;

  const { isAdd, isEdit, onCloseActionModal, onOpenAddModal, onOpenEditModal } =
    useActionWithURL({
      closeModalPath: commonPath
    });

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (id: string) => {
      const archiveMutation = isArchiving
        ? organizationArchive
        : organizationUnarchive;

      return archiveMutation({ id, command: {} });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateOrganizations(queryClient);
      mutation.reset();
    }
  });

  const onHandleArchive = useCallback(() => {
    if (selectedOrganization?.id) {
      mutation.mutate(selectedOrganization?.id);
    }
  }, [selectedOrganization, isArchiving, mutation]);

  const onHandleActions = useCallback(
    (organization: Organization, type: OrganizationActionsType) => {
      setSelectedOrganization(organization);
      if (type === 'EDIT') {
        return onOpenEditModal(`${PAGE_PATH_ORGANIZATIONS}/${organization.id}`);
      } else if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
        setIsArchiving(type === 'ARCHIVE');
        return onOpenConfirmModal();
      } else {
        onOpenAddModal();
      }
    },
    []
  );

  const handleOnCloseModal = useCallback(() => {
    onCloseActionModal();
    onCloseConfirmModal();
    setSelectedOrganization(undefined);
  }, []);

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      {(!!isAdd || !!isEdit) && (
        <OrganizationCreateUpdateModal
          isOpen={!!isAdd || !!isEdit}
          onClose={handleOnCloseModal}
          organization={selectedOrganization!}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={handleOnCloseModal}
          onSubmit={onHandleArchive}
          title={
            isArchiving
              ? t(`table:popover.archive-org`)
              : t(`table:popover.unarchive-org`)
          }
          description={
            <Trans
              i18nKey={
                isArchiving
                  ? 'table:organization.confirm-archive-desc'
                  : 'table:organization.confirm-unarchive-desc'
              }
              values={{ name: selectedOrganization?.name }}
              components={{ bold: <strong /> }}
            />
          }
          loading={mutation.isPending}
        />
      )}
    </>
  );
};

export default PageLoader;
