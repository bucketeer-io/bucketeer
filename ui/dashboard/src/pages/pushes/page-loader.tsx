import { useCallback, useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { pushUpdater } from '@api/push';
import { pushDelete } from '@api/push/push-delete';
import { useQueryPush } from '@queries/push-details';
import { invalidatePushes } from '@queries/pushes';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_PUSHES } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Push } from '@types';
import { useSearchParams } from 'utils/search-params';
import ConfirmModal from 'elements/confirm-modal';
import PageContent from './page-content';
import PushCreateUpdateModal from './push-modal/push-create-update-modal';
import { PushActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const { notify, errorNotify } = useToast();
  const { searchOptions } = useSearchParams();
  const pushEnvironmentId = searchOptions?.environmentId;

  const commonPath = useMemo(
    () => `/${currentEnvironment.urlCode}${PAGE_PATH_PUSHES}`,
    [currentEnvironment]
  );

  const {
    id: pushId,
    isAdd,
    isEdit,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal
  } = useActionWithURL({
    closeModalPath: commonPath
  });

  const [selectedPush, setSelectedPush] = useState<Push>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);
  const [isDeletePush, setIsDeletePush] = useState<boolean>(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const {
    data: pushCollection,
    isLoading: isLoadingPush,
    isError,
    error
  } = useQueryPush({
    params: {
      id: pushId as string,
      environmentId: pushEnvironmentId as string
    },
    enabled: !!isEdit && !!pushId && !selectedPush
  });

  const onHandleActions = useCallback(
    (push: Push, type: PushActionsType) => {
      setSelectedPush(push);
      if (type === 'EDIT') {
        return onOpenEditModal(
          `${commonPath}/${push.id}?environmentId=${push.environmentId}`
        );
      }
      onOpenConfirmModal();
      if (type === 'DELETE') return setIsDeletePush(true);
      setIsDisabling(type !== 'ENABLE');
    },
    [commonPath]
  );

  const handleOnCloseModal = useCallback((isRefresh?: boolean) => {
    const checkReset = isRefresh ?? true;
    onCloseActionModal();
    if (checkReset) {
      setSelectedPush(undefined);
    }
    setIsDeletePush(false);
    setIsDisabling(false);
    onCloseConfirmModal();
  }, []);

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
      handleOnCloseModal();
      invalidatePushes(queryClient);
      mutationState.reset();
    },
    onError: errorNotify
  });

  const onHandleConfirmSubmit = useCallback(() => {
    if (selectedPush?.id) {
      mutationState.mutate(selectedPush?.id);
    }
  }, [selectedPush]);

  useEffect(() => {
    if (pushCollection) {
      setSelectedPush(pushCollection.push);
    }
  }, [pushCollection]);

  useEffect(() => {
    if (isError && error) {
      errorNotify(error);
      handleOnCloseModal();
    }
  }, [isError, error]);

  return (
    <>
      <PageContent
        disabled={!editable}
        onAdd={onOpenAddModal}
        onHandleActions={onHandleActions}
      />
      {(!!isAdd || !!isEdit) && (
        <PushCreateUpdateModal
          disabled={!editable}
          isOpen={!!isAdd || !!isEdit}
          pushId={pushId}
          isLoadingPush={isLoadingPush}
          push={selectedPush}
          resetPush={() => setSelectedPush(undefined)}
          onClose={handleOnCloseModal}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={handleOnCloseModal}
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
          disabled={!editable}
        />
      )}
    </>
  );
};

export default PageLoader;
