import { useState } from 'react';
import { Trans } from 'react-i18next';
import { useLocation } from 'react-router-dom';
import { notificationUpdater } from '@api/notification';
import { invalidateNotifications } from '@queries/notifications';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { ID_NEW, PAGE_PATH_NOTIFICATIONS } from 'constants/routing';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { Notification } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchNotifications } from './collection-loader/use-fetch-notifications';
import AddNotificationModal from './notification-modal/add-notification-modal';
import EditNotificationModal from './notification-modal/edit-notification-modal';
import PageContent from './page-content';
import { NotificationActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const location = useLocation();

  const { isAdd, isEdit, onOpenAddModal, onOpenEditModal, onCloseActionModal } =
    useActionWithURL({
      idKey: '*',
      addPath: `${location.pathname}/${String(ID_NEW)}`,
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_NOTIFICATIONS}`
    });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchNotifications({
    pageSize: 1,
    organizationId: currentEnvironment.organizationId
  });
  const [isDisabling, setIsDisabling] = useState<boolean>(false);
  const [selectedNotification, setSelectedNotification] =
    useState<Notification>();

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (
    notification: Notification,
    type: NotificationActionsType
  ) => {
    if (type === 'EDIT')
      return onOpenEditModal(`${location.pathname}/${notification.id}`, {
        environmentId: notification.environmentId
      });
    setSelectedNotification(notification);
    if (type === 'ENABLE') {
      setIsDisabling(false);
      return onOpenConfirmModal();
    }
    setIsDisabling(true);
    onOpenConfirmModal();
  };

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
      mutationState.reset();
    }
  });

  const onHandleDisable = () => {
    if (selectedNotification?.id) {
      mutationState.mutate(selectedNotification);
    }
  };

  const isEmpty = collection?.subscriptions.length === 0;

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
        <AddNotificationModal isOpen={isAdd} onClose={onCloseActionModal} />
      )}
      {isEdit && (
        <EditNotificationModal isOpen={isEdit} onClose={onCloseActionModal} />
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
        />
      )}
    </>
  );
};

export default PageLoader;
