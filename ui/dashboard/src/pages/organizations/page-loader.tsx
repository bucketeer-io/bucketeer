import { useState } from 'react';
import { Trans } from 'react-i18next';
import { organizationArchive, organizationUnarchive } from '@api/organization';
import { invalidateOrganizations } from '@queries/organizations';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchOrganizations } from './collection-loader/use-fetch-organizations';
import AddOrganizationModal from './organization-modal/add-organization-modal';
import EditOrganizationModal from './organization-modal/edit-organization-modal';
import PageContent from './page-content';
import { OrganizationActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ pageSize: 1 });

  const [selectedOrganization, setSelectedOrganization] =
    useState<Organization>();

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
        ? organizationArchive
        : organizationUnarchive;

      return archiveMutation({ id });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateOrganizations(queryClient);
      mutation.reset();
    }
  });

  const onHandleArchive = () => {
    if (selectedOrganization?.id) {
      mutation.mutate(selectedOrganization?.id);
    }
  };

  const onHandleActions = (
    organization: Organization,
    type: OrganizationActionsType
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
    setSelectedOrganization(organization);
  };

  const isEmpty = collection?.Organizations.length === 0;

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
        <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      )}
      {isOpenAddModal && (
        <AddOrganizationModal
          isOpen={isOpenAddModal}
          onClose={onCloseAddModal}
        />
      )}
      {isOpenEditModal && (
        <EditOrganizationModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          organization={selectedOrganization!}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
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
