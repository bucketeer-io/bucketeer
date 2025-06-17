import { useState } from 'react';
import { Trans } from 'react-i18next';
import { pushUpdater } from '@api/push';
import { pushDelete } from '@api/push/push-delete';
import { invalidatePushes } from '@queries/pushes';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Push } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageContent from './page-content';
import AddPushModal from './push-modal/add-push-modal';
import EditPushModal from './push-modal/edit-push-modal';
import { PushActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { notify, errorNotify } = useToast();

  const [selectedPush, setSelectedPush] = useState<Push>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);
  const [isDeletePush, setIsDeletePush] = useState<boolean>(false);

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (push: Push, type: PushActionsType) => {
    setSelectedPush(push);
    if (type === 'EDIT') {
      return onOpenEditModal();
    }
    onOpenConfirmModal();
    if (type === 'DELETE') return setIsDeletePush(true);
    setIsDisabling(type !== 'ENABLE');
  };

  const mutationState = useMutation({
    mutationFn: async (id: string) => {
      return isDeletePush
        ? await pushDelete({
            id,
            environmentId: selectedPush?.environmentId
          })
        : await pushUpdater({
            id,
            environmentId: selectedPush?.environmentId,
            disabled: isDisabling
          });
    },
    onSuccess: () => {
      notify({
        message: t(`message:collection-action-success`, {
          collection: t('common:source-type.push'),
          action: t(isDeletePush ? 'common:deleted' : 'common:updated')
        })
      });
      onCloseConfirmModal();
      setIsDeletePush(false);
      invalidatePushes(queryClient);
      mutationState.reset();
    },
    onError: errorNotify
  });

  const onHandleConfirmSubmit = () => {
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
          onSubmit={onHandleConfirmSubmit}
          title={t(
            `popover.${isDeletePush ? 'delete' : isDisabling ? 'disable' : 'enable'}-push`
          )}
          description={
            <Trans
              i18nKey={`table:push.confirm-${isDeletePush ? 'delete' : isDisabling ? 'disable' : 'enable'}-desc`}
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
