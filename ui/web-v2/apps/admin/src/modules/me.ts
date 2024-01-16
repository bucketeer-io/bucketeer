import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { shallowEqual, useSelector } from 'react-redux';

import { getMe } from '../grpc/account';
import { AccountV2, ConsoleAccount } from '../proto/account/account_pb';
import { GetMeRequest } from '../proto/account/service_pb';
import { EnvironmentV2 } from '../proto/environment/environment_pb';
import {
  getCurrentEnvironmentId,
  setCurrentEnvironmentId,
} from '../storage/environment';

import { AppState } from '.';

const MODULE_NAME = 'me';

export interface Me {
  isAdmin: boolean;
  isLogin: boolean;
  consoleAccount: ConsoleAccount.AsObject;
}

export type MeState = Me;

export interface FetchMeParams {
  organizationId: string;
}

export const fetchMe = createAsyncThunk<
  Me,
  FetchMeParams | undefined,
  { state: AppState }
>('me/fetch', async (params) => {
  const request = new GetMeRequest();
  request.setOrganizationId(params.organizationId);
  const res = await getMe(request);
  return {
    isAdmin: res.response.toObject().account.isSystemAdmin,
    isLogin: true,
    consoleAccount: res.response.toObject().account,
  };
});

export const meSlice = createSlice({
  name: MODULE_NAME,
  initialState: {
    isAdmin: false,
    isLogin: false,
    consoleAccount: null,
  } as MeState,
  reducers: {
    clearMe(state) {
      return { isAdmin: false, isLogin: false, consoleAccount: null };
    },
    setCurrentEnvironment(state, action: PayloadAction<string>) {
      setCurrentEnvironmentId(action.payload);
      return state;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchMe.fulfilled, (_, action) => {
        // Since some old environments have empty id, so accept empty and reject only null or undefined
        const curEnvId =
          getCurrentEnvironmentId() != (null || undefined)
            ? getCurrentEnvironmentId()
            : action.payload.consoleAccount.environmentRolesList[0].environment
                .id;
        setCurrentEnvironmentId(curEnvId);
        return action.payload;
      })
      .addCase(fetchMe.rejected, () => {
        return { isAdmin: false, isLogin: false, consoleAccount: null };
      });
  },
});

export const useMe = (): MeState =>
  useSelector<AppState, MeState>((state) => state.me);

const currentEnvironmentRole = (
  state: AppState
): ConsoleAccount.EnvironmentRole.AsObject => {
  if ('consoleAccount' in state.me) {
    // Since some old environments have empty id, so accept empty and reject only null or undefined
    const curEnvId =
      getCurrentEnvironmentId() != (null || undefined)
        ? getCurrentEnvironmentId()
        : state.me.consoleAccount.environmentRolesList[0].environment.id;
    let curEnvRole = state.me.consoleAccount.environmentRolesList.find(
      (environmentRole) => environmentRole.environment.id === curEnvId
    );
    if (!curEnvRole) {
      curEnvRole = state.me.consoleAccount.environmentRolesList[0];
    }
    return curEnvRole;
  }
};

export const useCurrentEnvironment = (): EnvironmentV2.AsObject => {
  return useSelector<AppState, EnvironmentV2.AsObject>((state: AppState) => {
    return currentEnvironmentRole(state).environment;
  }, shallowEqual);
};

export const useEnvironments = (): Array<EnvironmentV2.AsObject> => {
  return useSelector<AppState, Array<EnvironmentV2.AsObject>>((state) => {
    if ('consoleAccount' in state.me) {
      return state.me.consoleAccount.environmentRolesList.map(
        (environmentRole) => environmentRole.environment
      );
    }
    return [];
  }, shallowEqual);
};

export const useIsEditable = (): boolean => {
  return useSelector<AppState, boolean>((state: AppState) => {
    if (state.me.consoleAccount.isSystemAdmin) return true;
    const envRole = currentEnvironmentRole(state);
    return envRole.role === AccountV2.Role.Environment.ENVIRONMENT_EDITOR;
  }, shallowEqual);
};

export const useIsOwner = (): boolean => {
  return useSelector<AppState, boolean>((state: AppState) => {
    if (state.me.consoleAccount.isSystemAdmin) return true;
    // TODO: Once we replace the console design, we should check if the user is owner/admin of the organization
    // We check if the user's environment role is Editor until the organization is introduced into the console design.
    // return (
    //   state.me.consoleAccount.organizationRole === AccountV2.Role.Organization.ORGANIZATION_OWNER ||
    //   state.me.consoleAccount.organizationRole === AccountV2.Role.Organization.ORGANIZATION_ADMIN
    // )
    const envRole = currentEnvironmentRole(state);
    return envRole.role === AccountV2.Role.Environment.ENVIRONMENT_EDITOR;
  }, shallowEqual);
};

export const { clearMe, setCurrentEnvironment } = meSlice.actions;
