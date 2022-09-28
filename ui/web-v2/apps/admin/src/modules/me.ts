import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { shallowEqual, useSelector } from 'react-redux';

import { getMe } from '../grpc/account';
import { Account, EnvironmentRole } from '../proto/account/account_pb';
import { GetMeRequest } from '../proto/account/service_pb';
import { Environment } from '../proto/environment/environment_pb';
import {
  getCurrentEnvironmentId,
  setCurrentEnvironmentId,
} from '../storage/environment';

import { AppState } from '.';

const MODULE_NAME = 'me';

export interface Me {
  isAdmin: boolean;
  environmentRoles: Array<EnvironmentRole.AsObject>;
  isLogin: boolean;
}

export type MeState = Me;

export const fetchMe = createAsyncThunk<Me>('me/fetch', async () => {
  const res = await getMe(new GetMeRequest());
  return {
    isAdmin: res.response.getIsAdmin(),
    environmentRoles: res.response.toObject().environmentRolesList,
    isLogin: true,
  };
});

export const meSlice = createSlice({
  name: MODULE_NAME,
  initialState: { isAdmin: false, environmentRoles: [], isLogin: false },
  reducers: {
    clearMe(state) {
      return { isAdmin: false, environmentRoles: [], isLogin: false };
    },
    setCurrentEnvironment(state, action: PayloadAction<string>) {
      setCurrentEnvironmentId(action.payload);
      return state;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchMe.fulfilled, (_, action) => {
        const curEnvId = getCurrentEnvironmentId()
          ? getCurrentEnvironmentId()
          : action.payload.environmentRoles[0].environment.id;
        setCurrentEnvironmentId(curEnvId);
        return action.payload;
      })
      .addCase(fetchMe.rejected, () => {
        return { isAdmin: false, environmentRoles: [], isLogin: false };
      });
  },
});

export const useMe = (): MeState =>
  useSelector<AppState, MeState>((state) => state.me);

const currentEnvironmentRole = (state: AppState): EnvironmentRole.AsObject => {
  if ('environmentRoles' in state.me) {
    const curEnvId = getCurrentEnvironmentId()
      ? getCurrentEnvironmentId()
      : state.me.environmentRoles[0].environment.id;
    const envRole = state.me.environmentRoles.find(
      (environmentRole) => environmentRole.environment.id === curEnvId
    );
    return envRole;
  }
  return new EnvironmentRole().toObject();
};

export const useCurrentEnvironmentRole = (): EnvironmentRole.AsObject => {
  return useSelector<AppState, EnvironmentRole.AsObject>(
    currentEnvironmentRole,
    shallowEqual
  );
};

export const useCurrentEnvironment = (): Environment.AsObject => {
  return useSelector<AppState, Environment.AsObject>((state: AppState) => {
    return currentEnvironmentRole(state).environment;
  }, shallowEqual);
};

export const useEnvironments = (): Array<Environment.AsObject> => {
  return useSelector<AppState, Array<Environment.AsObject>>((state) => {
    if ('environmentRoles' in state.me) {
      return state.me.environmentRoles.map(
        (environmentRole) => environmentRole.environment
      );
    }
    return [];
  }, shallowEqual);
};

export const useIsEditable = (): boolean => {
  return useSelector<AppState, boolean>((state: AppState) => {
    if (state.me.isAdmin) return true;
    const envRole = currentEnvironmentRole(state);
    return (
      envRole.role === Account.Role.EDITOR ||
      envRole.role === Account.Role.OWNER
    );
  }, shallowEqual);
};

export const useIsOwner = (): boolean => {
  return useSelector<AppState, boolean>((state: AppState) => {
    if (state.me.isAdmin) return true;
    const envRole = currentEnvironmentRole(state);
    return envRole.role === Account.Role.OWNER;
  }, shallowEqual);
};

export const { clearMe, setCurrentEnvironment } = meSlice.actions;
