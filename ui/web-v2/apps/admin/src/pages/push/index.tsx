import { listTags } from '@/modules/tags';
import { ListTagsRequest } from '@/proto/feature/service_pb';
import { yupResolver } from '@hookform/resolvers/yup';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { ConfirmDialog } from '../../components/ConfirmDialog';
import { Overlay } from '../../components/Overlay';
import { PushAddForm } from '../../components/PushAddForm';
import { PushList } from '../../components/PushList';
import { PushUpdateForm } from '../../components/PushUpdateForm';
import { PUSH_LIST_PAGE_SIZE } from '../../constants/push';
import {
  ID_NEW,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_PUSHES,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useCurrentEnvironment } from '../../modules/me';
import {
  selectById as selectPushById,
  listPushes,
  createPush,
  updatePush,
  deletePush,
  OrderBy,
  OrderDirection,
} from '../../modules/pushes';
import { Push } from '../../proto/push/push_pb';
import { ListPushesRequest } from '../../proto/push/service_pb';
import { AppDispatch } from '../../store';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
  SORT_OPTIONS_NAME_DESC,
} from '../../types/list';
import { PushSortOption, isPushSortOption } from '../../types/push';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: PushSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListPushesRequest.OrderBy.CREATED_AT,
        orderDirection: ListPushesRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListPushesRequest.OrderBy.CREATED_AT,
        orderDirection: ListPushesRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListPushesRequest.OrderBy.NAME,
        orderDirection: ListPushesRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_NAME_DESC:
      return {
        orderBy: ListPushesRequest.OrderBy.NAME,
        orderDirection: ListPushesRequest.OrderDirection.DESC,
      };
    default:
      return {
        orderBy: ListPushesRequest.OrderBy.CREATED_AT,
        orderDirection: ListPushesRequest.OrderDirection.DESC,
      };
  }
};

export const PushIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const currentEnvironment = useCurrentEnvironment();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const { pushId } = useParams<{ pushId: string }>();
  const isNew = pushId == ID_NEW;
  const isUpdate = pushId ? pushId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);
  const push = useSelector<AppState, Push.AsObject | undefined>(
    (state) => selectPushById(state.push, pushId),
    shallowEqual
  );
  const updatePushList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isPushSortOption(options && options.sort) ? options.sort : '-createdAt'
      );
      const cursor = (page - 1) * PUSH_LIST_PAGE_SIZE;
      dispatch(
        listPushes({
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
      updatePushList(options, 1);
    },
    [updateURL, updatePushList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updatePushList(searchOptions, page);
    },
    [updateURL, searchOptions, updatePushList]
  );

  const handleOnClickAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_PUSHES}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      name: '',
      fcmApiKey: '',
      tags: null,
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const add = useCallback(
    async (data) => {
      dispatch(
        createPush({
          environmentNamespace: currentEnvironment.id,
          name: data.name,
          fcmApiKey: data.fcmApiKey,
          tags: data.tags,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_PUSHES}`
        );
        updatePushList(null, 1);
      });
    },
    [dispatch, location]
  );

  const handleOnClickUpdate = useCallback(
    (p: Push.AsObject) => {
      setOpen(true);
      resetUpdate({
        name: p.name,
        fcmApiKey: p.fcmApiKey,
        tags: p.tagsList,
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_PUSHES}/${p.id}`,
        search: location.search,
      });
    },
    [setOpen, history, push, location]
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
      let tags: Array<string>;
      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.tags) {
        tags = data.tags;
      }
      dispatch(
        updatePush({
          environmentNamespace: currentEnvironment.id,
          id: pushId,
          name: name,
          currentTags: push.tagsList,
          tags: tags,
        })
      ).then(() => {
        updatePushList(
          searchOptions,
          searchOptions.page ? Number(searchOptions.page) : 1
        );
        handleOnClose();
      });
    },
    [dispatch, push, pushId, updatePushList]
  );

  const handleOnClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${PAGE_PATH_PUSHES}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd]);

  const deleteMethod = useForm({
    defaultValues: {
      push: null,
    },
    mode: 'onChange',
  });

  const { handleSubmit: deleteHandleSubmit, setValue: deleteSetValue } =
    deleteMethod;

  const handleOnClickDelete = useCallback(
    (p: Push.AsObject) => {
      deleteSetValue('push', p);
      setIsConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleDeletePush = useCallback(
    (data) => {
      dispatch(
        deletePush({
          environmentNamespace: currentEnvironment.id,
          id: data.push.id,
        })
      ).then(() => {
        updatePushList(null, 1);
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
    updatePushList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updatePushList]);

  useEffect(() => {
    dispatch(
      listTags({
        environmentNamespace: currentEnvironment.id,
        pageSize: 99999,
        cursor: '',
        orderBy: ListTagsRequest.OrderBy.DEFAULT,
        orderDirection: ListTagsRequest.OrderDirection.ASC,
        searchKeyword: null,
      })
    );
  }, []);

  return (
    <>
      <div className="m-10">
        <PushList
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
              <PushAddForm
                onSubmit={handleAddSubmit(add)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
          {isUpdate && (
            <FormProvider {...updateMethod}>
              <PushUpdateForm
                onSubmit={handleUpdateSubmit(update)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
        </Overlay>
      </div>
      <ConfirmDialog
        open={isConfirmDialogOpen}
        onConfirm={deleteHandleSubmit(handleDeletePush)}
        onClose={() => setIsConfirmDialogOpen(false)}
        title={f(messages.push.confirm.deleteTitle)}
        description={f(messages.push.confirm.deleteDescription, {
          pushName:
            deleteMethod.getValues().push && deleteMethod.getValues().push.name,
        })}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
