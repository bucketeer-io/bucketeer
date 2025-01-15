import { useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { useLocation } from 'react-router-dom';
import { accountDisable, accountEnable } from '@api/account';
import { accountDeleter } from '@api/account/account-deleter';
import { invalidateAccounts } from '@queries/accounts';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_MEMBERS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Account } from '@types';
import ConfirmModal from 'elements/confirm-modal';
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
  const { t } = useTranslation(['table']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);
  const location = useLocation();

  const {
    isAdd,
    isEdit,
    params,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal
  } = useActionWithURL({});

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
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const handleOnCloseActionModal = () =>
    onCloseActionModal(`/${params?.envUrlCode}${PAGE_PATH_MEMBERS}`);

  const [isOpenDetailsModal, onOpenDetailsModal, onCloseDetailsModal] =
    useToggleOpen(false);

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
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

  const mutationState = useMutation({
    mutationFn: async (email: string) => {
      const archiveMutation = isDisabling ? accountDisable : accountEnable;

      return archiveMutation({
        email,
        organizationId: currenEnvironment.organizationId,
        command: {}
      });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateAccounts(queryClient);
      mutationState.reset();
    }
  });

  const onHandleDisable = () => {
    if (selectedMember?.email) {
      mutationState.mutate(selectedMember.email);
    }
  };

  const onHandleActions = (member: Account, type: MemberActionsType) => {
    if (type === 'EDIT')
      return onOpenEditModal(`${location.pathname}/${member.email}`);
    setSelectedMember(member);
    if (type === 'DELETE') return onOpenDeleteModal();
    if (type === 'DETAILS') return onOpenDetailsModal();
    if (type === 'ENABLE') {
      setIsDisabling(false);
      return onOpenConfirmModal();
    }
    setIsDisabling(true);
    onOpenConfirmModal();
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
      {isAdd && (
        <AddMemberModal isOpen={isAdd} onClose={handleOnCloseActionModal} />
      )}
      {isEdit && (
        <EditMemberModal isOpen={isEdit} onClose={handleOnCloseActionModal} />
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
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleDisable}
          title={
            isDisabling
              ? t(`table:popover.disable-member`)
              : t(`table:popover.enable-member`)
          }
          description={
            <Trans
              i18nKey={
                isDisabling
                  ? 'table:members.confirm-disable-desc'
                  : 'table:members.confirm-enable-desc'
              }
              values={{ email: selectedMember?.email }}
              components={{ bold: <strong /> }}
            />
          }
          loading={mutationState.isPending}
        />
      )}
    </>
  );
};

export default PageLoader;
