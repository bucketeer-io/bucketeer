import { Nullable } from 'option-t/lib/Nullable';

import { Token } from '../proto/auth/token_pb';

const KEY = 'auth_token';

let cache: Nullable<Token.AsObject> = null;

export const getToken = (): Nullable<Token.AsObject> => {
  if (cache) {
    return cache;
  }
  try {
    const tokenStr = window.localStorage.getItem(KEY);
    if (tokenStr) {
      const token = JSON.parse(tokenStr);
      cache = token;
      return token;
    }
  } catch (error) {
    // ignore
  }
  return null;
};

export const setToken = (token: Token.AsObject): void => {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(token));
  } catch (error) {
    // ignore
  }
  cache = token;
};

export const clearToken = (): void => {
  cache = null;
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    // ignore
  }
};
