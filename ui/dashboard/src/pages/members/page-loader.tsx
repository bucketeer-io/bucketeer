import { useState } from 'react';
import { accountDeleter } from '@api/account/account-deleter';
import { invalidateAccounts } from '@queries/accounts';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Account } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchMembers } from './collection-loader/use-fetch-members';
import AddMemberModal from './member-modal/add-member-modal';
import DeleteMemberModal from './member-modal/delete-member-modal';
import EditMemberModal from './member-modal/edit-member-modal';
import MemberDetailsModal from './member-modal/member-details-modal';
import PageContent from './page-content';
import { MemberActionsType } from './types';

const PageLoader = () => {
  const { notify } = useToast();
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchMembers({
    pageSize: 1,
    organizationId: currenEnvironment.organizationId
  });

  const [selectedMember, setSelectedMember] = useState<Account>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenDetailsModal, onOpenDetailsModal, onCloseDetailsModal] =
    useToggleOpen(false);

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (member: Account) => {
      return accountDeleter({
        email: member.email,
        organizationId: member.organizationId,
        command: {}
      });
    },
    onSuccess: () => {
      onCloseDeleteModal();
      invalidateAccounts(queryClient);
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{selectedMember?.email}</b>
            {` has been deleted successfully!`}
          </span>
        )
      });
      mutation.reset();
    }
  });

  const onDeleteMember = () => {
    if (selectedMember) {
      mutation.mutate(selectedMember);
    }
  };

  const onHandleActions = (member: Account, type: MemberActionsType) => {
    if (type === 'EDIT') {
      onOpenEditModal();
    } else if (type === 'DELETE') {
      onOpenDeleteModal();
    } else if (type === 'DETAILS') {
      onOpenDetailsModal();
    }
    setSelectedMember(member);
  };

  const isEmpty = collection?.accounts.length === 0;

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
        <AddMemberModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenEditModal && (
        <EditMemberModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          member={selectedMember!}
        />
      )}
      {isOpenDeleteModal && (
        <DeleteMemberModal
          isOpen={isOpenDeleteModal}
          onClose={onCloseDeleteModal}
          member={selectedMember!}
          loading={mutation.isPending}
          onSubmit={onDeleteMember}
        />
      )}
      {isOpenDetailsModal && (
        <MemberDetailsModal
          isOpen={isOpenDetailsModal}
          onClose={onCloseDetailsModal}
          member={selectedMember!}
        />
      )}
    </>
  );
};

export default PageLoader;
