import { useState } from 'react';
import { Trans } from 'react-i18next';
import { pushUpdater } from '@api/push';
import { invalidatePushes } from '@queries/pushes';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Push } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageContent from './page-content';
import AddPushModal from './push-modal/add-push-modal';
import EditPushModal from './push-modal/edit-push-modal';
import { PushActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const queryClient = useQueryClient();

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
    if (type === 'EDIT') {
      onOpenEditModal();
    } else {
      setIsDisabling(type !== 'ENABLE');
      onOpenConfirmModal();
    }
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

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />

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
