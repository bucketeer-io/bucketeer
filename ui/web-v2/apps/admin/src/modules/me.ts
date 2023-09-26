import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { shallowEqual, useSelector } from 'react-redux';

import { getMeV2 } from '../grpc/account';
import { Account, EnvironmentRoleV2 } from '../proto/account/account_pb';
import { GetMeV2Request } from '../proto/account/service_pb';
import { EnvironmentV2 } from '../proto/environment/environment_pb';
import {
  getCurrentEnvironmentId,
  setCurrentEnvironmentId,
} from '../storage/environment';

import { AppState } from '.';

const MODULE_NAME = 'me';

export interface Me {
  isAdmin: boolean;
  environmentRoles: Array<EnvironmentRoleV2.AsObject>;
  isLogin: boolean;
}

export type MeState = Me;

export const fetchMe = createAsyncThunk<Me>('me/fetch', async () => {
  const res = await getMeV2(new GetMeV2Request());
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
        // Since some old environments have empty id, so accept empty and reject only null or undefined
        const curEnvId = (getCurrentEnvironmentId() != (null || undefined))
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

const currentEnvironmentRole = (state: AppState): EnvironmentRoleV2.AsObject => {
  if ('environmentRoles' in state.me) {
    // Since some old environments have empty id, so accept empty and reject only null or undefined
    const curEnvId = (getCurrentEnvironmentId() != (null || undefined))
      ? getCurrentEnvironmentId()
      : state.me.environmentRoles[0].environment.id;
    const envRole = state.me.environmentRoles.find(
      (environmentRole) => environmentRole.environment.id === curEnvId
    );
    if (!envRole) {
      return state.me.environmentRoles[0];
    }
    return envRole;
  }
  return new EnvironmentRoleV2().toObject();
};

export const useCurrentEnvironmentRole = (): EnvironmentRoleV2.AsObject => {
  return useSelector<AppState, EnvironmentRoleV2.AsObject>(
    currentEnvironmentRole,
    shallowEqual
  );
};

export const useCurrentEnvironment = (): EnvironmentV2.AsObject => {
  return useSelector<AppState, EnvironmentV2.AsObject>((state: AppState) => {
    return currentEnvironmentRole(state).environment;
  }, shallowEqual);
};

export const useEnvironments = (): Array<EnvironmentV2.AsObject> => {
  return useSelector<AppState, Array<EnvironmentV2.AsObject>>((state) => {
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
