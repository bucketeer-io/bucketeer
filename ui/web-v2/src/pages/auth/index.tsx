import { FC, useEffect, memo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLocation, useHistory } from 'react-router-dom';

import { AppState } from '../../modules';
import { exchangeTokenFromUrl, hasToken } from '../../modules/auth';
import { AppDispatch } from '../../store';

export const AuthCallbackPage: FC = memo(() => {
  const history = useHistory();
  const dispatch = useDispatch<AppDispatch>();
  const location = useLocation();
  const loading = useSelector<AppState, boolean>((state) => state.auth.loading);

  useEffect(() => {
    const query = location.search;
    dispatch(exchangeTokenFromUrl(query));
  }, [dispatch]);

  useEffect(() => {
    if (hasToken()) {
      history.push('/');
    }
  }, [loading]);

  return (
    <div className="flex justify-center items-center p-3">
      <div className="w-6 h-6 border-4 border-t-primary rounded-full animate-spin" />
    </div>
  );
});
