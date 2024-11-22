import { useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { APIKey } from '@types';
import PageLayout from 'elements/page-layout';
import AddAPIKeyModal from './api-key-modal/add-api-key-modal';
import EditAPIKeyModal from './api-key-modal/edit-api-key-modal';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchAPIKeys } from './collection-loader/use-fetch-apikey';
import PageContent from './page-content';
import { APIKeyActionsType } from './types';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchAPIKeys({
    pageSize: 1,
    environmentNamespace: currenEnvironment.id
  });

  const [selectedAPIKey, setSelectedAPIKey] = useState<APIKey>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const onHandleActions = (apiKey: APIKey, type: APIKeyActionsType) => {
    if (type === 'EDIT') {
      onOpenEditModal();
    }
    setSelectedAPIKey(apiKey);
  };

  const isEmpty = collection?.apiKeys.length === 0;

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
        <AddAPIKeyModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenEditModal && selectedAPIKey && (
        <EditAPIKeyModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          apiKey={selectedAPIKey}
        />
      )}
    </>
  );
};

export default PageLoader;
