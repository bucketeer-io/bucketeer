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

import { AccountAddForm } from '../../components/AccountAddForm';
import { AccountList } from '../../components/AccountList';
import { AccountUpdateForm } from '../../components/AccountUpdateForm';
import { ConfirmDialog } from '../../components/ConfirmDialog';
import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import { ACCOUNT_LIST_PAGE_SIZE } from '../../constants/account';
import {
  ID_NEW,
  PAGE_PATH_ACCOUNTS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectAccountById,
  listAccounts,
  disableAccount,
  enableAccount,
  getAccount,
  createAccount,
  updateAccount,
  OrderBy,
  OrderDirection,
} from '../../modules/accounts';
import { useCurrentEnvironment } from '../../modules/me';
import { Account } from '../../proto/account/account_pb';
import { ListAccountsRequest } from '../../proto/account/service_pb';
import { AppDispatch } from '../../store';
import { AccountSortOption, isAccountSortOption } from '../../types/account';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: AccountSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAccountsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAccountsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListAccountsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAccountsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListAccountsRequest.OrderBy.EMAIL,
        orderDirection: ListAccountsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListAccountsRequest.OrderBy.EMAIL,
        orderDirection: ListAccountsRequest.OrderDirection.DESC,
      };
  }
};

export const AccountIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const history = useHistory();
  const location = useLocation();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const { url } = useRouteMatch();
  const { accountId } = useParams<{ accountId: string }>();
  const isNew = accountId == ID_NEW;
  const isUpdate = accountId ? accountId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [account, getAccountError] = useSelector<
    AppState,
    [Account.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectAccountById(state.accounts, accountId),
      state.accounts.getAccountError,
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

  const updateAccountList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isAccountSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * ACCOUNT_LIST_PAGE_SIZE;
      const role =
        options && options.role != null ? Number(options.role) : null;
      const disabled =
        options && options.enabled ? options.enabled === 'false' : null;
      dispatch(
        listAccounts({
          environmentNamespace: currentEnvironment.id,
          pageSize: ACCOUNT_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          disabled: disabled,
          role: role,
        })
      );
    },
    [dispatch]
  );

  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

  const switchEnabledMethod = useForm({
    defaultValues: {
      accountId: '',
      enabled: false,
    },
    mode: 'onChange',
  });

  const {
    handleSubmit: switchEnableHandleSubmit,
    setValue: switchEnabledSetValue,
  } = switchEnabledMethod;

  const handleClickSwitchEnabled = useCallback(
    (accountId: string, enabled: boolean) => {
      switchEnabledSetValue('accountId', accountId);
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
            return enableAccount({
              environmentNamespace: currentEnvironment.id,
              id: data.accountId,
            });
          }
          return disableAccount({
            environmentNamespace: currentEnvironment.id,
            id: data.accountId,
          });
        })()
      ).then(() => {
        setIsConfirmDialogOpen(false);
        dispatch(
          getAccount({
            environmentNamespace: currentEnvironment.id,
            email: data.accountId,
          })
        );
      });
    },
    [dispatch, setIsConfirmDialogOpen]
  );

  const handleSearchOptionsChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateAccountList(options, 1);
    },
    [updateURL, updateAccountList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateAccountList(searchOptions, page);
    },
    [updateURL, updateAccountList, searchOptions]
  );

  const handleOpenAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleOpenUpdate = useCallback(
    (a: Account.AsObject) => {
      setOpen(true);
      resetUpdate({
        email: a.email,
        role: a.role.toString(),
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}/${a.id}`,
        search: location.search,
      });
    },
    [setOpen, history, account, location]
  );

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      email: null,
      role: null,
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const updateMethod = useForm({
    resolver: yupResolver(updateFormSchema),
    mode: 'onChange',
  });
  const { handleSubmit: handleUpdateSubmit, reset: resetUpdate } = updateMethod;

  const handleClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd, resetUpdate]);

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createAccount({
          environmentNamespace: currentEnvironment.id,
          email: data.email,
          role: data.role,
        })
      ).then(() => {
        resetAdd();
        setOpen(false);
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}`
        );
        updateAccountList(null, 1);
      });
    },
    [dispatch]
  );

  const handleUpdate = useCallback(
    async (data) => {
      dispatch(
        updateAccount({
          environmentNamespace: currentEnvironment.id,
          id: accountId,
          role: data.role,
        })
      ).then(() => {
        dispatch(
          getAccount({
            environmentNamespace: currentEnvironment.id,
            email: accountId,
          })
        );
        handleClose();
      });
    },
    [dispatch, accountId]
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
    updateAccountList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateAccountList]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.account.list.header.title)}
          description={f(messages.account.list.header.description)}
        />
      </div>
      <div className="m-10">
        <AccountList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onSwitchEnabled={handleClickSwitchEnabled}
          onAdd={handleOpenAdd}
          onUpdate={handleOpenUpdate}
          onChangeSearchOptions={handleSearchOptionsChange}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <AccountAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <AccountUpdateForm
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
            ? f(messages.account.confirm.enableTitle)
            : f(messages.account.confirm.disableTitle)
        }
        description={
          '「' +
          switchEnabledMethod.getValues().accountId +
          '」' +
          (switchEnabledMethod.getValues().enabled
            ? f(messages.account.confirm.enableDescription)
            : f(messages.account.confirm.disableDescription))
        }
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
