import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { parse } from 'query-string';

import { urls } from '../config';
import { getState, setState } from '../cookie';
import * as authGrpc from '../grpc/auth';
import {
  GetAuthenticationURLRequest,
  ExchangeTokenRequest,
  RefreshTokenRequest,
  SignInRequest,
  SwitchOrganizationRequest
} from '../proto/auth/service_pb';
import { Token } from '../proto/auth/token_pb';
import {
  clearToken as clearTokenFromStorage,
  setToken
} from '../storage/token';
import { PAGE_PATH_ROOT } from '../constants/routing';
import { history } from '../history';

const MODULE_NAME = 'auth';

interface SwitchOrganizationParams {
  accessToken: string;
  organizationId: string;
}
export const switchOrganization = createAsyncThunk<
  Token.AsObject,
  SwitchOrganizationParams
>(
  `${MODULE_NAME}/switchOrganization`,
  async ({ accessToken, organizationId }) => {
    try {
      const request = new SwitchOrganizationRequest();
      request.setAccessToken(accessToken);
      request.setOrganizationId(organizationId);

      const result = await authGrpc.switchOrganization(request);
      return result.response.getToken().toObject();
    } catch (error) {
      throw new Error('switch organization failed.');
    }
  }
);

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
        request.setType(2);
        const result = await authGrpc.exchangeToken(request);
        return result.response.getToken().toObject();
      }
    }
    throw new Error('exchange token failed.');
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

interface RefreshTokenParams {
  token: string;
}

export const refreshToken = createAsyncThunk<
  Token.AsObject,
  RefreshTokenParams
>(`${MODULE_NAME}/refreshToken`, async (params) => {
  const request = new RefreshTokenRequest();
  request.setRefreshToken(params.token);
  const result = await authGrpc.refreshToken(request);
  return result.response.getToken().toObject();
});

interface SignInParams {
  email: string;
  password: string;
}

export const signIn = createAsyncThunk<Token.AsObject, SignInParams>(
  `${MODULE_NAME}/signIn`,
  async (params) => {
    const request = new SignInRequest();
    request.setEmail(params.email);
    request.setPassword(params.password);
    const result = await authGrpc.signIn(request);
    return result.response.getToken().toObject();
  }
);

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
        clearTokenFromStorage();
        history.push(PAGE_PATH_ROOT);
      })
      .addCase(getAuthenticationURL.pending, (state) => {
        state.loading = true;
      })
      .addCase(getAuthenticationURL.fulfilled, (state, action) => {
        window.location.href = action.payload;
        state.loading = false;
      })
      .addCase(refreshToken.rejected, (state) => {
        state.loading = false;
        clearTokenFromStorage();
        history.push(PAGE_PATH_ROOT);
      })
      .addCase(refreshToken.pending, (state) => {
        state.loading = true;
      })
      .addCase(refreshToken.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
      })
      .addCase(exchangeTokenFromUrl.rejected, (state) => {
        state.loading = false;
        clearTokenFromStorage();
        history.push(PAGE_PATH_ROOT);
      })
      .addCase(exchangeTokenFromUrl.pending, (state) => {
        state.loading = true;
      })
      .addCase(exchangeTokenFromUrl.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
      })
      .addCase(signIn.rejected, (state) => {
        state.loading = false;
        clearTokenFromStorage();
      })
      .addCase(signIn.pending, (state) => {
        state.loading = true;
      })
      .addCase(signIn.fulfilled, (state, action) => {
        setToken(action.payload);
        state.loading = false;
      });
  }
});

export const { clearToken } = authSlice.actions;
