import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { useCallback, FC, memo, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import {
  useHistory,
  useRouteMatch,
  useLocation,
  useParams,
} from 'react-router-dom';

import { APIKeyAddForm } from '../../components/APIKeyAddForm';
import { APIKeyList } from '../../components/APIKeyList';
import { APIKeyUpdateForm } from '../../components/APIKeyUpdateForm';
import { ConfirmDialog } from '../../components/ConfirmDialog';
import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import { APIKEY_LIST_PAGE_SIZE } from '../../constants/apiKey';
import {
  ID_NEW,
  PAGE_PATH_APIKEYS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectAPIKeyById,
  listAPIKeys,
  disableAPIKey,
  enableAPIKey,
  getAPIKey,
  createAPIKey,
  updateAPIKey,
  OrderBy,
  OrderDirection,
} from '../../modules/apiKeys';
import { useCurrentEnvironment } from '../../modules/me';
import { APIKey } from '../../proto/account/api_key_pb';
import { ListAPIKeysRequest } from '../../proto/account/service_pb';
import { AppDispatch } from '../../store';
import { APIKeySortOption, isAPIKeySortOption } from '../../types/apiKey';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import {
  addApiKeyFormSchema,
  AddApiKeyFormSchema,
  updateApiKeyFormSchema,
  UpdateApiKeyFormSchema,
} from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: APIKeySortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAPIKeysRequest.OrderBy.CREATED_AT,
        orderDirection: ListAPIKeysRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListAPIKeysRequest.OrderBy.CREATED_AT,
        orderDirection: ListAPIKeysRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListAPIKeysRequest.OrderBy.NAME,
        orderDirection: ListAPIKeysRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListAPIKeysRequest.OrderBy.NAME,
        orderDirection: ListAPIKeysRequest.OrderDirection.DESC,
      };
  }
};

export const APIKeyIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const history = useHistory();
  const location = useLocation();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const { url } = useRouteMatch();
  const { apiKeyId } = useParams<{ apiKeyId: string }>();
  const isNew = apiKeyId == ID_NEW;
  const isUpdate = apiKeyId ? apiKeyId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [apiKey, getAPIKeyError] = useSelector<
    AppState,
    [APIKey.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectAPIKeyById(state.apiKeys, apiKeyId),
      state.apiKeys.getAPIKeyError,
    ],
    shallowEqual
  );
  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      history.replace(
        `${url}?${stringifySearchParams({
          ...options,
        })}`
      );
    },
    [history]
  );

  const updateAPIKeyList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isAPIKeySortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * APIKEY_LIST_PAGE_SIZE;
      const disabled =
        options && options.enabled ? options.enabled === 'false' : null;
      dispatch(
        listAPIKeys({
          environmentNamespace: currentEnvironment.id,
          pageSize: APIKEY_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          disabled: disabled,
        })
      );
    },
    [dispatch]
  );

  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

  const switchEnabledMethod = useForm({
    defaultValues: {
      apiKeyId: '',
      apiKeyName: '',
      enabled: false,
    },
    mode: 'onChange',
  });

  const {
    handleSubmit: switchEnableHandleSubmit,
    setValue: switchEnabledSetValue,
  } = switchEnabledMethod;

  const handleClickSwitchEnabled = useCallback(
    (apiKeyId: string, apiKeyName: string, enabled: boolean) => {
      switchEnabledSetValue('apiKeyId', apiKeyId);
      switchEnabledSetValue('apiKeyName', apiKeyName);
      switchEnabledSetValue('enabled', enabled);
      setIsConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleSwitchEnabled = useCallback(
    async (data) => {
      dispatch(
        (() => {
          if (data.enabled) {
            return enableAPIKey({
              environmentNamespace: currentEnvironment.id,
              id: data.apiKeyId,
            });
          }
          return disableAPIKey({
            environmentNamespace: currentEnvironment.id,
            id: data.apiKeyId,
          });
        })()
      ).then(() => {
        setIsConfirmDialogOpen(false);
        dispatch(
          getAPIKey({
            environmentNamespace: currentEnvironment.id,
            id: data.apiKeyId,
          })
        );
      });
    },
    [dispatch, setIsConfirmDialogOpen]
  );

  const handleSearchOptionsChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateAPIKeyList(options, 1);
    },
    [updateURL]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateAPIKeyList(searchOptions, page);
    },
    [updateURL, searchOptions]
  );

  const addMethod = useForm({
    resolver: yupResolver(addApiKeyFormSchema),
    defaultValues: {
      name: '',
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const updateMethod = useForm({
    resolver: yupResolver(updateApiKeyFormSchema),
    mode: 'onChange',
  });
  const { handleSubmit: handleUpdateSubmit, reset: resetUpdate } = updateMethod;

  const handleOpenAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleOpenUpdate = useCallback(
    (a: APIKey.AsObject) => {
      setOpen(true);
      resetUpdate({
        name: a.name,
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}/${a.id}`,
        search: location.search,
      });
    },
    [setOpen, history, apiKey, location]
  );

  const handleClose = useCallback(() => {
    setOpen(false);
    resetAdd();
    resetUpdate();
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd, resetUpdate]);

  const handleAdd = useCallback(
    async (data: AddApiKeyFormSchema) => {
      dispatch(
        createAPIKey({
          environmentNamespace: currentEnvironment.id,
          name: data.name,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`
        );
        updateAPIKeyList(null, 1);
      });
    },
    [dispatch, location]
  );

  const handleUpdate = useCallback(
    async (data: UpdateApiKeyFormSchema) => {
      dispatch(
        updateAPIKey({
          environmentNamespace: currentEnvironment.id,
          id: apiKeyId,
          name: data.name,
        })
      ).then(() => {
        dispatch(
          getAPIKey({
            environmentNamespace: currentEnvironment.id,
            id: apiKeyId,
          })
        );
        handleClose();
      });
    },
    [dispatch, getAPIKey, apiKeyId]
  );

  useEffect(() => {
    history.listen(() => {
      // Handle browser's back button
      if (history.action === 'POP') {
        if (open) {
          setOpen(false);
        }
      }
    });
  });

  useEffect(() => {
    updateAPIKeyList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [dispatch, searchOptions]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.apiKey.list.header.title)}
          description={f(messages.apiKey.list.header.description)}
        />
      </div>
      <div className="m-10">
        <APIKeyList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onSwitchEnabled={handleClickSwitchEnabled}
          onAdd={handleOpenAdd}
          onChangeSearchOptions={handleSearchOptionsChange}
          onUpdate={handleOpenUpdate}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <APIKeyAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <APIKeyUpdateForm
              onSubmit={handleUpdateSubmit(handleUpdate)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
      </Overlay>
      <ConfirmDialog
        open={isConfirmDialogOpen}
        onConfirm={switchEnableHandleSubmit(handleSwitchEnabled)}
        onClose={() => setIsConfirmDialogOpen(false)}
        title={
          switchEnabledMethod.getValues().enabled
            ? f(messages.apiKey.confirm.enableTitle)
            : f(messages.apiKey.confirm.disableTitle)
        }
        description={
          '「' +
          switchEnabledMethod.getValues().apiKeyName +
          '」' +
          (switchEnabledMethod.getValues().enabled
            ? f(messages.apiKey.confirm.enableDescription)
            : f(messages.apiKey.confirm.disableDescription))
        }
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
