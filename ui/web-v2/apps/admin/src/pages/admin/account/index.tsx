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

import { AdminAccountAddForm } from '../../../components/AdminAccountAddForm';
import { AdminAccountList } from '../../../components/AdminAccountList';
import { ConfirmDialog } from '../../../components/ConfirmDialog';
import { Overlay } from '../../../components/Overlay';
import { ACCOUNT_LIST_PAGE_SIZE } from '../../../constants/account';
import {
  ID_NEW,
  PAGE_PATH_ACCOUNTS,
  PAGE_PATH_NEW,
  PAGE_PATH_ADMIN,
} from '../../../constants/routing';
import { messages } from '../../../lang/messages';
import { AppState } from '../../../modules';
import {
  selectById as selectAccountById,
  listAccounts,
  disableAccount,
  enableAccount,
  getAccount,
  createAccount,
  OrderBy,
  OrderDirection,
} from '../../../modules/adminAccounts';
import { useCurrentEnvironment } from '../../../modules/me';
import { Account } from '../../../proto/account/account_pb';
import { ListAdminAccountsRequest } from '../../../proto/account/service_pb';
import { AppDispatch } from '../../../store';
import { AccountSortOption, isAccountSortOption } from '../../../types/account';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../../utils/search-params';

import { addFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: AccountSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAdminAccountsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAdminAccountsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListAdminAccountsRequest.OrderBy.CREATED_AT,
        orderDirection: ListAdminAccountsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListAdminAccountsRequest.OrderBy.EMAIL,
        orderDirection: ListAdminAccountsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListAdminAccountsRequest.OrderBy.EMAIL,
        orderDirection: ListAdminAccountsRequest.OrderDirection.DESC,
      };
  }
};

export const AdminAccountIndexPage: FC = memo(() => {
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
          pageSize: ACCOUNT_LIST_PAGE_SIZE,
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
              id: data.accountId,
            });
          }
          return disableAccount({
            id: data.accountId,
          });
        })()
      ).then(() => {
        setIsConfirmDialogOpen(false);
        dispatch(
          getAccount({
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
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_ACCOUNTS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      email: null,
      role: null,
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const handleClose = useCallback(() => {
    resetAdd();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_ACCOUNTS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd]);

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createAccount({
          email: data.email,
        })
      ).then(() => {
        resetAdd();
        setOpen(false);
        history.replace(`${PAGE_PATH_ADMIN}${PAGE_PATH_ACCOUNTS}`);
        updateAccountList(null, 1);
      });
    },
    [dispatch]
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
      <div className="flex items-stretch m-10 text-sm">
        <p className="text-gray-700">
          {f(messages.adminAccount.list.header.description)}
        </p>
        <a
          className="link"
          target="_blank"
          href="https://bucketeer.io/docs/#/managing-teams"
          rel="noreferrer"
        >
          {f(messages.readMore)}
        </a>
      </div>
      <div className="m-10">
        <AdminAccountList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onSwitchEnabled={handleClickSwitchEnabled}
          onAdd={handleOpenAdd}
          onChangeSearchOptions={handleSearchOptionsChange}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <AdminAccountAddForm
              onSubmit={handleAddSubmit(handleAdd)}
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
            ? f(messages.adminAccount.confirm.enableTitle)
            : f(messages.adminAccount.confirm.disableTitle)
        }
        description={f(
          switchEnabledMethod.getValues().enabled
            ? messages.adminAccount.confirm.enableDescription
            : messages.adminAccount.confirm.disableDescription,
          {
            accountId: switchEnabledMethod.getValues().accountId,
          }
        )}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
