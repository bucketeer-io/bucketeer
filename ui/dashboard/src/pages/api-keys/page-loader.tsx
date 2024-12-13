import { useState } from 'react';
import { Trans } from 'react-i18next';
import { apiKeyUpdater } from '@api/api-key';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { APIKey } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import AddAPIKeyModal from './api-key-modal/add-api-key-modal';
import EditAPIKeyModal from './api-key-modal/edit-api-key-modal';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchAPIKeys } from './collection-loader/use-fetch-apikey';
import PageContent from './page-content';
import { APIKeyActionsType } from './types';

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
  } = useFetchAPIKeys({
    pageSize: 1,
    organizationId: currenEnvironment.organizationId
  });

  const [selectedAPIKey, setSelectedAPIKey] = useState<APIKey>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (apiKey: APIKey, type: APIKeyActionsType) => {
    if (type === 'EDIT') {
      onOpenEditModal();
    } else if (type === 'ENABLE') {
      setIsDisabling(false);
      onOpenConfirmModal();
    } else if (type === 'DISABLE') {
      setIsDisabling(true);
      onOpenConfirmModal();
    }
    setSelectedAPIKey(apiKey);
  };

  const mutationState = useMutation({
    mutationFn: async (id: string) => {
      return apiKeyUpdater({
        id,
        environmentId: currenEnvironment.id,
        disabled: isDisabling
      });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateAPIKeys(queryClient);
      mutationState.reset();
    }
  });

  const onHandleDisable = () => {
    if (selectedAPIKey?.id) {
      mutationState.mutate(selectedAPIKey.id);
    }
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
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleDisable}
          title={
            isDisabling
              ? t(`table:popover.disable-api-key`)
              : t(`table:popover.enable-api-key`)
          }
          description={
            <Trans
              i18nKey={
                isDisabling
                  ? 'table:api-key.confirm-disable-desc'
                  : 'table:api-key.confirm-enable-desc'
              }
              values={{ name: selectedAPIKey?.name }}
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
