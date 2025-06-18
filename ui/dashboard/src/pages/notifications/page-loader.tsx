import { useCallback, useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { notificationUpdater } from '@api/notification';
import { useQueryNotification } from '@queries/notification-details';
import { invalidateNotifications } from '@queries/notifications';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_NOTIFICATIONS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Notification } from '@types';
import { useSearchParams } from 'utils/search-params';
import ConfirmModal from 'elements/confirm-modal';
import AddNotificationModal from './notification-modal/add-notification-modal';
import EditNotificationModal from './notification-modal/edit-notification-modal';
import PageContent from './page-content';
import { NotificationActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { notify, errorNotify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const commonPath = useMemo(
    () => `/${currentEnvironment.urlCode}${PAGE_PATH_NOTIFICATIONS}`,
    [currentEnvironment]
  );

  const {
    id: notificationId,
    isAdd,
    isEdit,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal
  } = useActionWithURL({
    closeModalPath: commonPath
  });

  const [isDisabling, setIsDisabling] = useState<boolean>(false);
  const [selectedNotification, setSelectedNotification] =
    useState<Notification>();
  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);
  const { searchOptions } = useSearchParams();
  const notificationEnvironmentId = searchOptions?.environmentId;

  const {
    data: notificationCollection,
    isLoading: isLoadingNotification,
    isError,
    error
  } = useQueryNotification({
    params: {
      id: notificationId as string,
      environmentId: notificationEnvironmentId as string
    },
    enabled: !!isEdit && !!notificationId && !selectedNotification
  });

  const onHandleActions = useCallback(
    (notification: Notification, type: NotificationActionsType) => {
      setSelectedNotification(notification);
      switch (type) {
        case 'EDIT':
          return onOpenEditModal(
            `${commonPath}/${notification.id}?environmentId=${notification.environmentId}`
          );
        case 'ENABLE':
        case 'DISABLE':
          setIsDisabling(type === 'DISABLE');
          return onOpenConfirmModal();
        default:
          break;
      }
    },
    [commonPath]
  );

  const mutationState = useMutation({
    mutationFn: async (notification: Notification) => {
      return notificationUpdater({
        id: notification.id,
        environmentId: notification.environmentId,
        disabled: isDisabling
      });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateNotifications(queryClient);
      setSelectedNotification(undefined);
      mutationState.reset();
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:notification'),
          action: t('common:updated')
        })
      });
    },
    onError: error => errorNotify(error)
  });

  const onHandleDisable = useCallback(() => {
    if (selectedNotification?.id) {
      mutationState.mutate(selectedNotification);
    }
  }, [selectedNotification]);

  useEffect(() => {
    if (notificationCollection) {
      setSelectedNotification(notificationCollection.subscription);
    }
  }, [notificationCollection]);

  useEffect(() => {
    if (isError && error) {
      errorNotify(error);
      onCloseActionModal();
    }
  }, [isError, error]);

  return (
    <>
      <PageContent
        disabled={!editable}
        onAdd={onOpenAddModal}
        onHandleActions={onHandleActions}
      />

      {isAdd && (
        <AddNotificationModal
          disabled={!editable}
          isOpen={isAdd}
          onClose={onCloseActionModal}
        />
      )}
      {isEdit && (
        <EditNotificationModal
          disabled={!editable}
          isOpen={isEdit}
          isLoadingNotification={isLoadingNotification}
          notification={selectedNotification}
          onClose={onCloseActionModal}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleDisable}
          title={
            isDisabling
              ? t(`table:popover.disable-notification`)
              : t(`table:popover.enable-notification`)
          }
          description={
            <Trans
              i18nKey={
                isDisabling
                  ? 'table:notification.confirm-disable-desc'
                  : 'table:notification.confirm-enable-desc'
              }
              values={{ name: selectedNotification?.name }}
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
