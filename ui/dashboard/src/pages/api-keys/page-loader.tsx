import { useCallback, useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { apiKeyUpdater } from '@api/api-key';
import { useQueryAPIKey } from '@queries/api-key-details';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_APIKEYS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { APIKey } from '@types';
import { onFormatEnvironments } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import ConfirmModal from 'elements/confirm-modal';
import AddAPIKeyModal from './api-key-modal/add-api-key-modal';
import EditAPIKeyModal from './api-key-modal/edit-api-key-modal';
import PageContent from './page-content';
import { APIKeyActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const commonPath = useMemo(
    () => `/${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`,
    [currentEnvironment]
  );

  const {
    id: apiKeyId,
    isAdd,
    isEdit,
    onCloseActionModal,
    onOpenAddModal,
    onOpenEditModal
  } = useActionWithURL({
    closeModalPath: commonPath
  });
  const { notify, errorNotify } = useToast();

  const [selectedAPIKey, setSelectedAPIKey] = useState<APIKey>();
  const [isDisabling, setIsDisabling] = useState<boolean>(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);
  const { searchOptions } = useSearchParams();

  const { data: collection, isLoading: isLoadingEnvs } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = collection?.environments || [];
  const { formattedEnvironments } = onFormatEnvironments(environments);

  const apiKeyEnvironmentId = searchOptions?.environmentId;

  const {
    data: apiKeyCollection,
    isLoading: isLoadingApiKey,
    isError,
    error
  } = useQueryAPIKey({
    params: {
      environmentId: apiKeyEnvironmentId as string,
      id: apiKeyId as string
    },
    enabled: !!isEdit && !!apiKeyId && !selectedAPIKey
  });

  const onHandleActions = useCallback(
    (apiKey: APIKey, type: APIKeyActionsType) => {
      setSelectedAPIKey(apiKey);
      const environment = formattedEnvironments.find(
        item => item.name === apiKey.environmentName
      );
      switch (type) {
        case 'EDIT':
          return onOpenEditModal(
            `${commonPath}/${apiKey.id}?environmentId=${environment?.id}`
          );
        case 'ENABLE':
        case 'DISABLE':
          setIsDisabling(type === 'DISABLE');
          return onOpenConfirmModal();
        default:
          break;
      }
    },
    [formattedEnvironments, commonPath]
  );

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
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:source-type.push'),
          action: t('common:updated')
        })
      });
    },
    onError: error => errorNotify(error)
  });

  const onHandleDisable = useCallback(() => {
    if (selectedAPIKey?.id) {
      mutationState.mutate(selectedAPIKey.id);
    }
  }, [selectedAPIKey]);

  useEffect(() => {
    if (apiKeyCollection) {
      setSelectedAPIKey(apiKeyCollection.apiKey);
    }
  }, [apiKeyCollection]);

  useEffect(() => {
    if (isError && error) {
      errorNotify(error);
      onCloseActionModal();
    }
  }, [isError, error]);

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      {isAdd && (
        <AddAPIKeyModal
          isOpen={isAdd}
          isLoadingEnvs={isLoadingEnvs}
          environments={formattedEnvironments}
          onClose={onCloseActionModal}
        />
      )}
      {isEdit && (
        <EditAPIKeyModal
          isOpen={isEdit}
          isLoadingApiKey={isLoadingApiKey}
          apiKey={selectedAPIKey}
          environments={formattedEnvironments}
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
