import { useState } from 'react';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Organization } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchOrganizations } from './collection-loader/use-fetch-organizations';
import AddOrganizationModal from './organization-modal/add-organization-modal';
import EditOrganizationModal from './organization-modal/edit-organization-modal';
import PageContent from './page-content';

const PageLoader = () => {
  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchOrganizations({ pageSize: 1 });

  const [selectedOrganization, setSelectedOrganization] =
    useState<Organization>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

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
        <>
          <PageContent
            onAdd={onOpenAddModal}
            onEdit={value => {
              setSelectedOrganization(value);
              onOpenEditModal();
            }}
          />
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
        </>
      )}
    </>
  );
};

export default PageLoader;
