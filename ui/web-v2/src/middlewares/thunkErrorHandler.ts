import { Middleware, AnyAction, MiddlewareAPI } from 'redux';

import { AppThunk } from '../modules';
import { addToast } from '../modules/toasts';
import { AppDispatch } from '../store';
import { clearToken } from '../modules/auth';
import { clearOrganizationId } from '../storage/organizationId';
import { clearMe } from '../modules/me';
import { PAGE_PATH_ROOT } from '../constants/routing';
import { clearCurrentEnvironmentId } from '../storage/environment';
import { history } from '../history';

export const UNAUTHENTICATED_ERROR = 'UNAUTHENTICATED_ERROR';

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
        if (action.error.message === UNAUTHENTICATED_ERROR) {
          history.push(PAGE_PATH_ROOT);
          dispatch(clearToken());
          dispatch(clearMe());
          clearOrganizationId();
          clearCurrentEnvironmentId();
        } else {
          dispatch(
            addToast({ message: action.error.message, severity: 'error' })
          );
        }
      }
    }

    return res;
  };
