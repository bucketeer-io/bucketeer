import { type Nullable } from 'option-t/nullable';
import { AuthToken } from '@types';

export const AUTH_TOKEN_KEY = 'auth_token';

export const getTokenStorage = (): Nullable<AuthToken> => {
  try {
    const tokenStr = window.localStorage.getItem(AUTH_TOKEN_KEY);
    if (tokenStr) {
      const token = JSON.parse(tokenStr);
      return token;
    }
  } catch (error) {
    console.error(error);
  }
  return null;
};

export const setTokenStorage = (token: AuthToken): void => {
  try {
    window.localStorage.setItem(AUTH_TOKEN_KEY, JSON.stringify(token));
  } catch (error) {
    console.error(error);
  }
};

export const clearTokenStorage = (): void => {
  try {
    window.localStorage.removeItem(AUTH_TOKEN_KEY);
  } catch (error) {
    console.error(error);
  }
};
