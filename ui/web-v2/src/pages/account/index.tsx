import { Option } from '../../components/Select';
import { intl } from '../../lang';
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
  useParams
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
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectAccountById,
  listAccounts,
  disableAccount,
  enableAccount,
  updateAccount,
  getAccount,
  createAccount,
  OrderBy,
  OrderDirection
} from '../../modules/accounts';
import { useCurrentEnvironment } from '../../modules/me';
import { AccountV2 } from '../../proto/account/account_pb';
import { ListAccountsV2Request } from '../../proto/account/service_pb';
import { AppDispatch } from '../../store';
import { AccountSortOption, isAccountSortOption } from '../../types/account';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC
} from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';
import { listTags } from '../../modules/tags';
import { ListTagsRequest } from '../../proto/tag/service_pb';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: AccountSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListAccountsV2Request.OrderBy.CREATED_AT,
        orderDirection: ListAccountsV2Request.OrderDirection.ASC
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListAccountsV2Request.OrderBy.CREATED_AT,
        orderDirection: ListAccountsV2Request.OrderDirection.DESC
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListAccountsV2Request.OrderBy.EMAIL,
        orderDirection: ListAccountsV2Request.OrderDirection.ASC
      };
    default:
      return {
        orderBy: ListAccountsV2Request.OrderBy.EMAIL,
        orderDirection: ListAccountsV2Request.OrderDirection.DESC
      };
  }
};

// TODO: Remove this when the console 3.0 is ready
enum AccountRoleV1 {
  VIEWER = 0,
  EDITOR = 1,
  OWNER = 2
}

// TODO: Remove this when the console 3.0 is ready
export const getRoleListV1 = (): Option[] => {
  return [
    {
      value: AccountRoleV1.VIEWER.toString(),
      label: intl.formatMessage(messages.account.role.viewer)
    },
    {
      value: AccountRoleV1.EDITOR.toString(),
      label: intl.formatMessage(messages.account.role.editor)
    },
    {
      value: AccountRoleV1.OWNER.toString(),
      label: intl.formatMessage(messages.account.role.owner)
    }
  ];
};

// TODO: Remove this when the console 3.0 is ready
export const getRoleV1 = (
  orgRole: AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap],
  envRole: AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap]
): Option => {
  // If it's a editor
  if (
    envRole == AccountV2.Role.Environment.ENVIRONMENT_EDITOR &&
    orgRole == AccountV2.Role.Organization.ORGANIZATION_MEMBER
  ) {
    return {
      value: AccountRoleV1.EDITOR.toString(),
      label: intl.formatMessage(messages.account.role.editor)
    };
    // If it's an admin
  } else if (
    envRole == AccountV2.Role.Environment.ENVIRONMENT_EDITOR &&
    orgRole == AccountV2.Role.Organization.ORGANIZATION_ADMIN
  ) {
    return {
      value: AccountRoleV1.OWNER.toString(),
      label: intl.formatMessage(messages.account.role.owner)
    };
    // If it's an onwer
  } else if (
    envRole == AccountV2.Role.Environment.ENVIRONMENT_EDITOR &&
    orgRole == AccountV2.Role.Organization.ORGANIZATION_OWNER
  ) {
    return {
      value: AccountRoleV1.OWNER.toString(),
      label: intl.formatMessage(messages.account.role.owner)
    };
  }
  // Anything else returns viewer
  return {
    value: AccountRoleV1.VIEWER.toString(),
    label: intl.formatMessage(messages.account.role.viewer)
  };
};

// TODO: Remove this when the console 3.0 is ready
export const convertToAccountV2Role = (
  roleV1: AccountRoleV1
): [
  AccountV2.Role.OrganizationMap[keyof AccountV2.Role.OrganizationMap],
  AccountV2.Role.EnvironmentMap[keyof AccountV2.Role.EnvironmentMap]
] => {
  if (roleV1 == AccountRoleV1.VIEWER) {
    return [
      AccountV2.Role.Organization.ORGANIZATION_MEMBER,
      AccountV2.Role.Environment.ENVIRONMENT_VIEWER
    ];
  }
  if (roleV1 == AccountRoleV1.EDITOR) {
    return [
      AccountV2.Role.Organization.ORGANIZATION_MEMBER,
      AccountV2.Role.Environment.ENVIRONMENT_EDITOR
    ];
  }
  return [
    AccountV2.Role.Organization.ORGANIZATION_ADMIN,
    AccountV2.Role.Environment.ENVIRONMENT_EDITOR
  ];
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
  const [open, setOpen] = useState(isNew || isUpdate);
  const [account] = useSelector<
    AppState,
    [AccountV2.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectAccountById(state.accounts, accountId),
      state.accounts.getAccountError
    ],
    shallowEqual
  );
  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      history.replace(
        `${url}?${stringifySearchParams({
          ...options
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
      const tags =
        options && Array.isArray(options.tagIds)
          ? options.tagIds
          : typeof options?.tagIds === 'string' && options?.tagIds.length > 0
            ? [options.tagIds]
            : [];

      dispatch(
        listAccounts({
          environmentId: currentEnvironment.id,
          organizationId: currentEnvironment.organizationId,
          pageSize: ACCOUNT_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          disabled: disabled,
          role: role,
          tags
        })
      );
    },
    [dispatch]
  );

  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

  const switchEnabledMethod = useForm({
    defaultValues: {
      accountId: '',
      enabled: false
    },
    mode: 'onChange'
  });

  const {
    handleSubmit: switchEnableHandleSubmit,
    setValue: switchEnabledSetValue
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
              organizationId: currentEnvironment.organizationId,
              environmentId: currentEnvironment.id,
              email: data.accountId
            });
          }
          return disableAccount({
            organizationId: currentEnvironment.organizationId,
            environmentId: currentEnvironment.id,
            email: data.accountId
          });
        })()
      ).then(() => {
        setIsConfirmDialogOpen(false);
        dispatch(
          getAccount({
            organizationId: currentEnvironment.organizationId,
            email: data.accountId
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
      search: location.search
    });
  }, [setOpen, history, location]);

  const handleOpenUpdate = useCallback(
    (a: AccountV2.AsObject) => {
      setOpen(true);
      const envRole = a.environmentRolesList.find(
        (e) => e.environmentId === currentEnvironment.id
      );
      resetUpdate({
        name: a.name,
        email: a.email,
        role: getRoleV1(a.organizationRole, envRole.role).value,
        tags: a.tagsList
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}/${a.email}`,
        search: location.search
      });
    },
    [setOpen, history, account, location]
  );

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      email: null,
      role: null,
      tags: []
    },
    mode: 'onChange'
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const updateMethod = useForm({
    resolver: yupResolver(updateFormSchema),
    mode: 'onChange'
  });
  const {
    handleSubmit: handleUpdateSubmit,
    reset: resetUpdate,
    formState: { dirtyFields }
  } = updateMethod;

  const handleClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}`,
      search: location.search
    });
  }, [setOpen, history, location, resetAdd, resetUpdate]);

  const handleAdd = useCallback(
    async (data) => {
      const [orgRole, envRole] = convertToAccountV2Role(data.role);
      dispatch(
        createAccount({
          organizationId: currentEnvironment.organizationId,
          name: data.name,
          email: data.email,
          environmentId: currentEnvironment.id,
          environmentRole: envRole,
          organizationRole: orgRole,
          tagsList: data.tags
        })
      ).then(() => {
        resetAdd();
        setOpen(false);
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_ACCOUNTS}`
        );
        updateAccountList(null, 1);
        fetchListTags();
      });
    },
    [dispatch]
  );

  const handleUpdate = useCallback(
    async (data) => {
      let name: string;
      let tagsList: string[];

      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.tags) {
        tagsList = data.tags;
      }
      const [orgRole, envRole] = convertToAccountV2Role(data.role);
      dispatch(
        updateAccount({
          organizationId: currentEnvironment.organizationId,
          environmentId: currentEnvironment.id,
          name: name,
          email: accountId,
          environmentRole: envRole,
          organizationRole: orgRole,
          tagsList: tagsList
        })
      ).then(() => {
        dispatch(
          getAccount({
            organizationId: currentEnvironment.organizationId,
            email: accountId
          })
        );
        handleClose();
        fetchListTags();
      });
    },
    [dispatch, accountId, dirtyFields]
  );

  const fetchListTags = useCallback(() => {
    dispatch(
      listTags({
        environmentId: currentEnvironment.id,
        pageSize: 0,
        cursor: '',
        orderBy: ListTagsRequest.OrderBy.DEFAULT,
        orderDirection: ListTagsRequest.OrderDirection.ASC,
        searchKeyword: null
      })
    );
  }, [dispatch]);

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
    fetchListTags();
  }, [updateAccountList]);

  useEffect(() => {
    if (isUpdate) {
      dispatch(
        getAccount({
          organizationId: currentEnvironment.organizationId,
          email: accountId
        })
      ).then((e) => {
        const payload = e.payload as AccountV2.AsObject;
        const envRole = payload.environmentRolesList.find(
          (e) => e.environmentId === currentEnvironment.id
        );
        resetUpdate({
          name: payload.name,
          email: payload.email,
          role: getRoleV1(payload.organizationRole, envRole.role).value,
          tags: payload.tagsList
        });
      });
    }
  }, []);

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
