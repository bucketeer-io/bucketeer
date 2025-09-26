import { useCallback, useEffect, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { apiKeyUpdater } from '@api/api-key';
import { useQueryAPIKey } from '@queries/api-key-details';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getAccountAccess,
  getCurrentEnvironment,
  getEditorEnvironments,
  useAuth
} from 'auth';
import { ID_NEW, PAGE_PATH_APIKEYS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { APIKey } from '@types';
import { useSearchParams } from 'utils/search-params';
import ConfirmModal from 'elements/confirm-modal';
import APIKeyCreateUpdateModal from './api-key-modal/api-key-create-update-modal';
import PageContent from './page-content';
import { APIKeyActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { editorEnvironments } = getEditorEnvironments(consoleAccount!);
  const navigate = useNavigate();

  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );
  const commonPath = useMemo(
    () => `/${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`,
    [currentEnvironment]
  );

  const isDisabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const {
    id: apiKeyId,
    isAdd,
    isEdit,
    onCloseActionModal,
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

  const handleOpenAddModal = useCallback(
    () => navigate(`${commonPath}/${ID_NEW}`),
    [commonPath]
  );

  const handleOnCloseModal = useCallback(() => {
    onCloseActionModal();
    onCloseConfirmModal();
  }, []);

  const onHandleActions = useCallback(
    (apiKey: APIKey, type: APIKeyActionsType) => {
      setSelectedAPIKey(apiKey);
      const environment = editorEnvironments.find(
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
    [editorEnvironments, commonPath]
  );

  const mutationState = useMutation({
    mutationFn: async (selectedAPIKey: APIKey) => {
      const environmentId = editorEnvironments.find(
        item => item.name === selectedAPIKey.environmentName
      )?.id;
      return apiKeyUpdater({
        id: selectedAPIKey.id,
        environmentId,
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
    if (selectedAPIKey) {
      mutationState.mutate(selectedAPIKey);
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
      handleOnCloseModal();
    }
  }, [isError, error]);

  return (
    <>
      <PageContent
        onAdd={handleOpenAddModal}
        onHandleActions={onHandleActions}
      />
      {(!!isAdd || !!isEdit) && (
        <APIKeyCreateUpdateModal
          isOpen={!!isAdd || !!isEdit}
          apiKeyEnvironmentId={apiKeyEnvironmentId as string}
          isLoadingApiKey={isLoadingApiKey}
          apiKey={selectedAPIKey}
          environments={editorEnvironments}
          resetApiKey={() => setSelectedAPIKey(undefined)}
          onClose={handleOnCloseModal}
        />
      )}
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={handleOnCloseModal}
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
          disabled={isDisabled}
        />
      )}
    </>
  );
};

export default PageLoader;
