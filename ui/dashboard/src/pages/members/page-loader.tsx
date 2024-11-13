import { useState } from 'react';
import { accountDeleter } from '@api/account/account-deleter';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Account } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchMembers } from './collection-loader/use-fetch-members';
import AddMemberModal from './member-modal/add-member-modal';
import EditMemberModal from './member-modal/edit-member-modal';
import PageContent from './page-content';
import { MemberActionsType } from './types';

const PageLoader = () => {
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

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const onHandleActions = (member: Account, type: MemberActionsType) => {
    if (type === 'EDIT') {
      onOpenEditModal();
      setSelectedMember(member);
    } else if (type === 'DELETE') {
      accountDeleter({
        email: member.email,
        organizationId: member.organizationId,
        command: {}
      });
    }
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
        <>
          <PageContent
            onAdd={onOpenAddModal}
            onHandleActions={onHandleActions}
          />
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
        </>
      )}
    </>
  );
};

export default PageLoader;
