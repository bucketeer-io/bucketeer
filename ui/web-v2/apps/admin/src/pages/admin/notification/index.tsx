import { yupResolver } from '@hookform/resolvers/yup';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { AdminNotificationAddForm } from '../../../components/AdminNotificationAddForm';
import { AdminNotificationList } from '../../../components/AdminNotificationList';
import { AdminNotificationUpdateForm } from '../../../components/AdminNotificationUpdateForm';
import { ConfirmDialog } from '../../../components/ConfirmDialog';
import { Overlay } from '../../../components/Overlay';
import { NOTIFICATION_LIST_PAGE_SIZE } from '../../../constants/adminNotification';
import {
  ID_NEW,
  PAGE_PATH_ADMIN,
  PAGE_PATH_NOTIFICATIONS,
  PAGE_PATH_NEW,
} from '../../../constants/routing';
import { messages } from '../../../lang/messages';
import { AppState } from '../../../modules';
import {
  createNotification,
  deleteNotification,
  disableNotification,
  enableNotification,
  getNotification,
  listNotification,
  OrderBy,
  OrderDirection,
  selectById as selectNotificationById,
  updateNotification,
} from '../../../modules/adminNotifications';
import { ListAdminSubscriptionsRequest } from '../../../proto/notification/service_pb';
import { Subscription } from '../../../proto/notification/subscription_pb';
import { AppDispatch } from '../../../store';
import {
  AdminNotificationSortOption,
  isAdminNotificationSortOption,
} from '../../../types/adminNotification';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../../types/list';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: AdminNotificationSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAdminSubscriptionsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAdminSubscriptionsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListAdminSubscriptionsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAdminSubscriptionsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListAdminSubscriptionsRequest.OrderBy.NAME,
        orderDirection: ListAdminSubscriptionsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_NAME_DESC:
      return {
        orderBy: ListAdminSubscriptionsRequest.OrderBy.NAME,
        orderDirection: ListAdminSubscriptionsRequest.OrderDirection.DESC,
      };
    default:
      return {
        orderBy: ListAdminSubscriptionsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAdminSubscriptionsRequest.OrderDirection.DESC,
      };
  }
};

export const AdminNotificationIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const { notificationId } = useParams<{ notificationId: string }>();
  const notification = useSelector<AppState, Subscription.AsObject | undefined>(
    (state) => selectNotificationById(state.adminNotification, notificationId),
    shallowEqual
  );
  const isNew = notificationId == ID_NEW;
  const isUpdate = notificationId ? notificationId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen] =
    useState(false);
  const [isEnableConfirmDialogOpen, setIsEnableConfirmDialogOpen] =
    useState(false);
  const updateNotificationList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isAdminNotificationSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * NOTIFICATION_LIST_PAGE_SIZE;
      const disabled =
        options && options.enabled ? options.enabled === 'false' : null;
      dispatch(
        listNotification({
          pageSize: NOTIFICATION_LIST_PAGE_SIZE,
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
      updateNotificationList(options, 1);
    },
    [updateURL, updateNotificationList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateNotificationList(searchOptions, page);
    },
    [updateURL, searchOptions, updateNotificationList]
  );

  const handleOnClickAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_NOTIFICATIONS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      name: '',
      webhookUrl: '',
      sourceTypes: null,
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const add = useCallback(
    async (data) => {
      dispatch(
        createNotification({
          name: data.name,
          sourceTypes: data.sourceTypes,
          webhookUrl: data.webhookUrl,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        history.replace(`${PAGE_PATH_ADMIN}${PAGE_PATH_NOTIFICATIONS}`);
        updateNotificationList(null, 1);
      });
    },
    [dispatch, location]
  );

  const handleOnClickUpdate = useCallback(
    (s: Subscription.AsObject) => {
      setOpen(true);
      resetUpdate({
        name: s.name,
        webhookUrl: s.recipient.slackChannelRecipient.webhookUrl,
        sourceTypes: [...s.sourceTypesList].sort(),
      });
      history.push({
        pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_NOTIFICATIONS}/${s.id}`,
        search: location.search,
      });
    },
    [setOpen, history, notification, location]
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
      let sourceTypes: Array<
        Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]
      >;
      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.sourceTypes) {
        sourceTypes = data.sourceTypes;
      }
      dispatch(
        updateNotification({
          id: notificationId,
          name: name,
          currentSourceTypes: notification.sourceTypesList,
          sourceTypes: sourceTypes,
        })
      ).then(() => {
        dispatch(
          getNotification({
            id: notificationId,
          })
        );
        handleOnClose();
      });
    },
    [dispatch, notification, notificationId]
  );

  const handleOnClose = useCallback(() => {
    resetAdd();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_NOTIFICATIONS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd]);

  const switchMethod = useForm({
    defaultValues: {
      notification: null,
    },
    mode: 'onChange',
  });

  const { handleSubmit: switchHandleSubmit, setValue: switchSetValue } =
    switchMethod;

  const handleOnSwitch = useCallback(
    (s: Subscription.AsObject) => {
      switchSetValue('notification', s);
      setIsEnableConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleSwitch = useCallback(
    (data) => {
      dispatch(
        data.notification.disabled
          ? enableNotification({
              id: data.notification.id,
            })
          : disableNotification({
              id: data.notification.id,
            })
      ).then(() => {
        updateNotificationList(
          searchOptions,
          searchOptions.page ? Number(searchOptions.page) : 1
        );
        setIsEnableConfirmDialogOpen(false);
      });
    },
    [dispatch, setIsEnableConfirmDialogOpen]
  );

  const deleteMethod = useForm({
    defaultValues: {
      notification: null,
    },
    mode: 'onChange',
  });

  const { handleSubmit: deleteHandleSubmit, setValue: deleteSetValue } =
    deleteMethod;

  const handleOnClickDelete = useCallback(
    (s: Subscription.AsObject) => {
      deleteSetValue('notification', s);
      setIsDeleteConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleDelete = useCallback(
    (data) => {
      dispatch(
        deleteNotification({
          id: data.notification.id,
        })
      ).then(() => {
        updateNotificationList(
          searchOptions,
          searchOptions.page ? Number(searchOptions.page) : 1
        );
        setIsDeleteConfirmDialogOpen(false);
      });
    },
    [dispatch, setIsDeleteConfirmDialogOpen]
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
    updateNotificationList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateNotificationList]);

  return (
    <>
      <div className="m-10">
        <AdminNotificationList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onChangeSearchOptions={handleSearchOptionsChange}
          onAdd={handleOnClickAdd}
          onUpdate={handleOnClickUpdate}
          onSwitch={handleOnSwitch}
          onDelete={handleOnClickDelete}
        />
      </div>
      <Overlay open={open} onClose={handleOnClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <AdminNotificationAddForm
              onSubmit={handleAddSubmit(add)}
              onCancel={handleOnClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <AdminNotificationUpdateForm
              onSubmit={handleUpdateSubmit(update)}
              onCancel={handleOnClose}
            />
          </FormProvider>
        )}
      </Overlay>
      <ConfirmDialog
        open={isDeleteConfirmDialogOpen}
        onConfirm={deleteHandleSubmit(handleDelete)}
        onClose={() => setIsDeleteConfirmDialogOpen(false)}
        title={f(messages.notification.confirm.deleteTitle)}
        description={f(messages.notification.confirm.deleteDescription, {
          notificationName:
            deleteMethod.getValues().notification &&
            deleteMethod.getValues().notification.name,
        })}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
      <ConfirmDialog
        open={isEnableConfirmDialogOpen}
        onConfirm={switchHandleSubmit(handleSwitch)}
        onClose={() => setIsEnableConfirmDialogOpen(false)}
        title={f(
          switchMethod.getValues().notification &&
            switchMethod.getValues().notification.disabled
            ? messages.notification.confirm.enableTitle
            : messages.notification.confirm.disableTitle
        )}
        description={f(
          switchMethod.getValues().notification &&
            switchMethod.getValues().notification.disabled
            ? messages.notification.confirm.enableDescription
            : messages.notification.confirm.disableDescription,
          {
            notificationName:
              switchMethod.getValues().notification &&
              switchMethod.getValues().notification.name,
          }
        )}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
