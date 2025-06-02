import { useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { accountDisable, accountEnable } from '@api/account';
import { accountDeleter } from '@api/account/account-deleter';
import { invalidateAccounts } from '@queries/accounts';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Account } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import AddMemberModal from './member-modal/add-member-modal';
import DeleteMemberModal from './member-modal/delete-member-modal';
import EditMemberModal from './member-modal/edit-member-modal';
import MemberDetailsModal from './member-modal/member-details-modal';
import PageContent from './page-content';
import { MemberActionsType } from './types';

const PageLoader = () => {
  const { notify } = useToast();
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedMember, setSelectedMember] = useState<Account>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenDetailsModal, onOpenDetailsModal, onCloseDetailsModal] =
    useToggleOpen(false);

  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (member: Account) => {
      return accountDeleter({
        email: member.email,
        organizationId: member.organizationId
      });
    },
    onSuccess: () => {
      onCloseDeleteModal();
      invalidateAccounts(queryClient);
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:member'),
          action: t('common:deleted')
        })
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
        organizationId: currentEnvironment.organizationId
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
    setSelectedMember(member);
    switch (type) {
      case 'EDIT':
        return onOpenEditModal();
      case 'DELETE':
        return onOpenDeleteModal();
      case 'DETAILS':
        return onOpenDetailsModal();
      case 'DISABLE':
      case 'ENABLE':
        setIsDisabling(type === 'DISABLE');
        return onOpenConfirmModal();
      default:
        return;
    }
  };

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
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
