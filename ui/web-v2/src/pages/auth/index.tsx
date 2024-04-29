import { FC, useEffect, memo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLocation, useHistory } from 'react-router-dom';

import { AppState } from '../../modules';
import { exchangeBucketeerTokenFromUrl, hasToken } from '../../modules/auth';
import { AppDispatch } from '../../store';

export const AuthCallbackPage: FC = memo(() => {
  const history = useHistory();
  const dispatch = useDispatch<AppDispatch>();
  const location = useLocation();
  const loading = useSelector<AppState, boolean>((state) => state.auth.loading);

  useEffect(() => {
    const query = location.search;
    dispatch(exchangeBucketeerTokenFromUrl(query));
  }, [dispatch]);

  useEffect(() => {
    if (hasToken()) {
      history.push('/');
    }
  }, [loading]);

  return <div className="spinner mx-auto mt-4" />;
});
