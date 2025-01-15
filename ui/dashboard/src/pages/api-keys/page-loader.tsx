import { useState } from 'react';
import { Trans } from 'react-i18next';
import { useLocation } from 'react-router-dom';
import { apiKeyUpdater } from '@api/api-key';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { ID_NEW, PAGE_PATH_APIKEYS } from 'constants/routing';
import useActionWithURL from 'hooks/use-action-with-url';
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
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const location = useLocation();

  const { isAdd, isEdit, onOpenAddModal, onOpenEditModal, onCloseActionModal } =
    useActionWithURL({
      idKey: '*',
      addPath: `${location.pathname}/${String(ID_NEW)}`,
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`
    });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchAPIKeys({
    pageSize: 1,
    organizationId: currentEnvironment.organizationId
  });

  const [selectedAPIKey, setSelectedAPIKey] = useState<APIKey>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const onHandleActions = (apiKey: APIKey, type: APIKeyActionsType) => {
    if (type === 'EDIT') {
      return onOpenEditModal(`${location.pathname}/${apiKey.id}`, {
        environmentName: apiKey.environmentName
      });
    }
    setSelectedAPIKey(apiKey);
    onOpenConfirmModal();
    if (type === 'ENABLE') {
      return setIsDisabling(false);
    }
    setIsDisabling(true);
  };

  const mutationState = useMutation({
    mutationFn: async (id: string) => {
      return apiKeyUpdater({
        id,
        environmentId: currentEnvironment.id,
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

      {isAdd && <AddAPIKeyModal isOpen={isAdd} onClose={onCloseActionModal} />}
      {isEdit && (
        <EditAPIKeyModal isOpen={isEdit} onClose={onCloseActionModal} />
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
