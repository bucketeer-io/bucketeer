import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { parse } from 'query-string';

import { urls } from '../config';
import { getState, setState } from '../cookie';
import * as authGrpc from '../grpc/auth';
import {
  ExchangeTokenRequest,
  GetAuthCodeURLRequest,
  RefreshTokenRequest,
} from '../proto/auth/service_pb';
import { Token } from '../proto/auth/token_pb';
import {
  clearToken as clearTokenFromStorage,
  getToken as getTokenFromStorage,
  setToken,
} from '../storage/token';

const MODULE_NAME = 'auth';

export const exchangeTokenFromUrl = createAsyncThunk<Token.AsObject, string>(
  `${MODULE_NAME}/exchangeTokenFromUrl`,
  async (query) => {
    const { code, state } = parse(query);
    const stateFromCookie = getState();
    if (!!code && state === stateFromCookie) {
      if (typeof code === 'string') {
        const request = new ExchangeTokenRequest();
        request.setCode(code);
        request.setRedirectUrl(urls.AUTH_REDIRECT);
        const result = await authGrpc.exchangeToken(request);
        return result.response.getToken().toObject();
      }
    }
    throw new Error('exchange token failed.');
  }
);

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
        refreshToken({
          token: token.refreshToken,
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
    thunkAPI.dispatch(getAuthCodeURL({ state }));
  }
);

interface GetAuthCodeURLParams {
  state: string;
}

export const getAuthCodeURL = createAsyncThunk<string, GetAuthCodeURLParams>(
  `${MODULE_NAME}/getAuthCodeURL`,
  async (params) => {
    const request = new GetAuthCodeURLRequest();
    request.setState(params.state);
    request.setRedirectUrl(urls.AUTH_REDIRECT);
    const result = await authGrpc.getAuthCodeURL(request);
    return result.response.getUrl();
  }
);

interface RefreshTokenParams {
  token: string;
}

export const refreshToken = createAsyncThunk<
  Token.AsObject,
  RefreshTokenParams
>(`${MODULE_NAME}/refreshToken`, async (params) => {
  const request = new RefreshTokenRequest();
  request.setRefreshToken(params.token);
  request.setRedirectUrl(urls.AUTH_REDIRECT);
  const result = await authGrpc.refreshToken(request);
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
  loading: false,
};

export type AuthState = typeof initialState;

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearToken(state) {
      clearTokenFromStorage();
      return initialState;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getAuthCodeURL.rejected, (state) => {
        state.loading = false;
        // retry
        location.reload();
      })
      .addCase(getAuthCodeURL.pending, (state) => {
        state.loading = true;
      })
      .addCase(getAuthCodeURL.fulfilled, (state, action) => {
        window.location.href = action.payload;
        state.loading = false;
      })
      .addCase(refreshToken.rejected, (state) => {
        state.loading = false;
        clearTokenFromStorage();
        location.reload();
      })
      .addCase(refreshToken.pending, (state) => {
        state.loading = true;
      })
      .addCase(refreshToken.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
        location.reload();
      })
      .addCase(exchangeTokenFromUrl.rejected, (state) => {
        state.loading = false;
      })
      .addCase(exchangeTokenFromUrl.pending, (state) => {
        state.loading = true;
      })
      .addCase(exchangeTokenFromUrl.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
      });
  },
});

export const { clearToken } = authSlice.actions;
