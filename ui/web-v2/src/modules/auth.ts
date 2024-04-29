import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { parse } from 'query-string';

import { urls } from '../config';
import { getState, setState } from '../cookie';
import * as authGrpc from '../grpc/auth';
import {
  GetAuthenticationURLRequest,
  ExchangeBucketeerTokenRequest,
  RefreshBucketeerTokenRequest
} from '../proto/auth/service_pb';
import { Token } from '../proto/auth/token_pb';
import {
  clearToken as clearTokenFromStorage,
  getToken as getTokenFromStorage,
  setToken
} from '../storage/token';

const MODULE_NAME = 'auth';

export const exchangeBucketeerTokenFromUrl = createAsyncThunk<
  Token.AsObject,
  string
>(`${MODULE_NAME}/exchangeBucketeerTokenFromUrl`, async (query) => {
  const { code, state } = parse(query);
  const stateFromCookie = getState();
  if (!!code && state === stateFromCookie) {
    if (typeof code === 'string') {
      const request = new ExchangeBucketeerTokenRequest();
      request.setCode(code);
      request.setRedirectUrl(urls.AUTH_REDIRECT);
      request.setType(2);
      const result = await authGrpc.exchangeBucketeerToken(request);
      return result.response.getToken().toObject();
    }
  }
  throw new Error('exchange token failed.');
});

export const setupAuthToken = createAsyncThunk<void>(
  `${MODULE_NAME}/setupAuthToken`,
  async (_, thunkAPI) => {
    const token = getTokenFromStorage();
    if (!token || !token.accessToken) {
      thunkAPI.dispatch(redirectToAuthUrl());
      throw new Error('token not found.');
    }

    if (isExpiredToken(token.expiry)) {
      thunkAPI.dispatch(
        refreshBucketeerToken({
          token: token.refreshToken
        })
      );
      return;
    }
    return;
  }
);

export const redirectToAuthUrl = createAsyncThunk<void>(
  `${MODULE_NAME}/redirecttoAuthUrl`,
  async (_, thunkAPI) => {
    const state = `${Date.now()}`;
    setState(state);
    thunkAPI.dispatch(getAuthenticationURL({ state }));
  }
);

interface GetAuthenticationURLParams {
  state: string;
}

export const getAuthenticationURL = createAsyncThunk<
  string,
  GetAuthenticationURLParams
>(`${MODULE_NAME}/getAuthenticationURL`, async (params) => {
  const request = new GetAuthenticationURLRequest();
  request.setState(params.state);
  request.setRedirectUrl(urls.AUTH_REDIRECT);
  request.setType(2); // google auth type
  const result = await authGrpc.getAuthenticationURL(request);
  return result.response.getUrl();
});

interface RefreshBucketeerTokenParams {
  token: string;
}

export const refreshBucketeerToken = createAsyncThunk<
  Token.AsObject,
  RefreshBucketeerTokenParams
>(`${MODULE_NAME}/refreshBucketeerToken`, async (params) => {
  const request = new RefreshBucketeerTokenRequest();
  request.setRefreshToken(params.token);
  request.setRedirectUrl(urls.AUTH_REDIRECT);
  request.setType(2);
  const result = await authGrpc.refreshBucketeerToken(request);
  return result.response.getToken().toObject();
});

const isExpiredToken = (expiry: number): boolean => {
  const now = Number(Date.now() / 1000);
  return now > expiry;
};

export const hasToken = (): boolean => {
  const token = getTokenFromStorage();

  if (!token || !token.accessToken) {
    return false;
  }
  if (isExpiredToken(token.expiry)) {
    return false;
  }
  return true;
};

const initialState = {
  loading: false
};

export type AuthState = typeof initialState;

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearToken() {
      clearTokenFromStorage();
      return initialState;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(getAuthenticationURL.rejected, (state) => {
        state.loading = false;
        // retry
        location.reload();
      })
      .addCase(getAuthenticationURL.pending, (state) => {
        state.loading = true;
      })
      .addCase(getAuthenticationURL.fulfilled, (state, action) => {
        window.location.href = action.payload;
        state.loading = false;
      })
      .addCase(refreshBucketeerToken.rejected, (state) => {
        state.loading = false;
        clearTokenFromStorage();
        location.reload();
      })
      .addCase(refreshBucketeerToken.pending, (state) => {
        state.loading = true;
      })
      .addCase(refreshBucketeerToken.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
        location.reload();
      })
      .addCase(exchangeBucketeerTokenFromUrl.rejected, (state) => {
        state.loading = false;
      })
      .addCase(exchangeBucketeerTokenFromUrl.pending, (state) => {
        state.loading = true;
      })
      .addCase(exchangeBucketeerTokenFromUrl.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
      });
  }
});

export const { clearToken } = authSlice.actions;
