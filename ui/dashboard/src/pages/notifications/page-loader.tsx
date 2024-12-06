import { useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { Notification } from '@types';
import PageLayout from 'elements/page-layout';
// import AddNotificationModal from './slack-modal/add-api-key-modal';
// import EditNotificationModal from './slack-modal/edit-api-key-modal';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchNotifications } from './collection-loader/use-fetch-notifications';
import PageContent from './page-content';
import { NotificationActionsType } from './types';

const PageLoader = () => {
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

  const [selectedNotification, setSelectedNotification] = useState<Notification>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const onHandleActions = (apiKey: Notification, type: NotificationActionsType) => {
    if (type === 'EDIT') {
      onOpenEditModal();
    }
    setSelectedNotification(apiKey);
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

      {/* {isOpenAddModal && (
        <AddNotificationModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenEditModal && selectedNotification && (
        <EditNotificationModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          apiKey={selectedNotification}
        />
      )} */}
    </>
  );
};

export default PageLoader;
