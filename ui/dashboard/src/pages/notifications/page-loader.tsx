import { useState } from 'react';
import { Trans } from 'react-i18next';
import { notificationUpdater } from '@api/notification';
import { invalidateNotifications } from '@queries/notifications';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
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
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchNotifications({
    pageSize: 1,
    organizationId: currenEnvironment.organizationId
  });
  const [isDisabling, setIsDisabling] = useState<boolean>(false);
  const [selectedNotification, setSelectedNotification] =
    useState<Notification>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (
    notification: Notification,
    type: NotificationActionsType
  ) => {
    if (type === 'EDIT') {
      onOpenEditModal();
    } else if (type === 'ENABLE') {
      setIsDisabling(false);
      onOpenConfirmModal();
    } else if (type === 'DISABLE') {
      setIsDisabling(true);
      onOpenConfirmModal();
    }
    setSelectedNotification(notification);
  };

  const mutationState = useMutation({
    mutationFn: async (notification: Notification) => {
      return notificationUpdater({
        id: notification.id,
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

      {isOpenAddModal && (
        <AddNotificationModal
          isOpen={isOpenAddModal}
          onClose={onCloseAddModal}
        />
      )}
      {isOpenEditModal && selectedNotification && (
        <EditNotificationModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          notification={selectedNotification}
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
        />
      )}
    </>
  );
};

export default PageLoader;
