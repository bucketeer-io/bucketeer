import { useState } from 'react';
import { Trans } from 'react-i18next';
import { pushUpdater } from '@api/push';
import { invalidatePushes } from '@queries/pushes';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Push } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchPushes } from './collection-loader/use-fetch-pushes';
import PageContent from './page-content';
import AddPushModal from './push-modal/add-push-modal';
import EditPushModal from './push-modal/edit-push-modal';
import { PushActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchPushes({
    pageSize: 1,
    environmentId: currenEnvironment.id
  });

  const [selectedPush, setSelectedPush] = useState<Push>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (push: Push, type: PushActionsType) => {
    setSelectedPush(push);

    if (type === 'EDIT') return onOpenEditModal();
    if (type === 'ENABLE') {
      setIsDisabling(false);
      return onOpenConfirmModal();
    }
    setIsDisabling(true);
    onOpenConfirmModal();
  };

  const mutationState = useMutation({
    mutationFn: async (id: string) => {
      return pushUpdater({
        id,
        environmentId: selectedPush?.environmentId,
        disabled: isDisabling
      });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidatePushes(queryClient);
      mutationState.reset();
    }
  });

  const onHandleDisable = () => {
    if (selectedPush?.id) {
      mutationState.mutate(selectedPush?.id);
    }
  };

  const isEmpty = collection?.pushes.length === 0;

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
        <AddPushModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenEditModal && selectedPush && (
        <EditPushModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          push={selectedPush}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleDisable}
          title={
            isDisabling
              ? t(`table:popover.disable-push`)
              : t(`table:popover.enable-push`)
          }
          description={
            <Trans
              i18nKey={
                isDisabling
                  ? 'table:push.confirm-disable-desc'
                  : 'table:push.confirm-enable-desc'
              }
              values={{ name: selectedPush?.name }}
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
