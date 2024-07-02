import { Middleware, AnyAction, MiddlewareAPI } from 'redux';

import { AppThunk } from '../modules';
import { addToast } from '../modules/toasts';
import { AppDispatch } from '../store';
import { getToken } from '../storage/token';
import { refreshBucketeerToken } from '../modules/auth';

export const TOKEN_IS_EXPIRED = 'token is expired';

function isPlainAction(action: AppThunk | AnyAction): action is AnyAction {
  return typeof action !== 'function';
}

export const thunkErrorHandler: Middleware =
  ({ dispatch }: MiddlewareAPI<AppDispatch>) =>
  (next) =>
  async (action) => {
    let res;
    try {
      res = await next(action);
    } catch (err) {
      if (process.env.NODE_ENV === 'development') {
        console.error(err);
      }

      if (err && err.code && err.message) {
        dispatch(addToast({ message: err.message, severity: 'error' }));
      } else {
        throw err;
      }
    }

    if (isPlainAction(action)) {
      if (action.type.includes('rejected')) {
        if (action.error.message === TOKEN_IS_EXPIRED) {
          const token = getToken();
          dispatch(refreshBucketeerToken({ token: token.refreshToken }));
        }
        dispatch(
          addToast({ message: action.error.message, severity: 'error' })
        );
      }
    }

    return res;
  };
