import { useCallback, useState } from 'react';
import { Trans } from 'react-i18next';
import { organizationArchive, organizationUnarchive } from '@api/organization';
import { invalidateOrganizations } from '@queries/organizations';
import { useMutation, useQueryClient } from '@tanstack/react-query';
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
    setSelectedOrganization(undefined);
  }, []);

  return (
    <>
      <PageContent
        onAdd={onOpenCreateUpdateModal}
        onHandleActions={onHandleActions}
      />
      {isOpenCreateUpdateModal && (
        <OrganizationCreateUpdateModal
          isOpen={isOpenCreateUpdateModal}
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
