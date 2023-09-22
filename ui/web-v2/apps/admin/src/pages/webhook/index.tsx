import { yupResolver } from '@hookform/resolvers/yup';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { ConfirmDialog } from '../../components/ConfirmDialog';
import { Overlay } from '../../components/Overlay';
import { WebhookAddForm } from '../../components/WebhookAddForm';
import { WebhookList } from '../../components/WebhookList';
import { WebhookUpdateForm } from '../../components/WebhookUpdateForm';
import { PUSH_LIST_PAGE_SIZE } from '../../constants/push';
import {
  ID_NEW,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
  PAGE_PATH_WEBHOOKS,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { useCurrentEnvironment } from '../../modules/me';
import {
  listWebhooks,
  createWebhook,
  updateWebhook,
  deleteWebhook,
  OrderBy,
  OrderDirection,
} from '../../modules/webhooks';
import { ListWebhooksRequest } from '../../proto/autoops/service_pb';
import { Webhook } from '../../proto/autoops/webhook_pb';
import { AppDispatch } from '../../store';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { WebhookSortOption, isWebhookSortOption } from '../../types/webhook';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: WebhookSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListWebhooksRequest.OrderBy.CREATED_AT,
        orderDirection: ListWebhooksRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListWebhooksRequest.OrderBy.CREATED_AT,
        orderDirection: ListWebhooksRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListWebhooksRequest.OrderBy.NAME,
        orderDirection: ListWebhooksRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_NAME_DESC:
      return {
        orderBy: ListWebhooksRequest.OrderBy.NAME,
        orderDirection: ListWebhooksRequest.OrderDirection.DESC,
      };
    default:
      return {
        orderBy: ListWebhooksRequest.OrderBy.CREATED_AT,
        orderDirection: ListWebhooksRequest.OrderDirection.DESC,
      };
  }
};

export const WebhookIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const currentEnvironment = useCurrentEnvironment();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const { webhookId } = useParams<{ webhookId: string }>();
  const isNew = webhookId == ID_NEW;
  const isUpdate = webhookId ? webhookId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

  const updateWebhookList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isWebhookSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * PUSH_LIST_PAGE_SIZE;
      dispatch(
        listWebhooks({
          environmentNamespace: currentEnvironment.id,
          pageSize: PUSH_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
        })
      );
    },
    [dispatch]
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

  const handleSearchOptionsChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateWebhookList(options, 1);
    },
    [updateURL, updateWebhookList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateWebhookList(searchOptions, page);
    },
    [updateURL, searchOptions, updateWebhookList]
  );

  const handleOnClickAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_WEBHOOKS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      name: '',
      description: '',
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const add = useCallback(
    async (data) => {
      dispatch(
        createWebhook({
          environmentNamespace: currentEnvironment.id,
          name: data.name,
          description: data.description,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_WEBHOOKS}`
        );
        updateWebhookList(null, 1);
      });
    },
    [dispatch, location]
  );

  const handleOnClickUpdate = useCallback(
    (w: Webhook.AsObject) => {
      setOpen(true);
      resetUpdate({
        name: w.name,
        description: w.description,
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_WEBHOOKS}/${w.id}`,
        search: location.search,
      });
    },
    [setOpen, history, location]
  );

  const updateMethod = useForm({
    resolver: yupResolver(updateFormSchema),
    mode: 'onChange',
  });

  const {
    handleSubmit: handleUpdateSubmit,
    formState: { dirtyFields },
    reset: resetUpdate,
  } = updateMethod;

  const update = useCallback(
    async (data) => {
      let name: string;
      let description: string;

      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.description) {
        description = data.description;
      }
      dispatch(
        updateWebhook({
          environmentNamespace: currentEnvironment.id,
          id: webhookId,
          name: name,
          description,
        })
      ).then(() => {
        updateWebhookList(
          searchOptions,
          searchOptions.page ? Number(searchOptions.page) : 1
        );
        handleOnClose();
      });
    },
    [dispatch, webhookId, updateWebhookList]
  );

  const handleOnClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_WEBHOOKS}`,
      search: location.search,
    });
  }, [setOpen, resetAdd, resetUpdate]);

  const deleteMethod = useForm({
    defaultValues: {
      webhook: null,
    },
    mode: 'onChange',
  });

  const { handleSubmit: deleteHandleSubmit, setValue: deleteSetValue } =
    deleteMethod;

  const handleOnClickDelete = useCallback(
    (w: Webhook.AsObject) => {
      deleteSetValue('webhook', w);
      setIsConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleDeleteWebhook = useCallback(
    (data) => {
      dispatch(
        deleteWebhook({
          environmentNamespace: currentEnvironment.id,
          id: data.webhook.id,
        })
      ).then(() => {
        updateWebhookList(null, 1);
        setIsConfirmDialogOpen(false);
      });
    },
    [dispatch, setIsConfirmDialogOpen]
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
    updateWebhookList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateWebhookList]);

  return (
    <>
      <div className="m-10">
        <WebhookList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onChangeSearchOptions={handleSearchOptionsChange}
          onAdd={handleOnClickAdd}
          onUpdate={handleOnClickUpdate}
          onDelete={handleOnClickDelete}
        />
        <Overlay open={open} onClose={handleOnClose}>
          {isNew && (
            <FormProvider {...addMethod}>
              <WebhookAddForm
                onSubmit={handleAddSubmit(add)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
          {isUpdate && (
            <FormProvider {...updateMethod}>
              <WebhookUpdateForm
                onSubmit={handleUpdateSubmit(update)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
        </Overlay>
      </div>
      <ConfirmDialog
        open={isConfirmDialogOpen}
        onConfirm={deleteHandleSubmit(handleDeleteWebhook)}
        onClose={() => setIsConfirmDialogOpen(false)}
        title={f(messages.webhook.confirm.deleteTitle)}
        description={f(messages.webhook.confirm.deleteDescription, {
          webhookName:
            deleteMethod.getValues().webhook &&
            deleteMethod.getValues().webhook.name,
        })}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
