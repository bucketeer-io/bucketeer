import { type Nullable } from 'option-t/nullable';
import { AuthToken } from '@types';

const KEY = 'demo_token';

export const getDemoTokenStorage = (): Nullable<AuthToken> => {
  try {
    const tokenStr = window.localStorage.getItem(KEY);
    if (tokenStr) {
      const token = JSON.parse(tokenStr);
      return token;
    }
  } catch (error) {
    console.error(error);
  }
  return null;
};

export const setDemoTokenStorage = (token: AuthToken): void => {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(token));
  } catch (error) {
    console.error(error);
  }
};

export const clearDemoTokenStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
